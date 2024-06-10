// SPDX-License-Identifier: AGPL-3.0-only

package ingester

import (
	"context"
	"flag"
	"time"

	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/grpcutil"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"
	"google.golang.org/grpc/codes"

	"github.com/grafana/mimir/pkg/mimirpb"
)

const (
	circuitBreakerResultSuccess      = "success"
	circuitBreakerResultError        = "error"
	circuitBreakerResultOpen         = "circuit_breaker_open"
	circuitBreakerDefaultPushTimeout = 2 * time.Second
	circuitBreakerDefaultReadTimeout = 30 * time.Second
	circuitBreakerPushRequestType    = "push"
	circuitBreakerReadRequestType    = "read"
)

type circuitBreakerMetrics struct {
	circuitBreakerTransitions *prometheus.CounterVec
	circuitBreakerResults     *prometheus.CounterVec
}

func newCircuitBreakerMetrics(r prometheus.Registerer, currentStateFn func(string) circuitbreaker.State, requestTypes []string) *circuitBreakerMetrics {
	cbMetrics := &circuitBreakerMetrics{
		circuitBreakerTransitions: promauto.With(r).NewCounterVec(prometheus.CounterOpts{
			Name: "cortex_ingester_circuit_breaker_transitions_total",
			Help: "Number of times the circuit breaker has entered a state.",
		}, []string{"request_type", "state"}),
		circuitBreakerResults: promauto.With(r).NewCounterVec(prometheus.CounterOpts{
			Name: "cortex_ingester_circuit_breaker_results_total",
			Help: "Results of executing requests via the circuit breaker.",
		}, []string{"request_type", "result"}),
	}
	circuitBreakerCurrentStateGaugeFn := func(requestType string, state circuitbreaker.State) prometheus.GaugeFunc {
		return promauto.With(r).NewGaugeFunc(prometheus.GaugeOpts{
			Name:        "cortex_ingester_circuit_breaker_current_state",
			Help:        "Boolean set to 1 whenever the circuit breaker is in a state corresponding to the label name.",
			ConstLabels: map[string]string{"request_type": requestType, "state": state.String()},
		}, func() float64 {
			if currentStateFn(requestType) == state {
				return 1
			}
			return 0
		})
	}
	for _, requestType := range requestTypes {
		for _, s := range []circuitbreaker.State{circuitbreaker.OpenState, circuitbreaker.HalfOpenState, circuitbreaker.ClosedState} {
			circuitBreakerCurrentStateGaugeFn(requestType, s)
			// We initialize all possible states for the given requestType and the circuitBreakerTransitions metrics
			cbMetrics.circuitBreakerTransitions.WithLabelValues(requestType, s.String())
		}

		for _, r := range []string{circuitBreakerResultSuccess, circuitBreakerResultError, circuitBreakerResultOpen} {
			// We initialize all possible results for the given requestType and the circuitBreakerResults metrics
			cbMetrics.circuitBreakerResults.WithLabelValues(requestType, r)
		}
	}
	return cbMetrics
}

type CircuitBreakerConfig struct {
	Enabled                    bool          `yaml:"enabled" category:"experimental"`
	FailureThresholdPercentage uint          `yaml:"failure_threshold_percentage" category:"experimental"`
	FailureExecutionThreshold  uint          `yaml:"failure_execution_threshold" category:"experimental"`
	ThresholdingPeriod         time.Duration `yaml:"thresholding_period" category:"experimental"`
	CooldownPeriod             time.Duration `yaml:"cooldown_period" category:"experimental"`
	InitialDelay               time.Duration `yaml:"initial_delay" category:"experimental"`
	RequestTimeout             time.Duration `yaml:"request_timeout" category:"experiment"`
	testModeEnabled            bool          `yaml:"-"`
}

func (cfg *CircuitBreakerConfig) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet, defaultRequestDuration time.Duration) {
	f.BoolVar(&cfg.Enabled, prefix+"enabled", false, "Enable circuit breaking when making requests to ingesters")
	f.UintVar(&cfg.FailureThresholdPercentage, prefix+"failure-threshold-percentage", 10, "Max percentage of requests that can fail over period before the circuit breaker opens")
	f.UintVar(&cfg.FailureExecutionThreshold, prefix+"failure-execution-threshold", 100, "How many requests must have been executed in period for the circuit breaker to be eligible to open for the rate of failures")
	f.DurationVar(&cfg.ThresholdingPeriod, prefix+"thresholding-period", time.Minute, "Moving window of time that the percentage of failed requests is computed over")
	f.DurationVar(&cfg.CooldownPeriod, prefix+"cooldown-period", 10*time.Second, "How long the circuit breaker will stay in the open state before allowing some requests")
	f.DurationVar(&cfg.InitialDelay, prefix+"initial-delay", 0, "How long the circuit breaker should wait between an activation request and becoming effectively active. During that time both failures and successes will not be counted.")
	f.DurationVar(&cfg.RequestTimeout, prefix+"request-timeout", defaultRequestDuration, "The maximum duration of an ingester's request before it triggers a circuit breaker. This configuration is used for circuit breakers only, and its timeouts aren't reported as errors.")
}

// circuitBreaker abstracts the ingester's server-side circuit breaker functionality.
// A nil *circuitBreaker is a valid noop implementation.
type circuitBreaker struct {
	cfg         CircuitBreakerConfig
	requestType string
	logger      log.Logger
	metrics     *circuitBreakerMetrics
	active      atomic.Bool
	cb          circuitbreaker.CircuitBreaker[any]

	// testRequestDelay is needed for testing purposes to simulate long-lasting requests
	testRequestDelay time.Duration
}

func newCircuitBreaker(cfg CircuitBreakerConfig, metrics *circuitBreakerMetrics, requestType string, logger log.Logger) *circuitBreaker {
	if !cfg.Enabled {
		return nil
	}
	active := atomic.NewBool(false)
	cb := circuitBreaker{
		cfg:         cfg,
		requestType: requestType,
		logger:      logger,
		metrics:     metrics,
		active:      *active,
	}

	circuitBreakerTransitionsCounterFn := func(metrics *circuitBreakerMetrics, requestType string, state circuitbreaker.State) prometheus.Counter {
		return metrics.circuitBreakerTransitions.WithLabelValues(requestType, state.String())
	}

	cbBuilder := circuitbreaker.Builder[any]().
		WithDelay(cfg.CooldownPeriod).
		OnClose(func(event circuitbreaker.StateChangedEvent) {
			circuitBreakerTransitionsCounterFn(cb.metrics, requestType, circuitbreaker.ClosedState).Inc()
			level.Info(logger).Log("msg", "circuit breaker is closed", "previous", event.OldState, "current", event.NewState)
		}).
		OnOpen(func(event circuitbreaker.StateChangedEvent) {
			circuitBreakerTransitionsCounterFn(cb.metrics, requestType, circuitbreaker.OpenState).Inc()
			level.Warn(logger).Log("msg", "circuit breaker is open", "previous", event.OldState, "current", event.NewState)
		}).
		OnHalfOpen(func(event circuitbreaker.StateChangedEvent) {
			circuitBreakerTransitionsCounterFn(cb.metrics, requestType, circuitbreaker.HalfOpenState).Inc()
			level.Info(logger).Log("msg", "circuit breaker is half-open", "previous", event.OldState, "current", event.NewState)
		})

	if cfg.testModeEnabled {
		// In case of testing purposes, we initialize the circuit breaker with count based failure thresholding,
		// since it is more deterministic, and therefore it is easier to predict the outcome.
		cbBuilder = cbBuilder.WithFailureThreshold(cfg.FailureThresholdPercentage)
	} else {
		// In case of production code, we prefer time based failure thresholding.
		cbBuilder = cbBuilder.WithFailureRateThreshold(cfg.FailureThresholdPercentage, cfg.FailureExecutionThreshold, cfg.ThresholdingPeriod)
	}

	cb.cb = cbBuilder.Build()
	return &cb
}

func isCircuitBreakerFailure(err error) bool {
	if err == nil {
		return false
	}

	// We only consider timeouts or ingester hitting a per-instance limit
	// to be errors worthy of tripping the circuit breaker since these
	// are specific to a particular ingester, not a user or request.

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	statusCode := grpcutil.ErrorToStatusCode(err)
	if statusCode == codes.DeadlineExceeded {
		return true
	}

	var ingesterErr ingesterError
	if errors.As(err, &ingesterErr) {
		return ingesterErr.errorCause() == mimirpb.INSTANCE_LIMIT
	}

	return false
}

func (cb *circuitBreaker) isActive() bool {
	return cb != nil && cb.active.Load()
}

func (cb *circuitBreaker) activate() {
	if cb == nil {
		return
	}
	if cb.cfg.InitialDelay == 0 {
		cb.active.Store(true)
	}
	time.AfterFunc(cb.cfg.InitialDelay, func() {
		cb.active.Store(true)
	})
}

// tryAcquirePermit tries to acquire a permit to use the circuit breaker and returns whether a permit was acquired.
// If it was possible to acquire a permit, success flag true and no error are returned. The acquired permit must be
// returned by a call to finishRequest.
// If it was not possible to acquire a permit, success flag false is returned. In this case no call to finishRequest
// is needed. If the permit was not acquired because of an error, that causing error is returned as well.
func (cb *circuitBreaker) tryAcquirePermit() (bool, error) {
	if !cb.isActive() {
		return false, nil
	}
	if !cb.cb.TryAcquirePermit() {
		cb.metrics.circuitBreakerResults.WithLabelValues(cb.requestType, circuitBreakerResultOpen).Inc()
		return false, newCircuitBreakerOpenError(cb.requestType, cb.cb.RemainingDelay())
	}
	return true, nil
}

// finishRequest completes a request executed upon a successfully acquired circuit breaker permit.
// It records the result of the request with the circuit breaker. Requests that lasted longer than
// the given maximumAllowedDuration are treated as a failure.
// The returned error is only used for testing purposes.
func (cb *circuitBreaker) finishRequest(actualDuration time.Duration, maximumAllowedDuration time.Duration, err error) error {
	if !cb.isActive() {
		return nil
	}
	if cb.cfg.testModeEnabled {
		actualDuration += cb.testRequestDelay
	}
	var deadlineErr error
	if maximumAllowedDuration < actualDuration {
		deadlineErr = context.DeadlineExceeded
	}
	return cb.recordResult(err, deadlineErr)
}

func (cb *circuitBreaker) recordResult(errs ...error) error {
	if !cb.isActive() {
		return nil
	}

	for _, err := range errs {
		if err != nil && isCircuitBreakerFailure(err) {
			cb.cb.RecordFailure()
			cb.metrics.circuitBreakerResults.WithLabelValues(cb.requestType, circuitBreakerResultError).Inc()
			return err
		}
	}
	cb.cb.RecordSuccess()
	cb.metrics.circuitBreakerResults.WithLabelValues(cb.requestType, circuitBreakerResultSuccess).Inc()
	return nil
}

type ingesterCircuitBreaker struct {
	push *circuitBreaker
	read *circuitBreaker
}

func newIngesterCircuitBreaker(pushCfg CircuitBreakerConfig, readCfg CircuitBreakerConfig, logger log.Logger, registerer prometheus.Registerer) *ingesterCircuitBreaker {
	prCB := &ingesterCircuitBreaker{}
	state := func(requestType string) circuitbreaker.State {
		switch requestType {
		case circuitBreakerPushRequestType:
			if prCB.push.isActive() {
				return prCB.push.cb.State()
			}
		case circuitBreakerReadRequestType:
			if prCB.read.isActive() {
				return prCB.read.cb.State()
			}
		}
		return -1
	}
	metrics := newCircuitBreakerMetrics(registerer, state, []string{circuitBreakerPushRequestType, circuitBreakerReadRequestType})
	prCB.push = newCircuitBreaker(pushCfg, metrics, circuitBreakerPushRequestType, logger)
	prCB.read = newCircuitBreaker(readCfg, metrics, circuitBreakerReadRequestType, logger)
	return prCB
}

func (cb *ingesterCircuitBreaker) activate() {
	if cb == nil {
		return
	}
	cb.push.activate()
	cb.read.activate()
}

// tryPushAcquirePermit tries to acquire a permit to use the push circuit breaker and returns whether a permit was acquired.
// If it was possible, tryPushAcquirePermit returns a function that should be called to release the acquired permit.
// If it was not possible, the causing error is returned.
func (cb *ingesterCircuitBreaker) tryPushAcquirePermit() (func(time.Duration, error), error) {
	if cb == nil {
		return func(time.Duration, error) {}, nil
	}

	pushAcquiredPermit, err := cb.push.tryAcquirePermit()
	if err != nil {
		return nil, err
	}
	return func(duration time.Duration, err error) {
		if pushAcquiredPermit {
			_ = cb.push.finishRequest(duration, cb.push.cfg.RequestTimeout, err)
		}
	}, nil
}

// tryReadAcquirePermit tries to acquire a permit to use the read circuit breaker and returns whether a permit was acquired.
// If it was possible, tryReadAcquirePermit returns a function that should be called to release the acquired permit.
// If it was not possible, the causing error is returned.
func (cb *ingesterCircuitBreaker) tryReadAcquirePermit() (func(time.Duration, error), error) {
	if cb == nil {
		return func(time.Duration, error) {}, nil
	}

	if cb.push.isActive() && cb.push.cb.State() == circuitbreaker.OpenState {
		return nil, newCircuitBreakerOpenError(cb.push.requestType, cb.push.cb.RemainingDelay())
	}

	readAcquiredPermit, err := cb.read.tryAcquirePermit()
	if err != nil {
		return nil, err
	}
	return func(duration time.Duration, err error) {
		if readAcquiredPermit {
			_ = cb.read.finishRequest(duration, cb.read.cfg.RequestTimeout, err)
		}
	}, nil
}

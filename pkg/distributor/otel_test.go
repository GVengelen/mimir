// SPDX-License-Identifier: AGPL-3.0-only

package distributor

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/dskit/user"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
	prometheustranslator "github.com/prometheus/prometheus/storage/remote/otlptranslator/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"

	"github.com/grafana/mimir/pkg/distributor/otlp"
	"github.com/grafana/mimir/pkg/mimirpb"
	"github.com/grafana/mimir/pkg/util/test"
	"github.com/grafana/mimir/pkg/util/validation"
)

func TestOTelMetricsToTimeSeries(t *testing.T) {
	const tenantID = "testTenant"
	discardedDueToOTelParseError := promauto.With(nil).NewCounterVec(prometheus.CounterOpts{
		Name: "discarded_due_to_otel_parse_error",
		Help: "Number of metrics discarded due to OTLP parse errors.",
	}, []string{tenantID, "group"})
	resourceAttrs := map[string]string{
		"service.name":        "service name",
		"service.instance.id": "service ID",
		"existent-attr":       "resource value",
		// This one is for testing conflict with metric attribute.
		"metric-attr": "resource value",
		// This one is for testing conflict with auto-generated job attribute.
		"job": "resource value",
		// This one is for testing conflict with auto-generated instance attribute.
		"instance": "resource value",
	}
	expTargetInfoLabels := []mimirpb.LabelAdapter{
		{
			Name:  labels.MetricName,
			Value: "target_info",
		},
	}
	for k, v := range resourceAttrs {
		switch k {
		case "service.name":
			k = "job"
		case "service.instance.id":
			k = "instance"
		case "job", "instance":
			// Ignore, as these labels are generated from service.name and service.instance.id
			continue
		default:
			k = prometheustranslator.NormalizeLabel(k)
		}
		expTargetInfoLabels = append(expTargetInfoLabels, mimirpb.LabelAdapter{
			Name:  k,
			Value: v,
		})
	}
	sort.Stable(otlp.ByLabelName(expTargetInfoLabels))

	md := pmetric.NewMetrics()
	{
		rm := md.ResourceMetrics().AppendEmpty()
		for k, v := range resourceAttrs {
			rm.Resource().Attributes().PutStr(k, v)
		}
		il := rm.ScopeMetrics().AppendEmpty()
		m := il.Metrics().AppendEmpty()
		m.SetName("test_metric")
		dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
		dp.SetIntValue(123)
		dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		dp.Attributes().PutStr("metric-attr", "metric value")
	}

	testCases := []struct {
		name                      string
		expectedLabels            []mimirpb.LabelAdapter
		promoteResourceAttributes []string
	}{
		{
			name:                      "Successful conversion without resource attribute promotion",
			promoteResourceAttributes: nil,
			expectedLabels: []mimirpb.LabelAdapter{
				{
					Name:  "__name__",
					Value: "test_metric",
				},
				{
					Name:  "instance",
					Value: "service ID",
				},
				{
					Name:  "job",
					Value: "service name",
				},
				{
					Name:  "metric_attr",
					Value: "metric value",
				},
			},
		},
		{
			name:                      "Successful conversion with resource attribute promotion",
			promoteResourceAttributes: []string{"non-existent-attr", "existent-attr"},
			expectedLabels: []mimirpb.LabelAdapter{
				{
					Name:  "__name__",
					Value: "test_metric",
				},
				{
					Name:  "instance",
					Value: "service ID",
				},
				{
					Name:  "job",
					Value: "service name",
				},
				{
					Name:  "metric_attr",
					Value: "metric value",
				},
				{
					Name:  "existent_attr",
					Value: "resource value",
				},
			},
		},
		{
			name:                      "Successful conversion with resource attribute promotion, conflicting resource attributes are ignored",
			promoteResourceAttributes: []string{"non-existent-attr", "existent-attr", "metric-attr", "job", "instance"},
			expectedLabels: []mimirpb.LabelAdapter{
				{
					Name:  "__name__",
					Value: "test_metric",
				},
				{
					Name:  "instance",
					Value: "service ID",
				},
				{
					Name:  "job",
					Value: "service name",
				},
				{
					Name:  "existent_attr",
					Value: "resource value",
				},
				{
					Name:  "metric_attr",
					Value: "metric value",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, testOld := range []bool{false, true} {
				t.Run(fmt.Sprintf("old function=%t", testOld), func(t *testing.T) {
					functionUnderTest := otelMetricsToTimeseries
					if testOld {
						functionUnderTest = otelMetricsToTimeseriesOld
					}
					mimirTS, err := functionUnderTest(
						tenantID, true, tc.promoteResourceAttributes, discardedDueToOTelParseError, log.NewNopLogger(), md,
					)
					require.NoError(t, err)
					require.Len(t, mimirTS, 2)
					var ts mimirpb.PreallocTimeseries
					var targetInfo mimirpb.PreallocTimeseries
					for i := range mimirTS {
						for _, lbl := range mimirTS[i].Labels {
							if lbl.Name != labels.MetricName {
								continue
							}

							if lbl.Value == "target_info" {
								targetInfo = mimirTS[i]
							} else {
								ts = mimirTS[i]
							}
						}
					}

					assert.ElementsMatch(t, ts.Labels, tc.expectedLabels)
					assert.ElementsMatch(t, targetInfo.Labels, expTargetInfoLabels)
				})
			}
		})
	}
}

func BenchmarkOTLPHandler(b *testing.B) {
	var samples []prompb.Sample
	for i := 0; i < 1000; i++ {
		ts := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(i) * time.Second)
		samples = append(samples, prompb.Sample{
			Value:     1,
			Timestamp: ts.UnixNano(),
		})
	}
	sampleSeries := []prompb.TimeSeries{
		{
			Labels: []prompb.Label{
				{Name: "__name__", Value: "foo"},
			},
			Samples: samples,
			Histograms: []prompb.Histogram{
				remote.HistogramToHistogramProto(1337, test.GenerateTestHistogram(1)),
			},
		},
	}
	// Sample metadata needs to correspond to every series in the sampleSeries
	sampleMetadata := []mimirpb.MetricMetadata{
		{
			Help: "metric_help",
			Unit: "metric_unit",
		},
	}
	exportReq := TimeseriesToOTLPRequest(sampleSeries, sampleMetadata)

	pushFunc := func(_ context.Context, pushReq *Request) error {
		if _, err := pushReq.WriteRequest(); err != nil {
			return err
		}

		pushReq.CleanUp()
		return nil
	}
	limits, err := validation.NewOverrides(
		validation.Limits{},
		validation.NewMockTenantLimits(map[string]*validation.Limits{}),
	)
	require.NoError(b, err)
	handler := OTLPHandler(100000, nil, nil, true, limits, RetryConfig{}, pushFunc, nil, nil, log.NewNopLogger(), true)

	b.Run("protobuf", func(b *testing.B) {
		req := createOTLPProtoRequest(b, exportReq, false)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			require.Equal(b, http.StatusOK, resp.Code)
			req.Body.(*reusableReader).Reset()
		}
	})

	b.Run("JSON", func(b *testing.B) {
		req := createOTLPJSONRequest(b, exportReq, false)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			require.Equal(b, http.StatusOK, resp.Code)
			req.Body.(*reusableReader).Reset()
		}
	})
}

func createOTLPProtoRequest(tb testing.TB, metricRequest pmetricotlp.ExportRequest, compress bool) *http.Request {
	tb.Helper()

	body, err := metricRequest.MarshalProto()
	require.NoError(tb, err)

	return createOTLPRequest(tb, body, compress, "application/x-protobuf")
}

func createOTLPJSONRequest(tb testing.TB, metricRequest pmetricotlp.ExportRequest, compress bool) *http.Request {
	tb.Helper()

	body, err := metricRequest.MarshalJSON()
	require.NoError(tb, err)

	return createOTLPRequest(tb, body, compress, "application/json")
}

func createOTLPRequest(tb testing.TB, body []byte, compress bool, contentType string) *http.Request {
	tb.Helper()

	if compress {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		_, err := gz.Write(body)
		require.NoError(tb, err)
		require.NoError(tb, gz.Close())

		body = b.Bytes()
	}

	// reusableReader is suitable for benchmarks
	req, err := http.NewRequest("POST", "http://localhost/", newReusableReader(body))
	require.NoError(tb, err)
	// Since http.NewRequest will deduce content length only from known io.Reader implementations,
	// define it ourselves
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", contentType)
	const tenantID = "test"
	req.Header.Set("X-Scope-OrgID", tenantID)
	ctx := user.InjectOrgID(context.Background(), tenantID)
	req = req.WithContext(ctx)
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
	}

	return req
}

type reusableReader struct {
	*bytes.Reader
	raw []byte
}

func newReusableReader(raw []byte) *reusableReader {
	return &reusableReader{
		Reader: bytes.NewReader(raw),
		raw:    raw,
	}
}

func (r *reusableReader) Close() error {
	return nil
}

func (r *reusableReader) Reset() {
	r.Reader.Reset(r.raw)
}

var _ io.ReadCloser = &reusableReader{}

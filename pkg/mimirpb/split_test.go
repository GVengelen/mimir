package mimirpb

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/require"
)

func prepareRequest() *WriteRequest {
	const numSeriesPerRequest = 100

	metrics := make([][]LabelAdapter, 0, numSeriesPerRequest)
	samples := make([]Sample, 0, numSeriesPerRequest)
	metadata := make([]*MetricMetadata, 0, numSeriesPerRequest)

	for i := 0; i < numSeriesPerRequest; i++ {
		metrics = append(metrics, []LabelAdapter{{Name: labels.MetricName, Value: "metric"}, {Name: "cardinality", Value: strconv.Itoa(i)}})
		samples = append(samples, Sample{Value: float64(i), TimestampMs: time.Now().UnixMilli()})
		metadata = append(metadata, &MetricMetadata{
			Type:             1 + MetricMetadata_MetricType(i)%STATESET,
			MetricFamilyName: fmt.Sprintf("metric_%d", i),
			Help:             fmt.Sprintf("help for metric_%d", i),
			Unit:             "unit",
		})
	}

	return ToWriteRequest(metrics, samples, nil, metadata, API)
}

func TestSplitRequestErrors(t *testing.T) {
	const timeseriesFieldSize = 100

	t.Run("unknown field", func(t *testing.T) {
		msg := []byte(nil)
		msg = append(msg, toVarint(tag(timeseriesField, wireTypeLen))...)
		msg = append(msg, toVarint(timeseriesFieldSize)...)
		msg = append(msg, make([]byte, timeseriesFieldSize)...)

		msg = append(msg, toVarint(tag(555, wireTypeVarint))...)

		_, err := SplitWriteRequestRequest(msg, len(msg)-1) // force splitting by using smaller maxSizeLimit
		require.ErrorContains(t, err, "unexpected tag 4440 (field: 555, type: 0)")
	})

	t.Run("wireTypeLen with too big length", func(t *testing.T) {
		msg := []byte(nil)
		msg = append(msg, toVarint(tag(timeseriesField, wireTypeLen))...)
		msg = append(msg, toVarint(math.MaxInt32+1)...)
		msg = append(msg, make([]byte, timeseriesFieldSize)...)

		_, err := SplitWriteRequestRequest(msg, len(msg)-1) // force splitting by using smaller maxSizeLimit
		require.ErrorContains(t, err, "invalid decoded length: 2147483648")
	})

	t.Run("short message", func(t *testing.T) {
		msg := []byte(nil)
		msg = append(msg, toVarint(tag(timeseriesField, wireTypeLen))...)
		msg = append(msg, toVarint(timeseriesFieldSize)...)
		msg = append(msg, make([]byte, timeseriesFieldSize-1)...)

		_, err := SplitWriteRequestRequest(msg, len(msg)-1) // force splitting by using smaller maxSizeLimit
		require.ErrorContains(t, err, "message too short, expected length: 100, remaining buffer: 99")
	})

	t.Run("invalid source value", func(t *testing.T) {
		msg := []byte(nil)
		msg = append(msg, toVarint(tag(timeseriesField, wireTypeLen))...)
		msg = append(msg, toVarint(timeseriesFieldSize)...)
		msg = append(msg, make([]byte, timeseriesFieldSize)...)

		msg = append(msg, toVarint(sourceFieldTag)...)
		msg = append(msg, toVarint(math.MaxInt32+1)...)

		_, err := SplitWriteRequestRequest(msg, len(msg)-1) // force splitting by using smaller maxSizeLimit
		require.ErrorContains(t, err, "invalid value 2147483648 for tag 16")
	})
}

func TestSplitRequestWithWeirdSource(t *testing.T) {
	testCases := map[string]func(t *testing.T) []byte{
		"simple": func(t *testing.T) []byte {
			return marshal(t, prepareRequest())
		},

		"weird source": func(t *testing.T) []byte {
			wr := prepareRequest()
			wr.Source = WriteRequest_SourceEnum(123456)
			return marshal(t, wr)
		},

		"weird source, set skipLabelNameValidation": func(t *testing.T) []byte {
			wr := prepareRequest()
			wr.Source = WriteRequest_SourceEnum(123456)
			wr.SkipLabelNameValidation = true
			return marshal(t, wr)
		},

		"request with prepended fields": func(t *testing.T) []byte {
			wr := prepareRequest()
			wr.Source = WriteRequest_SourceEnum(123456)
			wr.SkipLabelNameValidation = true
			m := marshal(t, wr)

			raw := []byte(nil)
			raw = append(raw, toVarint(sourceFieldTag)...)
			raw = append(raw, toVarint(32516)...)

			raw = append(raw, toVarint(skipLabelNameValidationFieldTag)...)
			raw = append(raw, toVarint(0)...)

			raw = append(raw, m...)

			// Verify that entire "raw" message can be unmarshaled.
			wr.Reset()
			err := wr.Unmarshal(raw)
			require.NoError(t, err)

			// Check that last values are honored.
			require.Equal(t, WriteRequest_SourceEnum(123456), wr.Source) // from wr
			require.True(t, wr.SkipLabelNameValidation)                  // from wr

			return raw
		},

		"request with repeated appended fields": func(t *testing.T) []byte {
			wr := prepareRequest()
			wr.Source = WriteRequest_SourceEnum(123456)
			wr.SkipLabelNameValidation = true
			raw := marshal(t, wr)

			// Add another "source" field
			raw = append(raw, toVarint(sourceFieldTag)...)
			raw = append(raw, toVarint(32516)...)

			// Add new "skipLabelNameValidation" field
			raw = append(raw, toVarint(skipLabelNameValidationFieldTag)...)
			raw = append(raw, toVarint(0)...)

			// One more "source"
			raw = append(raw, toVarint(sourceFieldTag)...)
			raw = append(raw, toVarint(555)...)

			// One more "skipLabelNameValidation" field. Technically bools should only have 0 or 1 values, but gogoproto accepts non-zero value as true.
			raw = append(raw, toVarint(skipLabelNameValidationFieldTag)...)
			raw = append(raw, toVarint(5)...)

			// Verify that we can unmarshal this, and both source and skipLabelNameValidation are set to last value in the message.
			wr.Reset()
			err := wr.Unmarshal(raw)
			require.NoError(t, err)

			// Check that last values are honored.
			require.Equal(t, WriteRequest_SourceEnum(555), wr.Source) // last appended source value
			require.True(t, wr.SkipLabelNameValidation)               // last appended skipLabelNameValidation value

			return raw
		},
	}

	for name, f := range testCases {
		marshalled := f(t)
		size := len(marshalled)

		for _, maxSize := range []int{size / 10, size / 5, size / 2, size - maxExtraBytes, size, size + maxExtraBytes} {
			var expectedSubrequests int
			if size <= maxSize {
				expectedSubrequests = 1
			} else {
				expectedSubrequests = int(math.Ceil(float64(size) / float64(maxSize-maxExtraBytes)))
			}

			t.Run(fmt.Sprintf("%s: total size: %d, max size: %d, expected requests: %d", name, size, maxSize, expectedSubrequests), func(t *testing.T) {
				reqs, err := SplitWriteRequestRequest(marshalled, maxSize)
				require.NoError(t, err)
				require.Equal(t, len(reqs), expectedSubrequests)
				verifyRequests(t, reqs, marshalled, maxSize)
			})
		}
	}
}

func marshal(t *testing.T, wr *WriteRequest) []byte {
	m, err := wr.Marshal()
	require.NoError(t, err)
	return m
}

func verifyRequests(t *testing.T, reqs [][]byte, original []byte, maxSize int) {
	wr := &WriteRequest{}
	err := wr.Unmarshal(original)
	require.NoError(t, err)

	// combine all subrequests together into single request, and verify that it's the same as original
	combined := WriteRequest{
		Source:                  wr.Source,                  // checked separately
		SkipLabelNameValidation: wr.SkipLabelNameValidation, // checked separately
	}

	for ix := range reqs {
		require.LessOrEqual(t, len(reqs[ix]), maxSize)

		p := WriteRequest{}
		// Check that all requests can be parsed.
		err := p.Unmarshal(reqs[ix])
		require.NoError(t, err)

		// Verify that all parsed requests have same source and skipValidation fields as original
		require.Equal(t, p.Source, wr.Source)
		require.Equal(t, p.SkipLabelNameValidation, wr.SkipLabelNameValidation)

		combined.Timeseries = append(combined.Timeseries, p.Timeseries...)
		combined.Metadata = append(combined.Metadata, p.Metadata...)
	}

	// Ideally we would use combined.Equal(wr), but PreallocWriteRequest uses PreallocTimeseries and that failes to compare with Timeseries
	// Instead, we just marshal both values and check serialized form.
	combinedSerialized, err := combined.Marshal()
	require.NoError(t, err)
	origSerialized, err := wr.Marshal()
	require.NoError(t, err)
	require.Equal(t, origSerialized, combinedSerialized)
}

func toVarint(val uint64) []byte {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], val)
	return buf[:n]
}

func BenchmarkWriteRequest_SplitByMaxMarshalSize(b *testing.B) {
	tests := map[string]struct {
		numSeries           int
		numLabelsPerSeries  int
		numSamplesPerSeries int
		numMetadata         int
	}{
		"write request with few series, few labels each, and no metadata": {
			numSeries:           50,
			numLabelsPerSeries:  10,
			numSamplesPerSeries: 1,
			numMetadata:         0,
		},
		"write request with few series, many labels each, and no metadata": {
			numSeries:           50,
			numLabelsPerSeries:  100,
			numSamplesPerSeries: 1,
			numMetadata:         0,
		},
		"write request with many series, few labels each, and no metadata": {
			numSeries:           1000,
			numLabelsPerSeries:  10,
			numSamplesPerSeries: 1,
			numMetadata:         0,
		},
		"write request with many series, many labels each, and no metadata": {
			numSeries:           1000,
			numLabelsPerSeries:  100,
			numSamplesPerSeries: 1,
			numMetadata:         0,
		},
		"write request with few metadata, and no series": {
			numSeries:           0,
			numLabelsPerSeries:  0,
			numSamplesPerSeries: 0,
			numMetadata:         50,
		},
		"write request with many metadata, and no series": {
			numSeries:           0,
			numLabelsPerSeries:  0,
			numSamplesPerSeries: 0,
			numMetadata:         1000,
		},
		"write request with both series and metadata": {
			numSeries:           500,
			numLabelsPerSeries:  25,
			numSamplesPerSeries: 1,
			numMetadata:         500,
		},
	}

	for testName, testData := range tests {
		b.Run(testName, func(b *testing.B) {
			req := generateWriteRequest(testData.numSeries, testData.numLabelsPerSeries, testData.numSamplesPerSeries, testData.numMetadata)
			reqSize := req.Size()

			// Test with different split size.
			splitScenarios := map[string]struct {
				maxSize              int
				expectedApproxSplits int
			}{
				"no splitting": {
					maxSize:              reqSize * 2,
					expectedApproxSplits: 1,
				},
				"split in few requests": {
					maxSize:              int(float64(reqSize) * 0.8),
					expectedApproxSplits: 2,
				},
				"split in many requests": {
					maxSize:              int(float64(reqSize) * 0.11),
					expectedApproxSplits: 10,
				},
			}

			for splitName, splitScenario := range splitScenarios {
				b.Run(splitName, func(b *testing.B) {
					// The actual number of splits may be slightly different then the expected, due to implementation
					// details (e.g. if a request both contain series and metadata, they're never mixed in the same split request).
					minExpectedSplits := splitScenario.expectedApproxSplits - 1
					maxExpectedSplits := splitScenario.expectedApproxSplits + 1
					if splitScenario.expectedApproxSplits == 0 {
						minExpectedSplits = 0
						maxExpectedSplits = 0
					}

					for n := 0; n < b.N; n++ {
						// Marshal the request each time. This is also offer a fair comparison with an alternative
						// implementation (we're considering) which does the splitting before marshalling.
						marshalled, err := req.Marshal()
						if err != nil {
							b.Fatal(err)
						}

						actualSplits, err := SplitWriteRequestRequest(marshalled, splitScenario.maxSize)
						if err != nil {
							b.Fatal(err)
						}

						// Ensure the number of splits match the expected ones.
						if numActualSplits := len(actualSplits); (numActualSplits < minExpectedSplits) || (numActualSplits > maxExpectedSplits) {
							b.Fatalf("expected between %d and %d splits but got %d", minExpectedSplits, maxExpectedSplits, numActualSplits)
						}

						for _, data := range actualSplits {
							if len(data) >= splitScenario.maxSize {
								b.Fatalf("the marshalled split request (%d bytes) is larger than max size (%d bytes)", len(data), splitScenario.maxSize)
							}
						}
					}
				})
			}
		})
	}
}

func generateWriteRequest(numSeries, numLabelsPerSeries, numSamplesPerSeries, numMetadata int) *WriteRequest {
	builder := labels.NewScratchBuilder(numLabelsPerSeries)

	// Generate timeseries.
	timeseries := make([]PreallocTimeseries, 0, numSeries)
	for i := 0; i < numSeries; i++ {
		curr := PreallocTimeseries{TimeSeries: &TimeSeries{}}

		// Generate series labels.
		builder.Reset()
		builder.Add(labels.MetricName, fmt.Sprintf("series_%d", i))
		for l := 1; l < numLabelsPerSeries; l++ {
			builder.Add(fmt.Sprintf("label_%d", l), fmt.Sprintf("this-is-the-value-of-label-%d", l))
		}
		curr.Labels = FromLabelsToLabelAdapters(builder.Labels())

		// Generate samples.
		curr.Samples = make([]Sample, 0, numSamplesPerSeries)
		for s := 0; s < numSamplesPerSeries; s++ {
			curr.Samples = append(curr.Samples, Sample{
				TimestampMs: int64(s),
				Value:       float64(s),
			})
		}

		// Add an exemplar.
		builder.Reset()
		builder.Add("trace_id", fmt.Sprintf("the-trace-id-for-%d", i))
		curr.Exemplars = []Exemplar{{
			Labels:      FromLabelsToLabelAdapters(builder.Labels()),
			TimestampMs: int64(i),
			Value:       float64(i),
		}}

		timeseries = append(timeseries, curr)
	}

	// Generate metadata.
	metadata := make([]*MetricMetadata, 0, numMetadata)
	for i := 0; i < numMetadata; i++ {
		metadata = append(metadata, &MetricMetadata{
			Type:             COUNTER,
			MetricFamilyName: fmt.Sprintf("series_%d", i),
			Help:             fmt.Sprintf("this is the help description for series %d", i),
			Unit:             "seconds",
		})
	}

	return &WriteRequest{
		Source:                  RULE,
		SkipLabelNameValidation: true,
		Timeseries:              timeseries,
		Metadata:                metadata,
	}
}

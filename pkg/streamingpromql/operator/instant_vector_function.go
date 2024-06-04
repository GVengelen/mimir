// SPDX-License-Identifier: AGPL-3.0-only
// Provenance-includes-location: https://github.com/prometheus/prometheus/blob/main/promql/engine.go
// Provenance-includes-location: https://github.com/prometheus/prometheus/blob/main/promql/functions.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: The Prometheus Authors

package operator

import (
	"context"
	"github.com/grafana/mimir/pkg/streamingpromql/functions"
	"github.com/grafana/mimir/pkg/streamingpromql/pooling"
	"github.com/grafana/mimir/pkg/streamingpromql/types"
)

// InstantVectorFunction performs a function over each series in an instant vector.
type InstantVectorFunction struct {
	Inner InstantVectorOperator
	Pool  *pooling.LimitingPool

	MetadataFunc   functions.SeriesMetadataFunction
	SeriesDataFunc functions.InstantVectorFunction
}

var _ InstantVectorOperator = &InstantVectorFunction{}

func (m *InstantVectorFunction) SeriesMetadata(ctx context.Context) ([]types.SeriesMetadata, error) {
	metadata, err := m.Inner.SeriesMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return m.MetadataFunc(metadata, m.Pool)
}

func (m *InstantVectorFunction) NextSeries(ctx context.Context) (types.InstantVectorSeriesData, error) {
	series, err := m.Inner.NextSeries(ctx)
	if err != nil {
		return types.InstantVectorSeriesData{}, err
	}

	return m.SeriesDataFunc(series, m.Pool)
}

func (m *InstantVectorFunction) Close() {
	m.Inner.Close()
}

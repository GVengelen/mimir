// SPDX-License-Identifier: AGPL-3.0-only

package streamingpromql

import (
	"flag"

	"github.com/prometheus/prometheus/promql"
)

type EngineOpts struct {
	CommonOpts     promql.EngineOpts
	FeatureToggles FeatureToggles

	// When operating in pedantic mode, we panic if memory consumption is > 0 after Query.Close()
	// (indicating something was not returned to a pool).
	Pedantic bool
}

type FeatureToggles struct {
	EnableAggregationOperations                  bool `yaml:"enable_aggregation_operations" category:"experimental"`
	EnableVectorVectorBinaryComparisonOperations bool `yaml:"enable_vector_vector_binary_comparison_operations" category:"experimental"`
	EnableVectorScalarBinaryComparisonOperations bool `yaml:"enable_vector_scalar_binary_comparison_operations" category:"experimental"`
	EnableScalarScalarBinaryComparisonOperations bool `yaml:"enable_scalar_scalar_binary_comparison_operations" category:"experimental"`
	EnableBinaryLogicalOperations                bool `yaml:"enable_binary_logical_operations" category:"experimental"`
	EnableScalars                                bool `yaml:"enable_scalars" category:"experimental"`
	EnableSubqueries                             bool `yaml:"enable_subqueries" category:"experimental"`
	EnableHistogramQuantileFunction              bool `yaml:"enable_histogram_quantile_function" category:"experimental"`
	EnableOneToManyAndManyToOneBinaryOperations  bool `yaml:"enable_one_to_many_and_many_to_one_binary_operations" category:"experimental"`
}

// EnableAllFeatures enables all features supported by MQE, including experimental or incomplete features.
var EnableAllFeatures = FeatureToggles{
	// Note that we deliberately use a keyless literal here to force a compilation error if we don't keep this in sync with new fields added to FeatureToggles.
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
}

func (t *FeatureToggles) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&t.EnableAggregationOperations, "querier.mimir-query-engine.enable-aggregation-operations", true, "Enable support for aggregation operations in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableVectorVectorBinaryComparisonOperations, "querier.mimir-query-engine.enable-vector-vector-binary-comparison-operations", true, "Enable support for binary comparison operations between two vectors in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableVectorScalarBinaryComparisonOperations, "querier.mimir-query-engine.enable-vector-scalar-binary-comparison-operations", true, "Enable support for binary comparison operations between a vector and a scalar in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableScalarScalarBinaryComparisonOperations, "querier.mimir-query-engine.enable-scalar-scalar-binary-comparison-operations", true, "Enable support for binary comparison operations between two scalars in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableBinaryLogicalOperations, "querier.mimir-query-engine.enable-binary-logical-operations", true, "Enable support for binary logical operations in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableScalars, "querier.mimir-query-engine.enable-scalars", true, "Enable support for scalars in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableSubqueries, "querier.mimir-query-engine.enable-subqueries", true, "Enable support for subqueries in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableHistogramQuantileFunction, "querier.mimir-query-engine.enable-histogram-quantile-function", true, "Enable support for the histogram_quantile function in the Mimir query engine. Only applies if the MQE is in use.")
	f.BoolVar(&t.EnableOneToManyAndManyToOneBinaryOperations, "querier.mimir-query-engine.enable-one-to-many-and-many-to-one-binary-operations", true, "Enable support for one-to-many and many-to-one binary operations (group_left/group_right) in the Mimir query engine. Only applies if the MQE is in use.")
}
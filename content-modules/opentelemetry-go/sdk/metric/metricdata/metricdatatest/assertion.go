// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metricdatatest provides testing functionality for use with the
// metricdata package.
package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Datatypes are the concrete data-types the metricdata package provides.
type Datatypes interface {
	metricdata.DataPoint[float64] |
		metricdata.DataPoint[int64] |
		metricdata.Gauge[float64] |
		metricdata.Gauge[int64] |
		metricdata.Histogram[float64] |
		metricdata.Histogram[int64] |
		metricdata.HistogramDataPoint[float64] |
		metricdata.HistogramDataPoint[int64] |
		metricdata.Extrema[int64] |
		metricdata.Extrema[float64] |
		metricdata.Metrics |
		metricdata.ResourceMetrics |
		metricdata.ScopeMetrics |
		metricdata.Sum[float64] |
		metricdata.Sum[int64] |
		metricdata.Exemplar[float64] |
		metricdata.Exemplar[int64]

	// Interface types are not allowed in union types, therefore the
	// Aggregation and Value type from metricdata are not included here.
}

type config struct {
	ignoreTimestamp bool
	ignoreExemplars bool
}

func newConfig(opts []Option) config {
	var cfg config
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// Option allows for fine grain control over how AssertEqual operates.
type Option interface {
	apply(cfg config) config
}

type fnOption func(cfg config) config

func (fn fnOption) apply(cfg config) config {
	return fn(cfg)
}

// IgnoreTimestamp disables checking if timestamps are different.
func IgnoreTimestamp() Option {
	return fnOption(func(cfg config) config {
		cfg.ignoreTimestamp = true
		return cfg
	})
}

// IgnoreExemplars disables checking if Exemplars are different.
func IgnoreExemplars() Option {
	return fnOption(func(cfg config) config {
		cfg.ignoreExemplars = true
		return cfg
	})
}

// AssertEqual asserts that the two concrete data-types from the metricdata
// package are equal.
func AssertEqual[T Datatypes](t *testing.T, expected, actual T, opts ...Option) bool {
	t.Helper()

	cfg := newConfig(opts)

	// Generic types cannot be type asserted. Use an interface instead.
	aIface := interface{}(actual)

	var r []string
	switch e := interface{}(expected).(type) {
	case metricdata.Exemplar[int64]:
		r = equalExemplars(e, aIface.(metricdata.Exemplar[int64]), cfg)
	case metricdata.Exemplar[float64]:
		r = equalExemplars(e, aIface.(metricdata.Exemplar[float64]), cfg)
	case metricdata.DataPoint[int64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[int64]), cfg)
	case metricdata.DataPoint[float64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[float64]), cfg)
	case metricdata.Gauge[int64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[int64]), cfg)
	case metricdata.Gauge[float64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[float64]), cfg)
	case metricdata.Histogram[float64]:
		r = equalHistograms(e, aIface.(metricdata.Histogram[float64]), cfg)
	case metricdata.Histogram[int64]:
		r = equalHistograms(e, aIface.(metricdata.Histogram[int64]), cfg)
	case metricdata.HistogramDataPoint[float64]:
		r = equalHistogramDataPoints(e, aIface.(metricdata.HistogramDataPoint[float64]), cfg)
	case metricdata.HistogramDataPoint[int64]:
		r = equalHistogramDataPoints(e, aIface.(metricdata.HistogramDataPoint[int64]), cfg)
	case metricdata.Extrema[int64]:
		r = equalExtrema(e, aIface.(metricdata.Extrema[int64]), cfg)
	case metricdata.Extrema[float64]:
		r = equalExtrema(e, aIface.(metricdata.Extrema[float64]), cfg)
	case metricdata.Metrics:
		r = equalMetrics(e, aIface.(metricdata.Metrics), cfg)
	case metricdata.ResourceMetrics:
		r = equalResourceMetrics(e, aIface.(metricdata.ResourceMetrics), cfg)
	case metricdata.ScopeMetrics:
		r = equalScopeMetrics(e, aIface.(metricdata.ScopeMetrics), cfg)
	case metricdata.Sum[int64]:
		r = equalSums(e, aIface.(metricdata.Sum[int64]), cfg)
	case metricdata.Sum[float64]:
		r = equalSums(e, aIface.(metricdata.Sum[float64]), cfg)
	default:
		// We control all types passed to this, panic to signal developers
		// early they changed things in an incompatible way.
		panic(fmt.Sprintf("unknown types: %T", expected))
	}

	if len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation, opts ...Option) bool {
	t.Helper()

	cfg := newConfig(opts)
	if r := equalAggregations(expected, actual, cfg); len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}

// AssertHasAttributes asserts that all Datapoints or HistogramDataPoints have all passed attrs.
func AssertHasAttributes[T Datatypes](t *testing.T, actual T, attrs ...attribute.KeyValue) bool {
	t.Helper()

	var reasons []string

	switch e := interface{}(actual).(type) {
	case metricdata.Exemplar[int64]:
		reasons = hasAttributesExemplars(e, attrs...)
	case metricdata.Exemplar[float64]:
		reasons = hasAttributesExemplars(e, attrs...)
	case metricdata.DataPoint[int64]:
		reasons = hasAttributesDataPoints(e, attrs...)
	case metricdata.DataPoint[float64]:
		reasons = hasAttributesDataPoints(e, attrs...)
	case metricdata.Gauge[int64]:
		reasons = hasAttributesGauge(e, attrs...)
	case metricdata.Gauge[float64]:
		reasons = hasAttributesGauge(e, attrs...)
	case metricdata.Sum[int64]:
		reasons = hasAttributesSum(e, attrs...)
	case metricdata.Sum[float64]:
		reasons = hasAttributesSum(e, attrs...)
	case metricdata.HistogramDataPoint[int64]:
		reasons = hasAttributesHistogramDataPoints(e, attrs...)
	case metricdata.HistogramDataPoint[float64]:
		reasons = hasAttributesHistogramDataPoints(e, attrs...)
	case metricdata.Extrema[int64], metricdata.Extrema[float64]:
		// Nothing to check.
	case metricdata.Histogram[int64]:
		reasons = hasAttributesHistogram(e, attrs...)
	case metricdata.Histogram[float64]:
		reasons = hasAttributesHistogram(e, attrs...)
	case metricdata.Metrics:
		reasons = hasAttributesMetrics(e, attrs...)
	case metricdata.ScopeMetrics:
		reasons = hasAttributesScopeMetrics(e, attrs...)
	case metricdata.ResourceMetrics:
		reasons = hasAttributesResourceMetrics(e, attrs...)
	default:
		// We control all types passed to this, panic to signal developers
		// early they changed things in an incompatible way.
		panic(fmt.Sprintf("unknown types: %T", actual))
	}

	if len(reasons) > 0 {
		t.Error(reasons)
		return false
	}

	return true
}

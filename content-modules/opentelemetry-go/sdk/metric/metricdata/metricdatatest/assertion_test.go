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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	attrA = attribute.NewSet(attribute.Bool("A", true))
	attrB = attribute.NewSet(attribute.Bool("B", true))

	fltrAttrA = []attribute.KeyValue{attribute.Bool("filter A", true)}
	fltrAttrB = []attribute.KeyValue{attribute.Bool("filter B", true)}

	startA = time.Now()
	startB = startA.Add(time.Millisecond)
	endA   = startA.Add(time.Second)
	endB   = startB.Add(time.Second)

	spanIDA  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDB  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDA = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDB = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	exemplarInt64A = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              -10,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarFloat64A = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              -10.0,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarInt64B = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrB,
		Time:               endB,
		Value:              12,
		SpanID:             spanIDB,
		TraceID:            traceIDB,
	}
	exemplarFloat64B = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrB,
		Time:               endB,
		Value:              12.0,
		SpanID:             spanIDB,
		TraceID:            traceIDB,
	}
	exemplarInt64C = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrA,
		Time:               endB,
		Value:              -10,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarFloat64C = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrA,
		Time:               endB,
		Value:              -10.0,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}

	dataPointInt64A = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	dataPointFloat64A = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64A},
	}
	dataPointInt64B = metricdata.DataPoint[int64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64B},
	}
	dataPointFloat64B = metricdata.DataPoint[float64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64B},
	}
	dataPointInt64C = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64C},
	}
	dataPointFloat64C = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64C},
	}

	minFloat64A              = metricdata.NewExtrema(-1.)
	minInt64A                = metricdata.NewExtrema[int64](-1)
	minFloat64B, maxFloat64B = metricdata.NewExtrema(3.), metricdata.NewExtrema(99.)
	minInt64B, maxInt64B     = metricdata.NewExtrema[int64](3), metricdata.NewExtrema[int64](99)
	minFloat64C              = metricdata.NewExtrema(-1.)
	minInt64C                = metricdata.NewExtrema[int64](-1)

	histogramDataPointInt64A = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minInt64A,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	histogramDataPointFloat64A = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minFloat64A,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64A},
	}
	histogramDataPointInt64B = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxInt64B,
		Min:          minInt64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64B},
	}
	histogramDataPointFloat64B = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxFloat64B,
		Min:          minFloat64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64B},
	}
	histogramDataPointInt64C = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrA,
		StartTime:    startB,
		Time:         endB,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minInt64C,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64C},
	}
	histogramDataPointFloat64C = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrA,
		StartTime:    startB,
		Time:         endB,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minFloat64C,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64C},
	}

	gaugeInt64A = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64A},
	}
	gaugeFloat64A = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64A},
	}
	gaugeInt64B = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64B},
	}
	gaugeFloat64B = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64B},
	}
	gaugeInt64C = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64C},
	}
	gaugeFloat64C = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64C},
	}

	sumInt64A = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64A},
	}
	sumFloat64A = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64A},
	}
	sumInt64B = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64B},
	}
	sumFloat64B = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64B},
	}
	sumInt64C = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64C},
	}
	sumFloat64C = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64C},
	}

	histogramInt64A = metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64A},
	}
	histogramFloat64A = metricdata.Histogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64A},
	}
	histogramInt64B = metricdata.Histogram[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64B},
	}
	histogramFloat64B = metricdata.Histogram[float64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64B},
	}
	histogramInt64C = metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64C},
	}
	histogramFloat64C = metricdata.Histogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64C},
	}

	metricsA = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        "1",
		Data:        sumInt64A,
	}
	metricsB = metricdata.Metrics{
		Name:        "B",
		Description: "B desc",
		Unit:        "By",
		Data:        gaugeFloat64B,
	}
	metricsC = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        "1",
		Data:        sumInt64C,
	}

	scopeMetricsA = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsA},
	}
	scopeMetricsB = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "B"},
		Metrics: []metricdata.Metrics{metricsB},
	}
	scopeMetricsC = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsC},
	}

	resourceMetricsA = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsA},
	}
	resourceMetricsB = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "B")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsB},
	}
	resourceMetricsC = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsC},
	}
)

type equalFunc[T Datatypes] func(T, T, config) []string

func testDatatype[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		r := f(a, b, newConfig(nil))
		assert.Greaterf(t, len(r), 0, "%v == %v", a, b)
	}
}

func testDatatypeIgnoreTime[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		c := newConfig([]Option{IgnoreTimestamp()})
		r := f(a, b, c)
		assert.Len(t, r, 0, "unexpected inequality")
	}
}

func testDatatypeIgnoreExemplars[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		c := newConfig([]Option{IgnoreExemplars()})
		r := f(a, b, c)
		assert.Len(t, r, 0, "unexpected inequality")
	}
}

func TestAssertEqual(t *testing.T) {
	t.Run("ResourceMetrics", testDatatype(resourceMetricsA, resourceMetricsB, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatype(scopeMetricsA, scopeMetricsB, equalScopeMetrics))
	t.Run("Metrics", testDatatype(metricsA, metricsB, equalMetrics))
	t.Run("HistogramInt64", testDatatype(histogramInt64A, histogramInt64B, equalHistograms[int64]))
	t.Run("HistogramFloat64", testDatatype(histogramFloat64A, histogramFloat64B, equalHistograms[float64]))
	t.Run("SumInt64", testDatatype(sumInt64A, sumInt64B, equalSums[int64]))
	t.Run("SumFloat64", testDatatype(sumFloat64A, sumFloat64B, equalSums[float64]))
	t.Run("GaugeInt64", testDatatype(gaugeInt64A, gaugeInt64B, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatype(gaugeFloat64A, gaugeFloat64B, equalGauges[float64]))
	t.Run("HistogramDataPointInt64", testDatatype(histogramDataPointInt64A, histogramDataPointInt64B, equalHistogramDataPoints[int64]))
	t.Run("HistogramDataPointFloat64", testDatatype(histogramDataPointFloat64A, histogramDataPointFloat64B, equalHistogramDataPoints[float64]))
	t.Run("DataPointInt64", testDatatype(dataPointInt64A, dataPointInt64B, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatype(dataPointFloat64A, dataPointFloat64B, equalDataPoints[float64]))
	t.Run("ExtremaInt64", testDatatype(minInt64A, minInt64B, equalExtrema[int64]))
	t.Run("ExtremaFloat64", testDatatype(minFloat64A, minFloat64B, equalExtrema[float64]))
	t.Run("ExemplarInt64", testDatatype(exemplarInt64A, exemplarInt64B, equalExemplars[int64]))
	t.Run("ExemplarFloat64", testDatatype(exemplarFloat64A, exemplarFloat64B, equalExemplars[float64]))
}

func TestAssertEqualIgnoreTime(t *testing.T) {
	t.Run("ResourceMetrics", testDatatypeIgnoreTime(resourceMetricsA, resourceMetricsC, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatypeIgnoreTime(scopeMetricsA, scopeMetricsC, equalScopeMetrics))
	t.Run("Metrics", testDatatypeIgnoreTime(metricsA, metricsC, equalMetrics))
	t.Run("HistogramInt64", testDatatypeIgnoreTime(histogramInt64A, histogramInt64C, equalHistograms[int64]))
	t.Run("HistogramFloat64", testDatatypeIgnoreTime(histogramFloat64A, histogramFloat64C, equalHistograms[float64]))
	t.Run("SumInt64", testDatatypeIgnoreTime(sumInt64A, sumInt64C, equalSums[int64]))
	t.Run("SumFloat64", testDatatypeIgnoreTime(sumFloat64A, sumFloat64C, equalSums[float64]))
	t.Run("GaugeInt64", testDatatypeIgnoreTime(gaugeInt64A, gaugeInt64C, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatypeIgnoreTime(gaugeFloat64A, gaugeFloat64C, equalGauges[float64]))
	t.Run("HistogramDataPointInt64", testDatatypeIgnoreTime(histogramDataPointInt64A, histogramDataPointInt64C, equalHistogramDataPoints[int64]))
	t.Run("HistogramDataPointFloat64", testDatatypeIgnoreTime(histogramDataPointFloat64A, histogramDataPointFloat64C, equalHistogramDataPoints[float64]))
	t.Run("DataPointInt64", testDatatypeIgnoreTime(dataPointInt64A, dataPointInt64C, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatypeIgnoreTime(dataPointFloat64A, dataPointFloat64C, equalDataPoints[float64]))
	t.Run("ExtremaInt64", testDatatypeIgnoreTime(minInt64A, minInt64C, equalExtrema[int64]))
	t.Run("ExtremaFloat64", testDatatypeIgnoreTime(minFloat64A, minFloat64C, equalExtrema[float64]))
	t.Run("ExemplarInt64", testDatatypeIgnoreTime(exemplarInt64A, exemplarInt64C, equalExemplars[int64]))
	t.Run("ExemplarFloat64", testDatatypeIgnoreTime(exemplarFloat64A, exemplarFloat64C, equalExemplars[float64]))
}

func TestAssertEqualIgnoreExemplars(t *testing.T) {
	hdpInt64 := histogramDataPointInt64A
	hdpInt64.Exemplars = []metricdata.Exemplar[int64]{exemplarInt64B}
	t.Run("HistogramDataPointInt64", testDatatypeIgnoreExemplars(histogramDataPointInt64A, hdpInt64, equalHistogramDataPoints[int64]))

	hdpFloat64 := histogramDataPointFloat64A
	hdpFloat64.Exemplars = []metricdata.Exemplar[float64]{exemplarFloat64B}
	t.Run("HistogramDataPointFloat64", testDatatypeIgnoreExemplars(histogramDataPointFloat64A, hdpFloat64, equalHistogramDataPoints[float64]))

	dpInt64 := dataPointInt64A
	dpInt64.Exemplars = []metricdata.Exemplar[int64]{exemplarInt64B}
	t.Run("DataPointInt64", testDatatypeIgnoreExemplars(dataPointInt64A, dpInt64, equalDataPoints[int64]))

	dpFloat64 := dataPointFloat64A
	dpFloat64.Exemplars = []metricdata.Exemplar[float64]{exemplarFloat64B}
	t.Run("DataPointFloat64", testDatatypeIgnoreExemplars(dataPointFloat64A, dpFloat64, equalDataPoints[float64]))
}

type unknownAggregation struct {
	metricdata.Aggregation
}

func TestAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, nil, nil)
	AssertAggregationsEqual(t, sumInt64A, sumInt64A)
	AssertAggregationsEqual(t, sumFloat64A, sumFloat64A)
	AssertAggregationsEqual(t, gaugeInt64A, gaugeInt64A)
	AssertAggregationsEqual(t, gaugeFloat64A, gaugeFloat64A)
	AssertAggregationsEqual(t, histogramInt64A, histogramInt64A)
	AssertAggregationsEqual(t, histogramFloat64A, histogramFloat64A)

	r := equalAggregations(sumInt64A, nil, config{})
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	r = equalAggregations(sumInt64A, gaugeInt64A, config{})
	assert.Len(t, r, 1, "should return with type mismatch only")

	r = equalAggregations(unknownAggregation{}, unknownAggregation{}, config{})
	assert.Len(t, r, 1, "should return with unknown aggregation only")

	r = equalAggregations(sumInt64A, sumInt64B, config{})
	assert.Greaterf(t, len(r), 0, "sums should not be equal: %v == %v", sumInt64A, sumInt64B)

	r = equalAggregations(sumInt64A, sumInt64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "sums should be equal: %v", r)

	r = equalAggregations(sumFloat64A, sumFloat64B, config{})
	assert.Greaterf(t, len(r), 0, "sums should not be equal: %v == %v", sumFloat64A, sumFloat64B)

	r = equalAggregations(sumFloat64A, sumFloat64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "sums should be equal: %v", r)

	r = equalAggregations(gaugeInt64A, gaugeInt64B, config{})
	assert.Greaterf(t, len(r), 0, "gauges should not be equal: %v == %v", gaugeInt64A, gaugeInt64B)

	r = equalAggregations(gaugeInt64A, gaugeInt64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "gauges should be equal: %v", r)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64B, config{})
	assert.Greaterf(t, len(r), 0, "gauges should not be equal: %v == %v", gaugeFloat64A, gaugeFloat64B)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "gauges should be equal: %v", r)

	r = equalAggregations(histogramInt64A, histogramInt64B, config{})
	assert.Greaterf(t, len(r), 0, "histograms should not be equal: %v == %v", histogramInt64A, histogramInt64B)

	r = equalAggregations(histogramInt64A, histogramInt64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "histograms should be equal: %v", r)

	r = equalAggregations(histogramFloat64A, histogramFloat64B, config{})
	assert.Greaterf(t, len(r), 0, "histograms should not be equal: %v == %v", histogramFloat64A, histogramFloat64B)

	r = equalAggregations(histogramFloat64A, histogramFloat64C, config{ignoreTimestamp: true})
	assert.Len(t, r, 0, "histograms should be equal: %v", r)
}

func TestAssertAttributes(t *testing.T) {
	AssertHasAttributes(t, minFloat64A, attribute.Bool("A", true)) // No-op, always pass.
	AssertHasAttributes(t, exemplarInt64A, attribute.Bool("filter A", true))
	AssertHasAttributes(t, exemplarFloat64A, attribute.Bool("filter A", true))
	AssertHasAttributes(t, dataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, dataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramDataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramDataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, metricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, scopeMetricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("A", true))

	r := hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "gaugeInt64A has A=True")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "gaugeFloat64A has A=True")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "sumInt64A has A=True")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "sumFloat64A has A=True")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "histogramInt64A has A=True")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "histogramFloat64A has A=True")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "gaugeInt64A does not have A=False")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "gaugeFloat64A does not have A=False")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "sumInt64A does not have A=False")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "sumFloat64A does not have A=False")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "histogramInt64A does not have A=False")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "histogramFloat64A does not have A=False")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "gaugeInt64A does not have Attribute B")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "gaugeFloat64A does not have Attribute B")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "sumInt64A does not have Attribute B")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "sumFloat64A does not have Attribute B")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "histogramIntA does not have Attribute B")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "histogramFloatA does not have Attribute B")
}

func TestAssertAttributesFail(t *testing.T) {
	fakeT := &testing.T{}
	assert.False(t, AssertHasAttributes(fakeT, dataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, dataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, exemplarInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, exemplarFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, sumInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, sumFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("B", true)))

	sum := metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []metricdata.DataPoint[int64]{
			dataPointInt64A,
			dataPointInt64B,
		},
	}
	assert.False(t, AssertHasAttributes(fakeT, sum, attribute.Bool("A", true)))
}

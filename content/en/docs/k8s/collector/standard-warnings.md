# 标准的警告

某些组件具有可能导致问题的场景。有些组件需要以特定的方式与收集器进行交互，以确保
组件按预期工作。本文档描述了可能影响收集器中组件的常见警告。

访问组件的 README，看看它是否受到这些标准警告的影响。

## Unsound Transformations

Incorrect usage of the component may lead to telemetry data that is unsound i.e.
not spec-compliant/meaningless. This would most likely be caused by converting
metric data types or creating new metrics from existing metrics.

## Statefulness

The component keeps state related to telemetry data and therefore needs all data
from a producer to be sent to the same Collector instance to ensure a correct
behavior. Examples of scenarios that require state would be computing/exporting
delta metrics, tail-based sampling and grouping telemetry.

## Identity Conflict

The component may change the
[identity of a metric](https://github.com/open-telemetry/opentelemetry-specification/blob/main//specification/metrics/data-model.md#opentelemetry-protocol-data-model-producer-recommendations)
or the
[identity of a timeseries](https://github.com/open-telemetry/opentelemetry-specification/blob/main//specification/metrics/data-model.md#timeseries-model).
This could be done by modifying the metric/timeseries's name, attributes, or
instrumentation scope. Modifying a metric/timeseries's identity could result in
a metric/timeseries identity conflict, which caused by two metrics/timeseries
sharing the same name, attributes, and instrumentation scope.

## Orphaned Telemetry

The component modifies the incoming telemetry in such a way that a span becomes
orphaned, that is, it contains a `trace_id` or `parent_span_id` that does not
exist. This may occur because the component can modify `span_id`, `trace_id`, or
`parent_span_id` or because the component can delete telemetry.

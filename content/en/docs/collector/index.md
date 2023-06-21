---
title: 收集器
description: Vendor-agnostic way to receive, process and export telemetry data.
aliases: [/docs/collector/about]
cascade:
  collectorVersion: 0.79.0
weight: 10
---

![OpenTelemetry Collector diagram with Jaeger, OTLP and Prometheus integration](img/otel-collector.svg)

## 介绍

The OpenTelemetry Collector offers a vendor-agnostic implementation of how to
receive, process and export telemetry data. It removes the need to run, operate,
and maintain multiple agents/collectors. This works with improved scalability
and supports open source observability data formats (e.g. Jaeger, Prometheus,
Fluent Bit, etc.) sending to one or more open source or commercial back-ends.
The local Collector agent is the default location to which instrumentation
libraries export their telemetry data.

## 目标

- _Usability_: 合理的默认配置，支持流行的协议，开箱即用的运行和收集。
- _Performance_: Highly stable and performant under varying loads and
  configurations.
- _Observability_: An exemplar of an observable service.
- _Extensibility_: Customizable without touching the core code.
- _Unification_: Single codebase, deployable as an agent or collector with
  support for traces, metrics, and logs (future).

## 何时使用收集器

对于大多数特定于语言的工具库，您都有针对流行后端和OTLP的导出器。你可能会想，

> 在什么情况下使用收集器发送数据，而不是让每个服务直接发送到后端?

对于尝试和开始使用OpenTelemetry，将数据直接发送到后端是快速获取价值的好方法。
此外，在开发或小规模环境中，您可以在没有收集器的情况下获得不错的结果。

但是，通常我们建议在服务旁边使用收集器，因为它允许您的服务快速卸载数据，并且收集器可以处理
额外的处理，如重试，批处理，加密，甚至敏感数据过滤。

[设置收集器](./getting-started.md)也比您想象的要容易:每种语言的默认OTLP导出器都假定有一个本地收集器端点，因此您启动收集器并开始进行遥测。

## 状态和发布

The **collector** status is: [mixed][], since core collector components
currently have mixed [stability levels][].

**Collector components** differ in their maturity levels. An effort is underway
to ensure that every component has its stability documented. To track the
progress of this effort, see `opentelemetry-collector-contrib` [issue #10116][].

{{% latest_release "collector-releases" /%}}

[issue #10116]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/10116
[mixed]: /docs/specs/otel/document-status/#mixed
[stability levels]:
  https://github.com/open-telemetry/opentelemetry-collector#stability-levels

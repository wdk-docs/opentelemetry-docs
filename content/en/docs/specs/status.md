---
title: 规范状态摘要
linkTitle: Status
aliases: [/docs/specs/otel/status]
weight: 10
---

OpenTelemetry 是在一个信号一个信号的基础上开发的。跟踪、度量、行李和日志记录都是
信号的例子。信号建立在上下文传播之上，上下文传播是一种跨分布式系统关联数据的共享
机制。

每个信号由四个[核心成分](../concepts/components.md)组成:

- APIs
- SDKs
- [OpenTelemetry Protocol](/docs/specs/otlp/README.md) (OTLP)
- [Collector](/docs/collector/)

信号也有贡献组件，一个插件和仪器的生态系统。所有检测工具都共享相同的语义约定，以
确保它们在观察常见操作(如 HTTP 请求)时生成相同的数据。

要了解有关信号和组件的更多信息，请参阅 OTel 规范[概述](./otel/overview.md).

## 组件生命周期

组件遵循开发生命周期:草稿、实验、稳定、弃用、删除。

- **Draft** components are under design, and have not been added to the
  specification.
- **Experimental** components are released and available for beta testing.
- **Stable** components are backwards compatible and covered under long term
  support.
- **Deprecated** components are stable but may eventually be removed.

有关生命周期和长期支持的完整定义，请参
见[版本控制和稳定性]。(/docs/specs/otel/versioning-and-stability/).

## 当前的状态

以下是当前可用信号的高级状态报告。请注意，虽然 OpenTelemetry 客户端遵循共享规范
，但它们是独立开发的。

Checking the current status for each client in the README of its
[github repo](https://github.com/open-telemetry) is recommended. Client support
for specific features can be found in the
[specification compliance tables](https://github.com/open-telemetry/opentelemetry-specification/blob/main/spec-compliance-matrix.md).

Note that, for each of the following sections, the **Collector** status is the
same as the **Protocol** status.

### [Tracing][]

- {{% spec_status "API" "otel/trace/api" "Status" %}}
- {{% spec_status "SDK" "otel/trace/sdk" "Status" %}}
- {{% spec_status "Protocol" "otlp" "Status" %}}
- Notes:
  - The tracing specification is now completely stable, and covered by long term
    support.
  - The tracing specification is still extensible, but only in a backwards
    compatible manner.
  - OpenTelemetry clients are versioned to v1.0 once their tracing
    implementation is complete.

### [Metrics][]

- {{% spec_status "API" "otel/metrics/api" "Status" %}}
- {{% spec_status "SDK" "otel/metrics/sdk" "Status" %}}
- {{% spec_status "Protocol" "otlp" "Status" %}}
- Notes:
  - OpenTelemetry Metrics is currently under active development.
  - The data model is stable and released as part of the OTLP protocol.
  - Experimental support for metric pipelines is available in the Collector.
  - Collector support for Prometheus is under development, in collaboration with
    the Prometheus community.

### [Baggage][]

- {{% spec_status "API" "otel/baggage/api" "Status" %}}
- **SDK:** stable
- **Protocol:** N/A
- Notes:
  - OpenTelemetry Baggage is now completely stable.
  - Baggage is not an observability tool, it is a system for attaching arbitrary
    keys and values to a transaction, so that downstream services may access
    them. As such, there is no OTLP or Collector component to baggage.

### [Logging][]

- {{% spec_status "Bridge API" "otel/logs/bridge-api" "Status" %}}
- {{% spec_status "SDK" "otel/logs/sdk" "Status" %}}
- {{% spec_status "Event API" "otel/logs/event-api" "Status" %}}
- {{% spec_status "Protocol" "otlp" "Status" %}}
- Notes:
  - The [logs data model][] is released as part of the OpenTelemetry Protocol.
  - Log processing for many data formats has been added to the Collector, thanks
    to the donation of Stanza to the OpenTelemetry project.
  - The OpenTelemetry Log Bridge API allows for writing appenders which bridge
    logs from existing log frameworks into OpenTelemetry. The Logs Bridge API is
    not meant to be called directly by end users. Log appenders are under
    development in many languages.
  - The OpenTelemetry Log SDK is the standard implementation of the Log Bridge
    API. Applications configure the SDK to indicate how logs are processed and
    exported (e.g. using OTLP).
  - The OpenTelemetry Event API allows log records to be emitted which conform
    to the [event semantic conventions][]. In contrast to the Log Bridge API,
    the Event API is intended to be called by end users. The Event API is under
    active development.

[baggage]: /docs/specs/otel/baggage/
[event semantic conventions]: /docs/specs/otel/logs/semantic_conventions/events/
[logging]: /docs/specs/otel/logs/
[logs data model]: /docs/specs/otel/logs/data-model/
[metrics]: /docs/specs/otel/metrics/
[tracing]: /docs/specs/otel/trace/

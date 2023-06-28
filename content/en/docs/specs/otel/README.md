<!--- Hugo front matter used to generate the website version of this page:
linkTitle: OTel spec
no_list: true
cascade:
  body_class: otel-docs-spec
  github_repo: &repo https://github.com/open-telemetry/opentelemetry-specification
  github_subdir: specification
  path_base_for_github_subdir: content/en/docs/specs/otel/
  github_project_repo: *repo
--->

# 规范

## 内容

- [概述](overview.md)
- [术语](glossary.md)
- [OpenTelemetry 客户端的版本控制和稳定性](versioning-and-stability.md)
- [图书馆的指导方针](library-guidelines.md)
  - [包/库布局](library-layout.md)
  - [一般错误处理准则](error-handling.md)
- API 规范
  - [上下文](context/README.md)
    - [Propagators](context/api-propagators.md)
  - [Baggage](baggage/api.md)
  - [追踪](trace/api.md)
  - [指标](metrics/api.md)
  - 日之惠
    - [Bridge API](logs/bridge-api.md)
    - [Event API](logs/event-api.md)
- SDK 规范
  - [追踪](trace/sdk.md)
  - [指标](metrics/sdk.md)
  - [日志](logs/sdk.md)
  - [资源](resource/sdk.md)
  - [配置](configuration/sdk-configuration.md)
- Data 规范
  - [语义约定](overview.md#semantic-conventions)
  - [协议](protocol/README.md)
    - [指标](metrics/data-model.md)
    - [日志](logs/data-model.md)
  - 兼容
    - [OpenCensus](compatibility/opencensus.md)
    - [OpenTracing](compatibility/opentracing.md)
    - [Prometheus 和 OpenMetrics](compatibility/prometheus_and_openmetrics.md)
    - [以非 otlp 日志格式跟踪上下文](compatibility/logging_trace_context.md)

## 符号约定和遵从性

[规范][]中的关键字"MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL"
应按照[BCP 14](https://tools.ietf.org/html/bcp14)
[[RFC2119](https://tools.ietf.org/html/rfc2119)]
[[RFC8174](https://tools.ietf.org/html/rfc8174)]中所描述的解释，当且仅当它们以大
写字母出现时，如下所示。

如果[规范][]的实现不能满足[规范][]中定义的"MUST", "MUST NOT", "REQUIRED",
"SHALL", or "SHALL NOT"中的一个或多个要求，则该[规范][]的实现不合规。相反，如
果[规范][]的实现满足[规范][]中定义的所有"MUST", "MUST NOT", "REQUIRED", "SHALL",
and "SHALL NOT"需求，则该实现是兼容的。

## 项目命名

- 官方项目名称是“OpenTelemetry”(“Open”和“Telemetry”之间没有空格)。
- OpenTelemetry 项目使用的官方缩写是“OTel”。避免使用“OT”，以避免与现在已弃用的
  “OpenTracing”项目混淆。
- 子项目的官方名称，如语言特定实现，遵循“OpenTelemetry {编程语言，运行时或组件的
  名称}”的模式，例如，“OpenTelemetry Python”，“OpenTelemetry .NET”或
  “OpenTelemetry Collector”。

## 项目简介

请参阅[项目存储库][]了解以下信息，以及更多:

- [Change / contribution process](../README.md#change--contribution-process)
- [Project timeline](../README.md#project-timeline)
- [Versioning the specification](../README.md#versioning-the-specification)
- [License](../README.md#license)

[项目存储库]: https://github.com/open-telemetry/opentelemetry-specification
[规范]: overview.md

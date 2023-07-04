---
title: 发行版
weight: 25
---

OpenTelemetry 项目目前提供了该收集器的[预构建发行版][]。 [发行版][]中包含的组件
可以在每个发行版的`manifest.yaml`中找到。

[预构建发行版]:
  https://github.com/open-telemetry/opentelemetry-collector-releases/releases
[发行版]:
  https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions

## 自定义的发行版

由于各种原因，OpenTelemetry 项目提供的现有发行版可能无法满足您的需求。无论您是想
要更小的版本，还是需要实现自定义功能，如[自定义身份验证器](../custom-auth)、接收
器、处理器或导出器。用于构建发行版的工具[ocb][] (OpenTelemetry Collector
Builder)可用于构建您自己的发行版。

[ocb]:
  https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder

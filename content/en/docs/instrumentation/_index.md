---
title: 仪表
description: >-
  许多流行的编程语言都支持OpenTelemetry代码插装
weight: 2
---

OpenTelemetry代码[植入][]支持下面列出的语言。
根据语言的不同，所涵盖的主题将包括以下部分或全部:

- 自动仪表
- 手动工具
- 导出数据

如果你正在使用Kubernetes，你可以使用[OpenTelemetry Operator for Kubernetes][otel-op]来为Java, Node.js和Python注入[自动检测库][auto]。

## 状态和发布

OpenTelemetry主要功能组件的现状如下:

{{% telemetry_support_table " " %}}

[auto]:
  https://github.com/open-telemetry/opentelemetry-operator#opentelemetry-auto-instrumentation-injection
[植入]: /docs/concepts/instrumentation/
[otel-op]: https://github.com/open-telemetry/opentelemetry-operator

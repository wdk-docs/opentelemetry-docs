---
title: Instrumentation 示例
title: 示例
aliases: [/docs/instrumentation/js/instrumentation_examples]
weight: 220
---

下面是一些 OpenTelemetry 仪器示例的资源。

## 核心的例子

OpenTelemetry 的[JavaScript 版本][repo]的存储库包含一些[示例]如何运行真实的应用
程序与 OpenTelemetry JavaScript。

[repo]: https://github.com/open-telemetry/opentelemetry-js
[示例]: https://github.com/open-telemetry/opentelemetry-js/tree/main/examples

## 插件和包示例

OpenTelemetry JavaScript 的许多包和插件在[贡献库]提供了一个使用示例。你可以
在[示例文件夹]中找到它们。

[贡献库]: https://github.com/open-telemetry/opentelemetry-js-contrib
[示例文件夹]:
  https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/examples

## 社区资源

[nodejs-opentelemetry-tempo][tempo]项目说明了 OpenTelemetry 的使用(通过自动和手
动检测)，涉及带有 DB 交互的微服务。它使用以下内容:

- [Prometheus](https://prometheus.io), 用于监视和警报
- [Loki](https://grafana.com/oss/loki/), 分布式日志
- [Tempo](https://grafana.com/oss/tempo/), 用于分布式跟踪
- [Grafana](https://grafana.com/grafana/) 对可视化

要了解更多细节，请访问[项目存储库][tempo]。

[tempo]: https://github.com/mnadeem/nodejs-opentelemetry-tempo

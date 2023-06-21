---
title: Instrumentation 示例
title: 示例
aliases: [/docs/instrumentation/js/instrumentation_examples]
weight: 220
---

下面是一些OpenTelemetry仪器示例的资源。

## 核心的例子

OpenTelemetry的[JavaScript版本][repo]的存储库包含一些[示例][]如何运行真实的应用程序与OpenTelemetry JavaScript。

[repo]: https://github.com/open-telemetry/opentelemetry-js
[examples]: https://github.com/open-telemetry/opentelemetry-js/tree/main/examples

## 插件和包示例

OpenTelemetry JavaScript的许多包和插件在[贡献库][]提供了一个使用示例。
你可以在[examples文件夹][]中找到它们。

[contributions repository]: https://github.com/open-telemetry/opentelemetry-js-contrib
[examples folder]: https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/examples

## 社区资源

[nodejs-opentelemetry-tempo][tempo]项目说明了OpenTelemetry的使用(通过自动和手动检测)，涉及带有DB交互的微服务。
它使用以下内容:

- [Prometheus](https://prometheus.io), for monitoring and alerting
- [Loki](https://grafana.com/oss/loki/), for distributed logging
- [Tempo](https://grafana.com/oss/tempo/), for distributed tracing
- [Grafana](https://grafana.com/grafana/) for visualization

要了解更多细节，请访问[项目存储库][tempo]。

[tempo]: https://github.com/mnadeem/nodejs-opentelemetry-tempo

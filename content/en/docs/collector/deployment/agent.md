---
title: 代理
description: 为什么以及如何向收集器发送信号，然后再从收集器发送到后端
weight: 2
---

代理收集器部署模式由应用程序(使用 [OpenTelemetry 协议(OTLP)][otlp]使用
OpenTelemetry SDK 进行检测)或其他收集器(使用 OTLP 导出器)[组
成][instrumentation]，这些收集器将遥测信号发送到与应用程序一起运行的[收集
器][collector]实例或与应用程序在同一台主机上运行的收集器实例(例如 sidecar 或守护
进程)。

每个客户端 SDK 或下游收集器都配置了一个收集器位置:

![Decentralized collector deployment concept](../../img/otel_agent_sdk_v2.svg)

1. 在应用程序中，SDK 被配置为将 OTLP 数据发送到收集器。
2. 采集器配置为将遥测数据发送到一个或多个后端。

## 示例

代理收集器部署模式的具体示例如下:您可以使用 OpenTelemetry Java SDK 手动设置一
个[用于导出指标的 Java 应用程序][instrument-java-metrics]。在应用程序的上下文中
，你可以将`OTEL_METRICS_EXPORTER`设置为`otlp`(这是默认值)，并使用收集器的地址配
置[otlp 导出器][otlp-exporter]，例如(在 Bash 或 `zsh` shell 中):

```
export OTEL_EXPORTER_OTLP_ENDPOINT=http://collector.example.com:4318
```

服务于`collector.example.com:4318`的收集器将被配置如下:

=== "Traces"

    ```yml
    receivers:
      otlp: # the OTLP receiver the app is sending traces to
        protocols:
          grpc:

    processors:
      batch:

    exporters:
      jaeger: # the Jaeger exporter, to ingest traces to backend
        endpoint: 'https://jaeger.example.com:14250'
        tls:
          insecure: true

    service:
      pipelines:
        traces/dev:
          receivers: [otlp]
          processors: [batch]
          exporters: [jaeger]
    ```

=== "Metrics"

    ```yml
    receivers:
      otlp: # the OTLP receiver the app is sending metrics to
        protocols:
          grpc:

    processors:
      batch:

    exporters:
      prometheusremotewrite: # the PRW exporter, to ingest metrics to backend
        endpoint: 'https://prw.example.com/v1/api/remote_write'

    service:
      pipelines:
        metrics/prod:
          receivers: [otlp]
          processors: [batch]
          exporters: [prometheusremotewrite]
    ```

=== "Logs"

    ```yml
    receivers:
      otlp: # the OTLP receiver the app is sending logs to
        protocols:
          grpc:

    processors:
      batch:

    exporters:
      file: # the File Exporter, to ingest logs to local file
        path: "./app42_example.log"
        rotation:

    service:
      pipelines:
        logs/dev:
          receivers: [otlp]
          processors: [batch]
          exporters: [file]
    ```

如果您想亲自尝试一下，可以看看端到端
的[Java][java-otlp-example]或[Python][py-otlp-example]示例。

## 权衡

有点:

- 入门很简单
- 清除应用程序与采集器 1:1 的映射关系

缺点:

- 可伸缩性(人力和负载方面)
- 僵化的

[instrumentation]: ../../instrumentation/index.md
[otlp]: ../../specs/otel/protocol/index.md
[collector]: ../index.md
[instrument-java-metrics]: ../../instrumentation/java/manual.md#metrics
[otlp-exporter]: ../../specs/otel/protocol/exporter.md
[java-otlp-example]:
  https://github.com/open-telemetry/opentelemetry-java-docs/tree/main/otlp
[py-otlp-example]:
  https://opentelemetry-python.readthedocs.io/en/stable/examples/metrics/instruments/README.html
[lb-exporter]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/loadbalancingexporter
[spanmetrics-processor]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/spanmetricsprocessor

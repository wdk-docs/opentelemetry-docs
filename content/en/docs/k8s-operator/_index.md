---
title: Kubernetes开放遥测操作器
linkTitle: K8s Operator
weight: 11
description: Kubernetes操作符的实现，它使用开放的遥测仪器库管理收集器和工作负载的自动检测。
spelling: cSpell:ignore Otel
aliases: [/docs/operator]
---

## 介绍

OpenTelemetry算子是[Kubernetes算子](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)的一个实现。

运营商管理:

- OpenTelemetry 收集器
- [使用OpenTelemetry工具库自动检测工作负载](https://github.com/open-telemetry/opentelemetry-operator#opentelemetry-auto-instrumentation-injection)

## 开始

要在现有集群中安装操作符，请确保安装了cert-manager并运行:

```bash
$ kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
```

一旦`opentelemetry-operator`部署就绪，创建一个开放遥测采集器(otelcol)实例，如下所示:

```bash
$ kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: simplest
spec:
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:
    processors:

    exporters:
      logging:

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: []
          exporters: [logging]
EOF
```

要了解更多配置选项，以及使用OpenTelemetry工具库设置工作负载的自动检测注入，请继续阅读[这里](https://github.com/open-telemetry/opentelemetry-operator/blob/main/README.md).

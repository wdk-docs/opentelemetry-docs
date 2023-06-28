---
title: 注入自动插装
linkTitle: Auto-instrumentation
weight: 11
description: 使用开放式遥测操作器的自动插装实现。
spelling: cSpell:ignore Otel
---

开放遥测操作器支持为.NET, Java, Nodejs 和 Python 服务注入和配置自动仪器库。

## 安装

首先，
将[OpenTelemetry Operator](https://github.com/open-telemetry/opentelemetry-operator)安
装到集群中。

你可以通
过[操作员释放清单](https://github.com/open-telemetry/opentelemetry-operator#getting-started)，
[操作员掌舵图](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator#opentelemetry-operator-helm-chart)，
或与[操作员中心](https://operatorhub.io/operator/opentelemetry-operator)。

在大多数情况下，您需要安
装[cert-manager](https://cert-manager.io/docs/installation/)。如果使用 helm 图表
，则可以选择生成自签名证书。

## 创建收集器(可选)

将遥测数据从容器发送到[收集器](../../collector/index.md)而不是直接发送到后端是最
佳实践。

收集器有助于简化秘密管理，从应用程序中解耦数据导出问题( 例如需要重试)，并允许您
向遥测添加额外的数据，例如使用[k8sattributesprocessor]组件。

!!! note

    如果您选择不使用收集器，则可以跳到下一节。

[k8sattributesprocessor]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/k8sattributesprocessor

Operator 为收集器提供一个[自定义资源定义(CRD)][CRD]，用于创建一个由 Operator 管
理的收集器实例。

下面的示例将收集器部署为`Deployment`(默认设置)，但也可以使用其他[部署模式]。

[CRD]:
  https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md#opentelemetrycollector
[部署模式]:
  https://github.com/open-telemetry/opentelemetry-operator#deployment-modes

当使用`Deployment`模式时，Operator 还将创建一个可用于与 Collector 交互的服务。服
务的名称是附加在`-collector`前面的`OpenTelemetryCollector`资源的名称。在我们的例
子中，这将是`demo-collector`。

```bash
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: demo
spec:
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:
    processors:
      memory_limiter:
        check_interval: 1s
        limit_percentage: 75
        spike_limit_percentage: 15
      batch:
        send_batch_size: 10000
        timeout: 10s

    exporters:
      logging:

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [memory_limiter, batch]
          exporters: [logging]
        metrics:
          receivers: [otlp]
          processors: [memory_limiter, batch]
          exporters: [logging]
        logs:
          receivers: [otlp]
          processors: [memory_limiter, batch]
          exporters: [logging]
EOF
```

上面的命令会执行 Collector 的部署，您可以将其用作 Pod 中自动检测的端点。

## 配置自动插装

为了能够管理自动插装，操作人员需要进行配置，以了解要对哪些 Pod 进行插装以及对这
些 Pod 使用哪种自动插装。这是通过[插装 CRD]完成的。

[插装 CRD]:
  https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md#instrumentation

正确创建插装资源对于使自动检测工作起着至关重要的作用。要使自动检测正常工作，需要
确保所有端点和环境变量都正确。

### .NET

下面的命令将创建一个基本的插装资源，该资源是专门为插装 .NET 服务配置的。

```bash
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4318
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: "1"
EOF
```

默认情况下，auto-instruments .NET 服务使用`otlp`和`http/protobuf`协议。这意味着
配置的端点必须能够通过`http/protobuf`接收 OTLP。因此，该示例使
用`http://demo-collector:4318`，它将连接到在上一步中创建的收集器的 otlreceiver
的`http`端口。

默认情况下，.NET 自动检测附带了[许多检测库]。这使检测变得容易，但可能导致过多或
不需要的数据。如果有任何你不想使用的库，你可以设
置`OTEL_DOTNET_AUTO_[SIGNAL]_[NAME]_INSTRUMENTATION_ENABLED=false` ，其
中`[SIGNAL]`是信号的类型，`[NAME]` 是区分大小写的库名。

[许多检测库]:
  https://github.com/open-telemetry/opentelemetry-dotnet-instrumentation/blob/main/docs/config.md#instrumentations

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4318
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: '1'
  dotnet:
    env:
      - name: OTEL_DOTNET_AUTO_TRACES_GRPCNETCLIENT_INSTRUMENTATION_ENABLED
        value: false
      - name: OTEL_DOTNET_AUTO_METRICS_PROCESS_INSTRUMENTATION_ENABLED
        value: false
```

有关更多详细信息，请参
阅[.NET Auto Instrumentation 文档](/docs/instrumentation/net/automatic/).

### Java

下面的命令创建一个基本的 Instrumentation 资源，该资源被配置为检测 Java 服务。

```bash
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4317
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: "1"
EOF
```

By default, the Instrumentation resource that auto-instruments Java services
uses `otlp` with the `grpc` protocol. This means that the configured endpoint
must be able to receive OTLP over `grpc`. Therefore, the example uses
`http://demo-collector:4317`, which connects to the `grpc` port of the
otlpreceiver of the Collector created in the previous step.

By default, the Java auto-instrumentation ships with
[many instrumentation libraries](/docs/instrumentation/java/automatic/#supported-libraries-frameworks-application-services-and-jvms).
This makes instrumentation easy, but could result in too much or unwanted data.
If there are any libraries you do not want to use you can set the
`OTEL_INSTRUMENTATION_[NAME]_ENABLED=false` where `[NAME]` is the name of the
library. If you know exactly which libraries you want to use, you can disable
the default libraries by setting
`OTEL_INSTRUMENTATION_COMMON_DEFAULT_ENABLED=false` and then use
`OTEL_INSTRUMENTATION_[NAME]_ENABLED=true` where `[NAME]` is the name of the
library. For more details, see
[Suppressing specific auto-instrumentation](/docs/instrumentation/java/automatic/agent-config/#suppressing-specific-auto-instrumentation).

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4317
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: '1'
  java:
    env:
      - name: OTEL_INSTRUMENTATION_KAFKA_ENABLED
        value: false
      - name: OTEL_INSTRUMENTATION_REDISCALA_ENABLED
        value: false
```

有关详细信息，请参
见[Java 代理配置](/docs/instrumentation/java/automatic/agent-config/).

### Node.js

下面的命令创建了一个基本的插装资源，该资源配置用于检测`Node.js`服务。

```bash
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4317
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: "1"
EOF
```

默认情况下，自动检测 Node.js 服务的`Instrumentation`资源使用 `otlp` 和 `grpc` 协
议。这意味着配置的端点必须能够通过 `grpc` 接收 `otlp` 。因此，本例使用
`http://demo-collector:4317`，它连接到在上一步中创建的收集器的`otlreceiver`的
`grpc` 端口。

默认情况下，Node.js 自动检测附带了许多[检测库]。目前，还没有办法选择只加入特定的
软件包或禁用特定的软件包。如果您不想使用默认映像包含的包，那么您必须提供自己的映
像，该映像只包含您想要的包，或者使用手动检测。

更多细节请参见[Node.js 自动插装].

[检测库]:
  https://github.com/open-telemetry/opentelemetry-js-contrib/blob/main/metapackages/auto-instrumentations-node/README.md#supported-instrumentations
[Node.js 自动插装]:
  ../../instrumentation/js/libraries.md#node-autoinstrumentation-package

### Python

下面的命令将创建一个基本的 Instrumentation 资源，该资源是专门为 Instrumentation
Python 服务配置的。

```bash
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4318
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: "1"
EOF
```

By default, the Instrumentation resource that auto-instruments python services
uses `otlp` with the `http/protobuf` protocol. This means that the configured
endpoint must be able to receive OTLP over `http/protobuf`. Therefore, the
example uses `http://demo-collector:4318`, which will connect to the `http` port
of the otlpreceiver of the Collector created in the previous step.

> As of operator v0.67.0, the Instrumentation resource automatically sets
> `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` and `OTEL_EXPORTER_OTLP_METRICS_PROTOCOL`
> to `http/protobuf` for Python services. If you use an older version of the
> Operator you **MUST** set these env variables to `http/protobuf`, or python
> auto-instrumentation will not work.

By default the Python auto-instrumentation will detect the packages in your
Python service and instrument anything it can. This makes instrumentation easy,
but can result in too much or unwanted data. If there are any packages you do
not want to instrument, you can set the `OTEL_PYTHON_DISABLED_INSTRUMENTATIONS`
environment variable

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: demo-instrumentation
spec:
  exporter:
    endpoint: http://demo-collector:4318
  propagators:
    - tracecontext
    - baggage
  sampler:
    type: parentbased_traceidratio
    argument: '1'
  python:
    env:
      - name: OTEL_PYTHON_DISABLED_INSTRUMENTATIONS
        value:
          <comma-separated list of package names to exclude from
          instrumentation>
```

有关更多细节，请参阅
[Python Agent Configuration 文档](/docs/instrumentation/python/automatic/agent-config/#disabling-specific-instrumentations)

---

既然已经创建了 Instrumentation 对象，那么集群就能够自动检测服务并将数据发送到端
点。然而，OpenTelemetry Operator 的自动检测遵循一个可选择的模型。为了激活自动检
测，您需要在部署中添加一个注释。

## 向现有部署添加注释

最后一步是选择自动检测服务。这是通过更新你的服务
的`spec.template.metadata.annotations`来包含一个特定于语言的注释来实现的:

- .NET: `instrumentation.opentelemetry.io/inject-dotnet: "true"`
- Java: `instrumentation.opentelemetry.io/inject-java: "true"`
- Node.js: `instrumentation.opentelemetry.io/inject-nodejs: "true"`
- Python: `instrumentation.opentelemetry.io/inject-python: "true"`

注释的可能值可以是

- `"true"` - 以当前命名空间的默认名称注入`Instrumentation`资源。
- `"my-instrumentation"` - 在当前命名空间中注入名为`"my-instrumentation"`的
  `Instrumentation` CR 实例。
- `"my-other-namespace/my-instrumentation"` - 从另一个命名空
  间`"my-other-namespace"`注入名为`"my-instrumentation"`的 `Instrumentation` CR
  实例.
- `"false"` - 不要注射

或者，可以将注释添加到名称空间中，这将导致该名称空间中的所有服务选择加入自动检测
。请参阅[操作员自动检测文档][otel-auto-instr]了解更多详细信息。

[otel-auto-instr]:
  https://github.com/open-telemetry/opentelemetry-operator/blob/main/README.md#opentelemetry-auto-instrumentation-injection

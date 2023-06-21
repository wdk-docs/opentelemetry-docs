---
title: 注入自动仪表
linkTitle: Auto-instrumentation
weight: 11
description:
  使用开放式遥测操作器的自动仪表实现。
spelling: cSpell:ignore Otel
---

开放遥测操作器支持为.NET, Java, Nodejs和Python服务注入和配置自动仪器库。

## 安装

首先，将[OpenTelemetry Operator](https://github.com/open-telemetry/opentelemetry-operator)安装到集群中。

你可以通过[操作员释放清单](https://github.com/open-telemetry/opentelemetry-operator#getting-started)，
[操作员掌舵图](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator#opentelemetry-operator-helm-chart)，或与[操作员中心](https://operatorhub.io/operator/opentelemetry-operator)。

在大多数情况下，您需要安装[cert-manager](https://cert-manager.io/docs/installation/)。
如果使用helm图表，则可以选择生成自签名证书。

## 创建OpenTelemetry收集器(可选)

将遥测数据从容器发送到[OpenTelemetry Collector](../../ Collector /)而不是直接发送到后端是最佳实践。
Collector有助于简化秘密管理，从应用程序中解耦数据导出问题(例如需要重试)，并允许您向遥测添加额外的数据，例如使用[k8sattributesprocessor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/k8sattributesprocessor)组件。
如果您选择不使用收集器，则可以跳到下一节。

Operator为OpenTelemetry Collector提供一个[自定义资源定义(CRD)](https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md#opentelemetrycollector)，用于创建一个由Operator管理的Collector实例。
下面的示例将Collector部署为部署(默认设置)，但也可以使用其他[部署模式](https://github.com/open-telemetry/opentelemetry-operator#deployment-modes)。

当使用`Deployment`模式时，操作员还将创建一个可用于与收集器交互的服务。
服务的名称是附加在`-collector`前面的`OpenTelemetryCollector`资源的名称。
在我们的例子中，这将是`demo-collector`。

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

上面的命令会导致Collector的部署，您可以将其用作pod中自动检测的端点。

## 配置Autoinstrumentation

为了能够管理自动仪表，操作人员需要进行配置，以了解要对哪些吊舱进行仪表以及对这些吊舱使用哪种自动仪表。
这是通过[仪表CRD](https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md#instrumentation)完成的。

正确创建Instrumentation资源对于使自动检测工作起着至关重要的作用。
要使自动检测正常工作，需要确保所有端点和环境变量都正确。

### .NET

下面的命令将创建一个基本的Instrumentation资源，该资源是专门为Instrumentation .NET服务配置的。

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

By default, the Instrumentation resource that auto-instruments .NET services
uses `otlp` with the `http/protobuf` protocol. This means that the configured
endpoint must be able to receive OTLP over `http/protobuf`. Therefore, the
example uses `http://demo-collector:4318`, which will connect to the `http` port
of the otlpreceiver of the Collector created in the previous step.

By default, the .NET auto-instrumentation ships with
[many instrumentation libraries](https://github.com/open-telemetry/opentelemetry-dotnet-instrumentation/blob/main/docs/config.md#instrumentations).
This makes instrumentation easy, but could result in too much or unwanted data.
If there are any libraries you do not want to use you can set the
`OTEL_DOTNET_AUTO_[SIGNAL]_[NAME]_INSTRUMENTATION_ENABLED=false` where
`[SIGNAL]` is the type of the signal and `[NAME]` is the case-sensitive name of
the library.

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

For more details, see
[.NET Auto Instrumentation docs](/docs/instrumentation/net/automatic/).

### Java

The following command creates a basic Instrumentation resource that is
configured for instrumenting Java services.

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

For more details, see
[Java Agent Configuration](/docs/instrumentation/java/automatic/agent-config/).

### Node.js

下面的命令创建了一个基本的仪表资源，该资源配置用于检测`Node.js`服务。

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

默认情况下，自动检测Node.js服务的`Instrumentation`资源使用 `otlp` 和 `grpc` 协议。
这意味着配置的端点必须能够通过 `grpc` 接收 `OTLP` 。
因此，本例使用 `http://demo-collector:4317` ，它连接到在上一步中创建的收集器的`otlreceiver`的 `grpc` 端口。

默认情况下，Node.js自动检测附带了[许多检测库](https://github.com/open-telemetry/opentelemetry-js-contrib/blob/main/metapackages/auto-instrumentations-node/README.md#supported-instrumentations)。
目前，还没有办法选择只加入特定的软件包或禁用特定的软件包。
如果您不想使用默认映像包含的包，那么您必须提供自己的映像，该映像只包含您想要的包，或者使用手动检测。

更多细节请参见[Node.js auto-instrumentation](../../instrumentation/js/libraries/#node-autoinstrumentation-package).

### Python

The following command will create a basic Instrumentation resource that is
configured specifically for instrumenting Python services.

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

[See the Python Agent Configuration docs for more details.](/docs/instrumentation/python/automatic/agent-config/#disabling-specific-instrumentations)

---

Now that your Instrumentation object is created, your cluster has the ability to
auto-instrument services and send data to an endpoint. However,
auto-instrumentation with the OpenTelemetry Operator follows an opt-in model. In
order to activate autoinstrumentation, you'll need to add an annotation to your
deployment.

## 向现有部署添加注释

最后一步是选择自动检测服务。
这是通过更新你的服务的`spec.template.metadata.annotations`来包含一个特定于语言的注释来实现的:

- .NET: `instrumentation.opentelemetry.io/inject-dotnet: "true"`
- Java: `instrumentation.opentelemetry.io/inject-java: "true"`
- Node.js: `instrumentation.opentelemetry.io/inject-nodejs: "true"`
- Python: `instrumentation.opentelemetry.io/inject-python: "true"`

注释的可能值可以是

- `"true"` - to inject `Instrumentation` resource with default name from the
  current namespace.
- `"my-instrumentation"` - to inject `Instrumentation` CR instance with name
  `"my-instrumentation"` in the current namespace.
- `"my-other-namespace/my-instrumentation"` - to inject `Instrumentation` CR
  instance with name `"my-instrumentation"` from another namespace
  `"my-other-namespace"`.
- `"false"` - do not inject

或者，可以将注释添加到名称空间中，这将导致该名称空间中的所有服务选择加入自动检测。
请参阅[操作员自动检测文档](https://github.com/open-telemetry/opentelemetry-operator/blob/main/README.md#opentelemetry-auto-instrumentation-injection)了解更多详细信息。

---
title: Kubernetes 开发
linkTitle: Kubernetes
aliases: [/docs/demo/kubernetes_deployment]
---

我们提供了一
个[OpenTelemetry Demo Helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo)来
帮助将 Demo 部署到现有的 Kubernetes 集群上。

[Helm](https://helm.sh)必须安装才能使用海图。请参考 Helm
的[文档](https://helm.sh/docs/)开始。

## 先决条件

- Kubernetes 1.23+
- 为应用程序提供 4 GB 的空闲 RAM
- Helm 3.9+ (仅适用于舵机的安装方法

## 使用 Helm 安装(推荐)

添加 OpenTelemetry Helm 存储库:

```shell
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
```

要安装版本名为 my-otel-demo 的图表，请运行以下命令:

```shell
helm install my-otel-demo open-telemetry/opentelemetry-demo
```

> **Note** OpenTelemetry Demo Helm chart 0.11.0 或更高版本需要执行下面提到的所有
> 使用方法。

## 使用 kubectl 安装

下面的命令将演示应用程序安装到 Kubernetes 集群。

```shell
kubectl create namespace otel-demo
kubectl apply --namespace otel-demo -f https://raw.githubusercontent.com/open-telemetry/opentelemetry-demo/main/kubernetes/opentelemetry-demo.yaml
```

> **Note** 这些清单是从 Helm 图表生成的，提供这些清单是为了方便。建议使用 Helm
> 图进行安装。

## 使用演示

演示应用程序将需要在 Kubernetes 集群外部公开的服务才能使用它们。您可以使
用`kubectl port-forward`命令或通过配置服务类型(例如:LoadBalancer)配置可选部署的
入口资源，将服务公开给本地系统。

### 使用 kubectl 端口转发公开服务

要公开 frontendproxy 服务，使用以下命令(将 `my-otel-demo` 替换为相应的 Helm 图表
发布名称):

```shell
kubectl port-forward svc/my-otel-demo-frontendproxy 8080:8080
```

为了从浏览器中正确收集跨度，您还需要公开 OpenTelemetry Collector 的 OTLP/HTTP 端
口(将`my-otel-demo`替换为相应的 Helm 图表发布名称):

```shell
kubectl port-forward svc/my-otel-demo-otelcol 4318:4318
```

> **Note**: `kubectl port-forward`将代理该端口，直到进程终止。您可能需要为每次使
> 用`kubectl port-forward`,，并使用<kbd>Ctrl-C</kbd>在完成时终止进程。

设置了 frontendproxy 和 Collector 端口转发后，您可以访问:

- Webstore: <http://localhost:8080/>
- Grafana: <http://localhost:8080/grafana/>
- Feature Flags UI: <http://localhost:8080/feature/>
- Load Generator UI: <http://localhost:8080/loadgen/>
- Jaeger UI: <http://localhost:8080/jaeger/ui/>

### 使用服务类型配置公开服务

> **Note** Kubernetes 集群可能没有适当的基础设施组件来启用 LoadBalancer 服务类型
> 或入口资源。在使用这些配置选项之前，请验证您的集群具有适当的支持。

每个演示服务(例如:frontendproxy)都提供了一种配置 Kubernetes 服务类型的方法。默认
情况下，这些将是 `ClusterIP`，但您可以使用 `serviceType` 属性为每个服务更改每个
。

要配置 frontendproxy 服务使用 LoadBalancer 服务类型，你需要在你的值文件中指定以
下内容:

```yaml
components:
  frontendProxy:
    service:
      type: LoadBalancer
```

> **Note** 建议在安装 Helm 图表时使用 values 文件，以便指定其他配置选项。

Helm 图不提供创建入口资源的工具。如果需要，这些将需要在安装 Helm 图表后手动创建
。一些 Kubernetes 提供程序需要特定的服务类型才能被入口资源使用(例如:EKS ALB 入口
，需要 NodePort 服务类型)。

为了正确地收集来自浏览器的跨度，您还需要公开 OpenTelemetry Collector 的
OTLP/HTTP 端口，以便用户 web 浏览器可以访问。还必须使用
`PUBLIC_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`环境变量将公开 OpenTelemetry
Collector 的位置传递给前端服务。你可以在你的 values 文件中使用以下代码:

```yaml
components:
  frontend:
    env:
      - name: PUBLIC_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
        value: http://otel-demo-collector.mydomain.com:4318/v1/traces
```

To install the Helm chart with a custom `my-values-file.yaml` values file use:

```shell
helm install my-otel-demo open-telemetry/opentelemetry-demo --values my-values-file.yaml
```

With the frontendproxy and Collector exposed, you can access the demo UI at the
base path for the frontendproxy. Other demo components can be accessed at the
following sub-paths:

- Webstore: `/` (base)
- Grafana: `/grafana`
- Feature Flags UI: `/feature`
- Load Generator UI: `/loadgen/` (must include trailing slash)
- Jaeger UI: `/jaeger/ui`

## 自带后端

Likely you want to use the Webstore as a demo application for an observability
backend you already have (e.g. an existing instance of Jaeger, Zipkin, or one of
the [vendor of your choice](/ecosystem/vendors/).

The OpenTelemetry Collector's configuration is exposed in the Helm chart. Any
additions you do will be merged into the default configuration. You can use this
to add your own exporters, and add them to the desired pipeline(s)

```yaml
opentelemetry-collector:
  config:
    exporters:
      otlphttp/example:
        endpoint: <your-endpoint-url>

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [batch]
          exporters: [otlphttp/example]
```

> **Note** 当将 YAML 值与 Helm 合并时，对象被合并，数组被替换。

Vendor backends might require you to add additional parameters for
authentication, please check their documentation. Some backends require
different exporters, you may find them and their documentation available at
[opentelemetry-collector-contrib/exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter).

要用自定义的 `my-values-file.yaml` 值文件安装 Helm 图表，使用:

```shell
helm install my-otel-demo open-telemetry/opentelemetry-demo --values my-values-file.yaml
```

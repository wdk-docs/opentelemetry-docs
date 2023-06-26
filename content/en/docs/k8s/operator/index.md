[![Continuous Integration][github-workflow-img]][github-workflow]
[![Go Report Card][goreport-img]][goreport] [![GoDoc][godoc-img]][godoc]

# OpenTelemetry Operator for Kubernetes

OpenTelemetry Operator
是[Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)的
一个实现。

Operator 管理:

- [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)
- 使用 OpenTelemetry 工具库自动检测工作负载

## 文档

- [API docs](./apis/index.md)

## Helm Charts

您可以通
过[Helm Chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator)从
opentelemetry-helm-charts 存储库安装 Opentelemetry Operator。更多信息请访
问[这里](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator)。

## 入门

要在现有集群中安装操作器，请确保安装
了[`cert-manager`](https://cert-manager.io/docs/installation/)并运行:

```bash
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
```

一旦`opentelemetry-operator`部署就绪就可以创建一个 OpenTelemetry Collector
(otelcol)实例，如下所示:

```yaml
kubectl apply -f - <<EOF
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
          processors: []
          exporters: [logging]
EOF
```

!!! WARNING

    在 OpenTelemetry Collector 格式稳定之前，可能需要在上面的示例中进行更改，
    以保持与所引用的 OpenTelemetry Collector 映像的最新版本兼容。

这将创建一个名为`simplest`的 OpenTelemetry Collector 实例，暴露一
个`jaeger-grpc`端口，以从插装化的应用程序中消费 `span` ，并通过`logging`导出这些
`span`，这将`span`写入接收`span`的 OpenTelemetry Collector 实例的控制台
(`stdout`)。

`config`节点保存应该按原样传递给底层 OpenTelemetry Collector 实例的`YAML`。请参
阅[OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)文
档以获取可能条目的参考。

此时，Operator 并不验证配置文件的内容:如果配置无效，则仍然会创建实例，但是底层的
OpenTelemetry Collector 可能会崩溃。

操作员检查配置文件以发现已配置的接收器及其端口。如果它找到具有端口的接收器，它将
创建一对 kubernetes 服务，其中一个是无头的，在集群中公开这些端口。无头服务包含一
个`service.beta.openshift.io/serving-cert-secret-name`注释，它将导致 OpenShift
创建一个包含证书和密钥的秘密。这个秘密可以作为卷、证书和密钥挂载在这些接收者的
TLS 配置中。

### 更新

如上所述，OpenTelemetry Collector 格式正在继续发展。但是，尽最大努力尝试升级所有
托管的`OpenTelemetryCollector`资源。

在某些情况下，可能希望防止操作符升级某些`OpenTelemetryCollector`资源。例如，当一
个资源被配置为自定义的`.Spec.Image`时，最终用户可能希望自己管理配置，而不是让操
作员升级它。这可以通过公开的属性`.Spec.UpgradeStrategy`在资源的基础上配置。

通过将资源的`.Spec.UpgradeStrategy`配置为 none，操作符将在升级例程中跳过给定的实
例。

`.Spec.UpgradeStrategy`的默认值和唯一可接受的值是`automatic`。

### 部署模式

`OpenTelemetryCollector`的`CustomResource`暴露了一个名为`.Spec.Mode`的属性，该属
性可用于指定收集器是否应该作为`DaemonSet`, `Sidecar`, or `Deployment`(默认)运行
。请看这个[例子]作为参考。

[例子]:
  https://github.com/open-telemetry/opentelemetry-operator/blob/main/tests/e2e/daemonset-features/01-install.yaml

#### 附接盒注入

通过将 Pod 注释`sidecar.opentelemetry.io/inject`设置为`"true"`，或者设置为具体
的`OpenTelemetryCollector`的名称，可以将带有 OpenTelemetryCollector 的 sidecar
注入到基于 Pod 的工作负载中，如下所示:

```yaml
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: sidecar-for-my-app
spec:
  mode: sidecar
  config: |
    receivers:
      jaeger:
        protocols:
          thrift_compact:
    processors:

    exporters:
      logging:

    service:
      pipelines:
        traces:
          receivers: [jaeger]
          processors: []
          exporters: [logging]
EOF

kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: myapp
  annotations:
    sidecar.opentelemetry.io/inject: "true"
spec:
  containers:
  - name: myapp
    image: jaegertracing/vertx-create-span:operator-e2e-tests
    ports:
      - containerPort: 8080
        protocol: TCP
EOF
```

当在同一个命名空间中有多个模式设置为`Sidecar`的`OpenTelemetryCollector` 资源时，
应该使用一个具体的名称。当同一个命名空间中只有一个`Sidecar`实例时，当注释被设置
为`"true"`时使用该实例。

注释值可以来自名称空间，也可以来自 pod。最具体的注释胜出，顺序如下:

- pod 注释被设置为具体的实例名或`"false"`时使用。
- 当 pod 注释不存在或设置为`"true"`，而命名空间设置为具体实例或设置为`"false"`时
  ，使用命名空间注释。

注释的可能值可以是:

- "true" - inject `OpenTelemetryCollector` resource from the namespace.
- "sidecar-for-my-app" - name of `OpenTelemetryCollector` CR instance in the
  current namespace.
- "my-other-namespace/my-instrumentation" - name and namespace of
  `OpenTelemetryCollector` CR instance in another namespace.
- "false" - do not inject

当使用基于 pod 的工作负载时，例如`Deployment` 或 `Statefulset`，请确保将注释添加
到`PodTemplate`部分。如:

```yaml
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  labels:
    app: my-app
  annotations:
    sidecar.opentelemetry.io/inject: "true" # WRONG
spec:
  selector:
    matchLabels:
      app: my-app
  replicas: 1
  template:
    metadata:
      labels:
        app: my-app
      annotations:
        sidecar.opentelemetry.io/inject: "true" # CORRECT
    spec:
      containers:
      - name: myapp
        image: jaegertracing/vertx-create-span:operator-e2e-tests
        ports:
          - containerPort: 8080
            protocol: TCP
EOF
```

当使用 sidecar 模式时，OpenTelemetry 收集器容器将使用 Kubernetes 资源属性设置环
境变量`OTEL_RESOURCE_ATTRIBUTES`，准备
由[resourcedetection](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor)处
理器使用。

### OpenTelemetry 自动插装注入

操作员可以注入和配置 OpenTelemetry 自动仪器库。目前支持 Apache HTTPD, DotNet,
Go, Java, NodeJS 和 Python。

要使用自动检测，请配置一个带有 SDK 和检测配置的 `Instrumentation` 资源。

```yaml
kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: my-instrumentation
spec:
  exporter:
    endpoint: http://otel-collector:4317
  propagators:
    - tracecontext
    - baggage
    - b3
  sampler:
    type: parentbased_traceidratio
    argument: "0.25"
  python:
    env:
      # Required if endpoint is set to 4317.
      # Python autoinstrumentation uses http/proto by default
      # so data must be sent to 4318 instead of 4317.
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: http://otel-collector:4318
  dotnet:
    env:
      # Required if endpoint is set to 4317.
      # Dotnet autoinstrumentation uses http/proto by default
      # See https://github.com/open-telemetry/opentelemetry-dotnet-instrumentation/blob/888e2cd216c77d12e56b54ee91dafbc4e7452a52/docs/config.md#otlp
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: http://otel-collector:4318
EOF
```

`propagators`的值被添加到`OTEL_PROPAGATORS`环境变量中。 `传播器`的有效值
由[OpenTelemetry Specification for OTEL_PROPAGATORS](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/#otel_propagators)定
义。

`sampler.type`的值被添加到`OTEL_TRACES_SAMPLER`环境变量中。 `sampler.type`的有效
值
由[OTEL_TRACES_SAMPLER 的 OpenTelemetry 规范](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/#otel_traces_sampler)定
义。 `sampler.argument`的值被添加到`OTEL_TRACES_SAMPLER_ARG`环境变量中。
`sampler.argument`的有效值将取决于所选择的采样器。请参
阅[OTEL_TRACES_SAMPLER_ARG 的 OpenTelemetry 规范](https://opentelemetry.io/docs/concepts/sdk-configuration/general-sdk-configuration/#otel_traces_sampler_arg)了
解更多详细信息。

以上 CR 可以通过 `kubectl get otelinst` 查询。

然后向 pod 添加注释以启用注入。可以将注释添加到名称空间中，以便该名称空间中的所
有 pod 都将获得检测，或者将注释添加到单独的 PodSpec 对象中，这些对象可以作为
Deployment、Statefulset 和其他资源的一部分使用。

Java:

```bash
instrumentation.opentelemetry.io/inject-java: "true"
```

NodeJS:

```bash
instrumentation.opentelemetry.io/inject-nodejs: "true"
```

Python:

```bash
instrumentation.opentelemetry.io/inject-python: "true"
```

DotNet:

```bash
instrumentation.opentelemetry.io/inject-dotnet: "true"
```

Go:

Go auto-instrumentation 还支持用于设
置[OTEL_GO_AUTO_TARGET_EXE env var](https://github.com/open-telemetry/opentelemetry-go-instrumentation/blob/main/docs/how-it-works.md)的
注释。这个 env 变量也可以通过 Instrumentation 资源设置，注释优先。由于 Go 自动检
测需要设置`OTEL_GO_AUTO_TARGET_EXE`，因此您必须通过注释或 Instrumentation 资源提
供有效的可执行路径。设置此值失败将导致仪器注入中止，使原始 pod 保持不变。

```bash
instrumentation.opentelemetry.io/inject-go: "true"
instrumentation.opentelemetry.io/otel-go-auto-target-exe: "/path/to/container/executable"
```

Apache HTTPD:

```bash
instrumentation.opentelemetry.io/inject-apache-httpd: "true"
```

OpenTelemetry SDK 环境变量:

```bash
instrumentation.opentelemetry.io/inject-sdk: "true"
```

注释的可能值可以是

- `"true"` - 从命名空间注入和 `Instrumentation` 资源。
- `"my-instrumentation"` - 当前命名空间中`Instrumentation`CR 实例的名称。
- `"my-other-namespace/my-instrumentation"` - `Instrumentation` CR 实例在另一个
  名称空间中的名称和名称空间。
- `"false"` - 不要注射

#### 多容器 Pod

如果没有指定其他内容，则在 pod 规范中可用的第一个容器上执行检测。在某些情况下(例
如在注入 Istio sidecar 的情况下)，有必要指定必须在哪个容器上执行此注入。

为此，有可能对将进行注射的吊舱进行微调。

为此，我们将使用`instrumentation.opentelemetry.io/container-names`注释，我们将为
其指定一个或多个必须进行注入的 pod 名称(`.spec.containers.name`):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment-with-multiple-containers
spec:
  selector:
    matchLabels:
      app: my-pod-with-multiple-containers
  replicas: 1
  template:
    metadata:
      labels:
        app: my-pod-with-multiple-containers
      annotations:
        instrumentation.opentelemetry.io/inject-java: 'true'
        instrumentation.opentelemetry.io/container-names: 'myapp,myapp2'
    spec:
      containers:
        - name: myapp
          image: myImage1
        - name: myapp2
          image: myImage2
        - name: myapp3
          image: myImage3
```

在上述情况下，`myapp`和`myapp2`容器将被检测，`myapp3`不会。

!!! NOTE

    Go 的自动检测 **不** 支持多容器 pod。当注入 Go 自动检测时，第一个 pod 应该是你想要检测的唯一 pod。

#### 使用定制的或供应商的工具

默认情况下，操作符使用上游自动插装库。可以通过覆盖 CR 中的映像字段来配置自定义自
动检测。

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: my-instrumentation
spec:
  java:
    image: your-customized-auto-instrumentation-image:java
  nodejs:
    image: your-customized-auto-instrumentation-image:nodejs
  python:
    image: your-customized-auto-instrumentation-image:python
  dotnet:
    image: your-customized-auto-instrumentation-image:dotnet
  go:
    image: your-customized-auto-instrumentation-image:go
  apacheHttpd:
    image: your-customized-auto-instrumentation-image:apache-httpd
```

自动检测的 Dockerfiles 可以在[autoinstrumentation 目录](./autoinstrumentation)中
找到。按照 Dockerfiles 中的说明来构建自定义容器映像。

#### 使用 Apache HTTPD 自动检测

对于`Apache HTTPD` 自动检测，默认情况下，检测假设 HTTPD 版本 2.4 和 HTTPD 配置目
录`/usr/local/apache2/conf`为它在官方的 `Apache HTTPD` 镜像中(参见
docker.io/httpd:latest)。如果您需要使用 2.2 版本，或者您的 HTTPD 配置目录不同，
或者您需要调整代理属性，请根据以下示例自定义工具规范:

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: my-instrumentation
  apache:
    image: your-customized-auto-instrumentation-image:apache-httpd
    version: 2.2
    configPath: /your-custom-config-path
    attrs:
      - name: ApacheModuleOtelMaxQueueSize
        value: '4096'
      - name: ...
        value: ...
```

所有可用属性的列表可以
在[otel-webserver-module](https://github.com/open-telemetry/opentelemetry-cpp-contrib/tree/main/instrumentation/otel-webserver-module)找
到

#### 只注入 OpenTelemetry SDK 环境变量

你可以为目前不能自动检测的应用配置 OpenTelemetry SDK，通过使用`inject-sdk`代替(
例如)`inject-python`或`inject-java`。这将注入环境变量，
如`OTEL_RESOURCE_ATTRIBUTES`，
`OTEL_TRACES_SAMPLER`和`OTEL_EXPORTER_OTLP_ENDPOINT`，您可以
在`Instrumentation`中配置，但实际上不会提供 SDK。

```bash
instrumentation.opentelemetry.io/inject-sdk: "true"
```

#### 控制插装功能

操作符允许通过特征门指定 Instrumentation 资源可以检测的语言。这些特征门必须通过
`--feature-gates`标志传递给操作符。该标志允许以逗号分隔的特征门标识符列表。在门
的前面加上'-'来禁用对相应语言的支持。用'+'作为门的前缀或不加前缀将启用对相应语言
的支持。如果一种语言在默认情况下是启用的，它的门只需要在禁用门时提供。

| Language    | Gate                                        | Default Value |
| ----------- | ------------------------------------------- | ------------- |
| Java        | `operator.autoinstrumentation.java`         | enabled       |
| NodeJS      | `operator.autoinstrumentation.nodejs`       | enabled       |
| Python      | `operator.autoinstrumentation.python`       | enabled       |
| DotNet      | `operator.autoinstrumentation.dotnet`       | enabled       |
| ApacheHttpD | `operator.autoinstrumentation.apache-httpd` | enabled       |
| Go          | `operator.autoinstrumentation.go`           | disabled      |

始终支持未在表中指定的语言，并且不能禁用。

### 目标分配程序

OpenTelemetry Operator 带有一个可选组件，即目标分配器(Target Allocator, TA)。当
创建 OpenTelemetryCollector 自定义资源(CR)并将 TA 设置为启用时，操作员将创建一个
新的部署和服务，以作为该 CR 的一部分为每个 Collector pod 提供特定
的`http_sd_config`指令。它还将更改 CR 中的 Prometheus 接收器配置，以便它使用来自
TA 的[http_sd_config](https://prometheus.io/docs/prometheus/latest/http_sd/)。下
面的例子展示了如何开始使用 Target Allocator:

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: collector-with-ta
spec:
  mode: statefulset
  targetAllocator:
    enabled: true
  config: |
    receivers:
      prometheus:
        config:
          scrape_configs:
          - job_name: 'otel-collector'
            scrape_interval: 10s
            static_configs:
            - targets: [ '0.0.0.0:8888' ]
            metric_relabel_configs:
            - action: labeldrop
              regex: (id|name)
              replacement: $$1
            - action: labelmap
              regex: label_(.+)
              replacement: $$1 

    exporters:
      logging:

    service:
      pipelines:
        metrics:
          receivers: [prometheus]
          processors: []
          exporters: [logging]
```

在上面的例子中，替换键中`$$`的用法是基于 Prometheus 接收
器[README](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/prometheusreceiver/README.md)文
档中提供的信息，该文档说明:
`Note: 由于收集器配置支持环境变量替换，普罗米修斯配置中的$字符被解释为环境变量。如果想在prometheus配置中使用$字符，必须使用$$转义。`

在幕后，OpenTelemetry 操作符将在对账后将 Collector 的配置转换为以下内容:

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: otel-collector
          scrape_interval: 10s
          http_sd_configs:
            - url: http://collector-with-ta-targetallocator:80/jobs/otel-collector/targets?collector_id=$POD_NAME
          metric_relabel_configs:
            - action: labeldrop
              regex: (id|name)
              replacement: $$1
            - action: labelmap
              regex: label_(.+)
              replacement: $$1

exporters:
  logging:

service:
  pipelines:
    metrics:
      receivers: [prometheus]
      processors: []
      exporters: [logging]
```

注意 Operator 如何从' `scrape_configs` '部分删除任何现有的服务发现配置(例如，'
`static_configs` '， ' `file_sd_configs` '等)，并添加一个' `http_sd_configs` '配
置，指向它所提供的 Target Allocator 实例。

OpenTelemetry 操作符还将在对账后将目标分配器的 promethueus 配置转换为以下内容:

```yaml
config:
  scrape_configs:
    - job_name: otel-collector
      scrape_interval: 10s
      static_configs:
        - targets: ['0.0.0.0:8888']
      metric_relabel_configs:
        - action: labeldrop
          regex: (id|name)
          replacement: $1
        - action: labelmap
          regex: label_(.+)
          replacement: $1
```

注意，在本例中，Operator 将替换键中的`$$`替换为单个`$`。这是因为收集器支持环境变
量替换，而 TA(目标分配器)不支持。因此，为了确保兼容性，TA 配置应该只包含一个`$`
符号。

更多关于 TargetAllocator 的信息可以在[这里](cmd/otel-allocator/README.md)找到.

#### 目标分配器配置重写

Prometheus 接收器现在显式支持从目标分配器获取抓取目标。因此，现在可以让 Operator
自动添加必要的目标分配器配置。此功能目前需要启
用`operator.collector.rewritetargetallocator`功能标志。启用该标志后，上一节中的
配置将呈现为:

```yaml
receivers:
  prometheus:
    config:
      global:
        scrape_interval: 1m
        scrape_timeout: 10s
        evaluation_interval: 1m
    target_allocator:
      endpoint: http://collector-with-ta-targetallocator:80
      interval: 30s
      collector_id: $POD_NAME

exporters:
  logging:

service:
  pipelines:
    metrics:
      receivers: [prometheus]
      processors: []
      exporters: [logging]
```

这还允许使用普罗米修斯操作符 crd 进行更直接的目标发现收集器配置。下面是一个最小
的例子:

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: collector-with-ta-prometheus-cr
spec:
  mode: statefulset
  targetAllocator:
    enabled: true
    serviceAccount: everything-prometheus-operator-needs
    prometheusCR:
      enabled: true
  config: |
    receivers:
      prometheus:

    exporters:
      logging:

    service:
      pipelines:
        metrics:
          receivers: [prometheus]
          processors: []
          exporters: [logging]
```

## 兼容性矩阵

### OpenTelemetry 操作符与 OpenTelemetry 收集器

OpenTelemetry Operator 遵循与操作数(OpenTelemetry Collector)相同的版本控制，直到
版本的次要部分。例如，OpenTelemetry Operator v0.18.1 跟踪 OpenTelemetry
Collector 0.18.0。版本的补丁部分表示操作符本身的补丁级别，而不是 OpenTelemetry
Collector 的补丁级别。每当 OpenTelemetry Collector 的新补丁版本发布时，我们将发
布操作符的新补丁版本。

By default, the OpenTelemetry Operator ensures consistent versioning between
itself and the managed `OpenTelemetryCollector` resources. That is, if the
OpenTelemetry Operator is based on version `0.40.0`, it will create resources
with an underlying OpenTelemetry Collector at version `0.40.0`.

When a custom `Spec.Image` is used with an `OpenTelemetryCollector` resource,
the OpenTelemetry Operator will not manage this versioning and upgrading. In
this scenario, it is best practice that the OpenTelemetry Operator version
should match the underlying core version. Given a `OpenTelemetryCollector`
resource with a `Spec.Image` configured to a custom image based on underlying
OpenTelemetry Collector at version `0.40.0`, it is recommended that the
OpenTelemetry Operator is kept at version `0.40.0`.

### OpenTelemetry Operator vs. Kubernetes vs. Cert Manager

我们努力与尽可能广泛的 Kubernetes 版本兼容，但是 Kubernetes 本身的一些更改需要我
们打破与旧 Kubernetes 版本的兼容性，可能是因为代码不兼容，或者以可维护性的名义。
每个已发布的操作符都将支持特定范围的 Kubernetes 版本，最晚在发布期间确定。

我们使用`cert-manager`来实现这个操作符的一些特性，第三列显示了已知与这个操作符的
版本一起工作的`cert-manager`的版本。

OpenTelemetry 操作符 _可能_ 在给定范围之外的版本上工作，但当打开新问题时，请确保
在受支持的版本上测试您的场景。

| OpenTelemetry Operator | Kubernetes     | Cert-Manager |
| ---------------------- | -------------- | ------------ |
| v0.79.0                | v1.19 to v1.27 | v1           |
| v0.78.0                | v1.19 to v1.27 | v1           |
| v0.77.0                | v1.19 to v1.26 | v1           |
| v0.76.1                | v1.19 to v1.26 | v1           |
| v0.75.0                | v1.19 to v1.26 | v1           |
| v0.74.0                | v1.19 to v1.26 | v1           |
| v0.73.0                | v1.19 to v1.26 | v1           |
| v0.72.0                | v1.19 to v1.26 | v1           |
| v0.71.0                | v1.19 to v1.25 | v1           |
| v0.70.0                | v1.19 to v1.25 | v1           |
| v0.69.0                | v1.19 to v1.25 | v1           |
| v0.68.0                | v1.19 to v1.25 | v1           |
| v0.67.0                | v1.19 to v1.25 | v1           |
| v0.66.0                | v1.19 to v1.25 | v1           |
| v0.64.1                | v1.19 to v1.25 | v1           |
| v0.63.1                | v1.19 to v1.25 | v1           |
| v0.62.1                | v1.19 to v1.25 | v1           |
| v0.61.0                | v1.19 to v1.25 | v1           |
| v0.60.0                | v1.19 to v1.25 | v1           |
| v0.59.0                | v1.19 to v1.24 | v1           |
| v0.58.0                | v1.19 to v1.24 | v1           |
| v0.57.2                | v1.19 to v1.24 | v1           |
| v0.56.0                | v1.19 to v1.24 | v1           |

## 贡献与开发

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

In addition to the
[core responsibilities](https://github.com/open-telemetry/community/blob/main/community-membership.md)
the operator project requires approvers and maintainers to be responsible for
releasing the project. See [RELEASE.md](./RELEASE.md) for more information and
release schedule.

Approvers
([@open-telemetry/operator-approvers](https://github.com/orgs/open-telemetry/teams/operator-approvers)):

- [Benedikt Bongartz](https://github.com/frzifus), Red Hat
- [Tyler Helmuth](https://github.com/TylerHelmuth), Honeycomb
- [Yuri Oliveira Sa](https://github.com/yuriolisa), Red Hat

Emeritus Approvers:

- [Anthony Mirabella](https://github.com/Aneurysm9), AWS
- [Dmitrii Anoshin](https://github.com/dmitryax), Splunk
- [Jay Camp](https://github.com/jrcamp), Splunk
- [James Bebbington](https://github.com/james-bebbington), Google
- [Owais Lone](https://github.com/owais), Splunk
- [Pablo Baeyens](https://github.com/mx-psi), DataDog

Target Allocator Maintainers
([@open-telemetry/operator-ta-maintainers](https://github.com/orgs/open-telemetry/teams/operator-ta-maintainers)):

- [Anthony Mirabella](https://github.com/Aneurysm9), AWS
- [Kristina Pathak](https://github.com/kristinapathak), Lightstep
- [Sebastian Poxhofer](https://github.com/secustor)

Maintainers
([@open-telemetry/operator-maintainers](https://github.com/orgs/open-telemetry/teams/operator-maintainers)):

- [Jacob Aronoff](https://github.com/jaronoff97), Lightstep
- [Pavol Loffay](https://github.com/pavolloffay), Red Hat
- [Vineeth Pothulapati](https://github.com/VineethReddy02), Timescale

Emeritus Maintainers

- [Alex Boten](https://github.com/codeboten), Lightstep
- [Bogdan Drutu](https://github.com/BogdanDrutu), Splunk
- [Juraci Paixão Kröhling](https://github.com/jpkrohling), Grafana Labs
- [Tigran Najaryan](https://github.com/tigrannajaryan), Splunk

Learn more about roles in the
[community repository](https://github.com/open-telemetry/community/blob/main/community-membership.md).

Thanks to all the people who already contributed!

[![Contributors][contributors-img]][contributors]

## License

[Apache 2.0 License](./LICENSE).

[github-workflow]:
  https://github.com/open-telemetry/opentelemetry-operator/actions
[github-workflow-img]:
  https://github.com/open-telemetry/opentelemetry-operator/workflows/Continuous%20Integration/badge.svg
[goreport-img]:
  https://goreportcard.com/badge/github.com/open-telemetry/opentelemetry-operator
[goreport]:
  https://goreportcard.com/report/github.com/open-telemetry/opentelemetry-operator
[godoc-img]:
  https://godoc.org/github.com/open-telemetry/opentelemetry-operator?status.svg
[godoc]:
  https://godoc.org/github.com/open-telemetry/opentelemetry-operator/pkg/apis/opentelemetry/v1alpha1#OpenTelemetryCollector
[contributors]:
  https://github.com/open-telemetry/opentelemetry-operator/graphs/contributors
[contributors-img]:
  https://contributors-img.web.app/image?repo=open-telemetry/opentelemetry-operator

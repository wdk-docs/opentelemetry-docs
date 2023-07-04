---
title: 扩容采集器
weight: 26
---

在使用 Opentelemetry Collector 规划可观察性管道时，应该考虑随着遥测收集的增加而
扩展管道的方法。

以下部分将指导您完成计划阶段，讨论要扩展哪些组件、如何确定何时扩展以及如何执行计
划。

## 如何扩展

虽然 OpenTelemetry Collector 在一个二进制文件中处理所有遥测信号类型，但实际情况
是每种类型可能有不同的缩放需求，并且可能需要不同的缩放策略。首先查看您的工作负载
，以确定哪种信号类型预计会占用最大的负载份额，以及哪种格式预计会被 Collector 接
收。例如，扩展抓取集群与扩展日志接收器有很大的不同。还要考虑工作负载的弹性:您是
在一天中的特定时间达到高峰，还是 24 小时内的负载都是相似的? 一旦您收集了这些信息
，您将了解需要扩展的内容。

例如，假设您要抓取数百个 Prometheus 端点，每分钟有来自 fluentd 实例的 tb 级日志
，以及一些应用程序指标和跟踪以 OTLP 格式从最新的微服务到达。在这种情况下，您将需
要一个可以单独扩展每个信号的体系结构:扩展 Prometheus 接收器需要 scraper 之间的协
调，以决定哪个 scraper 到哪个端点。相反，我们可以根据需要水平扩展无状态日志接收
器。在第三个收集器集群中设置度量和跟踪的 OTLP 接收器将允许我们隔离故障并更快地迭
代，而不必担心重新启动繁忙的管道。假设 OTLP 接收器允许摄取所有遥测类型，我们可以
将应用程序指标和跟踪保持在同一个实例上，并在需要时水平扩展它们。

## 何时扩展

同样，我们应该了解我们的工作负载，以决定何时扩大或缩小规模，但是收集器发出的一些
指标可以给您很好的提示，告诉您何时采取行动。

当 memory_limititer 处理器是管道的一部分时，Collector 可以给您的一个有用的提示是
度量' otelcol_processor_refused_span '。此处理器允许您限制收集器可以使用的内存量
。虽然 Collector 消耗的数据可能比该处理器中配置的最大数据量多一点，但是
memory_limiter 最终将阻止新数据通过管道，它将在此度量中记录这一事实。对于所有其
他遥测数据类型都存在相同的度量。如果数据经常被拒绝进入管道，您可能需要扩展
Collector 集群。一旦跨节点的内存消耗明显低于该处理器中设置的限制，您就可以缩小规
模。

另一组需要注意的指标是与出口商队列大小相关的指标
:`otelcol_exporter_queue_capacity` 和 `otelcol_exporter_queue_size`。收集器将在
等待工作线程可用以发送数据时将数据排在内存中。如果没有足够的工人或后端太慢，数据
就会开始堆积在队列中。一旦队列达到其容量(`otelcol_exporter_queue_size` >
`otelcol_exporter_queue_capacity`)，它拒绝数据
(`otelcol_exporter_enqueue_failed_spans`)。添加更多的 worker 通常会使 Collector
导出更多的数据，这可能不一定是您想要的(参
见[什么时候不按比例](#when-not-to-scale))。

熟悉您打算使用的组件也是值得的，因为不同的组件可能会产生其他指标。例如
，[负载平衡导出器将记录有关导出操作的定时信息](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/loadbalancingexporter#metrics)，
将其作为直方图`otelcol_loadbalancer_backend_latency`的一部分公开。您可以提取此信
息，以确定所有后端处理请求的时间是否相同:单个后端缓慢可能表明问题出在收集器外部
。

对于执行刮擦的接收器，例如 Prometheus 接收器，一旦完成刮擦所有目标所需的时间通常
非常接近刮擦间隔，则应该对刮擦进行缩放或分片。当这种情况发生时，是时候添加更多的
抓取器了，通常是新的 Collector 实例。

### 什么时候不能扩展

也许与知道何时进行扩展同样重要的是，了解哪些迹象表明扩展操作不会带来任何好处。一
个例子是遥测数据库无法跟上负载:在不扩展数据库的情况下，将 collector 添加到集群中
将没有帮助。类似地，当 Collector 和后端之间的网络连接饱和时，添加更多的
Collector 可能会导致有害的副作用。

同样，捕捉这种情况的一种方法是查看指标`otelcol_exporter_queue_size` 和
`otelcol_exporter_queue_capacity`。如果队列大小一直接近队列容量，则表明导出数据
比接收数据慢。您可以尝试增加队列大小，这将导致 Collector 消耗更多内存，但它也将
为后端提供一些喘息的空间，而不会永久丢失遥测数据。但是，如果您不断增加队列容量，
并且队列大小以相同的比例增长，则表明您可能需要查看收集器之外的情况。同样重要的是
要注意，在这里添加更多的工作人员是没有帮助的:您只会给已经承受高负载的系统施加更
大的压力。

后端可能出现问题的另一个迹象是`otelcol_exporter_send_failed_spans`指标的增加:这
表明向后端发送数据永久失败。当这种情况持续发生时，扩大收集器可能只会使情况变得更
糟。

## 如何扩展

此时，我们知道管道的哪些部分需要缩放。关于伸缩方面，我们有三种类型的组件:无状态
、刮削和有状态。

大多数 Collector 组件都是无状态的。即使它们在内存中保留了一些状态，也与扩展目的
无关。

像普罗米修斯接收器一样，scraper 被配置为从外部位置获取遥测数据。然后，接收器将逐
个抓取目标，将数据放入管道中。

像尾部采样处理器这样的组件不容易扩展，因为它们在内存中保留了一些与业务相关的状态
。在扩大规模之前，需要仔细考虑这些组成部分。

### 伸缩无状态收集器

好消息是，大多数情况下，扩展 Collector 很容易，因为只需添加新的副本并使用现成的
负载平衡器即可。当使用 gRPC 接收数据时，我们建议使用理解 gRPC 的负载平衡器。否则
，客户端总是会碰到同一个后备收集器。

您仍然应该考虑拆分收集管道并考虑可靠性。例如，当您的工作负载在 Kubernetes 上运行
时，您可能希望使用 DaemonSets 在与您的工作负载相同的物理节点上拥有一个
Collector，并在将数据发送到存储之前负责对数据进行预处理。当节点数量较少而 pod 数
量较多时，Sidecars 可能更有意义，因为您将在 Collector 层之间获得更好的 gRPC 连接
负载平衡，而不需要特定于 gRPC 的负载平衡器。使用 Sidecar 还可以避免在一个
DaemonSet pod 出现故障时导致节点中所有 pod 的关键组件瘫痪。

侧车模式包括将容器添加到工作负载 pod 中。
[OpenTelemetry Operator](/docs/k8s-operator/)可以自动为您添加。要做到这一点，你
需要一个 OpenTelemetry Collector CR，你需要注释你的 PodSpec 或 Pod，告诉操作员注
入一个 sidecar:

```yaml
---
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: sidecar-for-my-workload
spec:
  mode: sidecar
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
    processors:

    exporters:
      logging:

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: []
          exporters: [logging]
---
apiVersion: v1
kind: Pod
metadata:
  name: my-microservice
  annotations:
    sidecar.opentelemetry.io/inject: 'true'
spec:
  containers:
    - name: my-microservice
      image: my-org/my-microservice:v0.0.0
      ports:
        - containerPort: 8080
          protocol: TCP
```

In case you prefer to bypass the operator and add a sidecar manually, here’s an
example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-microservice
spec:
  containers:
    - name: my-microservice
      image: my-org/my-microservice:v0.0.0
      ports:
        - containerPort: 8080
          protocol: TCP
    - name: sidecar
      image: ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector:0.69.0
      ports:
        - containerPort: 8888
          name: metrics
          protocol: TCP
        - containerPort: 4317
          name: otlp-grpc
          protocol: TCP
      args:
        - --config=/conf/collector.yaml
      volumeMounts:
        - mountPath: /conf
          name: sidecar-conf
  volumes:
    - name: sidecar-conf
      configMap:
        name: sidecar-for-my-workload
        items:
          - key: collector.yaml
            path: collector.yaml
```

### 缩放刮刀

一些接收器正在积极获取遥测数据，并将其放置在管道中，如 hostmetrics 和 prometheus
接收器。虽然获取主机指标并不是我们通常会扩展的内容，但我们可能需要将为
Prometheus 接收器抓取数千个端点的工作分开。而且我们不能简单地使用相同的配置添加
更多的实例，因为每个 Collector 都会尝试在集群中抓取与其他 Collector 相同的端点，
从而导致更多的问题，比如乱序采样。

解决方案是按 Collector 实例对端点进行分片，这样，如果我们添加另一个 Collector 的
副本，每个副本将作用于不同的端点集。

这样做的一种方法是为每个 Collector 提供一个配置文件，以便每个 Collector 只发现该
Collector 的相关端点。例如，每个 Collector 可以负责一个 Kubernetes 名称空间或工
作负载上的特定标签。

扩展 Prometheus 接收器的另一种方法是使
用[Target Allocator](https://github.com/open-telemetry/opentelemetry-operator#target-allocator):)
它是一个额外的二进制文件，可以作为 OpenTelemetry Operator 的一部分进行部署，并使
用一致的散列算法将给定配置的 Prometheus 作业分配到收集器集群中。你可以像下面这样
使用自定义资源(CR)来使用目标分配器:

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

    exporters:
      logging:

    service:
      pipelines:
        traces:
          receivers: [prometheus]
          processors: []
          exporters: [logging]
```

在对账之后，OpenTelemetry 操作符将把 Collector 的配置转换为以下格式:

```yaml
   exporters:
      logging: null
    receivers:
      prometheus:
        config:
          global:
            scrape_interval: 1m
            scrape_timeout: 10s
            evaluation_interval: 1m
          scrape_configs:
          - job_name: otel-collector
            honor_timestamps: true
            scrape_interval: 10s
            scrape_timeout: 10s
            metrics_path: /metrics
            scheme: http
            follow_redirects: true
            http_sd_configs:
            - follow_redirects: false
              url: http://collector-with-ta-targetallocator:80/jobs/otel-collector/targets?collector_id=$POD_NAME
    service:
      pipelines:
        traces:
          exporters:
          - logging
          processors: []
          receivers:
          - prometheus
```

注意，Operator 是如何在“otel-collector”刮擦配置中添加“global”部分和“new
http_sd_configs”的，指向它所提供的 Target Allocator 实例。现在，要扩展收集器，请
更改 CR 的“replicas”属性，Target Allocator 将通过为每个收集器实例(pod)提供自定义
的“http_sd_config”来相应地分配负载。

### 伸缩状态收集器

某些组件可能在内存中保存数据，当按比例扩展时产生不同的结果。尾部采样处理器就是这
种情况，它在给定的时间段内保持内存中的跨度，仅在认为跟踪完成时才评估采样决策。通
过添加更多副本来扩展 Collector 集群意味着不同的收集器将接收给定跟踪的跨度，从而
导致每个收集器评估是否应该对该跟踪进行采样，可能会得到不同的答案。这种行为导致跟
踪范围丢失，错误地表示事务中发生的事情。

当使用跨指标处理器生成服务指标时，也会出现类似的情况。当不同的收集器接收到与同一
服务相关的数据时，基于服务名称的聚合将不准确。

为了克服这个问题，您可以在 collector 前面部署一个包含负载平衡导出器的 collector
层，以执行尾采样或跨度到度量的处理。负载平衡导出程序将一致地散列跟踪 ID 或服务名
称，并确定哪个收集器后端应该接收该跟踪的范围。您可以配置负载平衡导出器，以使用给
定 DNS a 条目后面的主机列表，例如 Kubernetes 无头服务。当支持该服务的部署向上或
向下扩展时，负载平衡导出器最终将看到更新的主机列表。或者，您可以指定负载平衡导出
程序要使用的静态主机列表。您可以通过增加副本的数量来扩展配置了负载平衡导出器的
collector 层。请注意，每个 Collector 可能会在不同的时间运行 DNS 查询，导致集群视
图出现片刻的差异。我们建议降低间隔值，以便在高弹性环境中集群视图只在短时间内不同
。

下面是一个使用 DNS a 记录(Kubernetes service otelcol 在 observability 命名空间上
)作为后端信息输入的示例配置:

```yaml
receivers:
  otlp:
    protocols:
      grpc:

processors:

exporters:
  loadbalancing:
    protocol:
      otlp:
    resolver:
      dns:
        hostname: otelcol.observability.svc.cluster.local

service:
  pipelines:
    traces:
      receivers:
        - otlp
      processors: []
      exporters:
        - loadbalancing
```

---
title: 配置
weight: 20
spelling: cSpell:ignore pprof zpages zipkin fluentforward hostmetrics opencensus
spelling: cSpell:ignore prometheus loglevel otlphttp upsert spanevents OIDC
spelling: cSpell:ignore prometheusremotewrite prodevent spanmetrics servicegraph
spelling: cSpell:ignore oidc cfssl genkey initca cfssljson gencert
---

假设熟悉以下页面:

- [数据收集](../../concepts/data-collection.md)以便了解适用于收集器的存储库。
- [安全指导](https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/security-best-practices.md)

## 基本

采集器由四个组件组成，用于访问遥测数据:

- ![](../../assets/img/logos/32x32/Receivers.svg){ width="32" }
  [Receivers](#receivers)
- ![](../../assets/img/logos/32x32/Processors.svg){ width="32" }
  [Processors](#processors)
- ![](../../assets/img/logos/32x32/Exporters.svg){ width="32" }
  [Exporters](#exporters)
- ![](../../assets/img/logos/32x32/Load_Balancer.svg){ width="32" }
  [Connectors](#connectors)

这些组件一旦配置好，就必须通过[service](#service)部分中的管道启用。

其次，还有[扩展](#extensions)，它们提供了可以添加到 Collector 的功能，但不需要直
接访问遥测数据，也不是管道的一部分。它们也在[service](#service)部分中启用。

一个示例配置如下:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:

exporters:
  otlp:
    endpoint: otelcol:4317

extensions:
  health_check:
  pprof:
  zpages:

service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

请注意，接收器、处理器、输出器和/或管道是通过组件标识符以`type[/name]`格式定义的
(例如: `otlp` 或 `otlp/2`)。只要标识符是唯一的，给定类型的组件就可以定义多次。例
如:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:
  otlp/2:
    protocols:
      grpc:
        endpoint: 0.0.0.0:55690

processors:
  batch:
  batch/test:

exporters:
  otlp:
    endpoint: otelcol:4317
  otlp/2:
    endpoint: otelcol2:4317

extensions:
  health_check:
  pprof:
  zpages:

service:
  extensions: [health_check, pprof, zpages]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    traces/2:
      receivers: [otlp/2]
      processors: [batch/test]
      exporters: [otlp/2]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

配置也可以引入其他文件，收集器可以将两个文件合并到 YAML 配置的单个文件中:

```yaml
receivers:
  otlp:
    protocols:
      grpc:

exporters: ${file:exporters.yaml}

service:
  extensions: []
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlp]
```

`exporters.yaml` 文件为:

```yaml
otlp:
  endpoint: otelcol.observability.svc.cluster.local:443
```

内存中的最终结果将是:

```yaml
receivers:
  otlp:
    protocols:
      grpc:

exporters:
  otlp:
    endpoint: otelcol.observability.svc.cluster.local:443

service:
  extensions: []
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlp]
```

<a name="receivers"></a>

## Receivers - 接收器

![](../../assets/img/logos/32x32/Receivers.svg){ width="35" }

接收器可以是基于推或拉的，它是数据进入收集器的方式。接收器可以支持一个或多
个[数据源](/docs/concepts/signals/)。

`receivers:`部分介绍如何配置接收器的。许多接收器都带有默认设置，因此只需指定接收
器的名称就足以配置它(例如，`zipkin:`)。如果需要配置或者用户希望更改默认配置，则
必须在本节中定义这种配置。将覆盖接收器为其提供默认配置的指定配置参数。

> 配置接收器不会启用它。接收器通过[service](#service)节中的管道启用。

必须配置一个或多个接收器。缺省情况下，没有配置接收器。下面提供了一个接收器的基本
示例。

> 有关接收器的详细配置，请参阅[receiver] README.md.

[receiver]: ./receiver.md

```yaml
receivers:
  # Data sources: logs
  fluentforward:
    endpoint: 0.0.0.0:8006

  # Data sources: metrics
  hostmetrics:
    scrapers:
      cpu:
      disk:
      filesystem:
      load:
      memory:
      network:
      process:
      processes:
      paging:

  # Data sources: traces
  jaeger:
    protocols:
      grpc:
      thrift_binary:
      thrift_compact:
      thrift_http:

  # Data sources: traces
  kafka:
    protocol_version: 2.0.0

  # Data sources: traces, metrics
  opencensus:

  # Data sources: traces, metrics, logs
  otlp:
    protocols:
      grpc:
      http:

  # Data sources: metrics
  prometheus:
    config:
      scrape_configs:
        - job_name: otel-collector
          scrape_interval: 5s
          static_configs:
            - targets: [localhost:8888]

  # Data sources: traces
  zipkin:
```

<a name="processors"></a>

## Processors - 处理器

![](./img/logos/32x32/Processors.svg){width="35"}

处理器在接收和导出之间的数据上运行。处理器是可选
的[有些是推荐的](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#recommended-processors).

`processors:`部分是如何配置处理器的。处理器可能带有默认设置，但许多处理器需要配
置。处理器的任何配置都必须在本节中完成。处理器为其提供默认配置的配置参数将被覆盖
。

> 配置处理器不能启用它。处理器是通过[service](#service)节中的管道启用的。

下面提供了默认处理器的一个基本示例。完整的处理器列表可以通过组
合[这里](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor)和[这里](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor)找
到。

> 有关详细的处理器配置，请参
> 阅[processor README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/README.md).

```yaml
processors:
  # Data sources: traces
  attributes:
    actions:
      - key: environment
        value: production
        action: insert
      - key: db.statement
        action: delete
      - key: email
        action: hash

  # Data sources: traces, metrics, logs
  batch:

  # Data sources: metrics
  filter:
    metrics:
      include:
        match_type: regexp
        metric_names:
          - prefix/.*
          - prefix_.*

  # Data sources: traces, metrics, logs
  memory_limiter:
    check_interval: 5s
    limit_mib: 4000
    spike_limit_mib: 500

  # Data sources: traces
  resource:
    attributes:
      - key: cloud.zone
        value: zone-1
        action: upsert
      - key: k8s.cluster.name
        from_attribute: k8s-cluster
        action: insert
      - key: redundant-attribute
        action: delete

  # Data sources: traces
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 15

  # Data sources: traces
  span:
    name:
      to_attributes:
        rules:
          - ^\/api\/v1\/document\/(?P<documentId>.*)\/update$
      from_attributes: [db.svc, operation]
      separator: '::'
```

<a name="exporters"></a>

## Exporters - 导出器

![](./img/logos/32x32/Exporters.svg){width="35"}

导出器可以是基于推或拉的，它是将数据发送到一个或多个后端/目的地的方式。导出器可
以支持一个或多个[数据源](../concepts/signals/index.md)。

`exporters:`部分是如何配置导出器的。导出程序可能带有默认设置，但许多导出程序需要
配置以至少指定目标和安全设置。导出器的任何配置都必须在本节中完成。将覆盖导出器为
其提供默认配置的指定配置参数。

> 配置导出器不会启用它。导出器是通过[service](#service)部分中的管道启用的。

必须配置一个或多个导出器。缺省情况下，没有配置导出器。下面提供了一个导出器的基本
示例。为了保证安全，某些导出器配置需要创建 x.509 证书，
如[设置证书](#setting-up-certificates)中所述。

> 有关导出器的详细配置，请参
> 阅[导出器 README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/README.md).

```yaml
exporters:
  # Data sources: traces, metrics, logs
  file:
    path: ./filename.json

  # Data sources: traces
  jaeger:
    endpoint: jaeger-all-in-one:14250
    tls:
      cert_file: cert.pem
      key_file: cert-key.pem

  # Data sources: traces
  kafka:
    protocol_version: 2.0.0

  # Data sources: traces, metrics, logs
  logging:
    loglevel: debug

  # Data sources: traces, metrics
  opencensus:
    endpoint: otelcol2:55678

  # Data sources: traces, metrics, logs
  otlp:
    endpoint: otelcol2:4317
    tls:
      cert_file: cert.pem
      key_file: cert-key.pem

  # Data sources: traces, metrics
  otlphttp:
    endpoint: https://example.com:4318

  # Data sources: metrics
  prometheus:
    endpoint: prometheus:8889
    namespace: default

  # Data sources: metrics
  prometheusremotewrite:
    endpoint: http://some.url:9411/api/prom/push
    # For official Prometheus (e.g. running via Docker)
    # endpoint: 'http://prometheus:9090/api/v1/write'
    # tls:
    #   insecure: true

  # Data sources: traces
  zipkin:
    endpoint: http://localhost:9411/api/v2/spans
```

<a name="connectors"></a>

## Connectors - 连接器

连接器既是输出器又是接收器。顾名思义，连接器连接两个管道:它在一个管道的末端作为
输出者使用数据，在另一个管道的开始作为接收者发出数据。它可以使用和发出相同数据类
型或不同数据类型的数据。连接器可以生成和发出数据来总结所使用的数据，也可以简单地
复制或路由数据。

`connectors:`部分是如何配置连接器的。

> 配置连接器不会启用它。连接器是通过[service](#service)节中的管道启用的。

可以配置一个或多个连接器。默认情况下，没有配置连接器。下面提供了一个连接器的基本
示例。

> 有关连接器的详细配置，请参
> 阅[connector README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector/README.md).

```yaml
connectors:
  forward:

  count:
    spanevents:
      my.prod.event.count:
        description: The number of span events from my prod environment.
        conditions:
          - 'attributes["env"] == "prod"'
          - 'name == "prodevent"'

  spanmetrics:
    histogram:
      explicit:
        buckets: [100us, 1ms, 2ms, 6ms, 10ms, 100ms, 250ms]
    dimensions:
      - name: http.method
        default: GET
      - name: http.status_code
    dimensions_cache_size: 1000
    aggregation_temporality: 'AGGREGATION_TEMPORALITY_CUMULATIVE'

  servicegraph:
    latency_histogram_buckets: [1, 2, 3, 4, 5]
    dimensions:
      - dimension-1
      - dimension-2
    store:
      ttl: 1s
      max_items: 10
```

<a name="extensions"></a>

## Extensions - 扩展

扩展主要用于不涉及处理遥测数据的任务。扩展的示例包括运行状况监视、服务发现和数据
转发。扩展是可选的。

`extensions:`部分是如何配置扩展的。许多扩展都带有默认设置，因此只需指定扩展的名
称就足以配置它(例如，`health_check:`)。如果需要配置或者用户希望更改默认配置，则
必须在本节中定义这种配置。扩展为其提供默认配置的指定配置参数将被覆盖。

> 配置扩展并不启用它。扩展在[service](#service)部分中启用。

默认情况下，没有配置任何扩展。下面提供了一个基本的扩展示例。

> 有关详细的扩展配置，请参
> 阅[extension README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension/README.md).

```yaml
extensions:
  health_check:
  pprof:
  zpages:
  memory_ballast:
    size_mib: 512
```

<a name="service"></a>

## Service - 服务

service 部分用于根据接收器、处理器、导出器和扩展部分中的配置配置在 Collector 中
启用哪些组件。如果一个组件被配置了，但没有在服务部分中定义，那么它是不启用的。服
务部分由三个部分组成:

- extensions
- pipelines
- telemetry

扩展包含要启用的所有扩展的列表。例如:

```yaml
service:
  extensions: [health_check, pprof, zpages]
```

管道可以是以下类型:

- traces: 收集和处理跟踪数据。
- metrics: 收集和处理度量数据。
- logs: 收集和处理日志数据。

管道由一组接收器、处理器和输出器组成。每个接收器/处理器/输出器必须在服务部分之外
的配置中定义，以便包含在管道中。

!!! note

    每个接收器/处理器/输出器可以在多个管道中使用。
    对于在多个管道中引用的处理器，每个管道将获得该处理器的一个单独实例。
    这与在多个管道中引用的接收器/导出器形成对比，在多个管道中，只有一个接收器/导出器实例用于所有管道。
    还要注意，处理器的顺序决定了处理数据的顺序。

管道配置示例如下:

```yaml
service:
  pipelines:
    metrics:
      receivers: [opencensus, prometheus]
      exporters: [opencensus, prometheus]
    traces:
      receivers: [opencensus, jaeger]
      processors: [batch]
      exporters: [opencensus, zipkin]
```

遥测是可以配置收集器本身的遥测的地方。它有两个子部分:“日志”和“指标”。

`logs`小节允许对收集器生成的日志进行配置。默认情况下，收集器将其日志以`INFO`的日
志级别写入 stderr。您还可以使用`initial_fields` 部分向所有日志添加静态键值对。
[在这里查看`logs`选项的完整列表](https://github.com/open-telemetry/opentelemetry-collector/blob/7666eb04c30e5cfd750db9969fe507562598f0ae/config/service.go#L41-L97)

`metrics`小节允许配置收集器生成的指标。默认情况下，收集器将生成关于自身的基本指
标，并将其暴露
在`localhost:8888/metrics`.[查看完整的“指标”选项列表](https://github.com/open-telemetry/opentelemetry-collector/blob/7666eb04c30e5cfd750db9969fe507562598f0ae/config/service.go#L99-L111)

遥测配置举例如下:

```yaml
service:
  telemetry:
    logs:
      level: debug
      initial_fields:
        service: my-instance
    metrics:
      level: detailed
      address: 0.0.0.0:8888
```

## 其他信息

### 配置环境变量

在 Collector 配置中支持环境变量的使用和扩展。例如，要使用存储
在`DB_KEY`和`OPERATION`环境变量上的值，你可以这样写:

```yaml
processors:
  attributes/example:
    actions:
      - key: ${env:DB_KEY}
        action: ${env:OPERATION}
```

使用`$$`表示文字`$`。例如，表示`$DataVisualization`看起来像这样:

```yaml
exporters:
  prometheus:
    endpoint: prometheus:8889
    namespace: $$DataVisualization
```

### 代理支持

利用`net/http`包的出口商(今天都是这样)尊重以下代理环境变量:

- HTTP_PROXY
- HTTPS_PROXY
- NO_PROXY

如果在 Collector 启动时间设置，那么无论协议如何，导出程序将会或不会按照这些环境
变量定义代理流量。

### 身份验证

大多数暴露 HTTP 或 gRPC 端口的接收器都能够使用收集器的身份验证机制来保护，并且大
多数使用 HTTP 或 gRPC 客户端的导出器都能够向传出请求添加身份验证数据。

收集器中的身份验证机制使用扩展机制，允许将自定义身份验证器插入收集器发行版中。如
果您对开发自定义身份验证器感兴趣，请查
看[构建自定义身份验证器](./custom-auth.md)文档。

每个身份验证扩展都有两种可能的用法:作为导出者的客户端身份验证器，向传出请求添加
身份验证数据;作为接收方的服务器身份验证器，对传入连接进行身份验证。请参考身份验
证扩展以了解其功能列表，但通常情况下，身份验证扩展只能实现其中一个特征。有关已知
身份验证器的列表，请使用本网站提供
的[注册表](/ecosystem/registry/?s=authenticator&component=extension)。

若要将服务器身份验证器添加到收集器中的接收器，请确保:

1. 在`.extensions`下添加验证器扩展及其配置
2. 向`.services.extensions`中添加对验证器的引用，以便收集器加载它
3. 在`.receivers.<your-receiver>.<http-or-grpc-config>.auth`下添加对验证器的引用

下面是一个在接收端使用 OIDC 验证器的示例，使其适用于从充当代理的 OpenTelemetry
collector 接收数据的远程收集器:

```yaml
extensions:
  oidc:
    issuer_url: http://localhost:8080/auth/realms/opentelemetry
    audience: collector

receivers:
  otlp/auth:
    protocols:
      grpc:
        auth:
          authenticator: oidc

processors:

exporters:
  logging:

service:
  extensions:
    - oidc
  pipelines:
    traces:
      receivers:
        - otlp/auth
      processors: []
      exporters:
        - logging
```

在代理端，这是一个使 OTLP 导出器获得 OIDC 令牌的示例，并将它们添加到远程收集器的
每个 RPC 中:

```yaml
extensions:
  oauth2client:
    client_id: agent
    client_secret: some-secret
    token_url: http://localhost:8080/auth/realms/opentelemetry/protocol/openid-connect/token

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317

processors:

exporters:
  otlp/auth:
    endpoint: remote-collector:4317
    auth:
      authenticator: oauth2client

service:
  extensions:
    - oauth2client
  pipelines:
    traces:
      receivers:
        - otlp
      processors: []
      exporters:
        - otlp/auth
```

### 设置证书

对于生产设置，我们强烈建议使用 TLS 证书，用于安全通信或 mTLS 用于相互身份验证。
请参见以下步骤生成本示例中使用的自签名证书。您可能希望使用当前的证书供应过程来获
取用于生产的证书。

安装[cfssl](https://github.com/cloudflare/cfssl)，并创建如下 `csr.json` 文件:

```json
{
  "hosts": ["localhost", "127.0.0.1"],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "OpenTelemetry Example"
    }
  ]
}
```

现在，运行以下命令:

```bash
cfssl genkey -initca csr.json | cfssljson -bare ca
cfssl gencert -ca ca.pem -ca-key ca-key.pem csr.json | cfssljson -bare cert
```

这将创建两个证书;首先，`ca.pem`中的"OpenTelemetry 示例"证书颁发机构(CA)
和`ca-key.pem`中的关联密钥客户端证书`cert.pem`(由 OpenTelemetry 示例 CA 签名)和
关联密钥`cert-key.pem`。

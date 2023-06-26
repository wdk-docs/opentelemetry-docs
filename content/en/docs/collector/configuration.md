---
title: 配置
weight: 20
spelling: cSpell:ignore pprof zpages zipkin fluentforward hostmetrics opencensus
spelling: cSpell:ignore prometheus loglevel otlphttp upsert spanevents OIDC
spelling: cSpell:ignore prometheusremotewrite prodevent spanmetrics servicegraph
spelling: cSpell:ignore oidc cfssl genkey initca cfssljson gencert
---

假设熟悉以下页面:

- [数据收集概念](/docs/concepts/data-collection/)以便了解适用于 OpenTelemetry
  Collector 的存储库。
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
(例如: `otlp` or `otlp/2`)。只要标识符是唯一的，给定类型的组件就可以定义多次。例
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

配置也可以包括其他文件，导致 Collector 将两个文件合并到 YAML 配置的单个内存表示
中:

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

## Receivers - 接收器

![](../../assets/img/logos/32x32/Receivers.svg){ width="35" }

接收器可以是基于推或拉的，它是数据进入收集器的方式。接收器可以支持一个或多
个[数据源](/docs/concepts/signals/)。

The `receivers:` section is how receivers are configured. Many receivers come
with default settings so simply specifying the name of the receiver is enough to
configure it (for example, `zipkin:`). If configuration is required or a user
wants to change the default configuration then such configuration must be
defined in this section. Configuration parameters specified for which the
receiver provides a default configuration are overridden.

`receivers:`部分是如何配置接收器的。许多接收器都带有默认设置，因此只需指定接收器
的名称就足以配置它(例如，`zipkin:`)。如果需要配置或者用户希望更改默认配置，则必
须在本节中定义这种配置。将覆盖接收器为其提供默认配置的指定配置参数。

> 配置接收器不会启用它。接收器通过[service](#service)节中的管道启用。

必须配置一个或多个接收器。缺省情况下，没有配置接收器。下面提供了一个接收器的基本
示例。

> 有关接收器的详细配置，请参
> 阅[receiver README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/README.md).

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

## Processors - 处理器

<img width="35" class="img-initial" src="/img/logos/32x32/Processors.svg"></img>

Processors are run on data between being received and being exported. Processors
are optional though
[some are recommended](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#recommended-processors).

The `processors:` section is how processors are configured. Processors may come
with default settings, but many require configuration. Any configuration for a
processor must be done in this section. Configuration parameters specified for
which the processor provides a default configuration are overridden.

> Configuring a processor does not enable it. Processors are enabled via
> pipelines within the [service](#service) section.

A basic example of the default processors is provided below. The full list of
processors can be found by combining the list found
[here](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor)
and
[here](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor).

> For detailed processor configuration, please see the
> [processor README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/README.md).

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

## Exporters - 导出器

<img width="35" class="img-initial" src="/img/logos/32x32/Exporters.svg"></img>

An exporter, which can be push or pull based, is how you send data to one or
more backends/destinations. Exporters may support one or more
[data sources](/docs/concepts/signals/).

The `exporters:` section is how exporters are configured. Exporters may come
with default settings, but many require configuration to specify at least the
destination and security settings. Any configuration for an exporter must be
done in this section. Configuration parameters specified for which the exporter
provides a default configuration are overridden.

> Configuring an exporter does not enable it. Exporters are enabled via
> pipelines within the [service](#service) section.

One or more exporters must be configured. By default, no exporters are
configured. A basic example of exporters is provided below. Certain exporter
configurations require x.509 certificates to be created in order to be secure,
as described in [setting up certificates](#setting-up-certificates).

> For detailed exporter configuration, see the
> [exporter README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/README.md).

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

## Connectors - 连接器

A connector is both an exporter and receiver. As the name suggests a Connector
connects two pipelines: It consumes data as an exporter at the end of one
pipeline and emits data as a receiver at the start of another pipeline. It may
consume and emit data of the same data type, or of different data types. A
connector may generate and emit data to summarize the consumed data, or it may
simply replicate or route data.

The `connectors:` section is how connectors are configured.

> Configuring a connectors does not enable it. Connectors are enabled via
> pipelines within the [service](#service) section.

One or more connectors may be configured. By default, no connectors are
configured. A basic example of connectors is provided below.

> For detailed connector configuration, please see the
> [connector README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/connector/README.md).

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

## Extensions - 扩展

Extensions are available primarily for tasks that do not involve processing
telemetry data. Examples of extensions include health monitoring, service
discovery, and data forwarding. Extensions are optional.

The `extensions:` section is how extensions are configured. Many extensions come
with default settings so simply specifying the name of the extension is enough
to configure it (for example, `health_check:`). If configuration is required or
a user wants to change the default configuration then such configuration must be
defined in this section. Configuration parameters specified for which the
extension provides a default configuration are overridden.

> Configuring an extension does not enable it. Extensions are enabled within the
> [service](#service) section.

By default, no extensions are configured. A basic example of extensions is
provided below.

> For detailed extension configuration, please see the
> [extension README.md](https://github.com/open-telemetry/opentelemetry-collector/blob/main/extension/README.md).

```yaml
extensions:
  health_check:
  pprof:
  zpages:
  memory_ballast:
    size_mib: 512
```

## Service - 服务

The service section is used to configure what components are enabled in the
Collector based on the configuration found in the receivers, processors,
exporters, and extensions sections. If a component is configured, but not
defined within the service section then it is not enabled. The service section
consists of three sub-sections:

- extensions
- pipelines
- telemetry

Extensions consist of a list of all extensions to enable. For example:

```yaml
service:
  extensions: [health_check, pprof, zpages]
```

Pipelines can be of the following types:

- traces: collects and processes trace data.
- metrics: collects and processes metric data.
- logs: collects and processes log data.

A pipeline consists of a set of receivers, processors and exporters. Each
receiver/processor/exporter must be defined in the configuration outside of the
service section to be included in a pipeline.

_Note:_ Each receiver/processor/exporter can be used in more than one pipeline.
For processor(s) referenced in multiple pipelines, each pipeline will get a
separate instance of that processor(s). This is in contrast to
receiver(s)/exporter(s) referenced in multiple pipelines, where only one
instance of a receiver/exporter is used for all pipelines. Also note that the
order of processors dictates the order in which data is processed.

The following is an example pipeline configuration:

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

Telemetry is where the telemetry for the collector itself can be configured. It
has two subsections: `logs` and `metrics`.

The `logs` subsection allows configuration of the logs generated by the
collector. By default the collector will write its logs to stderr with a log
level of `INFO`. You can also add static key-value pairs to all logs using the
`initial_fields` section.
[View the full list of `logs` options here.](https://github.com/open-telemetry/opentelemetry-collector/blob/7666eb04c30e5cfd750db9969fe507562598f0ae/config/service.go#L41-L97)

The `metrics` subsection allows configuration of the metrics generated by the
collector. By default the collector will generate basic metrics about itself and
expose them for scraping at `localhost:8888/metrics`
[View the full list of `metrics` options here.](https://github.com/open-telemetry/opentelemetry-collector/blob/7666eb04c30e5cfd750db9969fe507562598f0ae/config/service.go#L99-L111)

The following is an example telemetry configuration:

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

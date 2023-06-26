---
title: 网关
description:
  Why and how to send signals to a single OTLP end-point and from there to
  backends
weight: 3
---

网关收集器部署模式由应用程序(或其他收集器)组成，这些应用程序(或其他收集器)将遥测
信号发送到作为独立服务运行的一个或多个收集器实例提供的单个 OTLP 端点(例如
，Kubernetes 中的部署)，通常是每个集群、每个数据中心或每个区域。

一般情况下，您可以使用开箱即用的负载均衡器来在收集器之间分配负载:

![Gateway deployment concept](../../img/otel_gateway_sdk.svg)

对于遥测数据处理的处理必须在特定收集器中进行的用例，您可以使用两层设置，其中收集
器的管道在第一层配置了[跟踪 ID/服务名称感知的负载平衡导出程序][lb-exporter]，而
收集器在第二层处理向外扩展。例如，在使用[Tail Sampling 处理
器][tailsample-processor]时，您将需要使用负载平衡导出器，以便给定跟踪的所有跨度
到达应用尾部抽样策略的同一收集器实例。

让我们来看看这样一个使用负载均衡导出器的例子:

![Gateway deployment with load-balancing exporter](../../img/gateway-lb-sdk.svg)

1. 在应用程序中，SDK 被配置为将 OTLP 数据发送到中心位置。
2. 使用负载平衡导出器配置的收集器，它将信号分发到一组收集器。
3. 采集器配置为将遥测数据发送到一个或多个后端。

!!! warning

    目前，负载平衡导出器只支持`traces`类型的管道。

## Example

对于集中式收集器部署模式的具体示例，我们首先需要仔细研究负载平衡导出器。它有两个
主要配置字段:

- `resolver`，它决定在哪里找到下游收集器(或:后端)。如果在这里使用 `static`子键，
  则必须手动枚举收集器的 url。另一个支持的解析器是 DNS 解析器，它将定期检查更新
  和解析 IP 地址。对于这种解析器类型，`hostname`子键指定要查询的主机名，以便获得
  IP 地址列表。
- 使用`routing_key`字段，您可以告诉负载平衡导出器将 spans 路由到特定的下游收集器
  。如果您将此字段设置为`traceID` (默认)，则负载平衡导出程序将根据其`traceID` 导
  出 spans。否则，如果你使用`service`作为`routing_key`的值，它会根据它们的服务名
  称导出 spans，这在使用像[Span Metrics 连接器][spanmetrics-connector]这样的连接
  器时很有用，所以一个服务的所有 spans 将被发送到相同的下游收集器进行度量收集，
  保证准确的聚合。

服务于 OTLP 端点的第一层收集器将按照如下所示进行配置:

=== "Static"

    ```yml
    receivers:
      otlp:
        protocols:
          grpc:

    exporters:
      loadbalancing:
        protocol:
          otlp:
            insecure: true
        resolver:
          static:
            hostnames:
              - collector-1.example.com:4317
              - collector-2.example.com:5317
              - collector-3.example.com

    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [loadbalancing]
    ```

=== "DNS"

    ```yml
    receivers:
      otlp:
        protocols:
          grpc:

    exporters:
      loadbalancing:
        protocol:
          otlp:
            insecure: true
        resolver:
          dns:
            hostname: collectors.example.com

    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [loadbalancing]
    ```

=== "DNS with service"

    ```yml
    receivers:
      otlp:
        protocols:
          grpc:

    exporters:
      loadbalancing:
        routing_key: 'service'
        protocol:
          otlp:
            insecure: true
        resolver:
          dns:
            hostname: collectors.example.com
            port: 5317

    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [loadbalancing]
    ```

负载平衡导出程序发出的指标包括`otelcol_loadbalancer_num_backends` 和
`otelcol_loadbalancer_backend_latency`，您可以使用这些指标监视 OTLP 端点收集器的
运行状况和性能。

## 利弊

优点:

- 关注点分离，例如集中管理的凭据
- 集中策略管理(例如，过滤某些日志或采样)

缺点:

- 它又多了一个需要维护的东西，而且可能会失败(复杂性)
- 增加了级联收集器情况下的延迟
- 更高的整体资源使用(成本)

[lb-exporter]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/loadbalancingexporter
[tailsample-processor]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/tailsamplingprocessor
[spanmetrics-connector]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/connector/spanmetricsconnector

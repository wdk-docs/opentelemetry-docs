---
title: 管理
description: 如何大规模管理OpenTelemetry收集器部署
weight: 23
---

本文档描述了如何大规模管理 OpenTelemetry 收集器部署。

要充分利用此页，您应该了解如何安装和配置收集器。这些主题在其他地方有介绍:

- [入门][otel-collector-getting-started]了解如何安装 OpenTelemetry 收集器。
- [配置][otel-collector-configuration]关于如何配置 OpenTelemetry 收集器，设置遥
  测管道。

## 基础

大规模的遥测收集需要一种结构化的方法来管理代理。典型的代理管理任务包括:

1. 查询座席信息和配置。代理信息可以包括其版本、操作系统相关信息或功能。代理的配
   置指的是它的遥测收集设置，例如，OpenTelemetry 收集器[配
   置][otel-collector-configuration].
2. 升级/降级代理和管理特定于代理的包，包括基本代理功能和插件。
3. 将新配置应用于代理。由于环境的变化或政策的变化，可能需要这样做。
4. 代理的运行状况和性能监视，通常是 CPU 和内存使用情况，以及特定于代理的指标，例
   如，处理速率或与背压相关的信息。
5. 控制平面与代理之间的连接管理，例如处理 TLS 证书(吊销和轮换)。

并非每个用例都需要支持上述所有代理管理任务。在 OpenTelemetry 的上下文中，任务
_4. 健康状况和性能监视_ 最好使用 OpenTelemetry 完成。

## OpAMP

可观察性供应商和云提供商为代理管理提供专有解决方案。在开源可观察性领域，有一个新
兴的标准可以用于代理管理: 开放代理管理协议(OpAMP)。

[OpAMP 规范][opamp-spec]定义了如何管理一组遥测数据代理。这些代理可以
是[OpenTelemetry 收集器][otel-collector]、Fluent Bit 或其他任意组合的代理。

!!! Note

    这里使用术语 "agent" “代理”作为响应OpAMP的OpenTelemetry组件的统称，它可以是收集器，也可以是SDK组件。

OpAMP 是一个客户端/服务器协议，支持通过 HTTP 和 WebSockets 进行通信。

- OpAMP 服务器是控制平面的一部分，充当协调器，管理一组遥测代理。
- **OpAMP 客户端** 是数据平面的一部分。 OpAMP 的客户端可以在进程中实现，例
  如[OpenTelemetry 收集器中的 OpAMP 支持][opamp-in-otel-collector]中的情况。
  OpAMP 的客户端也可以在进程外实现。对于后一种选项，您可以使用管理器，该管理器负
  责与 OpAMP 服务器进行特定于 OpAMP 的通信，同时控制遥测代理，例如应用配置或升级
  它。请注意，主管/遥测通信不是 OpAMP 的一部分。

让我们来看一个具体的设置:

![OpAMP example setup](../img/opamp.svg)

1. OpenTelemetry 收集器，配置了管道，以:
   - (A) 接收来自下游源的信号
   - (B) 输出信号到上游目的地，可能包括有关采集器本身的遥测(由 OpAMP `own_xxx`连
     接设置表示)。
2. 实现服务器端 OpAMP 部分的控制平面与实现客户端 OpAMP 的收集器(或控制收集器的主
   管)之间的双向 OpAMP 控制流。

您可以通过使用[OpAMP 协议在 Go 中的实现][opamp-go]自己尝试一个简单的 OpAMP 设置
。对于下面的演练，您需要使用 1.19 或更高版本的 Go。

我们将设置一个简单的 OpAMP 控制平面，由示例 OpAMP 服务器组成，并让 OpenTelemetry
收集器通过示例 OpAMP 管理器连接到它。

首先，克隆`open-telemetry/opamp-go`repo:

```sh
git clone https://github.com/open-telemetry/opamp-go.git
```

接下来，我们需要 OpAMP 管理器可以管理的 OpenTelemetry 收集器二进制文件。为此，安
装[OpenTelemetry Collector Contrib][otelcolcontrib]发行版。收集器二进制文件的路
径(您将其安装到的位置)在下面被称为`$OTEL_COLLECTOR_BINARY`。

在`. /opamp-go/internal/examples/server`目录，启动 OpAMP 服务器:

```console
$ go run .
2023/02/08 13:31:32.004501 [MAIN] OpAMP Server starting...
2023/02/08 13:31:32.004815 [MAIN] OpAMP Server running...
```

在`./opamp-go/internal/examples/supervisor`目录中创建一个名为`supervisor.yaml`
的文件，其中包含以下内容(告诉 supervisor 在哪里找到服务器以及要管理的
OpenTelemetry 收集器二进制文件):

```yaml
server:
  endpoint: ws://127.0.0.1:4320/v1/opamp

agent:
  executable: $OTEL_COLLECTOR_BINARY
```

!!! Note

    确保将`$OTEL_COLLECTOR_BINARY`替换为实际的文件路径。
    例如，在Linux或macOS中，如果您将收集器安装在`/usr/local/bin/`中，则将`$OTEL_COLLECTOR_BINARY` 替换为`/usr/local/bin/otelcol`。

接下来，按如下方式创建一个收集器配置(将其保存
在`./opamp-go/internal/examples/supervisor`目录下的`effective.yaml`文件中):

```yaml
receivers:
  prometheus/own_metrics:
    config:
      scrape_configs:
        - job_name: otel-collector
          scrape_interval: 10s
          static_configs:
            - targets: [0.0.0.0:8888]
  hostmetrics:
    collection_interval: 10s
    scrapers:
      load:
      filesystem:
      memory:
      network:

exporters:
  logging:
    verbosity: detailed

service:
  pipelines:
    metrics:
      receivers: [hostmetrics, prometheus/own_metrics]
      exporters: [logging]
```

现在是时候启动管理器了(它将依次启动 OpenTelemetry 收集器):

```console
$ go run .
2023/02/08 13:32:54 Supervisor starting, id=01GRRKNBJE06AFVGQT5ZYC0GEK, type=io.opentelemetry.collector, version=1.0.0.
2023/02/08 13:32:54 Starting OpAMP client...
2023/02/08 13:32:54 OpAMP Client started.
2023/02/08 13:32:54 Starting agent /usr/local/bin/otelcol
2023/02/08 13:32:54 Connected to the server.
2023/02/08 13:32:54 Received remote config from server, hash=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855.
2023/02/08 13:32:54 Agent process started, PID=13553
2023/02/08 13:32:54 Effective config changed.
2023/02/08 13:32:54 Enabling own metrics pipeline in the config<F11>
2023/02/08 13:32:54 Effective config changed.
2023/02/08 13:32:54 Config is changed. Signal to restart the agent.
2023/02/08 13:32:54 Agent is not healthy: Get "http://localhost:13133": dial tcp [::1]:13133: connect: connection refused
2023/02/08 13:32:54 Stopping the agent to apply new config.
2023/02/08 13:32:54 Stopping agent process, PID=13553
2023/02/08 13:32:54 Agent process PID=13553 successfully stopped.
2023/02/08 13:32:54 Starting agent /usr/local/bin/otelcol
2023/02/08 13:32:54 Agent process started, PID=13554
2023/02/08 13:32:54 Agent is not healthy: Get "http://localhost:13133": dial tcp [::1]:13133: connect: connection refused
2023/02/08 13:32:55 Agent is not healthy: health check on http://localhost:13133 returned 503
2023/02/08 13:32:55 Agent is not healthy: health check on http://localhost:13133 returned 503
2023/02/08 13:32:56 Agent is not healthy: health check on http://localhost:13133 returned 503
2023/02/08 13:32:57 Agent is healthy.
```

如果一切都解决了，你现在应该能够访问 http://localhost:4321/ ，并访问 OpAMP 服务
器 UI，你应该看到你的收集器列出，由主管管理:

![OpAMP example setup](../img/opamp-server-ui.png)

您还可以查询采集器导出的指标(注意标签值):

```console
$ curl localhost:8888/metrics
...
# HELP otelcol_receiver_accepted_metric_points Number of metric points successfully pushed into the pipeline.
# TYPE otelcol_receiver_accepted_metric_points counter
otelcol_receiver_accepted_metric_points{receiver="prometheus/own_metrics",service_instance_id="01GRRKNBJE06AFVGQT5ZYC0GEK",service_name="io.opentelemetry.collector",service_version="1.0.0",transport="http"} 322
# HELP otelcol_receiver_refused_metric_points Number of metric points that could not be pushed into the pipeline.
# TYPE otelcol_receiver_refused_metric_points counter
otelcol_receiver_refused_metric_points{receiver="prometheus/own_metrics",service_instance_id="01GRRKNBJE06AFVGQT5ZYC0GEK",service_name="io.opentelemetry.collector",service_version="1.0.0",transport="http"} 0
```

## 其它信息

- 博客文章[使用 OpenTelemetry OpAMP 修改 go 语言上的服务遥
  测][blog-opamp-service-telemetry]
- YouTube 视频:
  - [闪电谈话:通过 OpAMP 协议管理 OpenTelemetry][opamp-lt]
  - [什么是 OpAMP 和什么是 BindPlane][opamp-bindplane]

[otel-collector]: ./index.md
[otel-collector-getting-started]: ./getting-started.md
[otel-collector-configuration]: ./configuration.md
[opamp-spec]:
  https://github.com/open-telemetry/opamp-spec/blob/main/specification.md
[opamp-in-otel-collector]:
  https://docs.google.com/document/d/1KtH5atZQUs9Achbce6LiOaJxLbksNJenvgvyKLsJrkc/edit#heading=h.ioikt02qpy5f
[opamp-go]: https://github.com/open-telemetry/opamp-go
[otelcolcontrib]:
  https://github.com/open-telemetry/opentelemetry-collector-releases/releases
[blog-opamp-service-telemetry]: ../../blog/2022/opamp/index.md
[opamp-lt]: https://www.youtube.com/watch?v=LUsfZFRM4yo
[opamp-bindplane]: https://www.youtube.com/watch?v=N18z2dOJSd8

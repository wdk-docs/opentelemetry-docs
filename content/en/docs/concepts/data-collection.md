---
title: 数据收集
description: >-
  OpenTelemetry项目通过OpenTelemetry Collector促进了遥测数据的收集。
weight: 50
---

OpenTelemetry项目通过OpenTelemetry Collector促进了遥测数据的收集。
OpenTelemetry Collector提供了一个与供应商无关的实现，用于接收、处理和导出遥测数据。
它消除了运行、操作和维护多个代理/收集器以支持向一个或多个开源或商业后端发送的开源可观察性数据格式(例如Jaeger、Prometheus等)的需要。
此外，Collector还为最终用户提供了对其数据的控制。
Collector是仪器库导出其遥测数据的默认位置。

> 收集器可以作为发行版提供，请参阅[这里](../distribution)了解更多信息。

## 部署

开放遥测采集器提供单一二进制和两种部署方法:

- **Agent:** A Collector instance running with the application or on the same
  host as the application (e.g. binary, sidecar, or daemonset).
- **Gateway:** One or more Collector instances running as a standalone service
  (e.g. container or deployment) typically per cluster, data center or region.

For information on how to use the Collector see the
[getting started documentation](/docs/collector/getting-started).

## 组件

The Collector is made up of the following components:

- <img width="32" class="img-initial" src="/img/logos/32x32/Receivers.svg"></img>
  `receivers`: How to get data into the Collector; these can be push or pull
  based
- <img width="32" class="img-initial" src="/img/logos/32x32/Processors.svg"></img>
  `processors`: What to do with received data
- <img width="32" class="img-initial" src="/img/logos/32x32/Exporters.svg"></img>
  `exporters`: Where to send received data; these can be push or pull based

These components are enabled through `pipelines`. Multiple instances of
components as well as pipelines can be defined via YAML configuration.

For more information about these components see the
[configuration documentation](/docs/collector/configuration).

## 存储库

The OpenTelemetry project provides two versions of the Collector:

- **[Core](https://github.com/open-telemetry/opentelemetry-collector/releases):**
  Foundational components such as configuration and generally applicable
  receivers, processors, exporters, and extensions.
- **[Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases):**
  All the components of core plus optional or possibly experimental components.
  Offers support for popular open source projects including Jaeger, Prometheus,
  and Fluent Bit. Also contains more specialized or vendor-specific receivers,
  processors, exporters, and extensions.

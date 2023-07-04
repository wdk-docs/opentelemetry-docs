---
title: 数据收集
description: >-
  OpenTelemetry项目通过OpenTelemetry Collector促进了遥测数据的收集。
weight: 50
---

OpenTelemetry 项目通过 Collector 促进了遥测数据的收集。 Collector 提供了一个与供
应商无关的实现，用于接收、处理和导出遥测数据。它消除了运行、操作和维护多个代理/
收集器以支持向一个或多个开源或商业后端发送的开源可观察性数据格式(例如
Jaeger、Prometheus 等)的需要。此外，Collector 还为最终用户提供了对其数据的控制。
Collector 是插装库导出其遥测数据的默认位置。

> 收集器可以作为发行版提供，请参阅[这里](./distribution.md)了解更多信息。

## 部署

开放遥测采集器提供单一二进制和两种部署方法:

- **代理:** 与应用程序一起运行或在与应用程序相同的主机上运行的 Collector 实例(
  例如 binary、sidecar 或 daemonset)。
- **网关:** 一个或多个 Collector 实例作为独立服务(例如容器或部署)运行，通常为每
  个集群、数据中心或区域。

有关如何使用收集器的信息，请参阅[入门文档](../collector/getting-started.md).

## 组件

收集器由以下组件组成:

- ![](../../assets/img/logos/32x32/Receivers.svg){ width="32" } `receivers`: 如
  何让数据进入收集器;这些可以是推或拉的基础
- ![](../../assets/img/logos/32x32/Processors.svg){ width="32" } `processors`:
  如何处理接收到的数据
- ![](../../assets/img/logos/32x32/Exporters.svg){ width="32" } `exporters`: 将
  接收到的数据发送到哪里;这些可以是推或拉的基础

这些组件是通过“管道”启用的。组件的多个实例以及管道可以通过 YAML 配置来定义。

有关这些组件的更多信息，请参阅[配置文档](../collector/configuration.md).

## 存储库

OpenTelemetry 项目提供了两个版本的 Collector:

- **[Core](https://github.com/open-telemetry/opentelemetry-collector/releases):**
  基本组件，如配置和一般适用的接收器、处理器、导出器和扩展。
- **[Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases):**
  所有核心组件加上可选或可能的实验组件。为流行的开源项目提供支持，包括
  Jaeger、Prometheus 和 Fluent Bit。还包含更专业或特定于供应商的接收器、处理器、
  导出器和扩展。

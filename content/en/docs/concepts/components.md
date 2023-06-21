---
title: 组件
description: 构成OpenTelemetry的主要组件
weight: 20
---

OpenTelemetry 目前由几个主要组件组成:

- [跨语言规范](/docs/specs/otel/)
- [OpenTelemetry 收集器](/docs/collector/)
- [每种语言 sdk](/docs/instrumentation/)
- [每种语言的工具库](/docs/concepts/instrumentation/libraries/)
- [按语言自动检测](/docs/concepts/instrumentation/automatic/)
- [K8s 操作器](/docs/k8s-operator/)

OpenTelemetry 允许您替换对特定于供应商的 sdk 和工具的需求，以生成和导出遥测数据
。

## 规范

描述所有实现的跨语言需求和期望。除了对术语的定义之外，该规范还定义了以下内容:

- **API:** 定义用于生成和关联跟踪、度量和日志数据的数据类型和操作。
- **SDK:** 定义特定于语言的 API 实现的需求。这里还定义了配置、数据处理和导出概念
  。
- **Data:** 定义遥测后端可以提供支持的 OpenTelemetry 协议(OTLP)和与供应商无关的
  语义约定。

有关更多信息，请参阅[规范](../specs/otel.md).

此外，API 概念的广泛注释的 protobuf 接口文件可以
在[原型存储库](https://github.com/open-telemetry/opentelemetry-proto)中找到。

## 收集器

OpenTelemetry Collector 是一个与供应商无关的代理，它可以接收、处理和导出遥测数据
。它支持以多种格式接收遥测数据(例如，OTLP、Jaeger、Prometheus 以及许多商业/专有
工具)，并将数据发送到一个或多个后端。它还支持在遥测数据导出之前对其进行处理和过
滤。收集器贡献包支持更多的数据格式和供应商后端。

有关更多信息，请参见[收集器](../collector/index.md).

## 语言 SDKs

OpenTelemetry 也有语言 SDKs，允许您使用 OpenTelemetry API 用您选择的语言生成遥测
数据，并将该数据导出到首选后端。这些 SDKs 还允许您为通用库和框架合并检测库，您可
以使用这些库和框架连接到应用程序中的手动检测。

有关详细信息，请参见[检测](/docs/concepts/instrumentation/).

## 工具库

OpenTelemetry 支持大量的组件，这些组件从受支持语言的流行库和框架中生成相关的遥测
数据。例如，来自 HTTP 库的入站和出站 HTTP 请求将生成关于这些请求的数据。

我们的长期目标是将流行的库编写为开箱即用的可观察库，这样就不需要引入单独的组件。

有关更多信息，请参见[检测库](/docs/concepts/instrumentation/libraries/).

## 自动仪表

如果适用，OpenTelemetry 的特定语言实现将提供一种方法来检测您的应用程序，而无需触
及您的源代码。虽然底层机制取决于语言，但至少这会将 OpenTelemetry API 和 SDK 功能
添加到您的应用程序中。此外，他们可能会添加一组 Instrumentation Libraries 和导出
器依赖项。

有关详细信息，请参见[检测](../concepts/instrumentation/automatic.md).

## K8s operator

OpenTelemetry Operator 是 Kubernetes Operator 的一个实现。Operator 使用
OpenTelemetry 管理 OpenTelemetry 收集器和工作负载的自动检测。

有关更多信息，请参阅[K8s Operator](../k8s-operator/index.md).

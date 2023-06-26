---
title: 手动插装
description: >-
  了解手动检测应用程序的基本步骤。
weight: 20
---

## 导入 OpenTelemetry API 和 SDK

首先需要将 OpenTelemetry 导入到服务代码中。如果您正在开发一个库或其他打算由可运
行二进制文件使用的组件，那么您只需要依赖于 API。如果您的工件是一个独立的流程或服
务，那么您将依赖于 API 和 SDK。有关 OpenTelemetry API 和 SDK 的更多信息，请参
见[specification](/docs/specs/otel/)。

## 配置 OpenTelemetry API

为了创建跟踪或度量，您需要首先创建跟踪程序和/或度量提供程序。通常，我们建议 SDK
应该为这些对象提供单个默认提供程序。然后，您将从该提供程序获得跟踪程序或插装实例
，并为其提供名称和版本。您在这里选择的名称应该确定要检测的确切内容——例如，如果您
正在编写一个库，那么您应该以您的库(例如`com.legitimatebusiness.myLibrary`)来命名
它，因为这个名称将命名生成的所有跨度或度量事件。还建议您提供与库或服务的当前版本
对应的版本字符串(即' semver:1.0.0 ')。

## 配置 OpenTelemetry SDK

如果您正在构建一个服务流程，您还需要为 SDK 配置适当的选项，以便将遥测数据导出到
某个分析后端。我们建议通过配置文件或其他机制以编程方式处理此配置。您可能还希望利
用不同语言的调优选项。

## 创建遥测数据

配置好 API 和 SDK 之后，您就可以通过从提供程序获得的跟踪器和度量对象自由地创建跟
踪和度量事件。为你的依赖项使用工具库——查看[registry](/ecosystem/registry/)或你的
语言的存储库，了解更多相关信息。

## 出口数据

一旦您创建了遥测数据，您将希望将其发送到某个地方。 OpenTelemetry 支持将数据从进
程导出到分析后端的两种主要方法，要么直接从进程导出，要么通
过[OpenTelemetry Collector](/docs/Collector)进行代理。

进程内导出需要您导入并依赖于一个或多个 _exporters_，这些库将 OpenTelemetry 的内
存跨度和度量对象转换为适合 Jaeger 或 Prometheus 等遥测分析工具的格式。此外
，OpenTelemetry 还支持一种名为“OTLP”的有线协议，所有 OpenTelemetry sdk 都支持该
协议。该协议可用于向 OpenTelemetry Collector 发送数据，OpenTelemetry Collector
是一个独立的二进制进程，可以作为服务实例的代理或侧车运行，也可以在单独的主机上运
行。然后可以配置 Collector 来转发和导出该数据到您选择的分析工具。

除了像 Jaeger 或 Prometheus 这样的开源工具之外，越来越多的公司支持从
OpenTelemetry 获取遥测数据。详情请参见[vendor](/ecosystem/vendor/)。

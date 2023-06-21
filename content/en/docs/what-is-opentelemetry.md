---
title: 什么是遥测?
description: 一个简短的解释什么是OpenTelemetry，什么不是。
aliases: [/about, /docs/concepts/what-is-opentelemetry, /otel]
weight: -1
---

微服务架构使开发人员能够更快、更独立地构建和发布软件，因为他们不再受制于与单片架构相关的复杂发布过程。

随着这些分布式系统的扩展，开发人员越来越难以看到他们自己的服务如何依赖或影响其他服务，特别是在部署之后或中断期间，速度和准确性至关重要。

> [可观察性](/docs/concepts/observability-primer/#what-is-observability)使开发人员和操作人员都可以获得对其系统的可见性。

## 那又怎样?

为了使系统可观察，必须对其进行检测。
也就是说，代码必须发出[trace](/docs/concepts/observability-primer/#distributed-traces)、[metrics](/docs/concepts/observability-primer/#reliability- metrics)和[logs](/docs/concepts/observability-primer/#logs)。
然后必须将检测的数据发送到可观察性后端。
现在有很多可观察性后端，从自托管的开源工具(例如[Jaeger](https://www.jaegertracing.io/)和[Zipkin](https://zipkin.io/))到商业SaaS产品。

在过去，检测代码的方式各不相同，因为每个可观察的后端都有自己的检测库和代理，用于向工具发送数据。

这意味着没有标准的数据格式来将数据发送到可观察性后端。
此外，如果一家公司选择切换Observability后端，这意味着他们必须重新编写代码并配置新的代理，以便能够向所选择的新工具发送遥测数据。

> 由于缺乏标准化，最终的结果是缺乏数据可移植性，并且增加了用户维护仪器库的负担。

认识到标准化的需要，云社区聚集在一起，两个开源项目诞生了:[OpenTracing](https://opentracing.io)(一个[cloud Native Computing Foundation (CNCF)](https://www.cncf.io)项目)和[OpenCensus](https://opencensus.io)(一个[Google opensource](https://opensource.google)社区项目)。

**OpenTracing** 提供了一个厂商中立的API，用于将遥测数据发送到可观察的后端;然而，它依赖于开发人员实现他们自己的库来满足规范。

**OpenCensus** 提供了一组特定于语言的库，开发人员可以使用它们来检测他们的代码并将其发送到他们支持的任何一个后端。

## 你好,OpenTelemetry !

为了实现单一标准，OpenCensus和OpenTracing于[2019年5月][cncf-incubating-project]合并为OpenTelemetry(简称OTel)。
作为CNCF的孵化项目，OpenTelemetry兼收并举。

OTel的目标是提供一套标准化的、与供应商无关的sdk、api和工具(/docs/collector)，用于摄取、转换和发送数据到可观察性后端(即开源或商业供应商)。

## OpenTelemetry能为我做什么?

OTel拥有广泛的行业支持和云提供商、[厂商](/ecosystem/vendors/)和终端用户的采用。它为您提供:

- 一个单一的、与厂商无关的检测库[每种语言](/docs/instrumentation)，支持自动和手动检测。
- 一个独立于供应商的[收集器](/docs/collector)二进制文件，可以以多种方式部署。
- 生成、发出、收集、处理和导出遥测数据的端到端实现。
- 完全控制您的数据，并能够通过配置将数据并行发送到多个目的地。
- 开放标准语义约定，确保与供应商无关的数据收集
- 能够并行支持多种[上下文传播](/docs/specs/otel/overview/#context-propagation)格式，以帮助随着标准的发展而迁移。
- 无论你在可观察性之旅的哪个位置，都是一条前进的道路。

通过对各种[开源和商业协议][otel-collector-contrib]，格式和上下文传播机制的支持，以及为OpenTracing和OpenCensus项目提供shims，采用OpenTelemetry很容易。

## OpenTelemetry不是什么

OpenTelemetry不像Jaeger或Prometheus那样是一个可观察的后端。
相反，它支持将数据导出到各种开源和商业后端。
它提供了一个可插拔的架构，因此可以很容易地添加其他技术协议和格式。

## 下一个什么?

- [入门](/docs/ Getting -started/) &mdash;跳进去吧!
- 了解[OpenTelemetry概念](/docs/concepts/).

[cncf-incubating-project]: https://www.cncf.io/blog/2021/08/26/opentelemetry-becomes-a-cncf-incubating-project/
[otel-collector-contrib]: https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver

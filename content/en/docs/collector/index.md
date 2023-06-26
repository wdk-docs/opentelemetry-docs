---
title: 收集器
description: Vendor-agnostic way to receive, process and export telemetry data.
aliases: [/docs/collector/about]
cascade:
  collectorVersion: 0.79.0
weight: 10
---

![OpenTelemetry Collector diagram with Jaeger, OTLP and Prometheus integration](img/otel-collector.svg)

## 介绍

OpenTelemetry Collector 提供了一个与供应商无关的实现，用于接收、处理和导出遥测数
据。它消除了运行、操作和维护多个代理/收集器的需要。它具有改进的可伸缩性，并支持
将开源可观察数据格式(例如 Jaeger、Prometheus、Fluent Bit 等)发送到一个或多个开源
或商业后端。本地 Collector 代理是仪器库将其遥测数据导出到的默认位置。

## 目标

- _可用性_: 合理的默认配置，支持流行的协议，开箱即用的运行和收集。
- _表演_: 在不同的负载和配置下高度稳定和高性能。
- _可观察性_: 可观察服务的范例。
- _可扩展性_: 无需触及核心代码即可自定义。
- _统一_: 单个代码库，可作为支持跟踪、度量和日志(未来)的代理或收集器进行部署。

## 何时使用收集器

对于大多数特定于语言的工具库，您都有针对流行后端和 OTLP 的导出器。你可能会想，

> 在什么情况下使用收集器发送数据，而不是让每个服务直接发送到后端?

对于尝试和开始使用 OpenTelemetry，将数据直接发送到后端是快速获取价值的好方法。此
外，在开发或小规模环境中，您可以在没有收集器的情况下获得不错的结果。

但是，通常我们建议在服务旁边使用收集器，因为它允许您的服务快速卸载数据，并且收集
器可以处理额外的处理，如重试，批处理，加密，甚至敏感数据过滤。

[设置收集器](./getting-started.md)也比您想象的要容易:每种语言的默认 OTLP 导出器
都假定有一个本地收集器端点，因此您启动收集器并开始进行遥测。

## 状态和发布

**收集器** 状态为:[mixed][]，因为核心收集器组件当前具有混合的[稳定级别][]。

**收集器组件** 成熟度不同。正在努力确保每个组件都有其稳定性文档。要跟踪此工作的
进展，请参阅`opentelemetry-collector-contrib` [issue #10116][].

[issue #10116]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/10116
[mixed]: /docs/specs/otel/document-status/#mixed
[稳定级别]:
  https://github.com/open-telemetry/opentelemetry-collector#stability-levels

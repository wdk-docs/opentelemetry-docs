---
title: 没有收集器
description: Why and how to send signals directly from app to backends
weight: 1
---

最简单的模式是根本不使用收集器。该模式由带有 OpenTelemetry SDK 的应用程
序[instrumented][instrumentation]组成，该 SDK 将遥测信号(跟踪，度量，日志)直接导
出到后端:

![没有收集器部署概念](../../img/otel_sdk.svg)

## 例子

请参阅[编程语言的代码插装][instrumentation]了解如何将信号从应用程序直接导出到后
端具体的端到端示例。

## 权衡

优点:

- 使用简单(特别是在开发/测试环境中)
- 没有额外的移动部件需要操作(在生产环境中)

缺点:

- 如果收集、处理或摄取发生变化，则需要更改代码
- 应用程序代码和后端之间的强耦合
- 每种语言实现的导出器数量有限

[instrumentation]: ../../instrumentation/index.md

---
title: 分发
description: >-
  分发版(不要与fork混淆)是OpenTelemetry组件的定制版本。
weight: 90
---

OpenTelemetry 项目由多个支持多个[信号](../signals)的[组件](../components) 组成。
OpenTelemetry 的参考实现如下:

- [特定于语言的工具库](../instrumentation)
- [收集器二进制文件](../data-collection)

可以从任何参考实现创建一个发行版。

## 什么是分发?

分发版(不要与 fork 混淆)是 OpenTelemetry 组件的定制版本。发行版是上游
OpenTelemetry 存储库的包装器，带有一些定制。发行版中的自定义可能包括:

- 为特定后端或供应商简化使用或自定义使用的脚本
- 更改后端、供应商或最终用户所需的默认设置
- 可能是供应商或最终用户特定的附加包装选项
- 测试、性能和安全覆盖超出了 OpenTelemetry 提供的范围
- OpenTelemetry 提供的功能之外的其他功能
- OpenTelemetry 提供的功能更少

分发将大致分为以下几类:

- **"Pure":** 这些发行版提供与上游版本相同的功能，并且 100%兼容。定制通常是为了
  便于使用或打包。这些定制可能是特定于后端、供应商或最终用户的。
- **"Plus":** 这些发行版提供了与上游版本相同的功能。除了在纯发行版中发现的定制之
  外，还包括其他组件。这方面的例子包括没有上溯到 OpenTelemetry 项目的仪器库或供
  应商导出程序。
- **"Minus":** 这些发行版提供了来自上游的一组简化的功能。这方面的例子包括移除
  OpenTelemetry Collector 项目中的仪器库或接收器/处理器/导出器/扩展。提供这些发
  行版可能是为了增加可支持性和安全性考虑。

## 谁会创造一个分发?

任何人都可以创建一个发行版。今天，一些[供应商](/ecosystem/vendors/)提供发行版。
此外，如果终端用户希望在[Registry](/ecosystem/registry/) 中使用没有上行到
OpenTelemetry 项目的组件，他们可能会考虑创建一个发行版。

## 贡献还是分发?

在你继续阅读并学习如何创建你自己的发行版之前，问问你自己，你在 OpenTelemetry 组
件之上添加的东西是否对每个人都有益，因此应该包含在参考实现中:

- Can your scripts for "ease of use" be generalized?
- Can your changes to default settings be the better option for everyone?
- Are your additional packaging options really specific?
- Might your test, performance & security coverage work with the reference
  implementation as well?
- Have you checked with the community if your additional capabilities could be
  part of the standard?

## 创建自己的发行版

### Collector

A guide on how to create your own distribution is available in this blog post:
["Building your own OpenTelemetry Collector distribution"](https://medium.com/p/42337e994b63)

If you are building your own distribution, the
[OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder)
might be a good starting point.

### 特定于语言的工具库

There are language specific extensibility mechanisms to customize the
instrumentation libraries:

- [Javaagent](../../instrumentation/java/automatic/extensions)

## 关于发行版，你应该知道些什么

When using OpenTelemetry project collateral such as logo and name for your
distribution, make sure that you are in line with the [OpenTelemetry Marketing
Guidelines for Contributing Organizations][guidelines].

The OpenTelemetry project does not certify distributions at this time. In the
future, OpenTelemetry may certify distributions and partners similarly to the
Kubernetes project. When evaluating a distribution, ensure using the
distribution does not result in vendor lock-in.

> Any support for a distribution comes from the distribution authors and not the
> OpenTelemetry authors.

[guidelines]:
  https://github.com/open-telemetry/community/blob/main/marketing-guidelines.md

---
title: 发行版
description: >-
  发行版(不要与fork混淆)是OpenTelemetry组件的定制版本。
weight: 90
---

OpenTelemetry 项目由多个支持多个[信号](../signals)的[组件](../components) 组成。
OpenTelemetry 的参考实现如下:

- [特定于语言的工具库](../instrumentation)
- [收集器二进制文件](../data-collection)

可以从任何参考实现创建一个发行版。

## 什么是发行版?

发行版(不要与 fork 混淆)是 OpenTelemetry 组件的定制版本。发行版是上游
OpenTelemetry 存储库的包装器，带有一些定制。发行版中的自定义可能包括:

- 为特定后端或供应商简化使用或自定义使用的脚本
- 更改后端、供应商或最终用户所需的默认设置
- 可能是供应商或最终用户特定的附加包装选项
- 测试、性能和安全覆盖超出了 OpenTelemetry 提供的范围
- OpenTelemetry 提供的功能之外的其他功能
- OpenTelemetry 提供的功能更少

发行版将大致分为以下几类:

- **"Pure":** 这些发行版提供与上游版本相同的功能，并且 100%兼容。定制通常是为了
  便于使用或打包。这些定制可能是特定于后端、供应商或最终用户的。
- **"Plus":** 这些发行版提供了与上游版本相同的功能。除了在纯发行版中发现的定制之
  外，还包括其他组件。这方面的例子包括没有上溯到 OpenTelemetry 项目的插装库或供
  应商导出程序。
- **"Minus":** 这些发行版提供了来自上游的一组简化的功能。这方面的例子包括移除
  OpenTelemetry Collector 项目中的插装库或接收器/处理器/导出器/扩展。提供这些发
  行版可能是为了增加可支持性和安全性考虑。

## 谁会创造一个发行版?

任何人都可以创建一个发行版。今天，一些[供应商](/ecosystem/vendors/)提供发行版。
此外，如果终端用户希望在[Registry](/ecosystem/registry/) 中使用没有上行到
OpenTelemetry 项目的组件，他们可能会考虑创建一个发行版。

## 贡献还是发行版本?

在你继续阅读并学习如何创建你自己的发行版之前，问问你自己，你在 OpenTelemetry 组
件之上添加的东西是否对每个人都有益，因此应该包含在参考实现中:

- 您的“易用性”脚本是否可以一般化?
- 更改默认设置对每个人来说都是更好的选择吗?
- 你的附加包装选项真的很具体吗?
- 您的测试、性能和安全覆盖是否也适用于参考实现?
- 您是否与社区确认过您的附加功能是否可以成为标准的一部分?

## 创建自己的发行版

### 收集器

关于如何创建自己的发行版的指南可以在这篇博客文章中找到
:[“构建自己的 OpenTelemetry 收集器发行版”](https://medium.com/p/42337e994b63)

如果您正在构建自己的发行版
，[OpenTelemetry 收集器构建器](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder)可
能是一个很好的起点。

### 特定于语言的工具库

有特定于语言的扩展机制来定制插装库:

- [Javaagent](../../instrumentation/java/automatic/extensions)

## 关于发行版，你应该知道些什么

在为您的分销使用 OpenTelemetry 项目附属品(如徽标和名称)时，请确保您符
合[OpenTelemetry 贡献组织营销指南][guidelines]。

OpenTelemetry 项目目前不认证发行版。在未来，OpenTelemetry 可能会像 Kubernetes 项
目一样认证发行版和合作伙伴。在评估发行版时，确保使用该发行版不会导致供应商锁定。

> 对发行版的任何支持都来自发行版的作者，而不是 OpenTelemetry 的作者。

[guidelines]:
  https://github.com/open-telemetry/community/blob/main/marketing-guidelines.md

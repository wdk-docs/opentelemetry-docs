---
title: 可观测入门指南
description: 核心可观察性概念。
weight: 9
spelling: cSpell:ignore KHTML
---

## 什么是可观察性?

可观察性让我们在不知道系统内部运作的情况下对系统提出问题，从而从外部理解系统。此
外，它允许我们轻松地排除故障并处理新问题(即“未知的未知”)，并帮助我们回答“为什么
会发生这种情况?”

为了能够对系统提出这些问题，必须对应用程序进行适当的检测。也就是说，应用程序代码
必须发出[信号](/docs/concepts/signals/)，例
如[trace](/docs/concepts/observability-primer/#distributed-traces)，
[metrics](/docs/concepts/observability-primer/#reliability- metrics)
和[logs](/docs/concepts/observability-primer/#logs)。当开发人员不需要添加更多的
检测来解决问题时，应用程序就被适当地检测了，因为他们已经拥有了所需的所有信息。

[OpenTelemetry](/docs/what-is-opentelemetry/)是应用程序代码被检测的机制，以帮助
使系统可观察。

## 可靠性和度量

遥测是指从系统发出的有关其行为的数据。数据可以
以[痕迹](#distributed-traces)、[指标](#reliability- metrics)和[日志](#logs)的形
式出现。

**可靠性** 回答的问题是:“服务是否在做用户期望它做的事情?”系统可以 100%正常运行，
但如果当用户点击“添加到购物车”将一条黑色裤子添加到购物车时，系统却一直在添加一条
红色裤子，那么系统就会被认为是**不**可靠的。

**指标** 是一段时间内关于基础设施或应用程序的数字数据的聚合。示例包括:系统错误率
、CPU 利用率、给定服务的请求率。

**SLI** ，即服务水平指标，表示对服务行为的度量。好的 SLI 从用户的角度来衡量您的
服务。例如，SLI 可以是网页加载的速度。

**SLO** ，即服务水平目标，是向组织/其他团队传达可靠性的手段。这可以通过将一个或
多个 sli 附加到业务值来实现。

## 理解分布式跟踪

为了理解分布式跟踪，让我们从一些基础知识开始。

### 日志

日志是由服务或其他组件发出的带有时间戳的消息。然而，
与[traces](#distributed-traces)不同，它们不一定与任何特定的用户请求或事务相关联
。它们在软件中几乎无处不在，并且在过去被开发人员和操作人员严重依赖，以帮助他们理
解系统行为。

示例日志:

```text
I, [2021-02-23T13:26:23.505892 #22473]  INFO -- : [6459ffe1-ea53-4044-aaa3-bf902868f730] Started GET "/" for ::1 at 2021-02-23 13:26:23 -0800
```

不幸的是，日志对于跟踪代码执行并不是非常有用，因为它们通常缺乏上下文信息，例如从
哪里调用它们。

当它们作为[span](#span)的一部分包含时，它们将变得更加有用。

### Spans

**span** 表示一个工作或操作单元。 它跟踪请求所做的特定操作，描绘出在执行该操作期
间发生的情况。

span 包含名称、与时间相关的数据
、[结构化日志消息](/docs/concepts/signals/traces/#span-events)和[其他元数据(即属性)](/docs/concepts/signals/traces/#
Attributes)，以提供关于它所跟踪的操作的信息。

#### Span 属性

下表包含了 span 属性的示例:

| Key              | Value                                                                                                                   |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------- |
| net.transport    | `IP.TCP`                                                                                                                |
| net.peer.ip      | `10.244.0.1`                                                                                                            |
| net.peer.port    | `10243`                                                                                                                 |
| net.host.name    | `localhost`                                                                                                             |
| http.method      | `GET`                                                                                                                   |
| http.target      | `/cart`                                                                                                                 |
| http.server_name | `frontend`                                                                                                              |
| http.route       | `/cart `                                                                                                                |
| http.scheme      | `http`                                                                                                                  |
| http.host        | `localhost`                                                                                                             |
| http.flavor      | `1.1`                                                                                                                   |
| http.status_code | `200`                                                                                                                   |
| http.user_agent  | `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36` |

For more on spans and how they pertain to OTel, see
[Spans](/docs/concepts/signals/traces/#spans).

### 分布式跟踪

分布式跟踪，通常称为跟踪，记录请求(由应用程序或最终用户发出)通过多服务架构(如微
服务和无服务器应用程序)传播时所采取的路径。

如果不进行跟踪，就很难确定分布式系统中性能问题的原因。

它提高了应用程序或系统运行状况的可见性，并允许我们调试难以在本地重现的行为。跟踪
对于分布式系统至关重要，因为分布式系统通常存在不确定性问题，或者过于复杂而无法在
本地重现。

跟踪通过分解请求流经分布式系统时所发生的事情，使调试和理解分布式系统变得不那么令
人生畏。

迹线由一个或多个跨度组成。第一个跨度表示根跨度。每个根跨度代表一个从头到尾的请求
。父类下面的跨提供了一个更深入的上下文，说明在请求期间发生了什么(或组成请求的步
骤)。

许多可观察性后端将轨迹可视化为瀑布图，如下图所示:

![Sample Trace](../../assets/img/waterfall_trace.png 'Trace waterfall diagram')

瀑布图显示了根跨度与其子跨度之间的父子关系。当一个 span 封装另一个 span 时，这也
表示嵌套关系。

有关 trace 及其与 OTel 的关系的更多信息，请参
见[trace](/docs/concepts/signals/traces/).

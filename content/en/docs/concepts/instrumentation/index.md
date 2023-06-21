---
title: 仪表
description: >-
  OpenTelemetry如何促进应用程序的自动和手动仪器。
aliases: [/docs/concepts/instrumenting]
weight: 15
---

为了使系统可观察，它必须 **仪器化** :也就是说，来自系统组件的代码必须发出[跟踪](/docs/concepts/observability-primer/#distributed-traces)，[度量](/docs/concepts/observability-primer/#reliability- metrics)和[日志](/docs/concepts/observability-primer/#logs)。

不需要修改源代码，您就可以使用[automatic instrumentation](automatic/)从应用程序收集遥测数据。
如果您以前使用APM代理从应用程序中提取遥测数据，那么自动检测将为您提供类似的开箱即用体验。

为了更方便地检测应用程序，您可以通过对OpenTelemetry api进行编码来[手动检测](manual/manual)应用程序。

为此，你不需要检测应用程序中使用的所有依赖项:

- 通过直接调用OpenTelemetry API，你的一些库将是开箱即用的可观察的。
  这些库有时被称为 **本机插装库**。
- 对于没有这种集成的库，OpenTelemetry项目提供了特定于语言的[仪器库][]

请注意，对于大多数语言，可以同时使用手动和自动插装:自动插装将允许您快速了解应用程序，而手动插装将使您能够将粒度可观察性嵌入到代码中。

[manual](manual/)和[automatic](automatic/)检测的确切安装机制因您所使用的开发语言而异，但下面几节将介绍一些相似之处。

[工具库]: /docs/specs/otel/overview/#instrumentation-libraries

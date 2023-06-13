---
title: 自动
description: >-
  Learn how Automatic Instrumentation can add observability to your application
  without the need to touch your code
weight: 10
---

如果适用，OpenTelemetry的特定语言实现将提供一种方法来检测您的应用程序，而无需触及您的源代码。
虽然底层机制取决于语言，但至少这会将OpenTelemetry API和SDK功能添加到您的应用程序中。
此外，他们可能会添加一组工具库和导出器依赖项。

可以通过环境变量和特定于语言的方式(如Java中的系统属性)进行配置。
至少，必须配置服务名称来标识被检测的服务。
各种其他配置选项可用，可能包括:

- 特定于数据源的配置
- 导出器配置
- 传播器配置
- 资源配置

自动仪表可用于以下语言:

- [.NET](/docs/instrumentation/net/automatic/)
- [Java](/docs/instrumentation/java/automatic/)
- [JavaScript](/docs/instrumentation/js/automatic/)
- [PHP](/docs/instrumentation/php/automatic/)
- [Python](/docs/instrumentation/python/automatic/)

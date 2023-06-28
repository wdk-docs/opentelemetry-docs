# 术语表

本文档定义了在本规范中使用的一些术语。

其他一些基本术语在[概述文档](overview.md)中有记录。

## 用户角色

### 应用程序所有者

应用程序或服务的维护者，负责配置和管理 OpenTelemetry SDK 的生命周期。

### 库作者

共享库的维护者，许多应用程序都依赖于这个库，OpenTelemetry 仪器也将其作为目标。

### 插装作者

OpenTelemetry 工具的维护者，根据 OpenTelemetry API 编写。这可能是在应用程序代码
、共享库或工具库中编写的工具。

### 插件作者

OpenTelemetry SDK 插件的维护者，针对 OpenTelemetry SDK 插件接口编写。

## 通用

### 信号

OpenTelemetry 是围绕信号或遥测分类构建的。指标、日志、轨迹和行李都是信号的例子。
每个信号代表一个连贯的、独立的功能集。每个信号遵循一个单独的生命周期，定义其当前
的稳定水平。

### 包

在本规范中，术语 **包** 描述了一组代表单个依赖的代码，这些依赖可以独立于其他包导
入到程序中。这个概念可以映射到某些语言中的不同术语，例如“模块”。请注意，在某些语
言中，术语“包”指的是不同的概念。

### ABI 兼容

ABI(应用程序二进制接口)是一种在机器代码级别定义软件组件之间交互的接口，例如在应
用程序可执行文件和共享对象库的编译二进制文件之间。 ABI 兼容性意味着库的新编译版
本可以正确地链接到目标可执行文件，而无需重新编译该可执行文件。

ABI 兼容性对于某些语言非常重要，尤其是那些提供某种形式的机器码的语言。对于其他语
言，ABI 兼容性可能不是相关需求。

### 带内和带外数据

> 在电信中，带内信令是在用于数据(如语音或视频)的同一频带或信道内发送控制信息。这
> 与带外信号形成对比，带外信号通过不同的通道发送，甚至通过单独的网络
> ([Wikipedia](https://en.wikipedia.org/wiki/In-band_signaling))。

在 opentelementetry 中，我们将带内数据称为作为业务消息的一部分在分布式系统组件之
间传递的数据，例如，当跟踪或包以 HTTP 头的形式包含在 HTTP 请求中时。这些数据通常
不包含遥测数据，但用于关联和连接由各个组件产生的遥测数据。遥测本身被称为带外数据
:它通过专用消息从应用程序传输，通常由后台例程异步传输，而不是从业务逻辑的关键路
径传输。导出到遥测后端的度量、日志和跟踪就是带外数据的例子。

### 手动插装

针对 OpenTelemetry API 编码，如[Tracing API](trace/api.md),
[Metrics API](metrics/api.md)，或其他从最终用户代码或共享框架(例如 MongoDB,
Redis 等)收集遥测数据。

### 自动插装

指不需要最终用户修改应用程序源代码的遥测收集方法。方法因编程语言而异，示例包括代
码操作(在编译期间或运行时)、猴子补丁或运行 eBPF 程序。

Synonym: _Auto-instrumentation_.

### 遥测 SDK

表示实现 _OpenTelemetry API_ 的库。

见[库指引](library-guidelines.md#sdk-implementation)
和[库资源语义约定](resource/semantic_conventions/README.md#telemetry-sdk).

### 构造函数

构造函数是应用程序所有者用来初始化和配置 OpenTelemetry SDK 和贡献包的公共代码。
构造函数的例子包括配置对象、环境变量和构造函数。

### SDK 插件

插件是扩展 OpenTelemetry SDK 的库。插件接口的例子有`SpanProcessor`, `Exporter`,
和 `Sampler` 接口。

### 导出库

导出器是 SDK 插件，它实现了“导出器”接口，并向消费者发送遥测信息。

### Instrumented Library

表示为其收集遥测信号(跟踪、度量、日志)的库。

对 OpenTelemetry API 的调用既可以由插装库本身完成，也可以由另一
个[插装库](#instrumentation-library)完成。

Example: `org.mongodb.client`.

### Instrumentation Library

表示为给定的[Instrumented library](#instrumented-library)提供检测的库。如果
_Instrumented Library_ 和 _Instrumentation Library_ 有内置的 OpenTelemetry
instrumentation，那么它们可能是同一个库。

请参阅[概述](overview.md#instrumentation-libraries)了解更详细的定义和命名指南。

Example: `io.opentelemetry.contrib.mongodb`.

Synonyms: _Instrumenting Library_.

### Instrumentation Scope

应用程序代码的逻辑单元，发出的遥测信息可以与它相关联。通常由开发人员决定什么表示
合理的检测范围。最常见的方法是使
用[instrumentation library](#instrumentation-library)作为范围，但是其他范围也很
常见，例如，可以选择一个模块、一个包或一个类作为检测范围。

如果代码单元有版本，则插装范围由(name,version)对定义，否则省略版本，只使用名称。
名称或(名称、版本)对唯一地标识发出遥测的代码的逻辑单元。确保唯一性的典型方法是使
用发出代码的完全限定名(例如，完全限定库名或完全限定类名)。

仪表范围用于获取[示踪剂或仪表](#tracer-name--meter-name).

检测范围可以有零个或多个附加属性，这些属性提供有关该范围的附加信息。例如，对于指
定工具库的作用域，可以记录一个附加属性来表示库的 URL。库的源代码存储在存储库的
URL 中。由于范围是构建时的概念，因此范围的属性不能在运行时更改。

### 示踪器名称/计算器名称

This refers to the `name` and (optional) `version` arguments specified when
creating a new `Tracer` or `Meter` (see
[Obtaining a Tracer](trace/api.md#tracerprovider)/[Obtaining a Meter](metrics/api.md#meterprovider)).
The name/version pair identifies the
[Instrumentation Scope](#instrumentation-scope), for example the
[Instrumentation Library](#instrumentation-library) or another unit of
application in the scope of which the telemetry is emitted.

这指的是 `name` 和(可选)在创建新的“跟踪器”或“计算器”时指定的“版本”参数(参
见[获取跟踪器](trace/api.md#tracerprovider)/[获取仪表](metrics/api.md#meterprovider))。
名称/版本对标识了[Instrumentation Scope](#instrumentation-scope)，例
如[Instrumentation Library](#instrumentation-library)或遥测发射范围内的另一个应
用单元。

### 执行单元

顺序代码执行的最小单元的总称，用于多任务的不同概念。例如线程、协程或纤维。

## 日志

### 日志记录

记录:事件的记录通常，记录包括一个时间戳，指示事件发生的时间，以及描述发生了什么
、发生在哪里等的其他数据。

Synonyms: _Log Entry_.

### 日志

有时用于指日志记录的集合。可能会有歧义，因为人们有时也会用“Log”来指代单一的“Log
Record”，因此这个术语应该谨慎使用，在可能产生歧义的上下文中，应该使用额外的限定
词(例如:“日志记录”)。

### 内嵌日志

“日志记录”嵌入在[Span](trace/api.md# Span)对象中，
在[Events](trace/api.md#add-events)列表中。

### 标准日志

“日志记录”没有嵌入到“Span”中，而是记录在其他地方。

### 日志属性

“日志记录”中包含的键/值对。

### 结构化的日志

日志记录的格式具有良好定义的结构，允许区分日志记录的不同元素(例如时间戳，属性等
)。例如，_Syslog 协议_ ([RFC 5424](https://tools.ietf.org/html/rfc5424))定义了“
结构化数据”格式。

### 平面文件日志

记录在文本文件中的日志，通常每条日志记录一行(尽管也可能有多行记录)。以更结构化的
格式(例如 JSON 文件)写入文本文件的日志是否被认为是平面文件日志，目前还没有通用的
行业协议。如果这样的区别是重要的，建议特别指出来。

### 日志附加/桥接

日志附加器或桥接器是一个组件，它使用[Log Bridge API](./logs/bridge-api.md)将日志
从现有的日志 API 桥接到 OpenTelemetry 中。术语“日志桥接器”和“日志 appender”可以
互换使用，这反映了这些组件将数据桥接到 OpenTelemetry 中，但在日志领域通常称为附
着器。

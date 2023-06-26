---
title: 术语表
description: >-
  您可能熟悉，也可能不熟悉OpenTelemetry项目使用的术语。
weight: 100
---

OpenTelemetry 项目使用的术语您可能不熟悉，也可能不熟悉。此外，该项目可能以不同于
其他项目的方式定义术语。本页包含项目中使用的术语及其含义。

## 通用术语

### **Aggregation**

在程序执行期间的一段时间内，将多个测量组合成有关测量的精确或估计统计信息的过程。
由[`Metric`](#metric) [`Data Source`](#data-source)使用。

### **API**

应用程序编程接口。在 OpenTelemetry 项目中，用于定义如何根
据[`Data Source`](#data-source)生成遥测数据。

### **Application**

为最终用户或其他应用程序设计的一个或多个[服务](#service)。

### **APM**

应用程序性能监视是关于监视软件应用程序及其性能(速度、可靠性、可用性等)，以检测问
题、警报和工具，以找到根本原因。

### **Attribute**

一个键值对。用于遥测信号-例如在[`Traces`](#trace)中将数据附加到[`Span`](#span)，
或在[`Metrics`](#metric)中。参见[属性规范][attribute]。

### **Automatic Instrumentation**

指不需要最终用户修改应用程序源代码的遥测收集方法。方法因编程语言而异，例如字节码
注入或猴子补丁。

### **Baggage**

一种传播名称/值对的机制，以帮助在事件和服务之间建立因果关系。参见[行李规
格][baggage]。

### **Client Library**

查看 [`插装库`](#instrumented-library).

### **Client-side App**

一个[应用程序](#application)的组件，它不在私有基础设施中运行，通常由最终用户直接
使用。客户端应用的例子有浏览器应用、移动应用和运行在物联网设备上的应用。

### **Collector**

关于如何接收、处理和导出遥测数据的与供应商无关的实现。可以作为代理或网关部署的单
个二进制文件。

也称为 OpenTelemetry 收集器。更多关于收集器的信息[在这里][收集器]。

### **Contrib**

几个[插装库](#instrumentation-library)和[收集器](#collector)提供了一组核心功能，
以及一个专用的贡献库，用于非核心功能，包括供应商的“出口器”。

### **Context Propagation**

允许所有[`Data Sources`](#data-source)共享一个底层上下文机制，用于
在[`Transaction`](#transaction)的生命周期内存储状态和访问数据。参见[上下文传播规
范][上下文传播]。

### **DAG**

[有向无环图][dag].

### **Data Source**

查看 [`Signal`](#signal)

### **Dimension**

查看 [`Label`](#label).

### **Distributed Tracing**

跟踪单个[`Request`](#request)的进程，称为跟踪，因为它是
由[`Services`](#service)处理的，组
成[`Application`](#application)。[`Distributed Trace`](#distributed-tracing)跨越
进程、网络和安全边界。

更多关于分布式跟踪的信息[在这里][distributed tracing].

### **Event**

在这种情况下，表示依赖于[`Data Source`](#data-source)。例如，[`Spans`](#span)。

### **Exporter**

提供向消费者发送遥测信息的功能。
由[`插装库`][spec-export-lib]和[`Collector`](/docs/collector/configuration#basics)使
用。导出器可以是 push-，也可以是 pull-based。

### **Field**

添加到[`Log Records`](#log-record)的名称/值对(类似
于[`Spans`](#span)的[`Attributes`](#attribute)和[`Metrics`](#metric)的[`Labels`](#label))。
参见[field spec][field]。

### **gRPC**

一个高性能、开源的通用[`RPC`](#rpc)框架。更多关于 gRPC 的信
息[在这里](https://grpc.io)。

### **HTTP**

[超文本传输协议][http]的简写。

### **Instrumented Library**

表示收集遥测信号([`Traces`](#trace), [`Metrics`](#metric), [`Logs`](#log))
的[`Library`](#library)。 [更多][spec-instrumented-lib]。

### **Instrumentation Library**

表示为给定的[`Instrumented Library`](#instrumented-library)提供检测
的[`Library`](#library)。
[`Instrumented Library`](#instrumented-library)和[`Instrumentation Library`](#instrumentation-library)可
能是相同的[`Library`](#library)，如果它有内置的 OpenTelemetry 仪器。 [更
多][spec-instrumentation-lib]。

### **JSON**

[JavaScript 对象表示法][json]的简写.

### **Label**

查看 [Attribute](#attribute).

### **Language**

编程语言。

### **Library**

由接口调用的特定于语言的行为集合。

### **Log**

有时用于指['日志记录'](#log-record)的集合。可能会有歧义，因为人们有时也会
用[`Log`](#log)来指代单个的[`Log Record`](#log-record)，因此这个术语应该谨慎使用
，在可能产生歧义的上下文中，应该使用额外的限定词(例如:“日志记录”)。查看[更
多][log]

### **Log Record**

['事件'](#e 的记录。通常，记录包括一个时间戳，表明[`Event`](#event)发生的时间，
以及描述发生了什么，发生在哪里等其他数据。查看[更多][log record]

### **Metadata**

A name/value pair added to telemetry data. OpenTelemetry calls this
[`Attributes`](#attribute) on [`Spans`](#span), [`Labels`](#label) on
[`Metrics`](#metric) and [`Fields`](#field) on [`Logs`](#log).

添加到遥测数据中的名称/值对。 OpenTelemetry 在[' span '](#span)上调
用[' Attributes '](#attribute)，在[' Metrics '](#metric)上调
用[' Labels '](#label)，在[' Logs '](#log)上调用[' Fields '](#field)。

### **Metric**

记录一个数据点，无论是原始测量或预定义的聚合，作为时间序列与['元数据'](#
Metadata)。 查看[更多][metric].

### **OC**

[' OpenCensus '](# OpenCensus)的缩写形式。

### **OpenCensus**

一组针对各种语言的库，允许您收集应用程序指标和分布式跟踪，然后将数据实时传输到您
选择的后端。 [OpenTelemetry 的前身](/docs/what-is-opentelemetry/#so-what)。 [更
多][opencensus]。

### **OpenTracing**

用于分布式跟踪的与供应商无关的 api 和工具。
[OpenTelemetry 的前身](/docs/what-is-opentelemetry/#so-what)。[更
多][opentracing]。

### **OT**

[`OpenTracing`](#opentracing)的简写.

### **OTel**

[OpenTelemetry](/docs/what-is-opentelemetry/)的简写.

### **OTelCol**

[OpenTelemetry Collector](#collector)的简写.

### **OTLP**

[OpenTelemetry Protocol](/docs/specs/otlp/)的简写.

### **Processor**

从接收数据到导出数据之间的操作。例如，批处理。
由['Instrumentation Libraries'](#instrumentation-library)和[Collector](/docs/collector/configuration/#processors)使
用。

### **Propagators**

用于序列化和反序列化遥测数据的特定部分，例如[`Spans`](#span)中的 span 上下文
和[`Baggage`](#baggage). 查看[更多][propagators].

### **Proto**

与语言无关的接口类型。 查看[更多][proto].

### **Receiver**

[`Collector`](/docs/collector/configuration/#receivers)用来定义如何接收遥测数据
的术语。接收器可以是推或拉为基础的。看到[更多][receiver]。

### **Request**

查看 [`Distributed Tracing`](#distributed-tracing).

### **Resource**

捕获有关记录遥测的实体的信息。例如，在 Kubernetes 上的容器中运行的产生遥测的进程
有一个 Pod 名称，它在一个命名空间中，可能是部署的一部分，也有一个名称。所有这三
个属性都可以包含在`Resource`中，并应用于任何数据源。

### **REST**

[Representational State Transfer][rest]的简写.

### **RPC**

[Remote Procedure Call][rpc]的简写.

### **Sampling**

控制导出数据量的机制。最常与[`Tracing`](#trace) [`Data Source`](#data-source)一
起使用. 查看[更多][sampling].

### **SDK**

软件开发工具包的简称。指遥测 SDK，表示实现 OpenTelemetry
[`API`](#api)的[`Library`](#library)

### **Semantic Conventions**

定义[' Metadata '](# Metadata)的标准名称和值，以便提供与供应商无关的遥测数据。

### **Service**

[`Application`](#application)的组件。一个[`Service`](#service)的多个实例通常是为
了高可用性和可扩展性而部署的。一个[`Service`](#service)可以部署在多个位置。

### **Signal**

[`Traces`](#trace), [`Metrics`](#metric) or [`Logs`](#log)之一。更多关于信号[在
这里][signals]。

### **Span**

表示[`Trace`](#trace)中的单个操作。查看[更多][span].

### **Span Link**

跨度链接是因果相关的跨度之间的链接。详情请参
见[跨间链接](/docs/specs/otel/overview#links-between-spans) 和
[指定链接](/docs/specs/otel/trace/api#specifying-links).。

### **Specification**

描述所有实现的跨语言需求和期望。查看[更多][specification].

### **Status**

操作的结果。通常用于指示是否发生错误。 查看[更多][status].

### **Tag**

查看 [`Metadata`](#metadata).

### **Trace**

[`Spans`](#span)的[`DAG`](#dag) ，其中[`Spans`](#span)之间的边定义为父/子关系。
查看[更多][trace].

### **Tracer**

负责创建[' span '](#span). 查看[更多][tracer].

### **Transaction**

查看 [`Distributed Tracing`](#distributed-tracing).

### **zPages**

在进程内替代外部导出程序。当包含时，它们在后台收集和汇总跟踪和度量信息;当被请求
时，这些数据被提供给网页。 查看[更多][zpages].

## 额外术语

### Traces

#### **[Trace API Terminology](/docs/specs/otel/trace/api)**

#### **[Trace SDK Terminology](/docs/specs/otel/trace/sdk)**

### Metrics

#### **[Metric API Terminology](/docs/specs/otel/metrics/api#overview)**

#### **[Metric SDK Terminology](/docs/specs/otel/metrics#specifications)**

### Logs

#### **[Trace Context Fields](/docs/specs/otel/logs/data-model#trace-context-fields)**

#### **[Severity Fields](/docs/specs/otel/logs/data-model#severity-fields)**

#### **[Log Record Fields](/docs/specs/otel/logs/data-model#log-and-event-record-definition)**

### Semantic Conventions

#### **[Resource Conventions](/docs/specs/otel/resource/semantic_conventions)**

#### **[Span Conventions](/docs/specs/otel/trace/semantic_conventions)**

#### **[Metric Conventions](/docs/specs/otel/metrics/semantic_conventions)**

[baggage]: /docs/specs/otel/baggage/api/
[attribute]: /docs/specs/otel/common/#attributes
[collector]: /docs/collector
[context propagation]: /docs/specs/otel/overview#context-propagation
[dag]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
[distributed tracing]: /docs/concepts/signals/traces/
[field]: /docs/specs/otel/logs/data-model#field-kinds
[http]: https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol
[json]: https://en.wikipedia.org/wiki/JSON
[log]: /docs/specs/otel/glossary#log
[log record]: /docs/specs/otel/glossary#log-record
[metric]: /docs/concepts/signals/metrics/
[opencensus]: https://opencensus.io
[opentracing]: https://opentracing.io
[propagators]: /docs/instrumentation/go/manual/#propagators-and-context
[proto]: https://github.com/open-telemetry/opentelemetry-proto
[receiver]: /docs/collector/configuration/#receivers
[rest]: https://en.wikipedia.org/wiki/Representational_state_transfer
[rpc]: https://en.wikipedia.org/wiki/Remote_procedure_call
[sampling]: /docs/specs/otel/trace/sdk#sampling
[signals]: /docs/concepts/signals/
[span]: /docs/specs/otel/trace/api#span
[spans]: /docs/specs/otel/trace/api#add-events
[spec-exporter-lib]: /docs/specs/otel/glossary/#exporter-library
[spec-instrumentation-lib]: /docs/specs/otel/glossary/#instrumentation-library
[spec-instrumented-lib]: /docs/specs/otel/glossary/#instrumented-library
[specification]: /docs/concepts/components/#specification
[status]: /docs/specs/otel/trace/api#set-status
[trace]: /docs/specs/otel/overview#traces
[tracer]: /docs/specs/otel/trace/api#tracer
[zpages]:
  https://github.com/open-telemetry/opentelemetry-specification/blob/main/experimental/trace/zpages.md

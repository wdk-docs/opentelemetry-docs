---
title: Traces
weight: 1
---

[**Traces**](/docs/concepts/observability-primer/#distributed-traces)向我们展示了当向应用程序发出请求时发生的情况。
无论您的应用程序是具有单个数据库的单体还是复杂的服务网格，跟踪对于理解请求在应用程序中的完整“路径”都是必不可少的。

考虑以下跟踪三个工作单元的示例:

```json
{
    "name": "Hello-Greetings",
    "context": {
        "trace_id": "0x5b8aa5a2d2c872e8321cf37308d69df2",
        "span_id": "0x5fb397be34d26b51",
    },
    "parent_id": "0x051581bf3cb55c13",
    "start_time": "2022-04-29T18:52:58.114304Z",
    "end_time": "2022-04-29T22:52:58.114561Z",
    "attributes": {
        "http.route": "some_route1"
    },
    "events": [
        {
            "name": "hey there!",
            "timestamp": "2022-04-29T18:52:58.114561Z",
            "attributes": {
                "event_attributes": 1
            }
        },
        {
            "name": "bye now!",
            "timestamp": "2022-04-29T18:52:58.114585Z",
            "attributes": {
                "event_attributes": 1
            }
        }
    ],
}
{
    "name": "Hello-Salutations",
    "context": {
        "trace_id": "0x5b8aa5a2d2c872e8321cf37308d69df2",
        "span_id": "0x93564f51e1abe1c2",
    },
    "parent_id": "0x051581bf3cb55c13",
    "start_time": "2022-04-29T18:52:58.114492Z",
    "end_time": "2022-04-29T18:52:58.114631Z",
    "attributes": {
        "http.route": "some_route2"
    },
    "events": [
        {
            "name": "hey there!",
            "timestamp": "2022-04-29T18:52:58.114561Z",
            "attributes": {
                "event_attributes": 1
            }
        }
    ],
}
{
    "name": "Hello",
    "context": {
        "trace_id": "0x5b8aa5a2d2c872e8321cf37308d69df2",
        "span_id": "0x051581bf3cb55c13",
    },
    "parent_id": null,
    "start_time": "2022-04-29T18:52:58.114201Z",
    "end_time": "2022-04-29T18:52:58.114687Z",
    "attributes": {
        "http.route": "some_route3"
    },
    "events": [
        {
            "name": "Guten Tag!",
            "timestamp": "2022-04-29T18:52:58.114561Z",
            "attributes": {
                "event_attributes": 1
            }
        }
    ],
}
```

这个示例跟踪输出有三个不同的类似日志的项目，称为[span](#span)，分别命名为“Hello-greetings”、“Hello-salutations”和“Hello”。
因为每个请求的上下文都有相同的“trace_id”，所以它们被认为是同一个Trace的一部分。

您将注意到的另一件事是，这个示例Trace的每个Span看起来都像一个结构化日志。
那是因为它确实是!
考虑trace的一种方法是，它们是包含上下文、相关性、层次结构等内容的结构化日志的集合。
但是，这些“结构化日志”可以来自不同的进程、服务、虚拟机、数据中心等。
这使得跟踪能够表示任何系统的端到端视图。

为了理解开放遥测中的跟踪是如何工作的，让我们看一下将在检测代码中发挥作用的组件列表:

- Tracer
- Tracer Provider
- Trace Exporter
- Trace Context

## Tracer Provider

Tracer Provider(有时称为“TracerProvider”)是“Tracer”的工厂。
在大多数应用程序中，跟踪程序提供程序初始化一次，其生命周期与应用程序的生命周期相匹配。
跟踪程序提供程序初始化还包括资源和导出程序初始化。
这通常是使用OpenTelemetry进行跟踪的第一步。
在某些语言sdk中，已经为您初始化了全局跟踪程序提供程序。

## Tracer

跟踪器创建的范围包含有关给定操作(例如服务中的请求)正在发生的事情的更多信息。
跟踪程序是从跟踪程序提供程序创建的。

## Trace Exporters

跟踪出口商将跟踪发送给消费者。
这个消费者可以是调试和开发时间的标准输出、OpenTelemetry Collector，或者您选择的任何开源或供应商后端。

## Context Propagation

上下文传播是支持分布式跟踪的核心概念。
使用上下文传播，无论在哪里生成span，都可以将span相互关联并组装到跟踪中。
我们通过两个子概念定义上下文传播:上下文和传播。

**上下文** 是一个对象，它包含发送和接收服务的信息，以便将一个跨度与另一个跨度相关联，并将其与整个跟踪相关联。
例如，如果服务A调用服务B，那么来自服务A的一个span(其ID在上下文中)将被用作在服务B中创建的下一个span的父span。

**传播** 是在服务和进程之间移动上下文的机制。
通过这样做，它组装了一个分布式跟踪。
它序列化或反序列化Span Context，并提供要从一个服务传播到另一个服务的相关跟踪信息。
我们现在有了我们所说的: **跟踪上下文** 。

上下文是一个抽象的概念——它需要一个具体的实现才能真正有用。
OpenTelemetry支持几种不同的上下文格式。
OpenTelemetry跟踪中使用的默认格式是W3C  `TraceContext`。
每个Context对象都与一个跨度相关联，并且可以在跨度上进行访问。
参见[Span Context](#Span-Context)。

通过结合上下文和传播，您现在可以组装跟踪。

> 有关更多信息，请参阅[跟踪规范][]

[跟踪规范]: /docs/specs/otel/overview/#tracing-signal

## Spans

[**span**](/docs/concepts/observability-primer/#span)表示一个工作或操作单元。
跨度是trace的构建块。
在OpenTelemetry中，它们包括以下信息:

- Name
- Parent span ID (empty for root spans)
- Start and End Timestamps
- [Span Context](#span-context)
- [Attributes](#attributes)
- [Span Events](#span-events)
- [Span Links](#span-links)
- [Span Status](#span-status)

Sample span:

```json
{
  "trace_id": "7bba9f33312b3dbb8b2c2c62bb7abe2d",
  "parent_id": "",
  "span_id": "086e83747d0e381e",
  "name": "/v1/sys/health",
  "start_time": "2021-10-22 16:04:01.209458162 +0000 UTC",
  "end_time": "2021-10-22 16:04:01.209514132 +0000 UTC",
  "status_code": "STATUS_CODE_OK",
  "status_message": "",
  "attributes": {
    "net.transport": "IP.TCP",
    "net.peer.ip": "172.17.0.1",
    "net.peer.port": "51820",
    "net.host.ip": "10.177.2.152",
    "net.host.port": "26040",
    "http.method": "GET",
    "http.target": "/v1/sys/health",
    "http.server_name": "mortar-gateway",
    "http.route": "/v1/sys/health",
    "http.user_agent": "Consul Health Check",
    "http.scheme": "http",
    "http.host": "10.177.2.152:26040",
    "http.flavor": "1.1"
  },
  "events": [
    {
      "name": "",
      "message": "OK",
      "timestamp": "2021-10-22 16:04:01.209512872 +0000 UTC"
    }
  ]
}
```

span可以嵌套，正如父span ID的存在所暗示的那样:子span表示子操作。
这允许span更准确地捕获应用程序中完成的工作。

### Span 上下文

Span context是一个不可变对象，包含以下内容:

- Trace ID表示该span是其中一部分的跟踪
- span的span ID
- 跟踪标志，一种二进制编码，包含有关跟踪的信息
- 跟踪状态，可以携带特定于供应商的跟踪信息的键值对列表

Span上下文是Span的一部分，它与[分布式上下文](#context-propagation)和[包袱](/docs/concepts/signals/baggage)一起被序列化和传播。

因为Span Context包含Trace ID，所以在创建[Span Links](#Span-Links)时使用它。

### 属性

属性是包含元数据的键值对，您可以使用这些元数据对Span进行注释，以携带有关它正在跟踪的操作的信息。

例如，如果span跟踪在电子商务系统中向用户的购物车中添加商品的操作，则可以捕获用户的ID、要添加到购物车中的商品的ID和购物车ID。

每种语言SDK实现的属性有以下规则:

- 键必须为非空字符串值
- 取值必须为非空字符串、布尔值、浮点数、整数或这些值的数组

此外，还有[语义属性](/docs/specs/otel/trace/semantic_conventions/)，它们是已知的元数据命名约定，通常出现在公共操作中。
尽可能使用语义属性命名是很有帮助的，这样可以跨系统标准化常见类型的元数据。

### Span 事件

可以将Span事件看作是Span上的结构化日志消息(或注释)，通常用于表示Span持续期间有意义的单一时间点。

例如，考虑web浏览器中的两种场景:

1. 跟踪页面加载
2. 表示页面何时具有交互性

Span最适合用于第一种场景，因为它是一个有开始和结束的操作。

Span Event最适合用于跟踪第二个场景，因为它代表了一个有意义的、奇异的时间点。

### Span 链接

存在链接，以便您可以将一个跨度与一个或多个跨度关联起来，这意味着存在因果关系。
例如，假设我们有一个分布式系统，其中一些操作由跟踪器跟踪。

为了响应其中的一些操作，一个额外的操作排队等待执行，但是它的执行是异步的。
我们也可以用trace来跟踪这个后续操作。

我们希望将后续操作的跟踪与第一个跟踪相关联，但是我们无法预测后续操作何时开始。
我们需要将这两条轨迹关联起来，因此我们将使用一个跨度链接。

您可以将第一个跟踪中的最后一个跨度链接到第二个跟踪中的第一个跨度。
现在，它们互为因果关系。

链接是可选的，但它是将跟踪跨度相互关联的好方法。

### Span 状态

一个状态将被附加到一个Span上。
通常，当应用程序代码中存在已知错误(例如异常)时，您将设置一个span状态。
Span Status将被标记为以下值之一:

- `Unset`
- `Ok`
- `Error`

处理异常时，可以将Span状态设置为Error。
否则，Span状态为Unset状态。
通过将Span状态设置为Unset，处理Span的后端现在可以分配最终状态。

### Span 类型

当创建一个span时，它是“Client”、“Server”、“Internal”、“Producer”或“Consumer”中的一个。
此跨度类型为跟踪后端提供了关于如何组装跟踪的提示。
根据OpenTelemetry规范，服务器跨度的父节点通常是远程客户端跨度，而客户端跨度的子节点通常是服务器跨度。
类似地，消费者span的父类始终是生产者，生产者span的子类始终是消费者。
如果没有提供，则假定span类型是内部的。

有关SpanKind的更多信息，请参见[SpanKind](/docs/specs/otel/trace/api/#spankind).

#### 客户端

客户端范围表示同步传出远程调用，例如传出HTTP请求或数据库调用。
请注意，在这个上下文中，“同步”不是指“async/await”，而是指它没有排队等待稍后的处理。

#### 服务器

服务器范围表示同步传入的远程调用，例如传入的HTTP请求或远程过程调用。

#### Internal

内部Span表示不跨越流程边界的操作。
插装函数调用或快速中间件之类的事情可能会使用内部Span。

#### Producer

生产者Span表示稍后可能异步处理的作业的创建。
它可以是远程作业，例如插入作业队列的作业，也可以是由事件侦听器处理的本地作业。

#### Consumer

消费者Span表示由生产者创建的作业的处理，并且可能在生产者Span已经结束很久之后才开始。

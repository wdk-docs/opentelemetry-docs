# 概述

本文档提供了 OpenTelemetry 项目的概述，并定义了重要的基本术语。

其他术语定义可在[词汇表](glossary.md)中找到.

## OpenTelemetry 客户端架构

![Cross cutting concerns](../internal/img/architecture.png)

在最高的体系结构级别，OpenTelemetry 客户端被组织
成[**信号**](glossary.md#signals)。每个信号都提供了一种特殊形式的可观测性。例如
，跟踪、度量和包袱是三个独立的信号。信号共享一个共同的子系统——上下文传播——但它们
彼此独立地工作。

每个信号都为软件提供了一种描述自身的机制。代码库，如 web 框架或数据库客户端，依
赖于各种信号来描述自己。然后可以将 OpenTelemetry 检测代码混合到该代码库中的其他
代码中。这使得 OpenTelemetry 成为一
个[**横切关注**](https://en.wikipedia.org/wiki/Cross-cutting_concern)——一个为了
提供价值而混合到许多其他软件中的软件。横切关注点，就其本质而言，违反了一个核心设
计原则——关注点分离。因此，OpenTelemetry 客户端设计需要格外小心，以避免为依赖于这
些横切 api 的代码库创建问题。

OpenTelemetry 客户端被设计成将每个信号中必须作为横切关注点导入的部分与可以独立管
理的部分分开。 OpenTelemetry 客户端也被设计成一个可扩展的框架。为了实现这些目标
，每个信号由四种类型的包组成:API、SDK、Semantic Conventions 和 Contrib。

### API

API 包由用于检测的横切公共接口组成。 OpenTelemetry 客户端中导入第三方库和应用程
序代码的任何部分都被认为是 API 的一部分。

### SDK

SDK 是 OpenTelemetry 项目提供的 API 的实现。在应用程序中，SDK
由[应用程序所有者](glossary.md#application-owner)安装和管理。请注意，SDK 包含额
外的公共接口，这些接口不被认为是 API 包的一部分，因为它们不是横切关注点。这些公
共接口被定义
为[构造函数](glossary.md#constructors)和[插件接口](glossary.md#sdk-plugins)。应
用程序所有者使用 SDK 构造函数;[插件作者](glossary.md#plugin-author)使用 SDK 插件
接口。 [工具作者](glossary.md#instrumentation-author)绝对不能直接引用任何 SDK 包
，只能引用 API。

### 语义约定

语义约定定义了键和值，这些键和值描述了应用程序使用的常见概念、协议和操作。

语义约定现在位于它们自己的存储库中:
https://github.com/open-telemetry/semantic-conventions

收集器和客户端库都应该自动将语义约定键和枚举值生成常量(或语言习惯等效)。在语义约
定稳定之前，生成的值不应该在稳定的包中分发。
[YAML](https://github.com/open-telemetry/semantic-conventions/tree/main/semantic_conventions)文
件必须作为生成的真实源。每种语言实现都应该
为[代码生成器](https://github.com/open-telemetry/build-tools/tree/main/semantic-conventions#code-generator)提
供特定于语言的支持。

此外，将列出规范所需的属性[在这里](semantic-conventions.md)

### 贡献包

OpenTelemetry 项目维护与流行的 OSS 项目的集成，这些项目对于观察现代 web 服务非常
重要。 API 集成的示例包括用于 web 框架、数据库客户端和消息队列的插装。示例 SDK
集成包括用于将遥测数据导出到流行的分析工具和遥测数据存储系统的插件。

一些插件，如 OTLP 导出器和 TraceContext 传播器，是 OpenTelemetry 规范所要求的。
这些必需的插件作为 SDK 的一部分包含。

可选的、独立于 SDK 的插件和工具包被称为 **Contrib** 包。 **API Contrib** 指的是
仅依赖于 API 的包; **SDK Contrib** 指的是同样依赖于 SDK 的包。

`Contrib` 一词特指由 OpenTelemetry 项目维护的插件和工具的集合;它不涉及第三方插件
托管在其他地方。

### 版本控制和稳定性

OpenTelemetry 重视稳定性和向后兼容性。请参
阅[版本控制和稳定性指南](./versioning-and-stability.md)了解详细信息。

## 跟踪信号

分布式跟踪是一组事件，由单个逻辑操作触发，并在应用程序的各个组件之间进行整合。分
布式跟踪包含跨进程、网络和安全边界的事件。当有人按下按钮以启动网站上的操作时，可
能会启动分布式跟踪—在本例中，跟踪将表示下游服务之间的调用，这些服务处理由按下按
钮发起的请求链。

### Traces

OpenTelemetry 中的 **Traces** 由它们的 **Spans** 隐式定义。特别地，一个
**Trace** 可以被认为是 **Spans** 的有向无环图(DAG)，其中 **Spans** 之间的边被定
义为父/子关系。

例如，下面是一个由 6 个 **Spans** 组成的 **Trace** 示例:

```
Causal relationships between Spans in a single Trace

        [Span A]  ←←←(the root span)
            |
     +------+------+
     |             |
 [Span B]      [Span C] ←←←(Span C is a `child` of Span A)
     |             |
 [Span D]      +---+-------+
               |           |
           [Span E]    [Span F]
```

有时，用时间轴来可视化 **Traces** 更容易，如下图所示:

```
Temporal relationships between Spans in a single Trace

––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–> time

 [Span A···················································]
   [Span B··········································]
      [Span D······································]
    [Span C····················································]
         [Span E·······]        [Span F··]
```

### Spans

一个 span 表示事务中的一个操作。每个 **Span** 封装了以下状态:

- 操作名称
- 开始和结束时间戳
- [**Attributes**](./common/README.md#attribute): 键值对列表。
- 一组零个或多个**事件**，每个事件本身是一个元组(时间戳，名称
  ，[**Attributes**](./common/README.md#attribute))。名称必须是字符串。
- 父节点的 **Span** 标识符。
- [**链接**](#links-between-spans)到零个或多个因果相关的 **Spans**(通过这些相关
  **Spans**的**SpanContext**)。
- 引用 Span 所需的 **SpanContext** 信息。见下文。

### SpanContext

表示在 **Trace** 中标识 **Span** 的所有信息，并且必须传播到子 Spans 和跨进程边界
。 **SpanContext** 包含跟踪标识符和从父 **Spans** 传播到子 **Spans** 的选项。

- **TraceId** 是跟踪的标识符。它是世界范围内独一无二的，几乎有足够的概率由 16 个
  随机生成的字节组成。TraceId 用于在所有进程中将特定跟踪的所有范围组合在一起。
- **SpanId** 是 span 的标识符。它是全局唯一的，实际上有足够的概率由 8 个随机生成
  的字节组成。当传递给子 Span 时，此标识符将成为子 **Span** 的父 span id。
- **TraceFlags** 表示跟踪的选项。它被表示为 1 字节(位图)。
  - 采样位-表示跟踪是否采样的位(掩码 `0x1`)。
- **Tracestate** 在键值对列表中携带跟踪系统特定的上下文。**Tracestate** 允许不同
  的供应商传播额外的信息并与他们的遗留 Id 格式进行互操作。欲了解更多详情，请参
  阅[本](https://w3c.github.io/trace-context/#tracestate-field).

### Spans 间链接

一个 **Span** 可以链接到零个或多个因果相关的其他 **Span** (由 **SpanContext** 定
义)。 **链接** 可以指向 **span** 内的单个 **Trace** 或跨不同的 **Trace** 。 **链
接** 可用于表示批处理操作，其中一个 **Span** 由多个初始化 **Span** 发起，每个
**Span** 表示在批处理中正在处理的单个传入项。

使用 **Link** 的另一个例子是声明起始跟踪和后续跟踪之间的关系。当 Trace 进入服务
的可信边界并且服务策略需要生成新的 Trace 而不是信任传入的 Trace 上下文时，可以使
用此方法。新的链接 Trace 还可以表示由众多快速传入请求之一发起的长时间运行的异步
数据处理操作。

当使用分散/收集(也称为 fork/join)模式时，根操作启动多个下游处理操作，并且所有这
些操作都聚合回单个 Span 中。最后一个 **Span** 链接到它聚合的许多操作。它们都是来
自同一个 Trace 的 **span** 。类似于 **Span** 的 Parent 字段。但是，建议不要在这
个场景中设置 **Span** 的 parent，因为在语义上，parent 字段代表一个单一的父场景，
在许多情况下，父 **Span** 完全包含子 **Span** 。在分散/收集和批处理场景中，情况
并非如此。

## 度量信号

OpenTelemetry 允许使用预定义的聚合和[一组属性](./common/README.md#attribute)记录
原始测量或度量.

使用 OpenTelemetry API 记录原始测量，允许最终用户决定应该为该度量应用哪种聚合算
法以及定义属性(维度)。它将在 gRPC 等客户端库中用于记录“server_latency”或
“received_bytes”的原始测量值。因此，最终用户将决定从这些原始测量中收集哪种类型的
聚合值。它可以是简单的平均或详细的直方图计算。

使用 OpenTelemetry API 用预定义的聚合记录度量也同样重要。它允许收集 cpu 和内存使
用等值，或者像“队列长度”这样的简单指标。

### 记录原始测量值

用于记录原始测量的主要类是“Measure”和“Measurement”。可以使用 OpenTelemetry API
记录“测量”列表以及附加上下文。因此，用户可以定义汇总这些“度量”，并使用传递的上下
文来定义结果度量的附加属性。

#### Measure

“Measure”描述了库记录的单个值的类型。它定义了公开度量的库和应用程序之间的契约，
应用程序将这些单独的度量聚合到一个“Measure”中。 “Measure”由名称、描述和一个值单
位来标识。

#### Measurement

“Measurement”描述要为“Measure”收集的单个值。 “Measurement”是 API 界面中的一个空
接口。该接口在 SDK 中定义。

### 使用预定义的聚合记录度量标准

所有类型的预聚合指标的基类称为“Metric”。它定义了基本的度量属性，如名称和属性。从
“Metric”继承的类定义了它们的聚合类型以及单个度量或点的结构。 API 定义了以下类型
的预聚合指标:

- 计数器报告瞬时测量。计数器值可以上升或保持不变，但永远不会下降。计数器值不能为
  负值。有两种类型的反度量值——`double` 和 `long`。
- 测量公制报告数值的瞬时测量值。仪表可以上下移动。仪表值可以是负的。有两种类型的
  测量公制值-`double` 和 `long`。

API 允许构造所选类型的`Metric`。 SDK 定义了查询要导出的`Metric`当前值的方式。

每种类型的`Metric`都有它的 API 来记录要聚合的值。 API 支持推拉两种模式设
置`Metric`值。

### 度量数据模型和 SDK

Metrics 数据模型[在这里指定](metrics/data-model.md)，并基
于[metrics.proto](https://github.com/open-telemetry/opentelemetry-proto/blob/master/opentelemetry/proto/metrics/v1/metrics.proto)。
该数据模型定义了三种语义:API 使用的 Event 模型、SDK 和 OTLP 使用的动态数据模型，
以及表示导出器应如何解释动态模型的 TimeSeries 模型。

不同的导出器有不同的功能(例如支持哪些数据类型)和不同的约束(例如属性键中允许哪些
字符)。指标旨在成为可能性的超集，而不是任何地方都支持的最低公分母。所有导出器都
通过 OpenTelemetry SDK 中定义的 Metric Producer 接口从 Metrics data Model 中使用
数据。

因此，Metrics 对数据的限制最小(例如，键中允许哪些字符)，处理 Metrics 的代码应该
避免验证和清理 Metrics 数据。相反，将数据传递给后端，依靠后端执行验证，并从后端
传回任何错误。

有关更多信息，请参阅[Metrics 数据模型规范](metrics/data-model.md)。

## 日志信号

### 数据模型

[日志数据模型](logs/data-model.md)定义了 OpenTelemetry 如何理解日志和事件。

## 包袱信号

除了跟踪传播之外，OpenTelemetry 还提供了一种简单的机制来传播名称/值对，称为“包袱
”。 “包袱”用于索引同一事务中具有先前服务提供的属性的服务中的可观察性事件。这有助
于在这些事件之间建立因果关系。

虽然“包袱”可以用于其他横切关注点的原型，但该机制主要是为了传达 OpenTelemetry 可
观察性系统的值。

这些值可以从“包袱”中使用，并用作度量的附加属性，或者用于日志和跟踪的附加上下文。
一些例子:

- web 服务可以从包含有关发送请求的服务的上下文中获益
- SaaS 提供程序可以包含有关负责该请求的 API 用户或令牌的上下文
- 确定特定的浏览器版本与图像处理服务中的故障相关联

为了与 OpenTracing 向后兼容，在使用 OpenTracing 桥时，包袱被传播为“包袱”。使用不
同标准的新关注点应该考虑创建一个新的横切关注点来覆盖它们的用例; 它们可能受益于
W3C 编码格式，但使用新的 HTTP 头在整个分布式跟踪中传输数据。

## 资源

“资源”捕获有关遥测记录的实体的信息。例如，Kubernetes 容器公开的指标可以链接到指
定集群、命名空间、pod 和容器名称的资源。

“资源”可以捕获实体标识的整个层次结构。它可以描述云中的主机和特定容器或进程中运行
的应用程序。

注意，一些进程识别信息可以通过 OpenTelemetry SDK 自动与遥测相关联。

## 上下文传播

OpenTelemetry 的所有横切关注点，如跟踪和度量，都共享一个底层的“上下文”机制，用于
在分布式事务的整个生命周期中存储状态和访问数据。

参见[背景信息](context/README.md)

## 传播器

OpenTelemetry 使用`Propagators`来序列化和反序列化横切关注点值，比如`Span`(通常只
有`SpanContext`部分)和`Baggage`。不同的“传播器”类型定义了特定传输所施加的限制，
并绑定到数据类型。

传播器 API 目前定义了一种“传播器”类型:

- `TextMapPropagator` 向载波中注入值并从载波中提取值作为文本。

## 收集器

OpenTelemetry 收集器是一组组件，可以从 OpenTelemetry 或其他监控/跟踪库(Jaeger,
Prometheus 等)测量的进程中收集跟踪、指标和最终的其他遥测数据(例如日志)，进行聚合
和智能采样，并将跟踪和指标导出到一个或多个监控/跟踪后端。收集器将允许丰富和转换
收集到的遥测数据(例如添加额外的属性或删除个人信息)。

OpenTelemetry 收集器有两种主要的操作模式: Agent(与应用程序一起在本地运行的守护进
程)和 Collector(独立运行的服务)。

在 OpenTelemetry Service 阅读更多内
容[长期愿景](https://github.com/open-telemetry/opentelemetry-collector/blob/master/docs/vision.md).

## 插装库

参见[插装库](glossary.md#instrumentation-library)

该项目的灵感是通过让每个库和应用程序直接调用 OpenTelemetry API，使它们成为开箱即
用的可观察对象。然而，许多库没有这样的集成，因此需要一个单独的库来注入这样的调用
，使用诸如包装接口、订阅特定于库的回调或将现有的遥测转换为 OpenTelemetry 模型等
机制。

使 OpenTelemetry 对另一个库具有可观察性的库称
为[插装库](glossary.md#instrumentation-library).

插装库的命名应该遵循插装库的任何命名约定(例如: 'middleware'用于 web 框架)。

如果没有确定的名称，建议使用"opentelemetry-instrumentation"作为包的前缀，后面跟
着被检测的库名称本身。例子包括:

- opentelemetry-instrumentation-flask (Python)
- @opentelemetry/instrumentation-grpc (Javascript)

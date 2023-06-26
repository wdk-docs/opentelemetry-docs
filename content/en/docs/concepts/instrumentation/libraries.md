---
title: 插装库
description: Learn how to add native instrumentation to your library.
aliases: [/docs/concepts/instrumenting-library]
weight: 40
---

OpenTelemetry 为许多库提供了[插装库][]，这通常是通过库钩子或猴子补丁库代码完成的
。

使用 OpenTelemetry 的本机库插装为用户提供了更好的可观察性和开发体验，消除了库暴
露和文档挂钩的需要:

- 自定义日志钩子可以被常见的和易于使用的 OpenTelemetry api 取代，用户将只与
  OpenTelemetry 交互
- 来自库和应用程序代码的跟踪、日志、指标是相关和一致的
- 通用约定允许用户在相同的技术和跨库和语言中获得相似和一致的遥测
- 遥测信号可以使用各种记录良好的 OpenTelemetry 扩展点对各种消费场景进行微调(过滤
  、处理、聚合)。

## 语义约定

查看可用的[语义约定](../../specs/otel/trace/semantic_conventions.md)，涵盖 web
框架、RPC 客户端、数据库、消息传递客户端、基础设施等!

如果您的库是其中之一-遵循惯例，它们是事实的主要来源，并告诉哪些信息应该包含在
spans 中。约定使检测保持一致:使用遥测技术的用户不必学习库的细节，而可观察性供应
商可以为各种各样的技术(例如数据库或消息传递系统)构建体验。当库遵循约定时，无需用
户输入或配置，许多场景就可以开箱即用。

如果您有任何反馈或想要添加一个新的会议-请来贡献!
[Instrumentation Slack](https://cloud-native.slack.com/archives/C01QZFGMLQ7)或[Specification repo](https://github.com/open-telemetry/opentelemetry-specification)是
一个很好的开始!

## 当 **不** 仪器

有些库是包装网络调用的瘦客户机。 OpenTelemetry 很可能有一个用于底层 RPC 客户端的
工具库(查看[registry](/ecosystem/registry/))。在这种情况下，可能没有必要检测包装
器库。

如果:

- 您的库是文档化或自解释 api 之上的瘦代理
- _和_ OpenTelemetry 有用于底层网络调用的工具
- _和_ 您的库不应该遵循任何惯例来丰富遥测技术

如果你有疑问-不要仪器-你可以在你看到需要的时候再做。

如果您选择不进行检测，那么提供一种方法为您的内部 RPC 客户端实例配置
OpenTelemetry 处理程序可能仍然是有用的。它在不支持全自动插装的语言中是必不可少的
，但在其他语言中仍然很有用。

如果您决定这样做，本文的其余部分将指导您使用什么以及如何使用。

## OpenTelemetry API

第一步是依赖于 OpenTelemetry API 包。

OpenTelemetry 有[两个主要模块](/docs/specs/otel/overview/)——API 和 SDK。
OpenTelemetry API 是一组抽象和非操作实现。除非您的应用程序导入 OpenTelemetry
SDK，否则您的检测工具不会做任何事情，也不会影响应用程序的性能。

**库应该只使用 OpenTelemetry API。**

你可能有理由担心添加新的依赖，这里有一些注意事项可以帮助你决定如何减少依赖地狱:

- OpenTelemetry Trace API 在 2021 年初达到稳定，它遵循[语义版本控制
  2.0](/docs/specs/otel/version -and-stability)和我们认真对待 API 稳定性。
- 当使用依赖时，请使用最早的稳定 OpenTelemetry API(1.0.\*)并避免更新它，除非您必
  须使用新功能。
- 当您的工具稳定下来时，请考虑将其作为一个单独的包发布，这样就不会给不使用它的用
  户带来问题。您可以将其保留在您的 repo 中，或
  者[将其添加到 OpenTelemetry](https://github.com/open-telemetry/oteps/blob/main/text/0155-external-modules.md#contrib-components)，
  这样它将与其他仪器包一起发布。
- 语义约定是[稳定的，但受制于演变][]:虽然这不会导致任何功能问题，但您可能需要每
  隔一段时间更新您的工具。将其放在预览插件或 opentelement_contrib_repo 中可能有
  助于保持惯例的最新，而不会破坏用户的更改。

[稳定的，但受制于演变]:
  ../../specs/otel/versioning-and-stability.md#semantic-conventions-stability

### 获取追踪器

所有应用程序配置都通过 Tracer API 对库隐藏。默认情况下，库应该
从[global `TracerProvider`](/docs/specs/otel/trace/api/#get-a-tracer)获取跟踪器
。

```java
private static final Tracer tracer = GlobalOpenTelemetry.getTracer("demo-db-client", "0.1.0-beta1");
```

对于库来说，有一个允许应用程序显式传递`TracerProvider`实例的 API 是很有用的，这
样可以更好地实现依赖注入并简化测试。

在获得跟踪程序时，提供您的库(或跟踪插件)名称和版本——它们显示在遥测数据上，帮助用
户处理和过滤遥测数据，了解它的来源，并调试/报告任何仪表问题。

## 仪器仪表

### 公共 api

公共 API 是很好的跟踪对象:为公共 API 调用创建的范围允许用户将遥测映射到应用程序
代码，了解库调用的持续时间和结果。调用 trace:

- 内部进行网络调用的公共方法或花费大量时间且可能失败的本地操作(例如 IO)
- 处理请求或消息的处理程序

**插装的例子:**

```java
private static final Tracer tracer = GlobalOpenTelemetry.getTracer("demo-db-client", "0.1.0-beta1");

private Response selectWithTracing(Query query) {
    // check out conventions for guidance on span names and attributes
    Span span = tracer.spanBuilder(String.format("SELECT %s.%s", dbName, collectionName))
            .setSpanKind(SpanKind.CLIENT)
            .setAttribute("db.name", dbName)
            ...
            .startSpan();

    // makes span active and allows correlating logs and nest spans
    try (Scope unused = span.makeCurrent()) {
        Response response = query.runWithRetries();
        if (response.isSuccessful()) {
            span.setStatus(StatusCode.OK);
        }

        if (span.isRecording()) {
           // populate response attributes for response codes and other information
        }
    } catch (Exception e) {
        span.recordException(e);
        span.setStatus(StatusCode.ERROR, e.getClass().getSimpleName());
        throw e;
    } finally {
        span.end();
    }
}
```

按照约定填充属性!如没有适用的规定，请参
阅[一般惯例](/docs/specs/otel/trace/semantic_conventions/span-general/).

### 嵌套网络和其他 spans

Network calls are usually traced with OpenTelemetry auto-instrumentations
through corresponding client implementation.

![Nested database and HTTP spans in Jaeger UI](../nested-spans.svg)

If OpenTelemetry does not support tracing your network client, use your best
judgement, here are some considerations to help:

- Would tracing network calls improve observability for users or your ability to
  support them?
- Is your library a wrapper on top of public, documented RPC API? Would users
  need to get support from the underlying service in case of issues?
  - instrument the library and make sure to trace individual network tries
- Would tracing those calls with spans be very verbose? or would it noticeably
  impact performance?
  - use logs with verbosity or span events: logs can be correlated to parent
    (public API calls), while span events should be set on public API span.
  - if they have to be spans (to carry and propagate unique trace context), put
    them behind a configuration option and disable them by default.

If OpenTelemetry already supports tracing your network calls, you probably don't
want to duplicate it. There may be some exceptions:

- to support users without auto-instrumentation (which may not work in certain
  environments or users may have concerns with monkey-patching)
- to enable custom (legacy) correlation and context propagation protocols with
  underlying service
- enrich RPC spans with absolutely essential library/service-specific
  information not covered by auto-instrumentation

WARNING: Generic solution to avoid duplication is under construction 🚧.

### Events

Traces are one kind of signal that your apps can emit. Events (or logs) and
traces complement, not duplicate, each other. Whenever you have something that
should have a verbosity, logs are a better choice than traces.

Chances are that your app uses logging or some similar module already. Your
module might already have OpenTelemetry integration -- to find out, see the
[registry](/ecosystem/registry/). Integrations usually stamp active trace
context on all logs, so users can correlate them.

If your language and ecosystem don't have common logging support, use [span
events][] to share additional app details. Events maybe more convenient if you
want to add attributes as well.

As a rule of thumb, use events or logs for verbose data instead of spans. Always
attach events to the span instance that your instrumentation created. Avoid
using the active span if you can, since you don't control what it refers to.

## 上下文传播

### 提取上下文

If you work on a library or a service that receives upstream calls, e.g. a web
framework or a messaging consumer, you should extract context from the incoming
request/message. OpenTelemetry provides the `Propagator` API, which hides
specific propagation standards and reads the trace `Context` from the wire. In
case of a single response, there is just one context on the wire, which becomes
the parent of the new span the library creates.

After you create a span, you should pass new trace context to the application
code (callback or handler), by making the span active; if possible, you should
do this explicitly.

```java
// extract the context
Context extractedContext = propagator.extract(Context.current(), httpExchange, getter);
Span span = tracer.spanBuilder("receive")
            .setSpanKind(SpanKind.SERVER)
            .setParent(extractedContext)
            .startSpan();

// make span active so any nested telemetry is correlated
try (Scope unused = span.makeCurrent()) {
  userCode();
} catch (Exception e) {
  span.recordException(e);
  span.setStatus(StatusCode.ERROR);
  throw e;
} finally {
  span.end();
}
```

Here're the full
[examples of context extraction in Java](/docs/instrumentation/java/manual/#context-propagation),
check out OpenTelemetry documentation in your language.

In the case of a messaging system, you may receive more than one message at
once. Received messages become
[_links_](/docs/instrumentation/java/manual/#create-spans-with-links) on the
span you create. Refer to
[messaging conventions](/docs/specs/otel/trace/semantic_conventions/messaging/)
for details (WARNING: messaging conventions are
[under constructions](https://github.com/open-telemetry/oteps/pull/173) 🚧).

### 注入上下文

When you make an outbound call, you will usually want to propagate context to
the downstream service. In this case, you should create a new span to trace the
outgoing call and use `Propagator` API to inject context into the message. There
may be other cases where you might want to inject context, e.g. when creating
messages for async processing.

```java
Span span = tracer.spanBuilder("send")
            .setSpanKind(SpanKind.CLIENT)
            .startSpan();

// make span active so any nested telemetry is correlated
// even network calls might have nested layers of spans, logs or events
try (Scope unused = span.makeCurrent()) {
  // inject the context
  propagator.inject(Context.current(), transportLayer, setter);
  send();
} catch (Exception e) {
  span.recordException(e);
  span.setStatus(StatusCode.ERROR);
  throw e;
} finally {
  span.end();
}
```

Here's the full
[example of context injection in Java](/docs/instrumentation/java/manual/#context-propagation).

There might be some exceptions:

- downstream service does not support metadata or prohibits unknown fields
- downstream service does not define correlation protocols. Is it possible that
  some future service version will support compatible context propagation?
  Inject it!
- downstream service supports custom correlation protocol.
  - best effort with custom propagator: use OpenTelemetry trace context if
    compatible.
  - or generate and stamp custom correlation ids on the span.

### 进程内的

- **Make your spans active** (aka current): it enables correlating spans with
  logs and any nested auto-instrumentations.
- If the library has a notion of context, support **optional** explicit trace
  context propagation _in addition_ to active spans
  - put spans (trace context) created by library in the context explicitly,
    document how to access it
  - allow users to pass trace context in your context
- Within the library, propagate trace context explicitly - active spans may
  change during callbacks!
  - capture active context from users on the public API surface as soon as you
    can, use it as a parent context for your spans
  - pass context around and stamp attributes, exceptions, events on explicitly
    propagated instances
  - this is essential if you start threads explicitly, do background processing
    or other things that can break due to async context flow limitations in your
    language

## Misc

### 设备注册

Please add your instrumentation library to the
[OpenTelemetry registry](/ecosystem/registry/), so users can find it.

### 表演

OpenTelemetry API is no-op and very performant when there is no SDK in the
application. When OpenTelemetry SDK is configured, it
[consumes bound resources](/docs/specs/otel/performance/).

Real-life applications, especially on the high scale, would frequently have
head-based sampling configured. Sampled-out spans are cheap and you can check if
the span is recording, to avoid extra allocations and potentially expensive
calculations, while populating attributes.

```java
// some attributes are important for sampling, they should be provided at creation time
Span span = tracer.spanBuilder(String.format("SELECT %s.%s", dbName, collectionName))
        .setSpanKind(SpanKind.CLIENT)
        .setAttribute("db.name", dbName)
        ...
        .startSpan();

// other attributes, especially those that are expensive to calculate
// should be added if span is recording
if (span.isRecording()) {
    span.setAttribute("db.statement", sanitize(query.statement()))
}
```

### 错误处理

OpenTelemetry API is
[forgiving at runtime](/docs/specs/otel/error-handling/#basic-error-handling-principles) -
does not fail on invalid arguments, never throws, and swallows exceptions. This
way instrumentation issues do not affect application logic. Test the
instrumentation to notice issues OpenTelemetry hides at runtime.

### 测试

Since OpenTelemetry has variety of auto-instrumentations, it's useful to try how
your instrumentation interacts with other telemetry: incoming requests, outgoing
requests, logs, etc. Use a typical application, with popular frameworks and
libraries and all tracing enabled when trying out your instrumentation. Check
out how libraries similar to yours show up.

For unit testing, you can usually mock or fake `SpanProcessor` and
`SpanExporter`.

```java
@Test
public void checkInstrumentation() {
  SpanExporter exporter = new TestExporter();

  Tracer tracer = OpenTelemetrySdk.builder()
           .setTracerProvider(SdkTracerProvider.builder()
              .addSpanProcessor(SimpleSpanProcessor.create(exporter)).build()).build()
           .getTracer("test");
  // run test ...

  validateSpans(exporter.exportedSpans);
}

class TestExporter implements SpanExporter {
  public final List<SpanData> exportedSpans = Collections.synchronizedList(new ArrayList<>());

  @Override
  public CompletableResultCode export(Collection<SpanData> spans) {
    exportedSpans.addAll(spans);
    return CompletableResultCode.ofSuccess();
  }
  ...
}
```

[插装库]: /docs/specs/otel/overview/#instrumentation-libraries
[span events]: /docs/specs/otel/trace/api/#add-events

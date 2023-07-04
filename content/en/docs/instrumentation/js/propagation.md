---
title: 传播
description: Context propagation for the JS SDK
aliases: [/docs/instrumentation/js/api/propagation]
weight: 65
spelling: cSpell:ignore traceparent tracestate
---

传播是在服务和进程之间移动数据的机制。虽然不限于跟踪，但它允许跟踪构建关于系统的
因果信息，这些信息可以跨任意分布在进程和网络边界上的服务。

## 使用插装库进行上下文传播

对于绝大多数用例，上下文传播是通过插装库完成的。

例如，如果你有几个通过 HTTP 通信的 Node.js 服务，你可以使
用[`express`]和[`http`]插装库在服务之间自动传播跟踪上下文。

[express]: https://www.npmjs.com/package/@opentelemetry/instrumentation-express
[http]: https://www.npmjs.com/package/@opentelemetry/instrumentation-http

**强烈建议您使用插装库来传播上下文。** 虽然可以手动传播上下文，但如果您的系统使
用插装库在服务之间进行通信，请使用匹配的插装库来传播上下文。

参考[插装库](./libraries.md)了解更多关于插装库以及如何使用它们的信息。

## 手动 W3C 跟踪上下文传播

在某些情况下，不可能使用检测库传播上下文。可能没有与您用来使服务相互通信的库匹配
的插装库。或者您可能有一些需求是插装库无法满足的，即使它们存在。

当你必须手动传播上下文时，你可以使用[上下文 api](./context.md)。

下面的通用示例演示了如何手动传播跟踪上下文。

首先，在发送服务中，你需要注入当前的`context`:

```js
// Sending service
import { context, propagation, trace } from '@opentelemetry/api';
const output = {};

// Serialize the traceparent and tracestate from context into
// an output object.
//
// This example uses the active trace context, but you can
// use whatever context is appropriate to your scenario.
propagation.inject(context.active(), output);

const { traceparent, tracestate } = output;
// You can then pass the traceparent and tracestate
// data to whatever mechanism you use to propagate
// across services.
```

在接收服务上，您需要提取`context`(例如，从解析过的 HTTP 标头中)，然后将它们设置
为当前跟踪上下文。

```js
// Receiving service
import { context, propagation, trace } from '@opentelemetry/api';

// Assume "input" is an object with 'traceparent' & 'tracestate' keys
const input = {};

// Extracts the 'traceparent' and 'tracestate' data into a context object.
//
// You can then treat this context as the active context for your
// traces.
let activeContext = propagation.extract(context.active(), input);

let tracer = trace.getTracer('app-name');

let span = tracer.startSpan(
  spanName,
  {
    attributes: {},
  },
  activeContext
);

// Set the created span as active in the deserialized context.
trace.setSpan(activeContext, span);
```

从那里，当您有一个反序列化的活动上下文时，您可以创建 spans，这些 spans 将成为来
自其他服务的相同跟踪的一部分。

你也可以使用[上下文](./context.md) API 以其他方式修改或设置反序列化的上下文。

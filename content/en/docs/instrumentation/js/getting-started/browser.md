---
title: Browser
aliases: [/docs/js/getting_started/browser]
weight: 20
---

While this guide uses the example application presented below, the steps to
instrument your own application should be similar.

## 先决条件

Ensure that you have the following installed locally:

- [Node.js](https://nodejs.org/en/download/)
- [TypeScript](https://www.typescriptlang.org/download), if you will be using
  TypeScript.

## 示例应用程序

This is a very simple guide, if you'd like to see more complex examples go to
[examples/opentelemetry-web](https://github.com/open-telemetry/opentelemetry-js/tree/main/examples/opentelemetry-web).

Copy the following file into an empty directory and call it `index.html`.

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Document Load Instrumentation Example</title>
    <base href="/" />
    <!--
      https://www.w3.org/TR/trace-context/
      Set the `traceparent` in the server's HTML template code. It should be
      dynamically generated server side to have the server's request trace Id,
      a parent span Id that was set on the server's request span, and the trace
      flags to indicate the server's sampling decision
      (01 = sampled, 00 = notsampled).
      '{version}-{traceId}-{spanId}-{sampleDecision}'
    -->
    <meta
      name="traceparent"
      content="00-ab42124a3c573678d4d8b21ba52df3bf-d21f7bc17caa5aba-01"
    />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    Example of using Web Tracer with document load instrumentation with console
    exporter and collector exporter
  </body>
</html>
```

### 安装

To create traces in the browser, you will need `@opentelemetry/sdk-trace-web`,
and the instrumentation `@opentelemetry/instrumentation-document-load`:

```shell
npm init -y
npm install @opentelemetry/api \
  @opentelemetry/sdk-trace-web \
  @opentelemetry/instrumentation-document-load \
  @opentelemetry/context-zone
```

### 初始化与配置

If you are coding in TypeScript, then run the following command:

```shell
tsc --init
```

Then acquire [parcel](https://parceljs.org/), which will (among other things)
let you work in Typescript.

```shell
npm install --save-dev parcel
```

Create an empty code file named `document-load` with a `.ts` or `.js` extension,
as appropriate, based on the language you've chosen to write your app in. Add
the following code to your HTML right before the `</body>` closing tag:

<!-- prettier-ignore-start -->
{{< tabpane lang=html persistLang=false >}}
{{< tab TypeScript >}}
<script type="module" src="document-load.ts"></script>
{{< /tab >}}
{{< tab JavaScript >}}
<script type="module" src="document-load.js"></script>
{{< /tab >}}
{{< /tabpane >}}
<!-- prettier-ignore-end -->

We will add some code that will trace the document load timings and output those
as OpenTelemetry Spans.

### 创建跟踪程序提供程序

Add the following code to the `document-load.ts|js` to create a tracer provider,
which brings the instrumentation to trace document load:

```js
/* document-load.ts|js file - the code snippet is the same for both the languages */
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { registerInstrumentations } from '@opentelemetry/instrumentation';

const provider = new WebTracerProvider();

provider.register({
  // Changing default contextManager to use ZoneContextManager - supports asynchronous operations - optional
  contextManager: new ZoneContextManager(),
});

// Registering instrumentations
registerInstrumentations({
  instrumentations: [new DocumentLoadInstrumentation()],
});
```

Now build the app with parcel:

```shell
npx parcel index.html
```

and open the development webserver (e.g. at `http://localhost:1234`) to see if
your code works.

There will be no output of traces yet, for this we need to add an exporter.

### 创建导出器

In the following example, we will use the `ConsoleSpanExporter` which prints all
spans to the console.

In order to visualize and analyze your traces, you will need to export them to a
tracing backend. Follow [these instructions](../../exporters) for setting up a
backend and exporter.

You may also want to use the `BatchSpanProcessor` to export spans in batches in
order to more efficiently use resources.

To export traces to the console, modify `document-load.ts|js` so that it matches
the following code snippet:

```js
/* document-load.ts|js file - the code is the same for both the languages */
import {
  ConsoleSpanExporter,
  SimpleSpanProcessor,
} from '@opentelemetry/sdk-trace-base';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { registerInstrumentations } from '@opentelemetry/instrumentation';

const provider = new WebTracerProvider();
provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));

provider.register({
  // Changing default contextManager to use ZoneContextManager - supports asynchronous operations - optional
  contextManager: new ZoneContextManager(),
});

// Registering instrumentations
registerInstrumentations({
  instrumentations: [new DocumentLoadInstrumentation()],
});
```

现在，重新构建应用程序并再次打开浏览器。
在开发人员工具栏的控制台中，您应该看到正在导出一些跟踪:

```json
{
  "traceId": "ab42124a3c573678d4d8b21ba52df3bf",
  "parentId": "cfb565047957cb0d",
  "name": "documentFetch",
  "id": "5123fc802ffb5255",
  "kind": 0,
  "timestamp": 1606814247811266,
  "duration": 9390,
  "attributes": {
    "component": "document-load",
    "http.response_content_length": 905
  },
  "status": {
    "code": 0
  },
  "events": [
    {
      "name": "fetchStart",
      "time": [1606814247, 811266158]
    },
    {
      "name": "domainLookupStart",
      "time": [1606814247, 811266158]
    },
    {
      "name": "domainLookupEnd",
      "time": [1606814247, 811266158]
    },
    {
      "name": "connectStart",
      "time": [1606814247, 811266158]
    },
    {
      "name": "connectEnd",
      "time": [1606814247, 811266158]
    },
    {
      "name": "requestStart",
      "time": [1606814247, 819101158]
    },
    {
      "name": "responseStart",
      "time": [1606814247, 819791158]
    },
    {
      "name": "responseEnd",
      "time": [1606814247, 820656158]
    }
  ]
}
```

### 添加的设备

如果你想检测AJAX请求、用户交互和其他，你可以为它们注册额外的检测:

```javascript
registerInstrumentations({
  instrumentations: [
    new UserInteractionInstrumentation(),
    new XMLHttpRequestInstrumentation(),
  ],
});
```

## 用于Web的元包

要利用最常见的仪器都在一个你可以简单地使用[开放遥测元包的web](https://www.npmjs.com/package/@opentelemetry/auto-instrumentations-web)

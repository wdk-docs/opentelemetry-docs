---
title: 浏览器
aliases: [/docs/js/getting_started/browser]
weight: 20
---

虽然本指南使用下面提供的示例应用程序，但是检测您自己的应用程序的步骤应该是类似的
。

## 先决条件

确保在本地安装了以下软件:

- [Node.js](https://nodejs.org/en/download/)
- [TypeScript](https://www.typescriptlang.org/download), 如果你要使用
  TypeScript。

## 示例应用程序

这是一个非常简单的指南，如果你想看更复杂的例子，请访问
[examples/opentelemetry-web](https://github.com/open-telemetry/opentelemetry-js/tree/main/examples/opentelemetry-web).

将以下文件复制到一个空目录中，并将其命名为`index.html`。

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

要在浏览器中创建跟踪，你需要`@opentelemetry/sdk-trace-web`和插
装`@opentelemetry/instrumentation-document-load`:

```shell
npm init -y
npm install @opentelemetry/api \
  @opentelemetry/sdk-trace-web \
  @opentelemetry/instrumentation-document-load \
  @opentelemetry/context-zone
```

### 初始化与配置

如果你在用 TypeScript 编码，那么运行下面的命令:

```shell
tsc --init
```

然后获取[parcel](https://parceljs.org/)，它可以让你用 Typescript 工作。

```shell
npm install --save-dev parcel
```

创建一个名为`document-load`的空代码文件，扩展名为`.ts`或`.js`，根据你选择的语言
编写应用程序。将以下代码添加到 HTML 中，在`</body>`结束标签之前:

=== "TypeScript"

    ```ts
    <script type="module" src="document-load.ts"></script>
    ```

=== "JavaScript"

    ```js
    <script type="module" src="document-load.js"></script>
    ```

我们将添加一些代码来跟踪文档加载时间并将其输出为 OpenTelemetry Spans.

### 创建跟踪程序提供程序

将以下代码添加到 `document-load.ts|js` 中，以创建一个跟踪程序提供程序，它将提供
跟踪文档加载的工具:

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

现在用包裹构建应用程序:

```shell
npx parcel index.html
```

并打开开发 web 服务器(例如在`http://localhost:1234`)，看看你的代码是否工作。

目前还没有痕迹输出，为此我们需要添加一个导出器。

### 创建导出器

在下面的示例中，我们将使用`ConsoleSpanExporter`将所有的 spans 打印到控制台。

为了可视化和分析您的跟踪，您需要将它们导出到跟踪后端。按照[这些说明](../../
exports)设置后端和导出器。

您可能还想使用`BatchSpanProcessor`来批量导出 spans，以便更有效地使用资源。

要将跟踪信息导出到控制台，请修改`document-load.ts|js`，以便与以下代码片段匹配:

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

现在，重新构建应用程序并再次打开浏览器。在开发人员工具栏的控制台中，您应该看到正
在导出一些跟踪:

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

如果你想检测 AJAX 请求、用户交互和其他，你可以为它们注册额外的检测:

```javascript
registerInstrumentations({
  instrumentations: [
    new UserInteractionInstrumentation(),
    new XMLHttpRequestInstrumentation(),
  ],
});
```

## 用于 Web 的元包

要利用最常见的插装都在一个你可以简单地使
用[开放遥测元包的 web](https://www.npmjs.com/package/@opentelemetry/auto-instrumentations-web)

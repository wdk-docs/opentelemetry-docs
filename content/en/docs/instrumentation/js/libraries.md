---
title: 使用工具库
linkTitle: Libraries
weight: 40
spelling: cSpell:ignore autoinstrumentation metapackage
---

您可以使用[插装库](../../specs/otel/glossary/#instrumentation-library)为库或框架
生成遥测数据。

例如
，[用于 Express 的 instrumentation 库](https://www.npmjs.com/package/@opentelemetry/instrumentation-express)将
根据入站 HTTP 请求自动创建[span](../../concepts/signals/traces/#span)。

## Setup

每个工具库都是一个 NPM 包，安装通常是这样完成的:

```console
npm install <name-of-package>
```

它通常在应用程序启动时注册，例如在创
建[TracerProvider](../../concepts/signals/traces/#tracer-provider)时.

## Node.js

### Node 自动插装包

OpenTelemetry 还定义
了[auto-instrumentations-node](https://www.npmjs.com/package/@opentelemetry/auto-instrumentations-node)元
包，它将所有基于 node .js 的仪器库捆绑到一个包中。这是一种方便的方法，可以为所有
库添加自动生成的遥测功能，而且工作量很小。

要使用这个包，首先安装它:

```shell
npm install @opentelemetry/auto-instrumentations-node
```

然后在你的跟踪初始化代码中，使用`registerInstrumentations`:

<!-- textlint-disable -->

<!-- prettier-ignore-start -->
{{< tabpane langEqualsHeader=true >}}

{{< tab TypeScript >}}
```ts
/* tracing.ts */

// Import dependencies
import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
import opentelemetry from "@opentelemetry/api";
import { Resource } from "@opentelemetry/resources";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { ConsoleSpanExporter, BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";

// This registers all instrumentation packages
registerInstrumentations({
  instrumentations: [
    getNodeAutoInstrumentations()
  ],
});

const resource =
  Resource.default().merge(
    new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: "service-name-here",
      [SemanticResourceAttributes.SERVICE_VERSION]: "0.1.0",
    })
  );

const provider = new NodeTracerProvider({
    resource: resource,
});
const exporter = new ConsoleSpanExporter();
const processor = new BatchSpanProcessor(exporter);
provider.addSpanProcessor(processor);

provider.register();
```
{{< /tab >}}

{{< tab JavaScript >}}
```js
/* tracing.js */

// Require dependencies
const { getNodeAutoInstrumentations } = require("@opentelemetry/auto-instrumentations-node");
const opentelemetry = require("@opentelemetry/api");
const { Resource } = require("@opentelemetry/resources");
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions");
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node");
const { registerInstrumentations } = require("@opentelemetry/instrumentation");
const { ConsoleSpanExporter, BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base");

// This registers all instrumentation packages
registerInstrumentations({
  instrumentations: [
    getNodeAutoInstrumentations()
  ],
});

const resource =
  Resource.default().merge(
    new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: "service-name-here",
      [SemanticResourceAttributes.SERVICE_VERSION]: "0.1.0",
    })
  );

const provider = new NodeTracerProvider({
    resource: resource,
});
const exporter = new ConsoleSpanExporter();
const processor = new BatchSpanProcessor(exporter);
provider.addSpanProcessor(processor);

provider.register();
```

{{< /tab >}}

{{< /tabpane >}}
<!-- prettier-ignore-end -->

<!-- textlint-enable -->

### 使用单独的插装包

如果您不希望使用元包，也许是为了减少依赖关系图的大小，您可以安装和注册单独的工具
包。

For example, here's how you can install and register only the
[instrumentation-express](https://www.npmjs.com/package/@opentelemetry/instrumentation-express)
and
[instrumentation-http](https://www.npmjs.com/package/@opentelemetry/instrumentation-http)
packages to instrument inbound and outbound HTTP traffic.

```shell
npm install --save @opentelemetry/instrumentation-http @opentelemetry/instrumentation-express
```

And then register each instrumentation library:

<!-- prettier-ignore-start -->
{{< tabpane langEqualsHeader=true >}}

{{< tab TypeScript >}}
```ts
/* tracing.ts */

// Import dependencies
import { HttpInstrumentation } from "@opentelemetry/instrumentation-http";
import { ExpressInstrumentation } from "@opentelemetry/instrumentation-express";
import opentelemetry from "@opentelemetry/api";
import { Resource } from "@opentelemetry/resources";
import { SemanticResourceAttributes } from "@opentelemetry/semantic-conventions";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { ConsoleSpanExporter, BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";

// This registers all instrumentation packages
registerInstrumentations({
  instrumentations: [
    // Express instrumentation expects HTTP layer to be instrumented
    new HttpInstrumentation(),
    new ExpressInstrumentation(),
  ],
});

const resource =
  Resource.default().merge(
    new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: "service-name-here",
      [SemanticResourceAttributes.SERVICE_VERSION]: "0.1.0",
    })
  );

const provider = new NodeTracerProvider({
    resource: resource,
});
const exporter = new ConsoleSpanExporter();
const processor = new BatchSpanProcessor(exporter);
provider.addSpanProcessor(processor);

provider.register();
```
{{< /tab >}}

{{< tab JavaScript >}}
```js
/* tracing.js */

// Require dependencies
const { HttpInstrumentation } = require("@opentelemetry/instrumentation-http");
const { ExpressInstrumentation } = require("@opentelemetry/instrumentation-express");
const opentelemetry = require("@opentelemetry/api");
const { Resource } = require("@opentelemetry/resources");
const { SemanticResourceAttributes } = require("@opentelemetry/semantic-conventions");
const { NodeTracerProvider } = require("@opentelemetry/sdk-trace-node");
const { registerInstrumentations } = require("@opentelemetry/instrumentation");
const { ConsoleSpanExporter, BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base");

// This registers all instrumentation packages
registerInstrumentations({
  instrumentations: [
    // Express instrumentation expects HTTP layer to be instrumented
    new HttpInstrumentation(),
    new ExpressInstrumentation(),
  ],
});

const resource =
  Resource.default().merge(
    new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]: "service-name-here",
      [SemanticResourceAttributes.SERVICE_VERSION]: "0.1.0",
    })
  );

const provider = new NodeTracerProvider({
    resource: resource,
});
const exporter = new ConsoleSpanExporter();
const processor = new BatchSpanProcessor(exporter);
provider.addSpanProcessor(processor);

provider.register();
```
{{< /tab >}}

{{< /tabpane >}}
<!-- prettier-ignore-end -->

## 配置工具库

一些工具库提供了额外的配置选项。

例如
，[Express instrumentation](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-express#express-instrumentation-options)提
供了忽略指定中间件或丰富使用请求钩子自动创建的范围的方法。

您需要参考每个仪器库的文档来进行高级配置。

## 可用的仪器库

OpenTelemetry 生成的仪器库的完整列表可
从[opentelemetry-js-contrib](https://github.com/open-telemetry/opentelemetry-js-contrib)存
储库获得。

您还可以
在[registry](../../../ecosystem/registry/index.md?language=js&component=instrumentation)中
找到更多可用的工具。

## 下一个步骤

在你设置好仪器库之后，你可能想添加[手动工具](./manual.md)来收集自定义的遥测数据
。

你还需要配置一个合适的导出器来[导出遥测数据](./exporters.md)到一个或多个遥测后端
。

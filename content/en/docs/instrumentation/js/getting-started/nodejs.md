---
title: Node.js
description: Get telemetry for your app in less than 5 minutes!
aliases: [/docs/js/getting_started/nodejs]
spelling: cSpell:ignore rolldice autoinstrumentation autoinstrumentations KHTML
weight: 10
---

本页将向您展示如何在 Node.js 中开始使用 OpenTelemetry。

您将学习如何自动检测一个简单的应用程序，以一种
将[traces][]、[metrics][]和[logs][]发送到控制台的方式。

## 先决条件

确保在本地安装了以下软件:

- [Node.js](https://nodejs.org/en/download/)
- [TypeScript](https://www.typescriptlang.org/download), 如果你要使用
  TypeScript。

## 示例应用程序

下面的示例使用一个基本的[Express](https://expressjs.com/)应用程序。如果你不使用
Express，没关系——你也可以将 OpenTelemetry JavaScript 与其他 web 框架一起使用，比
如 Koa 和 Nest.JS。要获得支持框架的完整库列表，请参
见[registry](/ecosystem/registry/?component=instrumentation&language=js)。

有关更详细的示例，请参见[示例](/docs/instrumentation/js/examples/).

### 依赖关系

首先，在新目录中设置一个空的`package.json`:

```shell
npm init -y
```

接下来，安装 Express 依赖项。

=== "TypeScript"

    ```sh
    npm install typescript \
      ts-node \
      @types/node \
      express \
      @types/express
    ```

=== "JavaScript"

    ```sh
    npm install express
    ```

### 创建并启动 HTTP 服务器

Create a file named `app.ts` (or `app.js` if not using typescript) and add the
following code to it:

=== "TypeScript"

    ```ts
    /*app.ts*/
    import express, { Express } from 'express';

    const PORT: number = parseInt(process.env.PORT || '8080');
    const app: Express = express();

    function getRandomNumber(min: number, max: number) {
      return Math.floor(Math.random() * (max - min) + min);
    }

    app.get('/rolldice', (req, res) => {
      res.send(getRandomNumber(1, 6).toString());
    });

    app.listen(PORT, () => {
      console.log(`Listening for requests on http://localhost:${PORT}`);
    });
    ```

=== "JavaScript"

    ```js
    /*app.js*/
    const express = require('express');

    const PORT = parseInt(process.env.PORT || '8080');
    const app = express();

    function getRandomNumber(min, max) {
      return Math.floor(Math.random() * (max - min) + min);
    }

    app.get('/rolldice', (req, res) => {
      res.send(getRandomNumber(1, 6).toString());
    });

    app.listen(PORT, () => {
      console.log(`Listening for requests on http://localhost:${PORT}`);
    });
    ```

Run the application with the following command and open
<http://localhost:8080/rolldice> in your web browser to ensure it is working.

=== "TypeScript"

    ```sh
    $ npx ts-node app.ts
    Listening for requests on http://localhost:8080
    ```

=== "JavaScript"

    ```sh
    $ node app.js
    Listening for requests on http://localhost:8080
    ```

## 插装

The following shows how to install, initialize, and run an application
instrumented with OpenTelemetry.

### 依赖

First, install the Node SDK and autoinstrumentations package.

The Node SDK lets you initialize OpenTelemetry with several configuration
defaults that are correct for the majority of use cases.

The `auto-instrumentations-node` package installs instrumentation packages that
will automatically create spans corresponding to code called in libraries. In
this case, it provides instrumentation for Express, letting the example app
automatically create spans for each incoming request.

```shell
npm install @opentelemetry/sdk-node \
  @opentelemetry/api \
  @opentelemetry/auto-instrumentations-node \
  @opentelemetry/sdk-metrics
```

To find all autoinstrumentation modules, you can look at the
[registry](/ecosystem/registry/?language=js&component=instrumentation).

### 设置

The instrumentation setup and configuration must be run _before_ your
application code. One tool commonly used for this task is the
[--require](https://nodejs.org/api/cli.html#-r---require-module) flag.

Create a file named `instrumentation.ts` (or `instrumentation.js` if not using
typescript) , which will contain your instrumentation setup code.

=== "TypeScript"

    ```ts
    /*instrumentation.ts*/
    import { NodeSDK } from '@opentelemetry/sdk-node';
    import { ConsoleSpanExporter } from '@opentelemetry/sdk-trace-node';
    import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
    import {
      PeriodicExportingMetricReader,
      ConsoleMetricExporter,
    } from '@opentelemetry/sdk-metrics';

    const sdk = new NodeSDK({
      traceExporter: new ConsoleSpanExporter(),
      metricReader: new PeriodicExportingMetricReader({
        exporter: new ConsoleMetricExporter(),
      }),
      instrumentations: [getNodeAutoInstrumentations()],
    });

    sdk.start();
    ```

=== "JavaScript"

    ```js
    /*instrumentation.js*/
    // Require dependencies
    const { NodeSDK } = require('@opentelemetry/sdk-node');
    const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-node');
    const {
      getNodeAutoInstrumentations,
    } = require('@opentelemetry/auto-instrumentations-node');
    const {
      PeriodicExportingMetricReader,
      ConsoleMetricExporter,
    } = require('@opentelemetry/sdk-metrics');

    const sdk = new NodeSDK({
      traceExporter: new ConsoleSpanExporter(),
      metricReader: new PeriodicExportingMetricReader({
        exporter: new ConsoleMetricExporter(),
      }),
      instrumentations: [getNodeAutoInstrumentations()],
    });

    sdk.start();
    ```

## 运行插装应用

现在你可以像平常一样运行你的应用程序了，但是你可以使用`--require`标志在加载应用
程序代码之前加载工具。

=== "TypeScript"

    ```sh
    $ npx ts-node --require ./instrumentation.ts app.ts
    Listening for requests on http://localhost:8080
    ```

=== "JavaScript"

    ```sh
    $ node --require ./instrumentation.js app.js
    Listening for requests on http://localhost:8080
    ```

在你的网页浏览器中打开<http://localhost:8080/rolldice>重新加载页面几次。过了一会
儿，你应该会看到控制台通过`ConsoleSpanExporter`打印出 spans。

<details>
<summary>查看示例输出</summary>

```json
{
  "traceId": "3f1fe6256ea46d19ec3ca97b3409ad6d",
  "parentId": "f0b7b340dd6e08a7",
  "name": "middleware - query",
  "id": "41a27f331c7bfed3",
  "kind": 0,
  "timestamp": 1624982589722992,
  "duration": 417,
  "attributes": {
    "http.route": "/",
    "express.name": "query",
    "express.type": "middleware"
  },
  "status": { "code": 0 },
  "events": []
}
{
  "traceId": "3f1fe6256ea46d19ec3ca97b3409ad6d",
  "parentId": "f0b7b340dd6e08a7",
  "name": "middleware - expressInit",
  "id": "e0ed537a699f652a",
  "kind": 0,
  "timestamp": 1624982589725778,
  "duration": 673,
  "attributes": {
    "http.route": "/",
    "express.name": "expressInit",
    "express.type": "middleware"
  },
  "status": { code: 0 },
  "events": []
}
{
  "traceId": "3f1fe6256ea46d19ec3ca97b3409ad6d",
  "parentId": "f0b7b340dd6e08a7",
  "name": "request handler - /",
  "id": "8614a81e1847b7ef",
  "kind": 0,
  "timestamp": 1624982589726941,
  "duration": 21,
  "attributes": {
    "http.route": "/",
    "express.name": "/",
    "express.type": "request_handler"
  },
  "status": { code: 0 },
  "events": []
}
{
  "traceId": "3f1fe6256ea46d19ec3ca97b3409ad6d",
  "parentId": undefined,
  "name": "GET /",
  "id": "f0b7b340dd6e08a7",
  "kind": 1,
  "timestamp": 1624982589720260,
  "duration": 11380,
  "attributes": {
    "http.url": "http://localhost:8080/",
    "http.host": "localhost:8080",
    "net.host.name": "localhost",
    "http.method": "GET",
    "http.route": "",
    "http.target": "/",
    "http.user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
    "http.flavor": "1.1",
    "net.transport": "ip_tcp",
    "net.host.ip": "::1",
    "net.host.port": 8080,
    "net.peer.ip": "::1",
    "net.peer.port": 61520,
    "http.status_code": 304,
    "http.status_text": "NOT MODIFIED"
  },
  "status": { "code": 1 },
  "events": []
}
```

</details>

生成的 span 跟踪到`/rolldice`路由的请求的生命周期。

再向端点发送几个请求。稍后，您将在控制台输出中看到指标，如下所示:

<details>
<summary>查看示例输出</summary>

```yaml
{
  descriptor: {
    name: 'http.server.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the inbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: [
    {
      attributes: [Object],
      startTime: [Array],
      endTime: [Array],
      value: [Object]
    }
  ]
}
{
  descriptor: {
    name: 'http.client.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the outbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: []
}
{
  descriptor: {
    name: 'db.client.connections.usage',
    type: 'UP_DOWN_COUNTER',
    description: 'The number of connections that are currently in the state referenced by the attribute "state".',
    unit: '{connections}',
    valueType: 1
  },
  dataPointType: 3,
  dataPoints: []
}
{
  descriptor: {
    name: 'http.server.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the inbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: [
    {
      attributes: [Object],
      startTime: [Array],
      endTime: [Array],
      value: [Object]
    }
  ]
}
{
  descriptor: {
    name: 'http.client.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the outbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: []
}
{
  descriptor: {
    name: 'db.client.connections.usage',
    type: 'UP_DOWN_COUNTER',
    description: 'The number of connections that are currently in the state referenced by the attribute "state".',
    unit: '{connections}',
    valueType: 1
  },
  dataPointType: 3,
  dataPoints: []
}
{
  descriptor: {
    name: 'http.server.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the inbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: [
    {
      attributes: [Object],
      startTime: [Array],
      endTime: [Array],
      value: [Object]
    }
  ]
}
{
  descriptor: {
    name: 'http.client.duration',
    type: 'HISTOGRAM',
    description: 'measures the duration of the outbound HTTP requests',
    unit: 'ms',
    valueType: 1
  },
  dataPointType: 0,
  dataPoints: []
}
{
  descriptor: {
    name: 'db.client.connections.usage',
    type: 'UP_DOWN_COUNTER',
    description: 'The number of connections that are currently in the state referenced by the attribute "state".',
    unit: '{connections}',
    valueType: 1
  },
  dataPointType: 3,
  dataPoints: []
}
```

</details>

## 下一个步骤

Enrich your instrumentation generated automatically with
[manual instrumentation](/docs/instrumentation/js/manual) of your own codebase.
This gets you customized observability data.

You'll also want to configure an appropriate exporter to
[export your telemetry data](/docs/instrumentation/js/exporters) to one or more
telemetry backends.

## 故障排除

出什么问题了吗?你可以启用诊断日志来验证 OpenTelemetry 是否被正确初始化:

=== "TypeScript"

    ```ts
    /*instrumentation.ts*/
    import { diag, DiagConsoleLogger, DiagLogLevel } from '@opentelemetry/api';

    // For troubleshooting, set the log level to DiagLogLevel.DEBUG
    diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.INFO);

    // const sdk = new NodeSDK({...
    ```

=== "JavaScript"

    ```js
    /*instrumentation.js*/
    // Require dependencies
    const { diag, DiagConsoleLogger, DiagLogLevel } = require('@opentelemetry/api');

    // For troubleshooting, set the log level to DiagLogLevel.DEBUG
    diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.INFO);

    // const sdk = new NodeSDK({...
    ```

[traces]: /docs/concepts/signals/traces/
[metrics]: /docs/concepts/signals/metrics/
[logs]: /docs/concepts/signals/logs/

---
title: 前端
---

前端负责为用户提供 UI，以及 UI 或其他客户端所利用的 API。该应用基
于[Next.JS](https://nextjs.org/)提供基于 React 的 UI 和 API 路由。

[前端源码](https://github.com/open-telemetry/opentelemetry-demo/blob/main/src/frontend/)

## 服务器插装

建议在启动 NodeJS 应用程序时使用 Node required 模块来初始化 SDK 和自动检测。在初
始化 OpenTelemetry Node.js SDK 时，您可以选择指定要利用哪些自动检测库，或者使
用`getNodeAutoInstrumentations()`函数，其中包括大多数流行的框架。
`utils/telemetry/Instrumentation.js`文件包含了初始化 SDK 和基于标
准[OpenTelemetry 环境变量](/docs/specs/otel/configuration/sdk-environment-variables/)用
于 OTLP 导出、资源属性和服务名称的自动检测所需的所有代码。

```javascript
const opentelemetry = require('@opentelemetry/sdk-node');
const {
  getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const {
  OTLPTraceExporter,
} = require('@opentelemetry/exporter-trace-otlp-grpc');
const {
  OTLPMetricExporter,
} = require('@opentelemetry/exporter-metrics-otlp-grpc');
const { PeriodicExportingMetricReader } = require('@opentelemetry/sdk-metrics');
const {
  alibabaCloudEcsDetector,
} = require('@opentelemetry/resource-detector-alibaba-cloud');
const {
  awsEc2Detector,
  awsEksDetector,
} = require('@opentelemetry/resource-detector-aws');
const {
  containerDetector,
} = require('@opentelemetry/resource-detector-container');
const { gcpDetector } = require('@opentelemetry/resource-detector-gcp');
const {
  envDetector,
  hostDetector,
  osDetector,
  processDetector,
} = require('@opentelemetry/resources');

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new OTLPTraceExporter(),
  instrumentations: [
    getNodeAutoInstrumentations({
      // only instrument fs if it is part of another trace
      '@opentelemetry/instrumentation-fs': {
        requireParentSpan: true,
      },
    }),
  ],
  metricReader: new PeriodicExportingMetricReader({
    exporter: new OTLPMetricExporter(),
  }),
  resourceDetectors: [
    containerDetector,
    envDetector,
    hostDetector,
    osDetector,
    processDetector,
    alibabaCloudEcsDetector,
    awsEksDetector,
    awsEc2Detector,
    gcpDetector,
  ],
});

sdk.start();
```

节点所需模块使用`--require`命令行参数加载。这可以
在`package.json`的`scripts.start`部分中完成，并使用`npm start`启动应用程序。

```json
  "scripts": {
    "start": "node --require ./Instrumentation.js server.js",
  },
```

## 追踪

### Span 异常和状态

You can use the span object's `recordException` function to create a span event
with the full stack trace of a handled error. When recording an exception also
be sure to set the span's status accordingly. You can see this in the catch
block of the `NextApiHandler` function in the
`utils/telemetry/InstrumentationMiddleware.ts` file.

```typescript
span.recordException(error as Exception);
span.setStatus({ code: SpanStatusCode.ERROR });
```

### 创建新 spans

New spans can be created and started using
`Tracer.startSpan("spanName", options)`. Several options can be used to specify
how the span can be created.

- `root: true` will create a new trace, setting this span as the root.
- `links` are used to specify links to other spans (even within another trace)
  that should be referenced.
- `attributes` are key/value pairs added to a span, typically used for
  application context.

```typescript
span = tracer.startSpan(`HTTP ${method}`, {
  root: true,
  kind: SpanKind.SERVER,
  links: [{ context: syntheticSpan.spanContext() }],
  attributes: {
    'app.synthetic_request': true,
    [SemanticAttributes.HTTP_TARGET]: target,
    [SemanticAttributes.HTTP_STATUS_CODE]: response.statusCode,
    [SemanticAttributes.HTTP_METHOD]: method,
    [SemanticAttributes.HTTP_USER_AGENT]: headers['user-agent'] || '',
    [SemanticAttributes.HTTP_URL]: `${headers.host}${url}`,
    [SemanticAttributes.HTTP_FLAVOR]: httpVersion,
  },
});
```

## 浏览器插装

The web-based UI that the frontend provides is also instrumented for web
browsers. OpenTelemetry instrumentation is included as part of the Next.js App
component in `pages/_app.tsx`. Here instrumentation is imported and initialized.

```typescript
import FrontendTracer from '../utils/telemetry/FrontendTracer';

if (typeof window !== 'undefined') FrontendTracer();
```

The `utils/telemetry/FrontendTracer.ts` file contains code to initialize a
TracerProvider, establish an OTLP export, register trace context propagators,
and register web specific auto-instrumentation libraries. Since the browser will
send data to an OpenTelemetry collector that will likely be on a separate
domain, CORS headers are also setup accordingly.

As part of the changes to carry over the `synthetic_request` attribute flag for
the backend services, the `applyCustomAttributesOnSpan` configuration function
has been added to the `instrumentation-fetch` library custom span attributes
logic that way every browser-side span will include it.

```typescript
import {
  CompositePropagator,
  W3CBaggagePropagator,
  W3CTraceContextPropagator,
} from '@opentelemetry/core';
import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { getWebAutoInstrumentations } from '@opentelemetry/auto-instrumentations-web';
import { Resource } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';

const FrontendTracer = async () => {
  const { ZoneContextManager } = await import('@opentelemetry/context-zone');

  const provider = new WebTracerProvider({
    resource: new Resource({
      [SemanticResourceAttributes.SERVICE_NAME]:
        process.env.NEXT_PUBLIC_OTEL_SERVICE_NAME,
    }),
  });

  provider.addSpanProcessor(new SimpleSpanProcessor(new OTLPTraceExporter()));

  const contextManager = new ZoneContextManager();

  provider.register({
    contextManager,
    propagator: new CompositePropagator({
      propagators: [
        new W3CBaggagePropagator(),
        new W3CTraceContextPropagator(),
      ],
    }),
  });

  registerInstrumentations({
    tracerProvider: provider,
    instrumentations: [
      getWebAutoInstrumentations({
        '@opentelemetry/instrumentation-fetch': {
          propagateTraceHeaderCorsUrls: /.*/,
          clearTimingResources: true,
          applyCustomAttributesOnSpan(span) {
            span.setAttribute('app.synthetic_request', 'false');
          },
        },
      }),
    ],
  });
};

export default FrontendTracer;
```

## 度量

TBD

## 日志

TBD

## Baggage

前端利用 OpenTelemetry 包袱来检查请求是否合成(来自负载生成器)。合成请求将强制创
建新的跟踪。来自新跟踪的根 span 将包含许多与 HTTP 请求检测的 span 相同的属性。

要确定是否设置了一个包袱项，您可以利用`propagation`API 来解析包袱标头，并利
用`baggage`API 来获取或设置条目。

```typescript
    const baggage = propagation.getBaggage(context.active());
    if (baggage?.getEntry("synthetic_request")?.value == "true") {...}
```

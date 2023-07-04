---
title: 导出器
weight: 50
spelling: cSpell:ignore proto nginx openzipkin
---

为了可视化和分析您的痕迹，您需要将它们导出到后端，
如[Jaeger](https://www.jaegertracing.io/)或[Zipkin](https://zipkin.io/)。
OpenTelemetry JS 为一些常见的开源后端提供了导出器。

下面您将找到一些关于如何设置后端和匹配的导出器的介绍。

## OTLP 端点

要将跟踪数据发送到 OTLP 端点(如[collector](/docs/collector)或 Jaeger)，您需要使
用导出包，例如`@opentelemetry/exporter-trace-otlp-proto`:

```shell
npm install --save @opentelemetry/exporter-trace-otlp-proto \
  @opentelemetry/exporter-metrics-otlp-proto
```

接下来，将导出器配置为指向 OTLP 端点。例如，你可以
从[入门](./getting-started/nodejs.md)中更新`instrumentation.ts|js`，如下所示:

=== "Typescript"

    ```Typescript
    /*tracing.ts*/
    import * as opentelemetry from "@opentelemetry/sdk-node";
    import {
      getNodeAutoInstrumentations,
    } from "@opentelemetry/auto-instrumentations-node";
    import {
      OTLPTraceExporter,
    } from "@opentelemetry/exporter-trace-otlp-proto";
    import {
      OTLPMetricExporter
    } from "@opentelemetry/exporter-metrics-otlp-proto";
    import {
      PeriodicExportingMetricReader
    } from "@opentelemetry/sdk-metrics";

    const sdk = new opentelemetry.NodeSDK({
      traceExporter: new OTLPTraceExporter({
        // optional - default url is http://localhost:4318/v1/traces
        url: "<your-otlp-endpoint>/v1/traces",
        // optional - collection of custom headers to be sent with each request, empty by default
        headers: {},
      }),
      metricReader: new PeriodicExportingMetricReader({
        exporter: new OTLPMetricExporter({
          url: '<your-otlp-endpoint>/v1/metrics', // url is optional and can be omitted - default is http://localhost:4318/v1/metrics
          headers: {}, // an optional object containing custom headers to be sent with each request
        }),
      }),
      instrumentations: [getNodeAutoInstrumentations()],
    });
    sdk.start();
    ```

=== "JavaScript"

    ```JavaScript
    /*tracing.js*/
    const opentelemetry = require("@opentelemetry/sdk-node");
    const {
      getNodeAutoInstrumentations,
    } = require("@opentelemetry/auto-instrumentations-node");
    const {
      OTLPTraceExporter,
    } = require("@opentelemetry/exporter-trace-otlp-proto");
    const {
      OTLPMetricExporter
    } = require("@opentelemetry/exporter-metrics-otlp-proto");
    const {
      PeriodicExportingMetricReader
    } = require('@opentelemetry/sdk-metrics');

    const sdk = new opentelemetry.NodeSDK({
      traceExporter: new OTLPTraceExporter({
        // optional - default url is http://localhost:4318/v1/traces
        url: "<your-otlp-endpoint>/v1/traces",
        // optional - collection of custom headers to be sent with each request, empty by default
        headers: {},
      }),
      metricReader: new PeriodicExportingMetricReader({
        exporter: new OTLPMetricExporter({
          url: '<your-otlp-endpoint>/v1/metrics', // url is optional and can be omitted - default is http://localhost:4318/v1/metrics
          headers: {}, // an optional object containing custom headers to be sent with each request
          concurrencyLimit: 1, // an optional limit on pending requests
        }),
      }),
      instrumentations: [getNodeAutoInstrumentations()],
    });
    sdk.start();
    ```

要快速试用`OTLPTraceExporter`，你可以在 docker 容器中运行 Jaeger:

```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

### 与 WebTracer 一起使用

当您在基于浏览器的应用程序中使用 OTLP 导出器时，您需要注意:

1. 不支持使用 gRPC 和 http/proto 进行导出
2. 您网站的[内容安全策略][](csp)可能会阻止您的导出
3. [跨域资源共享][](CORS)报头可能不允许发送导出文件
4. 您可能需要将您的收集器公开到公共互联网上

下面将介绍如何使用正确的导出器，如何配置 csp 和 CORS 标头，以及在暴露收集器时必
须采取的预防措施。

#### 使用 HTTP/JSON 的 OTLP 导出器

[OpenTelemetry Collector exporters with gRPC][]和[OpenTelemetry Collector
exporters with protobuf][]只能与 Node.js 一起工作，因此你只能使用[OpenTelemetry
Collector exporters with HTTP][]。

确保导出器的接收端(收集器或可观察性后端)支持`http/json`，并且将数据导出到正确的
端点，即确保端口设置为“4318”。

#### 配置 CSPs

如果您的网站正在使用内容安全策略(csp)，请确保包含 OTLP 端点的域。如果你的收集器
端点是`https://collector.example.com:4318/v1/traces`,添加以下指令:

```text
connect-src collector.example.com:4318/v1/traces
```

如果您的 CSP 不包括 OTLP 端点，您将看到一条错误消息，指出对您的端点的请求违反了
CSP 指令。

#### 配置 CORS 标头

如果您的网站和收集器托管在不同的来源，您的浏览器可能会阻止发送到收集器的请求。跨
域资源共享(Cross-Origin Resource Sharing, CORS)需要配置特殊的头信息。

OpenTelemetry 收集器为基于 http 的接收器提供了[一个特性][]来添加所需的标头，以允
许接收器接受来自 web 浏览器的跟踪:

```yaml
receivers:
  otlp:
    protocols:
      http:
        include_metadata: true
        cors:
          allowed_origins:
            - https://foo.bar.com
            - https://*.test.com
          allowed_headers:
            - Example-Header
          max_age: 7200
```

#### Securely expose your collector

To receive telemetry from a web application you need to allow the browsers of
your end-users to send data to your collector. If your web application is
accessible from the public internet, you also have to make your collector
accessible for everyone.

It is recommended that you do not expose your collector directly, but that you
put a reverse proxy (NGINX, Apache HTTP Server, ...) in front of it. The reverse
proxy can take care of SSL-offloading, setting the right CORS headers, and many
other features specific to web applications.

Below you will find a configuration for the popular NGINX web server to get you
started:

```nginx
server {
    listen 80 default_server;
    server_name _;
    location / {
        # Take care of preflight requests
        if ($request_method = 'OPTIONS') {
             add_header 'Access-Control-Max-Age' 1728000;
             add_header 'Access-Control-Allow-Origin' 'name.of.your.website.example.com' always;
             add_header 'Access-Control-Allow-Headers' 'Accept,Accept-Language,Content-Language,Content-Type' always;
             add_header 'Access-Control-Allow-Credentials' 'true' always;
             add_header 'Content-Type' 'text/plain charset=UTF-8';
             add_header 'Content-Length' 0;
             return 204;
        }

        add_header 'Access-Control-Allow-Origin' 'name.of.your.website.example.com' always;
        add_header 'Access-Control-Allow-Credentials' 'true' always;
        add_header 'Access-Control-Allow-Headers' 'Accept,Accept-Language,Content-Language,Content-Type' always;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass http://collector:4318;
    }
}
```

## Zipkin

To set up Zipkin as quickly as possible, run it in a docker container:

```shell
docker run --rm -d -p 9411:9411 --name zipkin openzipkin/zipkin
```

Install the exporter package as a dependency for your application:

```shell
npm install --save @opentelemetry/exporter-zipkin
```

Update your OpenTelemetry configuration to use the exporter and to send data to
your Zipkin backend:

=== "Typescript"

    ```Typescript
    import { ZipkinExporter } from "@opentelemetry/exporter-zipkin";
    import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";

    provider.addSpanProcessor(new BatchSpanProcessor(new ZipkinExporter()));
    ```

=== "JavaScript"

    ```JavaScript
    const { ZipkinExporter } = require("@opentelemetry/exporter-zipkin");
    const { BatchSpanProcessor } = require("@opentelemetry/sdk-trace-base");

    provider.addSpanProcessor(new BatchSpanProcessor(new ZipkinExporter()));
    ```

[content security policies]:
  https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/
[cross-origin resource sharing]:
  https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
[opentelemetry collector exporter with grpc]:
  https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-grpc
[opentelemetry collector exporter with protobuf]:
  https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-proto
[opentelemetry collector exporter with http]:
  https://www.npmjs.com/package/@opentelemetry/exporter-trace-otlp-http
[a feature]:
  https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/confighttp/README.md

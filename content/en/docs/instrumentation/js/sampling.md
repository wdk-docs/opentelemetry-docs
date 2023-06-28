---
title: 采样
weight: 80
---

[采样](../../concepts/sampling/)是一个限制系统产生的迹线数量的过程。 JavaScript
SDK 提供了几个[head-sampling](../../concepts/sampling#head-sampling)。

## 默认的行为

默认情况下，对所有跨度进行采样，因此对 100%的跟踪进行采样。如果您不需要管理数据
量，则不必设置采样器。

## TraceIDRatioBasedSampler

采样时，最常用的头部采样器是 TraceIdRatioBasedSampler。它确定地采样您作为参数传
入的跟踪的百分比。

### 环境变量

你可以用环境变量配置 TraceIdRatioBasedSampler:

```shell
export OTEL_TRACES_SAMPLER="traceidratio"
export OTEL_TRACES_SAMPLER_ARG="0.1"
```

This tells the SDK to sample spans such that only 10% of traces get created.

### Node.js

您还可以在代码中配置 TraceIdRatioBasedSampler。下面是 Node.js 的一个例子:

=== "TypeScript"

    ```ts
    import { TraceIdRatioBasedSampler } from '@opentelemetry/sdk-trace-node';

    const samplePercentage = 0.1;

    const sdk = new NodeSDK({
      // Other SDK configuration parameters go here
      sampler: new TraceIdRatioBasedSampler(samplePercentage),
    });
    ```

=== "JavaScript"

    ```js
    const { TraceIdRatioBasedSampler } = require('@opentelemetry/sdk-trace-node');

    const samplePercentage = 0.1;

    const sdk = new NodeSDK({
      // Other SDK configuration parameters go here
      sampler: new TraceIdRatioBasedSampler(samplePercentage),
    });
    ```

### Browser

You can also configure the TraceIdRatioBasedSampler in code. Here's an example
for browser apps:

=== "TypeScript"

    ```ts
    import {
      WebTracerProvider,
      TraceIdRatioBasedSampler,
    } from '@opentelemetry/sdk-trace-web';

    const samplePercentage = 0.1;

    const provider = new WebTracerProvider({
      sampler: new TraceIdRatioBasedSampler(samplePercentage),
    });
    ```

=== "JavaScript"

    ```js
    const {
      WebTracerProvider,
      TraceIdRatioBasedSampler,
    } = require('@opentelemetry/sdk-trace-web');

    const samplePercentage = 0.1;

    const provider = new WebTracerProvider({
      sampler: new TraceIdRatioBasedSampler(samplePercentage),
    });
    ```

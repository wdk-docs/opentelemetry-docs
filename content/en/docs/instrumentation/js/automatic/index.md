---
title: 自动 Instrumentation
linkTitle: 自动
weight: 20
---

自动检测提供了一种方法来检测任何 Node.js 应用程序，并从许多流行的库和框架中捕获
遥测数据，而无需任何代码更改。

## 设置

运行以下命令安装相应的包。

```shell
npm install --save @opentelemetry/api
npm install --save @opentelemetry/auto-instrumentations-node
```

`@opentelemetry/api` and `@opentelemetry/auto-instrumentations-node`包安装
api、SDK 和插装工具。

## 配置模块

该模块具有高度可配置性。

一种选择是通过使用`env` 从 CLI 设置环境变量来配置模块:

```shell
env OTEL_TRACES_EXPORTER=otlp OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=your-endpoint \
node --require @opentelemetry/auto-instrumentations-node/register app.js
```

或者，你可以使用 `export` 来设置环境变量:

```shell
export OTEL_TRACES_EXPORTER="otlp"
export OTEL_METRICS_EXPORTER="otlp"
export OTEL_EXPORTER_OTLP_ENDPOINT="your-endpoint"
export OTEL_NODE_RESOURCE_DETECTORS="env,host,os"
export OTEL_SERVICE_NAME="your-service-name"
export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"
node app.js
```

默认情况下，使用所有 SDK 资源检测器。您可以使用环境变
量`OTEL_NODE_RESOURCE_DETECTORS` 来只启用某些检测器，或者完全禁用它们。

要查看配置选项的全部范围，请参见[模块配置](module-config).

## 支持的库和框架

许多流行的 Node.js 库都是自动检测的。有关完整列表，请参
见[支持的仪器](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/metapackages/auto-instrumentations-node#supported-instrumentations).

## 故障排除

您可以通过将`OTEL_LOG_LEVEL`环境变量设置为以下其中一个来设置日志级别:

- `none`
- `error`
- `warn`
- `info`
- `debug`
- `verbose`
- `all`
- 默认级别是 `info`。

!!! NOTES

    - 在生产环境中，建议将 `OTEL_LOG_LEVEL` 设置为 `info`。
    - 无论环境或调试级别如何，日志总是被发送到控制台。
    - 调试日志非常冗长，可能会对应用程序的性能产生负面影响。只在需要时启用调试日志。

---
title: 自动插装配置
linkTitle: Configuration
weight: 10
---

这个模块可以通过设
置[环境变量](../../../../specs/otel/configuration/sdk-environment-variables/)进
行高度配置。自动检测行为的许多方面可以根据您的需要进行配置，例如资源检测器、导出
器、跟踪上下文传播头等等。

## 自动插装配置

SDK 和导出器配置可以使用环境变量进行设置。更多信息可以
在[这里](../../../../concepts/sdk-configuration/)找到。

## SDK 资源检测器配置

默认情况下，该模块将启用所有 SDK 资源检测器。你可以使
用`OTEL_NODE_RESOURCE_DETECTORS`环境变量来只启用某些检测器，或者完全禁用它们:

- `env`
- `host`
- `os`
- `process`
- `container`
- `alibaba`
- `aws`
- `gcp`
- `all` - 启用所有资源检测器
- `none` - 禁用资源检测

例如，要只启用`env` 和 `host`探测器，你可以设置:

```shell
OTEL_NODE_RESOURCE_DETECTORS=env,host
```

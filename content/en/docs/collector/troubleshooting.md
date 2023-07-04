---
title: 故障排除
description: 对采集器故障处理的建议
weight: 25
---

本页描述对 OpenTelemetry Collector 的运行状况或性能进行故障排除时的一些选项。
Collector 为调试问题提供了各种度量、日志和扩展。

## 发送测试数据

对于某些类型的问题，特别是验证配置和调试网络问题，将少量数据发送到配置为输出到本
地日志的收集器可能会有所帮助。详情请参
阅[本地出口商](https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/troubleshooting.md#local-exporters).

## 调试复杂管道的检查表

当遥测数据流经多个收集器和网络时，很难隔离问题。对于通过采集器或遥测管道中的其他
组件的遥测数据的每个"hop"，重要的是要验证以下内容:

- 采集器日志中是否有错误信息?
- 遥测技术是如何被吸收到这个组件中的?
- 该组件如何修改遥测(即采样，编辑)?
- 如何从该组件导出遥测数据?
- 遥测采用什么格式?
- 下一跳是如何配置的?
- 是否有阻止数据进出的网络策略?

### 更多的

有关详细建议(包括常见问题)，请参阅 Collector repo 中
的[故障排除](https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/troubleshooting.md)。

# 长期路线图

这份长期路线图(草案)是一份反映我们当前愿望的愿景文件。并不是承诺要实现此路线图中
列出的所有内容。本文档的主要目的是确保所有贡献者工作一致。随着时间的推移，我们的
愿景发生了变化，维护人员保留在路线图中添加、修改和删除项目的权利。

| Description                                                    | Status      | Links                                                                                                            |
| -------------------------------------------------------------- | ----------- | ---------------------------------------------------------------------------------------------------------------- |
| **测试**                                                       |
| 度量正确性测试                                                 | Done        | [#652](https://github.com/open-telemetry/opentelemetry-collector/issues/652)                                     |
|                                                                |
| **新格式**                                                     |
| 完整的 OTLP/HTTP 支持                                          | Done        | [#882](https://github.com/open-telemetry/opentelemetry-collector/issues/882)                                     |
| 为所有主核心处理器(属性、批处理、k8sattributes 等)添加日志支持 | Done        |
|                                                                |
| **5 最小值**                                                   |
| 针对大多数常见目标的发行包(例如 Docker、RPM、Windows 等)       | Done        | https://github.com/open-telemetry/opentelemetry-collector-releases/releases                                      |
| 检测和收集 AWS 上的环境指标和标签                              | Beta        | https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor |
| k8s 遥测检测与采集                                             | Beta        | https://pkg.go.dev/github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor    |
| 主机度量集合                                                   | Beta        | https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/hostmetricsreceiver         |
| 支持更多特定于应用程序的度量集合(例如 Kafka, Hadoop 等)        | In Progress | https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver                             |
|                                                                |
| **其他功能**                                                   |
| 安全停机(管道排水)                                             | Done        | [#483](https://github.com/open-telemetry/opentelemetry-collector/issues/483)                                     |
| 默认情况下，弃用队列重试处理器并启用每个导出器的队列           | Done        | [#1721](https://github.com/open-telemetry/opentelemetry-collector/issues/1721)                                   |

此时，OpenTelemetry Collector SIG 的大部分工作都集中在跨各种包实现 GA 状态上。请
参阅[GA 路线图](ga-roadmap.md)文档中的其他详细信息。

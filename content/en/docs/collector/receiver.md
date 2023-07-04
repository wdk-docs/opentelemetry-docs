# 接收器

接收器是数据进入收集器的方式。通常，接收器接受指定格式的数据，将其转换为内部格式
，并将其传递给适用管道中定义
的[处理器](./processor.md)和[导出器](./exporters.md)。

可用的跟踪接收器(按字母顺序排序):

- [OTLP Receiver](otlpreceiver/README.md) 可用的度量接收器(按字母顺序排序):
- [OTLP Receiver](otlpreceiver/README.md) 可用的日志接收器(按字母顺序排序):
- [OTLP Receiver](otlpreceiver/README.md)
  [contrib repository](https://github.com/open-telemetry/opentelemetry-collector-contrib)有
  更多的接收器，可以添加到自定义构建的收集器中。

## 配置接收器

接收器通过 YAML 在顶层的`receivers`标签下配置。必须至少有一个启用的接收器，才能
将配置视为有效。

下面是`examplereceiver`的配置示例。

```yaml
receivers:
  # Receiver 1.
  # <receiver type>:
  examplereceiver:
    # <setting one>: <value one>
    endpoint: 1.2.3.4:8080
    # ...
  # Receiver 2.
  # <receiver type>/<name>:
  examplereceiver/settings:
    # <setting two>: <value two>
    endpoint: 0.0.0.0:9211
```

接收器实例在配置的其他部分(如管道)中以其全名引用。全名由接收器类型、'/'和配置中
附加到接收器类型的名称组成。所有接收器的全名必须是唯一的。

对于上面的例子:

- 接收器 1 的全名为 `examplereceiver`.
- 接收器 2 的全名为 `examplereceiver/settings`.

接收器在添加到管道时启用。例如:

```yaml
service:
  pipelines:
    # Valid pipelines are: traces, metrics or logs
    # Trace pipeline 1.
    traces:
      receivers: [examplereceiver, examplereceiver/settings]
      processors: []
      exporters: [exampleexporter]
    # Trace pipeline 2.
    traces/another:
      receivers: [examplereceiver, examplereceiver/settings]
      processors: []
      exporters: [exampleexporter]
```

> 每个管道必须至少启用一个接收器才能成为有效配置。

---
title: 构建自定义收集器
weight: 29
---

如果您计划构建和调试自定义收集器接收器、处理器、扩展或导出程序，您将需要自己的收
集器实例。这将允许您在最喜欢的 Golang IDE 中直接启动和调试 OpenTelemetry
Collector 组件。

以这种方式进行组件开发的另一个有趣的方面是，您可以使用 IDE 中的所有调试特性(堆栈
跟踪是很好的老师!)来理解 Collector 本身如何与组件代码交互。

OpenTelemetry 社区开发了一个名为[OpenTelemetry Collector builder][ocb](或简称为
“ocb”)的工具来帮助人们组装他们自己的发行版，使构建一个包含他们自定义组件和公开可
用组件的发行版变得容易。

作为该过程的一部分，“构建器”将生成收集器的源代码，您可以使用它来帮助构建和调试您
自己的自定义组件，所以让我们开始吧。

## 步骤 1 - 安装构建器

`ocb`二进制文件可以从[OpenTelemetry Collector releases][releases]下载。您将在页
面底部找到资产列表。资产是根据操作系统和芯片组命名的，所以请下载适合您配置的资产
。

二进制文件有一个很长的名字，所以你可以简单地把它重命名为' ocb ';如果你运行的是
Linux 或 macOS，你还需要提供二进制文件的执行权限。

打开你的终端，输入以下命令来完成这两个操作:

```cmd
mv ocb_{{% param collectorVersion %}}_darwin_amd64 ocb
chmod 777 ocb
```

为了确保' ocb '已经准备好使用，进入你的终端并输入' ./ocb help '，一旦你按下
enter 键，你应该会在控制台中看到' help '命令的输出。

## 步骤 2 -创建构建器清单文件

The builder's `manifest` file is a `yaml` where you pass information about the
code generation and compile process combined with the components that you would
like to add to your Collector's distribution.

The `manifest` starts with a map named `dist` which contains tags to help you
configure the code generation and compile process. In fact, all the tags for
`dist` are the equivalent of the `ocb` command line `flags`.

Here are the tags for the `dist` map:

| Tag              | Description                                                                                        | Optional | Default Value                                                                     |
| ---------------- | -------------------------------------------------------------------------------------------------- | -------- | --------------------------------------------------------------------------------- |
| module:          | The module name for the new distribution, following Go mod conventions. Optional, but recommended. | Yes      | `go.opentelemetry.io/collector/cmd/builder`                                       |
| name:            | The binary name for your distribution                                                              | Yes      | `otelcol-custom`                                                                  |
| description:     | A long name for the application.                                                                   | Yes      | `Custom OpenTelemetry Collector distribution`                                     |
| otelcol_version: | The OpenTelemetry Collector version to use as base for the distribution.                           | Yes      | `{{% param collectorVersion %}}`                                                  |
| output_path:     | The path to write the output (sources and binary).                                                 | Yes      | `/var/folders/86/s7l1czb16g124tng0d7wyrtw0000gn/T/otelcol-distribution3618633831` |
| version:         | The version for your custom OpenTelemetry Collector.                                               | Yes      | `1.0.0`                                                                           |
| go:              | Which Go binary to use to compile the generated sources.                                           | Yes      | go from the PATH                                                                  |

As you can see on the table above, all the `dist` tags are optional, so you will
be adding custom values for them depending if your intentions to make your
custom Collector distribution available for consumption by other users or if you
are simply leveraging the `ocb` to bootstrap your component development and
testing environment.

For this tutorial, you will be creating a Collector's distribution to support
the development and testing of components.

Go ahead and create a manifest file named `builder-config.yaml` with the
following content:

> builder-config.yaml

```yaml
dist:
  name: otelcol-dev
  description: Basic OTel Collector distribution for Developers
  output_path: ./otelcol-dev
```

Now you need to add the modules representing the components you want to be
incorporated in this custom Collector distribution. Take a look at the
[ocb configuration documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder#configuration)
to understand the different modules and how to add the components.

We will be adding the following components to our development and testing
collector distribution:

- Exporters: Jaeger and Logging
- Receivers: OTLP
- Processors: Batch

Here is what my `builder-config.yaml` manifest file looks after adding the
modules for the components above:

<!-- prettier-ignore -->
```yaml
dist:
  name: otelcol-dev
  description: Basic OTel Collector distribution for Developers
  output_path: ./otelcol-dev
  otelcol_version: {{% param collectorVersion %}}

exporters:
  - gomod:
      go.opentelemetry.io/collector/exporter/loggingexporter v{{% param collectorVersion %}}
  - gomod:
      github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter
      v{{% param collectorVersion %}}

processors:
  - gomod:
      go.opentelemetry.io/collector/processor/batchprocessor v{{% param collectorVersion %}}

receivers:
  - gomod:
      go.opentelemetry.io/collector/receiver/otlpreceiver v{{% param collectorVersion %}}
```

## 步骤 3 -生成代码并构建收集器的分发。

All you need now is to let the `ocb` do it's job, so go to your terminal and
type the following command:

```cmd
./ocb --config builder-config.yaml
```

If everything went well, here is what the output of the command should look
like:

```nocode
2022-06-13T14:25:03.037-0500	INFO	internal/command.go:85	OpenTelemetry Collector distribution builder	{"version": "{{% param collectorVersion %}}", "date": "2023-01-03T15:05:37Z"}
2022-06-13T14:25:03.039-0500	INFO	internal/command.go:108	Using config file	{"path": "builder-config.yaml"}
2022-06-13T14:25:03.040-0500	INFO	builder/config.go:99	Using go	{"go-executable": "/usr/local/go/bin/go"}
2022-06-13T14:25:03.041-0500	INFO	builder/main.go:76	Sources created	{"path": "./otelcol-dev"}
2022-06-13T14:25:03.445-0500	INFO	builder/main.go:108	Getting go modules
2022-06-13T14:25:04.675-0500	INFO	builder/main.go:87	Compiling
2022-06-13T14:25:17.259-0500	INFO	builder/main.go:94	Compiled	{"binary": "./otelcol-dev/otelcol-dev"}
```

As defined in the `dist` section of your config file, you now have a folder
named `otelcol-dev` containing all the source code and the binary for your
Collector's distribution.

You can now use the generated code to bootstrap your component development
projects and easily build and distribute your own collector distribution with
your components.

[ocb]:
  https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder
[releases]: https://github.com/open-telemetry/opentelemetry-collector/releases

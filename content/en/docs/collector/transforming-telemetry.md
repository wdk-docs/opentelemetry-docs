---
title: 改变遥测
weight: 26
---

OpenTelemetry Collector 是在将数据发送到供应商或其他系统之前转换数据的方便场所。
这通常是出于数据质量、治理、成本和安全性的考虑。

从[Collector Contrib 存储库](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor)获
得的处理器支持对度量、跨度和日志数据进行数十种不同的转换。下面几节提供一些基本示
例，介绍如何开始使用一些常用的处理器。

处理器的配置，特别是高级转换，可能会对收集器性能产生重大影响。

## 基本过滤

**处理器**:
[filter processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/filterprocessor)

过滤处理器允许用户根据“包括”或“排除”规则过滤遥测数据。包含规则用于定义“允许列表
”，其中任何与包含规则不匹配的内容都将从收集器中删除。排除规则用于定义“拒绝列表
”，其中匹配规则的遥测数据将从收集器中删除。

例如，_只_ 允许来自 app1、app2 和 app3 服务的 span 数据，删除来自所有其他服务的
数据:

```yaml
processors:
  filter/allowlist:
    spans:
      include:
        match_type: strict
        services:
          - app1
          - app2
          - app3
```

为了只阻止来自名为 development 的服务的跨度，而允许所有其他跨度，使用了一个排除
规则:

```yaml
processors:
  filter/denylist:
    spans:
      exclude:
        match_type: strict
        services:
          - development
```

[filter processor docs](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/filterprocessor)有
更多的例子，包括对日志和指标进行过滤。

## 添加/删除属性

**Processor**:
[attributes processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor)
or
[resource processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor)

属性处理器可用于更新、插入、删除或替换度量或跟踪上的现有属性。例如，下面的配置将
一个名为 account_id 的属性添加到所有 span 中:

```yaml
processors:
  attributes/accountid:
    actions:
      - key: account_id
        value: 2245
        action: insert
```

资源处理器具有相同的配置，但只适用
于[资源属性](/docs/specs/otel/resource/semantic_conventions/)。使用资源处理器修
改与遥测相关的基础设施元数据。例如，插入 Kubernetes 集群名称:

```yaml
processors:
  resource/k8s:
    attributes:
      - key: k8s.cluster.name
        from_attribute: k8s-cluster
        action: insert
```

## 重命名度量或度量标签

**Processor:**
[metrics transform processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/metricstransformprocessor)

The
[metrics transform processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/metricstransformprocessor)
shares some functionality with the
[attributes processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor),
but also supports renaming and other metric-specific functionality.

```yaml
processors:
  metricstransform/rename:
    transforms:
      include: system.cpu.usage
      action: update
      new_name: system.cpu.usage_time
```

[度量转换处理器](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/metricstransformprocessor)还
支持正则表达式，以同时将转换规则应用于多个度量名称或度量标签。下面的示例将所有指
标的 cluster_name 重命名为 cluster-name:

```yaml
processors:
  metricstransform/clustername:
    transforms:
      - include: ^.*$
        match_type: regexp
        action: update
        operations:
          - action: update_label
            label: cluster_name
            new_label: cluster-name
```

## 用资源属性丰富遥测

**Processor**:
[资源检测处理器](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor)
和
[k8sattributes 处理器](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/k8sattributesprocessor)

这些处理器可用于用相关基础设施元数据丰富遥测，以帮助团队快速识别底层基础设施何时
影响服务运行状况或性能。

资源检测处理器将相关的云或主机级信息添加到遥测中:

```yaml
processors:
  resourcedetection/system:
    # Modify the list of detectors to match the cloud environment
    detectors: [env, system, gcp, ec2, azure]
    timeout: 2s
    override: false
```

类似地，K8s 处理器使用相关的 Kubernetes 元数据(如 pod 名称、节点名称或工作负载名
称)丰富遥测。收集器 pod 必须配置为具有对某些 Kubernetes RBAC api 的读访问权限，
文档中
有[这里](https://pkg.go.dev/github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor#hdr-RBAC)。
要使用默认选项，可以配置一个空块:

```yaml
processors:
  k8sattributes/default:
```

## 先进的转换

更高级的属性转换也可以
在[转换处理程序](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/transformprocessor)中
使用。转换处理器允许最终用户使
用[OpenTelemetry 转换语言](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/ottl)在
度量、日志和跟踪上指定转换。

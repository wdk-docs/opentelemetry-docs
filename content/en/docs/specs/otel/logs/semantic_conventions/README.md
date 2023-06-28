# 日志属性语义约定

!!! note

    语义约定正在转移到一个[新的位置](http://github.com/open-telemetry/semantic-conventions).

不允许对本文档进行任何修改。

**Status**: [Experimental](../../document-status.md)

定义了以下日志的语义约定:

- [General](general.md): 可用于描述日志记录的一般语义属性。
- [Log Media](media.md): 可用于描述日志源的语义属性。

定义了事件的以下语义约定:

- [Events](events.md): 使用日志数据模型表示事件时必须使用的语义属性。

除了日志的语义约定
，[trace](../../trace/semantic_conventions/README.md)和[metrics](../../metrics/semantic_conventions/README.md)，
OpenTelemetry 还用自己
的[资源语义约定](../../resource/semantic_conventions/README.md)定义了覆
盖[Resources](../../resource/sdk.md)的概念。

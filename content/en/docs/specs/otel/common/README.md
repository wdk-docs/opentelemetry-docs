<!--- Hugo front matter used to generate the website version of this page:
aliases: [/docs/reference/specification/common/common]
--->

# 通用规范概念

**Status**: [Stable, Feature-freeze](../document-status.md)

## Attribute

<a id="attributes"></a>

一个“属性”是一个键值对，它必须具有以下属性:

- 属性键必须是一个非`null`和非空字符串。
- 属性值为:
  - 基本类型:字符串、布尔值、双精度浮点数(IEEE 754-1985)或有符号 64 位整数。
  - 原始类型值的数组。数组必须是同构的，也就是说，它不能包含不同类型的值。

对于原生不支持非字符串值的协议，非字符串值应该表示为 json 编码的字符串。例如，表
达式`int64(100)`将被编码为`100`，`float64(1.5)`将被编码为`1.5`，任何类型的空数组
将被编码为`[]`。

Attribute values expressing a numerical value of zero, an empty string, or an
empty array are considered meaningful and MUST be stored and passed on to
processors / exporters.

Attribute values of `null` are not valid and attempting to set a `null` value is
undefined behavior.

`null` values SHOULD NOT be allowed in arrays. However, if it is impossible to
make sure that no `null` values are accepted (e.g. in languages that do not have
appropriate compile-time type checking), `null` values within arrays MUST be
preserved as-is (i.e., passed on to span processors / exporters as `null`). If
exporters do not support exporting `null` values, they MAY replace those values
by 0, `false`, or empty strings. This is required for map/dictionary structures
represented as two arrays with indices that are kept in sync (e.g., two
attributes `header_keys` and `header_values`, both containing an array of
strings to represent a mapping `header_keys[i] -> header_values[i]`).

See [Attribute Naming](attribute-naming.md) for naming guidelines.

See [Requirement Level](attribute-requirement-level.md) for requirement levels
guidelines.

See [this document](attribute-type-mapping.md) to find out how to map values
obtained outside OpenTelemetry into OpenTelemetry attribute values.

### 属性限制

执行错误的代码可能会产生意想不到的属性。如果对属性没有限制，它们会迅速耗尽可用内
存，导致难以安全恢复的崩溃。

默认情况下，SDK 应该按照下面的[可配置参数](#configurable-parameters)列表应用截断
。

如果 SDK 提供了一种方法:

- 设置一个属性值长度限制，这样对于每个属性值:
  - 如果它是一个字符串，如果它超过了这个限制(将其中的任何字符计数为 1)，sdk 必须
    截断该值，使其长度最多等于限制，
  - 如果它是字符串数组，则分别对每个值应用上述规则。
  - 否则一个值绝对不能被截断;
- 设置唯一属性键的限制如下:
  - 对于每个唯一的属性键，添加它将导致超过限制，SDK 必须丢弃该键/值对。

There MAY be a log emitted to indicate to the user that an attribute was
truncated or discarded. To prevent excessive logging, the log MUST NOT be
emitted more than once per record on which an attribute is set.

If the SDK implements the limits above, it MUST provide a way to change these
limits programmatically. Names of the configuration options SHOULD be the same
as in the list below.

An SDK MAY implement model-specific limits, for example
`SpanAttributeCountLimit` or `LogRecordAttributeCountLimit`. If both a general
and a model-specific limit are implemented, then the SDK MUST first attempt to
use the model-specific limit, if it isn't set, then the SDK MUST attempt to use
the general limit. If neither are defined, then the SDK MUST try to use the
model-specific limit default value, followed by the global limit default value.

#### 可配置参数

- `AttributeCountLimit` (Default=128) - 每条记录允许的最大属性计数;
- `AttributeValueLengthLimit` (Default=Infinity) - 允许的最大属性值长度;

#### 免税实体

资源属性应该不受上面描述的限制，因为资源不容易受到导致过多属性计数或大小的场景(
自动检测)的影响。资源每批只发送一次，而不是每个跨度发送一次，因此在资源上拥有更
多/更大的属性相对更便宜。资源在设计上也是不可变的，它们通常与限制一起传递给
TracerProvider。这使得为资源实现属性限制变得很尴尬。

属于度量的属性此时不受上述限制的限制，
如[度量属性限制](../metrics/sdk.md#attribute-limits)中所述。

## 属性集合

[Resources](../resource/sdk.md), Metrics
[data points](../metrics/data-model.md#metric-points),
[Spans](../trace/api.md#set-attributes), Span
[Events](../trace/api.md#add-events), Span
[Links](../trace/api.md#specifying-links) and
[Log Records](../logs/data-model.md) may contain a collection of attributes. The
keys in each such collection are unique, i.e. there MUST NOT exist more than one
key-value pair with the same key. The enforcement of uniqueness may be performed
in a variety of ways as it best fits the limitations of the particular
implementation.

Normally for the telemetry generated using OpenTelemetry SDKs the attribute
key-value pairs are set via an API that either accepts a single key-value pair
or a collection of key-value pairs. Setting an attribute with the same key as an
existing attribute SHOULD overwrite the existing attribute's value. See for
example Span's [SetAttribute](../trace/api.md#set-attributes) API.

A typical implementation of [SetAttribute](../trace/api.md#set-attributes) API
will enforce the uniqueness by overwriting any existing attribute values pending
to be exported, so that when the Span is eventually exported the exporters see
only unique attributes. The OTLP format in particular requires that exported
Resources, Spans, Metric data points and Log Records contain only unique
attributes.

Some other implementations may use a streaming approach where every
[SetAttribute](../trace/api.md#set-attributes) API call immediately results in
that individual attribute value being exported using a streaming wire protocol.
In such cases the enforcement of uniqueness will likely be the responsibility of
the recipient of this data.

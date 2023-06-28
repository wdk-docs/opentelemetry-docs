# 资源语义约定

!!! note

    语义约定正在转移到一个[新的位置](http://github.com/open-telemetry/semantic-conventions).

不允许对本文档进行任何修改。

**Status**: [Mixed](../../document-status.md)

本文档定义了资源的标准属性。这些属性通常在[资源](../sdk.md)中使用，也建议在需要
以一致的方式描述资源的任何其他地方使用。这些属性的大部分都继承
自[OpenCensus 资源标准](https://github.com/census-instrumentation/opencensus-specs/blob/master/resource/StandardResources.md).

## 待办事项

- 增加更多计算单元:AppEngine 单元等。
- 添加 Web 浏览器。
- 决定是否只使用小写字符串。
- 考虑为每个属性和属性组合添加可选/必需属性(例如，当提供 k8 资源时，可能需要所有
  k8)。

## Document Conventions

**Status**: [Stable](../../document-status.md)

Attributes are grouped logically by the type of the concept that they described.
Attributes in the same group have a common prefix that ends with a dot. For
example all attributes that describe Kubernetes properties start with "k8s."

See [Attribute Requirement Levels](../../common/attribute-requirement-level.md)
for details on when attributes should be included.

## Attributes with Special Handling

**Status**: [Stable](../../document-status.md)

Given their significance some resource attributes are treated specifically as
described below.

### Semantic Attributes with Dedicated Environment Variable

These are the attributes which MAY be configurable via a dedicated environment
variable as specified in
[OpenTelemetry Environment Variable Specification](../../configuration/sdk-environment-variables.md):

- [`service.name`](#service)

## Semantic Attributes with SDK-provided Default Value

These are the attributes which MUST be provided by the SDK as specified in the
[Resource SDK specification](../sdk.md#sdk-provided-resource-attributes):

- [`service.name`](#service)
- [`telemetry.sdk` group](#telemetry-sdk)

## Service

**Status**: [Stable](../../document-status.md)

**type:** `service`

**Description:** A service instance.

<!-- semconv service -->

| Attribute      | Type   | Description                      | Examples       | Requirement Level |
| -------------- | ------ | -------------------------------- | -------------- | ----------------- |
| `service.name` | string | Logical name of the service. [1] | `shoppingcart` | Required          |

**[1]:** MUST be the same for all instances of horizontally scaled services. If
the value was not specified, SDKs MUST fallback to `unknown_service:`
concatenated with [`process.executable.name`](process.md#process), e.g.
`unknown_service:bash`. If `process.executable.name` is not available, the value
MUST be set to `unknown_service`.

<!-- endsemconv -->

## Service (Experimental)

**Status**: [Experimental](../../document-status.md)

**type:** `service`

**Description:** Additions to service instance.

<!-- semconv service_experimental -->

| Attribute             | Type   | Description                                              | Examples                                                          | Requirement Level |
| --------------------- | ------ | -------------------------------------------------------- | ----------------------------------------------------------------- | ----------------- |
| `service.namespace`   | string | A namespace for `service.name`. [1]                      | `Shop`                                                            | Recommended       |
| `service.instance.id` | string | The string ID of the service instance. [2]               | `my-k8s-pod-deployment-1`; `627cc493-f310-47de-96bd-71410b7dec09` | Recommended       |
| `service.version`     | string | The version string of the service API or implementation. | `2.0.0`                                                           | Recommended       |

**[1]:** A string value having a meaning that helps to distinguish a group of
services, for example the team name that owns a group of services.
`service.name` is expected to be unique within the same namespace. If
`service.namespace` is not specified in the Resource then `service.name` is
expected to be unique for all services that have no explicit namespace defined
(so the empty/unspecified namespace is simply one more valid namespace).
Zero-length namespace string is assumed equal to unspecified namespace.

**[2]:** MUST be unique for each instance of the same
`service.namespace,service.name` pair (in other words
`service.namespace,service.name,service.instance.id` triplet MUST be globally
unique). The ID helps to distinguish instances of the same service that exist at
the same time (e.g. instances of a horizontally scaled service). It is
preferable for the ID to be persistent and stay the same for the lifetime of the
service instance, however it is acceptable that the ID is ephemeral and changes
during important lifetime events for the service (e.g. service restarts). If the
service has no inherent unique ID that can be used as the value of this
attribute it is recommended to generate a random Version 1 or Version 4 RFC 4122
UUID (services aiming for reproducible UUIDs may also use Version 5, see RFC
4122 for more recommendations).

<!-- endsemconv -->

Note: `service.namespace` and `service.name` are not intended to be concatenated
for the purpose of forming a single globally unique name for the service. For
example the following 2 sets of attributes actually describe 2 different
services (despite the fact that the concatenation would result in the same
string):

```
# Resource attributes that describes a service.
namespace = Company.Shop
service.name = shoppingcart
```

```
# Another set of resource attributes that describe a different service.
namespace = Company
service.name = Shop.shoppingcart
```

## Telemetry SDK

**Status**: [Stable](../../document-status.md)

**type:** `telemetry.sdk`

**Description:** The telemetry SDK used to capture data recorded by the
instrumentation libraries.

<!-- semconv telemetry -->

| Attribute                | Type   | Description                                         | Examples        | Requirement Level |
| ------------------------ | ------ | --------------------------------------------------- | --------------- | ----------------- |
| `telemetry.sdk.name`     | string | The name of the telemetry SDK as defined above. [1] | `opentelemetry` | Required          |
| `telemetry.sdk.language` | string | The language of the telemetry SDK.                  | `cpp`           | Required          |
| `telemetry.sdk.version`  | string | The version string of the telemetry SDK.            | `1.2.3`         | Required          |

**[1]:** The OpenTelemetry SDK MUST set the `telemetry.sdk.name` attribute to
`opentelemetry`. If another SDK, like a fork or a vendor-provided
implementation, is used, this SDK MUST set the `telemetry.sdk.name` attribute to
the fully-qualified class or module name of this SDK's main entry point or
another suitable identifier depending on the language. The identifier
`opentelemetry` is reserved and MUST NOT be used in this case. All custom
identifiers SHOULD be stable across different versions of an implementation.

`telemetry.sdk.language` has the following list of well-known values. If one of
them applies, then the respective value MUST be used, otherwise a custom value
MAY be used.

| Value    | Description |
| -------- | ----------- |
| `cpp`    | cpp         |
| `dotnet` | dotnet      |
| `erlang` | erlang      |
| `go`     | go          |
| `java`   | java        |
| `nodejs` | nodejs      |
| `php`    | php         |
| `python` | python      |
| `ruby`   | ruby        |
| `rust`   | rust        |
| `swift`  | swift       |
| `webjs`  | webjs       |

<!-- endsemconv -->

## Telemetry SDK (Experimental)

**Status**: [Experimental](../../document-status.md)

**type:** `telemetry.sdk`

**Description:** Additions to the telemetry SDK.

<!-- semconv telemetry_experimental -->

| Attribute                | Type   | Description                                                    | Examples | Requirement Level |
| ------------------------ | ------ | -------------------------------------------------------------- | -------- | ----------------- |
| `telemetry.auto.version` | string | The version string of the auto instrumentation agent, if used. | `1.2.3`  | Recommended       |

<!-- endsemconv -->

## Compute Unit

**Status**: [Experimental](../../document-status.md)

Attributes defining a compute unit (e.g. Container, Process, Function as a
Service):

- [Container](./container.md)
- [Function as a Service](./faas.md)
- [Process](./process.md)
- [Web engine](./webengine.md)

## Compute Instance

**Status**: [Experimental](../../document-status.md)

Attributes defining a computing instance (e.g. host):

- [Host](./host.md)

## Environment

**Status**: [Experimental](../../document-status.md)

Attributes defining a running environment (e.g. Operating System, Cloud, Data
Center, Deployment Service):

- [Operating System](./os.md)
- [Device](./device.md)
- [Cloud](./cloud.md)
- Deployment:
  - [Deployment Environment](./deployment_environment.md)
  - [Kubernetes](./k8s.md)
- [Browser](./browser.md)

## Version attributes

**Status**: [Stable](../../document-status.md)

Version attributes, such as `service.version`, are values of type `string`. They
are the exact version used to identify an artifact. This may be a semantic
version, e.g., `1.2.3`, git hash, e.g., `8ae73a`, or an arbitrary version
string, e.g., `0.1.2.20210101`, whatever was used when building the artifact.

## Cloud-Provider-Specific Attributes

**Status**: [Experimental](../../document-status.md)

Attributes that are only applicable to resources from a specific cloud provider.
Currently, these resources can only be defined for providers listed as a valid
`cloud.provider` in [Cloud](./cloud.md) and below. Provider-specific attributes
all reside in the `cloud_provider` directory. Valid cloud providers are:

- [Alibaba Cloud](https://www.alibabacloud.com/) (`alibaba_cloud`)
- [Amazon Web Services](https://aws.amazon.com/)
  ([`aws`](cloud_provider/aws/README.md))
- [Google Cloud Platform](https://cloud.google.com/)
  ([`gcp`](cloud_provider/gcp/README.md))
- [Microsoft Azure](https://azure.microsoft.com/) (`azure`)
- [Tencent Cloud](https://www.tencentcloud.com/) (`tencent_cloud`)
- [Heroku dyno](./cloud_provider/heroku.md)

# OpenTelemetry Collector

OpenTelemetry Collector 提供了一个与供应商无关的实现，用于接收、处理和导出遥测数
据。此外，为了支持开源遥测数据格式(例如 Jaeger、Prometheus 等)到多个开源或商业后
端，它消除了运行、操作和维护多个代理/收集器的需要。

目的:

- 可用性: 合理的默认配置，支持流行的协议，开箱即用的运行和收集。
- 高性能: 在各种负载和配置下都非常稳定和高性能。
- 可观测: 一个可观察服务的范例。
- 可扩展: 无需触及核心代码即可自定义。
- 统一性: 单个代码库，可作为支持跟踪、度量和日志的代理或收集器进行部署。

## 稳定的水平

收集器组件和实现处于稳定性的不同阶段，通常分为功能和配置两部分。每个组件的状态在
该组件的 README 文件中可用。虽然我们打算提供高质量的组件作为该存储库的一部分，但
我们承认，并非所有组件都已准备就绪。因此，每个组件应根据以下定义列出每个遥测信号
的当前稳定级别:

### Development

并非该组件的所有部分都已到位，而且它可能还不能作为任何发行版的一部分使用。应该报
告错误和性能问题，但组件所有者可能不会给予它们太多关注。你的反馈仍然是需要的，特
别是当它涉及到用户体验(配置选项，组件可观察性，技术实现细节，…)。根据情况的发展
，配置选项可能会经常失效。该组件不应在生产中使用。

### Alpha

The component is ready to be used for limited non-critical workloads and the
authors of this component would welcome your feedback. Bugs and performance
problems should be reported, but component owners might not work on them right
away. The configuration options might change often without backwards
compatibility guarantees.

### Beta

Same as Alpha, but the configuration options are deemed stable. While there
might be breaking changes between releases, component owners should try to
minimize them. A component at this stage is expected to have had exposure to
non-critical production workloads already during its **Alpha** phase, making it
suitable for broader usage.

### Stable

The component is ready for general availability. Bugs and performance problems
should be reported and there's an expectation that the component owners will
work on them. Breaking changes, including configuration options and the
component's output are not expected to happen without prior notice, unless under
special circumstances.

### Deprecated

The component is planned to be removed in a future version and no further
support will be provided. Note that new issues will likely not be worked on.
When a component enters "deprecated" mode, it is expected to exist for at least
two minor releases. See the component's readme file for more details on when a
component will cease to exist.

### Unmaintained

A component identified as unmaintained does not have an active code owner. Such
component may have never been assigned a code owner or a previously active code
owner has not responded to requests for feedback within 6 weeks of being
contacted. Issues and pull requests for unmaintained components will be labelled
as such. After 6 months of being unmaintained, these components will be removed
from official distribution. Components that are unmaintained are actively
seeking contributors to become code owners.

## 兼容性

当用作库时，OpenTelemetry Collector 尝试跟踪当前支持的 Go 版本，
如[由 Go 团队定义](https://go.dev/doc/devel/release#policy)。移除对不受支持的 Go
版本的支持并不会被认为是破坏性的改变。

从 Go 1.18 发布开始，OpenTelemetry Collector 对 Go 版本的支持将更新如下:

1. 新的 Go 小版本“N”发布后的第一个版本将为新的 Go 小版本添加构建和测试步骤。
2. Go 小版本“N”发布后的第一个版本将取消对 Go 小版本“N-2”的支持。

官方 OpenTelemetry Collector 发行版二进制文件可以使用任何受支持的 Go 版本构建。

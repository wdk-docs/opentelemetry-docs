---
title: 构建自定义身份验证器
weight: 30
---

OpenTelemetry Collector 允许接收方和导出方连接到身份验证器，提供了一种方法，既可
以在接收方一侧对传入连接进行身份验证，也可以在导出方一侧向传出请求添加身份验证数
据。

此机制是
在[' extensions '](https://pkg.go.dev/go.opentelemetry.io/collector/component#Extension)框
架之上实现的，本文档将指导您实现自己的身份验证器。如果您正在寻找有关如何使用现有
身份验证器的文档，请参阅入门页和您的身份验证器文档。您可以在本网站的注册表中找到
现有身份验证器的列表。

使用本指南了解如何构建自定义身份验证器的一般指导，并参考最新
的[API 参考指南](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth)了
解每种类型和函数的实际语义。

如果您在任何时候需要帮助，请加
入[CNCF Slack 工作区](https://slack.cncf.io)的[# opentelemetri -collector](https://cloud-native.slack.com/archives/C01N6P7KR6W)房
间。

## 体系结构

认证器是常规扩展，也满足与认证机制相关的一个或多个接口:

- [go.opentelemetry.io/collector/config/configauth/ServerAuthenticator](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth#ServerAuthenticator)
- [go.opentelemetry.io/collector/config/configauth/GRPCClientAuthenticator](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth#GRPCClientAuthenticator)
- [go.opentelemetry.io/collector/config/configauth/HTTPClientAuthenticator](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth#HTTPClientAuthenticator)

服务器身份验证器用于接收端，能够拦截 HTTP 和 gRPC 请求，而客户端身份验证器用于导
出端，能够向 HTTP 和 gRPC 请求添加身份验证数据。身份验证者可以同时实现这两个接口
，从而允许将扩展的单个实例用于传入和传出请求。请注意，用户可能仍然希望对传入和传
出请求使用不同的身份验证器，因此，不要让两端都需要使用您的身份验证器。

一旦认证器扩展在收集器发行版中可用，它就可以在配置文件中作为常规扩展引用:

```yaml
extensions:
  oidc:

receivers:
processors:
exporters:

service:
  extensions:
    - oidc
  pipelines:
    traces:
      receivers: []
      processors: []
      exporters: []
```

但是，验证器需要被消费组件引用才能有效。下面的示例显示了与上面相同的扩展名，现在
由名为' otlp/auth '的接收器使用:

```yaml
extensions:
  oidc:

receivers:
  otlp/auth:
    protocols:
      grpc:
        auth:
          authenticator: oidc

processors:
exporters:

service:
  extensions:
    - oidc
  pipelines:
    traces:
      receivers:
        - otlp/auth
      processors: []
      exporters: []
```

当需要给定验证器的多个实例时，它们可以有不同的名称:

```yaml
extensions:
  oidc/some-provider:
  oidc/another-provider:

receivers:
  otlp/auth:
    protocols:
      grpc:
        auth:
          authenticator: oidc/some-provider

processors:
exporters:

service:
  extensions:
    - oidc/some-provider
    - oidc/another-provider
  pipelines:
    traces:
      receivers:
        - otlp/auth
      processors: []
      exporters: []
```

### 服务器的身份验证器

服务器身份验证器本质上是带有“Authenticate”功能的扩展，接收有效负载头作为参数。如
果验证器能够验证传入的连接，它应该返回一个' nil '错误，如果不能，则返回具体的错
误。作为扩展，验证者应该确保
在[' Start '](https://pkg.go.dev/go.opentelemetry.io/collector/component#Component)阶
段初始化所需的所有资源，并期望在' Shutdown '时清理它们。

“Authenticate”调用是传入请求的热路径的一部分，将阻塞管道，因此请确保正确处理需要
进行的任何阻塞操作。具体地说，尊重上下文设定的最后期限，如果有的话。还要确保在扩
展中添加足够的可观察性，特别是以指标和跟踪的形式，这样用户就可以在错误率超过一定
水平时设置通知系统，并可以调试特定的故障。

### 客户端身份验证器

客户端验证器是实现以下一个或多个接口的验证器:

- [go.opentelemetry.io/collector/config/configauth/GRPCClientAuthenticator](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth#GRPCClientAuthenticator)
- [go.opentelemetry.io/collector/config/configauth/HTTPClientAuthenticator](https://pkg.go.dev/go.opentelemetry.io/collector/config/configauth#HTTPClientAuthenticator)

与服务器身份验证器类似，它们本质上是带有额外功能的扩展，每个扩展都接收一个对象，
该对象为身份验证器提供了将身份验证数据注入其中的机会。例如，HTTP 客户端验证器提
供一个[' http.RoundTripper '](https://pkg.go.dev/net/http#RoundTripper)，而 gRPC
客户端验证器可以产生一
个[' credentials.PerRPCCredentials '](https://pkg.go.dev/google.golang.org/grpc/credentials#PerRPCCredentials)。

## 向发行版添加自定义身份验证器

自定义身份验证器必须是与主收集器相同的二进制文件的一部分。在构建您自己的身份验证
器时，您可能还需要构建一个自定义发行版，或者为您的用户提供将您的扩展作为他们自己
的发行版的一部分来使用的方法。幸运的是，可以使
用[OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder)实
用程序来构建自定义发行版。

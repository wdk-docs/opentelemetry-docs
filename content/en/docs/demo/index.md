---
title: OpenTelemetry 演示文档
linkTitle: Demo
cascade:
  repo: https://github.com/open-telemetry/opentelemetry-demo
weight: 2
---

欢迎来到[OpenTelemetry Demo](/ecosystem/demo/)文档，它涵盖了如何安装和运行演示，
以及一些可以用来查看 OpenTelemetry 的场景。

## 运行 Demo

想要部署演示并查看实际情况吗?从这里开始。

- [Docker](docker-deployment/)
- [Kubernetes](kubernetes-deployment/)

## 语言特性参考

想要了解特定语言的检测是如何工作的?从这里开始。

| 语言          | 自动插装                             | 插装库                                                                                                      | 手动插装                                                                  |
| ------------- | ------------------------------------ | ----------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| .NET          |                                      | [Cart Service](services/cart/)                                                                              | [Cart Service](services/cart/)                                            |
| C++           |                                      |                                                                                                             | [货币服务](services/currency/)                                            |
| Erlang/Elixir |                                      | [特征标志服务](services/feature-flag/)                                                                      | [特征标志服务](services/feature-flag/)                                    |
| Go            |                                      | [会计服务](services/accounting/), [结账服务](services/checkout/), [产品目录服务](services/product-catalog/) | [结账服务](services/checkout/), [产品目录服务](services/product-catalog/) |
| Java          | [Ad 服务](services/ad/)              |                                                                                                             | [Ad 服务](services/ad/)                                                   |
| JavaScript    |                                      | [前端](services/frontend/)                                                                                  | [前端](services/frontend/), [支付服务](services/payment/)                 |
| Kotlin        |                                      | [欺诈侦测服务](services/fraud-detection/)                                                                   |                                                                           |
| PHP           |                                      | [报价服务](services/quote/)                                                                                 | [报价服务](services/quote/)                                               |
| Python        | [推荐服务](services/recommendation/) |                                                                                                             | [推荐服务](services/recommendation/)                                      |
| Ruby          |                                      | [电子邮件服务](services/email/)                                                                             | [电子邮件服务](services/email/)                                           |
| Rust          |                                      | [航运服务](services/shipping/)                                                                              | [航运服务](services/shipping/)                                            |

## 服务文档

关于 OpenTelemetry 如何在每个服务中部署的具体信息可以在这里找到:

- [Ad 服务](services/ad/)
- [Cart 服务](services/cart/)
- [结账服务](services/checkout/)
- [Email 服务](services/email/)
- [特征标志服务](services/feature-flag/)
- [前端](services/frontend/)
- [负载生成器](services/load-generator/)
- [支付服务](services/payment/)
- [产品目录服务](services/product-catalog/)
- [报价服务](services/quote/)
- [推荐服务](services/recommendation/)
- [航运服务](services/shipping/)

## 场景

如何用 OpenTelemetry 解决问题?这些场景引导您完成一些预配置的问题，并向您展示如何
解释 OpenTelemetry 数据以解决这些问题。

随着时间的推移，我们将添加更多的场景。

- 使用 Feature Flag 服务为产品 id: `OLJCESPC7Z`的 `GetProduct` 请求生成一
  个[产品目录错误](feature-flags)
- 发现内存泄漏并使用指标和跟踪对其进行诊断。(了解更多
  )(scenarios/recommendation-cache/)

## 参考

项目参考文档，如需求和特性矩阵。

- [体系结构](architecture/)
- [开发](development/)
- [特性标志参考](feature-flags/)
- [度量特征矩阵](metric-features/)
- [需求](./requirements/)
- [截图](screenshots/)
- [服务角色表](service-table/)
- [Span 属性参考](manual-span-attributes/)
- [测试](tests/)
- [追踪特征矩阵](trace-features/)

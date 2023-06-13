---
title: 按服务跟踪功能覆盖
linkTitle: Trace Feature Coverage
aliases: [/docs/demo/trace_service_features]
---

| 服务               | 语言            | 工具库 | 手创Span | Span数据充实 | RPC上下文传播 | Span链接 | 行李 | 资源发现 |
| ------------------ | --------------- | ------ | -------- | ------------ | ------------- | -------- | ---- | -------- |
| Accounting Service | Go              | 🚧      | 🚧        | 🚧            | 🚧             | 🚧        | 🚧    | ✅        |
| Ad                 | Java            | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Cart               | .NET            | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | ✅        |
| Checkout           | Go              | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | ✅        |
| Currency           | C++             | 🔕      | ✅        | ✅            | ✅             | 🔕        | 🔕    | 🚧        |
| Email              | Ruby            | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Feature Flag       | Erlang / Elixir | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Fraud Detection    | Kotlin          | ✅      | 🚧        | 🚧            | 🚧             | 🚧        | 🚧    | 🚧        |
| Frontend           | JavaScript      | ✅      | ✅        | ✅            | 🔕             | ✅        | ✅    | ✅        |
| Payment            | JavaScript      | ✅      | ✅        | ✅            | 🔕             | 🔕        | ✅    | ✅        |
| Product Catalog    | Go              | ✅      | 🔕        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Quote Service      | PHP             | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Recommendation     | Python          | ✅      | ✅        | ✅            | 🔕             | 🔕        | 🔕    | 🚧        |
| Shipping           | Rust            | 🔕      | ✅        | ✅            | ✅             | 🔕        | 🔕    | 🚧        |

Emoji传奇:

- 完成: ✅
- 不适用: 🔕
- 不在场(尚未): 🚧

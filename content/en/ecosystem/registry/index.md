---
title: 注册表
description: >-
  Find libraries, plugins, integrations, and other useful tools for extending
  OpenTelemetry.
aliases: [/registry/*]
type: default
layout: registry
outputs: [html, json]
body_class: registry
weight: 20
---

{{% blocks/lead color="white" %}}

# {{% param title %}}

{{% param description %}}

{{% /blocks/lead %}}

{{% blocks/section color="dark" %}}

## 你需要什么?

开放遥测注册表允许您在开放遥测生态系统中搜索仪器库、跟踪器实现、实用程序和其他有用的项目。

- Not able to find an exporter for your language? Remember, the   [OpenTelemetry Collector](/docs/collector) supports exporting to a variety of   systems and works with all OpenTelemetry Core Components!
- Are you a project maintainer? See,  [Adding a project to the OpenTelemetry Registry](adding).
- Check back regularly, the community and registry are growing!

{{% /blocks/section %}}

{{< blocks/section color="white" type="container-lg" >}}

{{<registry-search-form>}}

{{< /blocks/section >}}

---
title: 贡献组织开放遥测营销指南
linkTitle: Marketing Guidelines
---

OpenTelemetry(又名OTel)是终端用户、相邻OSS项目和最终销售基于OTel数据或组件的产品和服务的供应商之间的协作。
像许多面向标准的项目一样，在OTel上合作的供应商也在市场上竞争，因此，为贡献组织如何沟通和传递有关OTel的信息建立一些基本规则和期望是很重要的。



本文档分为两部分:

- **Goals and Guidelines:** 我们想要达到什么目标?我们的指导是什么?
- **Concerns and consequences:** 我们如何确定是否违反了指导方针?我们该怎么做呢?
 
## Goals and Guidelines

There are three high-level focus areas for these goals and guidelines.

### I: OpenTelemetry is a joint effort

- Do’s:
  - Use project collateral such as logo and name in line with the Linux
    Foundation’s branding and
    [trademark usage guidelines](https://www.linuxfoundation.org/trademark-usage/)
  - Emphasize that OTel would not be possible without collaboration from many
    contributors who happen to work for competing vendors and providers
  - Cite names of the other contributors and vendors involved with OTel efforts
  - Emphasize our common goals as a community to improve end user/developer
    experiences and empower them
- Don’ts:
  - Imply that a single provider is responsible for OTel itself, and/or one of
    its various component parts
  - Diminish the contributions of another organization or of another individual

### II: It’s not a competition

- Do’s:
  - Emphasize that all contributions are valuable, and that they come in many
    shapes and sizes, including:
  - Contributions to the core project code or to language- or framework-specific
    SDKs
  - Creating and sharing educational resources (videos, workshops, articles), or
    shared resources that can be used for educational purposes (e.g. a sample
    app using specific language/framework)
  - Community-building activities such as organizing an event or meetup group
  - Publicly recognize and thank other organizations for their contributions to
    OTel
- Don’ts:
  - Directly compare the volume or value of different contributors to OTel
    (E.g., via [CNCF devstats](https://devstats.cncf.io/))
  - Imply that infrequent or minor contributors to OTel are necessarily
    second-class citizens, and/or that their own OTel compatibility should be
    questioned as a result (in fact, there’s no reason that any provider needs
    to contribute to OTel in order to support it)

### III: Promote awareness of OTel interoperability and modularization

- Do’s:
  - “Shout from the rooftops” about OTel compatibility – the more that end-users
    understand what they can do with OTel data, the better
  - Emphasize the vendor-neutrality and portability of any OTel integration
- Don’ts:
  - Imply that an end-user isn’t “Using OTel” unless they’re using some specific
    set of components within OTel (OTel is a “wide” project with many decoupled
    components)
  - Publicly denigrate the OTel support of another provider, particularly
    without objective evidence

## Concerns and Consequences

Inevitably there will be instances where vendors (or at least their Marketing
departments) run afoul of these guidelines. To date, this hasn’t happened
frequently, so we don’t want to create an over-complicated process to handle
concerns.

Here is how we handle such circumstances:

1. Whomever notices the relevant public (marketing) content should write an
   email to cncf-opentelemetry-governance@lists.cncf.io and include an
   explanation of why the content is problematic, ideally referencing the
   [relevant guidelines above](#goals-and-guidelines).
1. The OTel Governance Committee (GC) will discuss the case during its next
   (weekly) meeting, or asynchronously via email if possible. The OTel GC
   guarantees a response via email **within two weeks** of the initial report.
1. If the GC agrees that there’s a problem, a corrective action will be
   recommended to the author of the content in question, and the GC will request
   that the organization that published the content train relevant employees on
   the content in this document as a further preventative measure.

If a pattern develops with a particular vendor, the GC will meet to discuss more
significant consequences – for instance, removing that vendor’s name from
OTel-maintained lists of compatible providers, or simply publicly documenting
the pattern of poor community behavior.

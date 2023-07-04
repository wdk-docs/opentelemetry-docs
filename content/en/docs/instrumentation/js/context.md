---
title: 上下文
description: OpenTelemetry JavaScript上下文API文档
aliases: [/docs/instrumentation/js/api/context]
weight: 60
---

创建并启动 HTTP 服务器为了使 OpenTelemetry 工作，它必须存储和传播重要的遥测数据
。例如，当接收到请求并启动 span 时，它必须对创建子 span 的组件可用。为了解决这个
问题，OpenTelemetry 将跨度存储在 Context 中。本文档描述了 JavaScript 的
OpenTelemetry 上下文 API 以及如何使用它。

更多信息:

- [上下文规范](../../specs/otel/context/README.md)
- [上下文 API 参考](https://open-telemetry.github.io/opentelemetry-js/classes/_opentelemetry_api.ContextAPI.html)

## 上下文管理器

上下文 API 依赖于上下文管理器来工作。本文档中的示例将假设您已经配置了上下文管理
器。通常情况下，上下文管理器是由你的 SDK 提供的，但是也可以像这样直接注册一个:

```typescript
import * as api from '@opentelemetry/api';
import { AsyncHooksContextManager } from '@opentelemetry/context-async-hooks';

const contextManager = new AsyncHooksContextManager();
contextManager.enable();
api.context.setGlobalContextManager(contextManager);
```

## 根上下文

`ROOT_CONTEXT`是空的上下文。如果没有活跃上下文，则`ROOT_CONTEXT`是活跃的。活跃上
下文将在下面[活跃上下文](#active-context)进行解释。

## 上下文的键

上下文条目是键值对。可以通过调用`api.createContextKey(description)`来创建键.

```typescript
import * as api from '@opentelemetry/api';

const key1 = api.createContextKey('My first key');
const key2 = api.createContextKey('My second key');
```

## 基本操作

### 获得条目

使用`context.getValue(key)`方法访问条目。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('some key');
// ROOT_CONTEXT is the empty context
const ctx = api.ROOT_CONTEXT;

const value = ctx.getValue(key);
```

### 设置条目

条目是通过使用`context.setValue(key, value)` 方法创建的。设置上下文条目将创建一
个新上下文，其中包含前一个上下文的所有条目，但包含新条目。设置上下文条目不会修改
之前的上下文。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('some key');
const ctx = api.ROOT_CONTEXT;

// add a new entry
const ctx2 = ctx.setValue(key, 'context 2');

// ctx2 contains the new entry
console.log(ctx2.getValue(key)); // "context 2"

// ctx is unchanged
console.log(ctx.getValue(key)); // undefined
```

### 删除条目

通过调用`context.deleteValue(key)`来删除条目。删除上下文条目将创建一个新的上下文
，其中包含前一个上下文的所有条目，但不包含由键标识的条目。删除上下文条目不会修改
前一个上下文。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('some key');
const ctx = api.ROOT_CONTEXT;
const ctx2 = ctx.setValue(key, 'context 2');

// remove the entry
const ctx3 = ctx.deleteValue(key);

// ctx3 does not contain the entry
console.log(ctx3.getValue(key)); // undefined

// ctx2 is unchanged
console.log(ctx2.getValue(key)); // "context 2"
// ctx is unchanged
console.log(ctx.getValue(key)); // undefined
```

## 活跃的上下文

!!! IMPORTANT

    这假定您已经配置了上下文管理器。如果没有，`api.context.active()`将总是返回`ROOT_CONTEXT`。

活跃上下文是由`api.context.active()`返回的上下文。上下文对象包含允许跟踪单个执行
线程的跟踪组件相互通信并确保成功创建跟踪的条目。例如，当创建一个 span 时，可以将
其添加到上下文中。稍后，当创建另一个 span 时，它可以使用上下文中的 span 作为其父
span。这是通过使用 node 中的[async_hooks]或[AsyncLocalStorage]或 web 中
的[zone.js]等机制来完成的，以便通过单次执行传播上下文。如果没有活跃上下文，则返
回`ROOT_CONTEXT`，它只是一个空的上下文对象。

[async_hooks]: https://nodejs.org/api/async_hooks.html
[AsyncLocalStorage]:
  https://nodejs.org/api/async_context.html#async_context_class_asynclocalstorage
[zone.js]: https://github.com/angular/angular/tree/main/packages/zone.js

### 获得活跃上下文

活跃上下文是由`api.context.active()`返回的上下文。

```typescript
import * as api from '@opentelemetry/api';

// Returns the active context
// If no context is active, the ROOT_CONTEXT is returned
const ctx = api.context.active();
```

### 设置活跃的上下文

上下文可以通过使用`api.context.with(ctx, callback)`来激活。在`callback`执行期间
，传递给`with`的上下文将由`context.active`返回。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('Key to store a value');
const ctx = api.context.active();

api.context.with(ctx.setValue(key, 'context 2'), async () => {
  // "context 2" is active
  console.log(api.context.active().getValue(key)); // "context 2"
});
```

`api.context.with(context, callback)` 的返回值是回调的返回值。回调总是同步调用的
。

```typescript
import * as api from '@opentelemetry/api';

const name = await api.context.with(api.context.active(), async () => {
  const row = await db.getSomeValue();
  return row['name'];
});

console.log(name); // name returned by the db
```

活跃上下文执行可能是嵌套的。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('Key to store a value');
const ctx = api.context.active();

// No context is active
console.log(api.context.active().getValue(key)); // undefined

api.context.with(ctx.setValue(key, 'context 2'), () => {
  // "context 2" is active
  console.log(api.context.active().getValue(key)); // "context 2"
  api.context.with(ctx.setValue(key, 'context 3'), () => {
    // "context 3" is active
    console.log(api.context.active().getValue(key)); // "context 3"
  });
  // "context 2" is active
  console.log(api.context.active().getValue(key)); // "context 2"
});

// No context is active
console.log(api.context.active().getValue(key)); // undefined
```

### 例子

这个更复杂的示例说明了如何不修改上下文，而是创建新的上下文对象。

```typescript
import * as api from '@opentelemetry/api';

const key = api.createContextKey('Key to store a value');

const ctx = api.context.active(); // Returns ROOT_CONTEXT when no context is active
const ctx2 = ctx.setValue(key, 'context 2'); // does not modify ctx

console.log(ctx.getValue(key)); //? undefined
console.log(ctx2.getValue(key)); //? "context 2"

const ret = api.context.with(ctx2, () => {
  const ctx3 = api.context.active().setValue(key, 'context 3');

  console.log(api.context.active().getValue(key)); //? "context 2"
  console.log(ctx.getValue(key)); //? undefined
  console.log(ctx2.getValue(key)); //? "context 2"
  console.log(ctx3.getValue(key)); //? "context 3"

  api.context.with(ctx3, () => {
    console.log(api.context.active().getValue(key)); //? "context 3"
  });
  console.log(api.context.active().getValue(key)); //? "context 2"

  return 'return value';
});

// The value returned by the callback is returned to the caller
console.log(ret); //? "return value"
```

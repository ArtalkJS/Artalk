# 插件开发

## 使用 Vite 开发

我们提供了 Vite 的集成插件，使用 Vite 开发能够开箱即用地构建 Artalk 插件，简化 Vite 配置。

参考文档：[@artalk/plugin-kit](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-kit/README.md)。

同时，我们提供了一个模版代码仓库，你可以直接 fork 来开发插件：[artalk-plugin-sample](https://github.com/ArtalkJS/artalk-plugin-sample)。

通过 Vite 结合前端生态，你可以任选你喜欢的技术栈，如 Vue、React、Svelte、SolidJS 等框架开发 Artalk 插件。

## `Artalk.use`

使用 `Artalk.use` 来扩展 Artalk。

```js
import Artalk from 'artalk'

Artalk.use((ctx) => {
  ctx.on('list-loaded', () => {
    console.log('评论列表加载完毕')
    ctx.getCommentNodes().forEach((node) => {
      node.getEl().style.background = 'red'
    })
  })
})

Artalk.use((ctx) => {
  let el = null

  ctx.on('mounted', () => {
    el = document.createElement('div')
    document.body.appendChild(el)
    console.log('Artalk 实例被挂载')
  })

  ctx.on('unmounted', () => {
    el.remove()
    console.log('Artalk 实例被销毁')
  })
})

const artalk = Artalk.init({ ... })
```

请注意：

- `Artalk.use` 必须在 `Artalk.init` 之前调用才生效。
- 请勿依赖插件的加载顺序，请监听事件来执行插件逻辑。
- 当所有的插件加载完毕后，会触发 `created` 事件。
- 当插件和配置加载完毕后，会触发 `mounted` 事件。
- 当 Artalk 实例被销毁时，会触发 `unmounted` 事件。
- 当配置发生变化时，会触发 `updated` 事件。

在 use 函数中，你可以通过 `ctx` 来访问 Artalk 的 Context 对象。

### 插件可配置选项

在 TypeScript 中，你可以通过 `ArtalkPlugin<T>` 来定义一个带选项的插件类型。

```ts
export interface DemoPluginOptions {
  foo?: string
}

export const ArtalkDemoPlugin: ArtalkPlugin<DemoPluginOptions> = (ctx, options = {}) => {
  console.log(options.foo)
}
```

在执行 `Artalk.use` 时，传入选项：

```ts
import { ArtalkDemoPlugin } from 'artalk-plugin-demo'

Artalk.use(ArtalkDemoPlugin, { foo: 'bar' })
```

## ContextAPI

Context 对象包含了 Artalk 的上下文信息。

| 成员                  | 说明                 |
| --------------------- | -------------------- |
| `ctx.getEl`           | 获取容器元素         |
| `ctx.getConf`         | 获取配置             |
| `ctx.updateConf`      | 更新配置             |
| `ctx.watchConf`       | 监听配置             |
| `ctx.setDarkMode`     | 设置夜间模式         |
| `ctx.getApi`          | 获取 API 客户端对象  |
| --                    | --                   |
| `ctx.fetch`           | 获取评论数据         |
| `ctx.reload`          | 重载评论列表         |
| `ctx.getComments`     | 获取所有评论数据对象 |
| `ctx.getCommentNodes` | 获取所有评论节点对象 |
| --                    | --                   |
| `ctx.on`              | 添加事件监听         |
| `ctx.off`             | 解除事件监听         |
| `ctx.trigger`         | 触发事件             |
| --                    | --                   |
| `ctx.get`             | 获取依赖             |
| `ctx.inject`          | 注入依赖             |

::: warning
Context API 目前仍不稳定，开发可能会有变动，升级请关注 CHANGELOG。
:::

参考：[@artalk/src/types/context.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/context.ts)

Artalk 包含完整的 TypeScript 类型定义，你可以通过编辑器的自动补全来查看 API。

## 示例插件

以下为 Artalk 官方维护的外置插件列表：

| 插件                                                                                         | 说明           |
| -------------------------------------------------------------------------------------------- | -------------- |
| [artalk-plugin-sample](https://github.com/ArtalkJS/artalk-plugin-sample)                     | 示例插件       |
| [@artalk/plugin-katex](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-katex)       | LaTeX 公式插件 |
| [@artalk/plugin-auth](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-auth)         | Auth 插件 (SolidJS) |
| [@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-lightbox) | 图片灯箱基础插件   |

同时 Artalk 内部也有很多插件的实现，你可以参考源码来开发插件：

[@ArtalkJS/Artalk - src/plugins](https://github.com/ArtalkJS/Artalk/tree/master/ui/artalk/src/plugins)

## 后端插件开发

【TODO】

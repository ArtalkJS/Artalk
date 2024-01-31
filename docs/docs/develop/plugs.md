# 插件开发

使用 `Artalk.use` 来扩展 Artalk。

```js
import Artalk from 'artalk'

Artalk.use((ctx) => {
  ctx.on('mounted', () => {
    ctx.get('editor').setContent("Hello World")
  })

  ctx.on('list-loaded', () => {
    console.log('评论列表加载完毕')
  })
})

const artalk = Artalk.init({ ... })
```

请注意：

- `Artalk.use` 必须在 `Artalk.init` 之前调用。
- 请勿依赖插件的加载顺序，请监听事件来执行插件逻辑。
- 当所有的插件加载完毕后，会触发 `inited` 事件。

在 use 函数中，你可以通过 `ctx` 来访问 Artalk 的上下文。

## Context

Context 对象包含了 Artalk 的上下文信息。

| 成员 | 说明 |
| --- | --- |
| `ctx.$root` | Artalk 容器元素 |
| `ctx.conf` | Artalk 配置 |
| `ctx.fetch` | 加载数据数据 |
| `ctx.reload` | 重载评论列表 |
| `ctx.on` | 添加事件监听 |
| `ctx.off` | 解除事件监听 |
| `ctx.trigger` | 触发事件 |
| `ctx.inject` | 依赖注入 |
| `ctx.get` | 获取依赖 |

::: warning
Context API 目前仍不稳定，开发可能会有变动，升级请关注 CHANGELOG。
:::

参考：[@artalk/src/types/context.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/context.ts)

Artalk 拥有完整的 TypeScript 类型定义，你可以通过编辑器的自动补全来查看 API。

## 示例插件

Artalk 有很多内置的插件，你可以参考它们的源码来开发自己的插件。

[@ArtalkJS/Artalk - src/plugins](https://github.com/ArtalkJS/Artalk/tree/master/ui/artalk/src/plugins)

我们还提供了一些外置的插件，同样可以参考：

| 插件 | 说明 |
| --- | --- |
| [@artalk/plugin-katex](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-katex) | LaTeX 公式插件 |
| [@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-lightbox) | 图片灯箱插件 |

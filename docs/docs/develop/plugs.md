# 插件开发

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

## Context

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

Artalk 有很多内置的插件，你可以参考它们的源码来开发自己的插件。

[@ArtalkJS/Artalk - src/plugins](https://github.com/ArtalkJS/Artalk/tree/master/ui/artalk/src/plugins)

我们还提供了一些外置的插件，同样可以参考：

| 插件                                                                                         | 说明           |
| -------------------------------------------------------------------------------------------- | -------------- |
| [@artalk/plugin-katex](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-katex)       | LaTeX 公式插件 |
| [@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/tree/master/ui/plugin-lightbox) | 图片灯箱插件   |

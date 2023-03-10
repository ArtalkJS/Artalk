# 插件开发

你可以使用 `Artalk.use` 来扩展 Artalk。

```js
Artalk.use((ctx) => {
  ctx.editor.setContent("Hello World")
})
```

## Context

::: warning
Context API 目前仍不稳定，开发可能会有变动，升级请关注 CHANGELOG。
:::

参考：[@artalk/types/context.d.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/packages/artalk/types/context.d.ts)

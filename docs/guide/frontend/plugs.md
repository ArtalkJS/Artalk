# 插件

::: warning
目前文档仍在陆续完善中...
:::

你可以通过 `Artalk.use` 来装载 Artalk 插件。

```js
Artalk.use((ctx) => {
  ctx.editor.setContent("Hello World")
})
```

## Context

参考：[@artalk/types/context.d.ts](https://github.com/ArtalkJS/Artalk/blob/master/packages/artalk/types/context.d.ts)

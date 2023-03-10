# Frontend API

## 初始化 `Artalk.init`

创建 Artalk 实例化对象。

```js
Artalk.init({
  el:        '#Comments',
  pageKey:   '/post/1',
  pageTitle: '关于引入 Artalk 的这档子事',
  server:    'http://your_domain:8080',
  site:      'Artalk 的博客',
})
```

## 更新配置 `Artalk.update`

修改 Artalk 当前配置。

```js
Artalk.update({
  // 新的配置...
})
```

该函数调用后，`conf-loaded` 事件将被触发。

更新配置不会自动刷新评论列表，可按需继续执行 `Artalk.reload()` 函数。

注：前端的一些配置将会被后端的配置覆盖，详情见：[在后端控制前端](../guide/backend/fe-control.md)

## 重新加载 `Artalk.reload`

刷新 Artalk 评论列表。

```js
Artalk.reload()
```

列表加载前将触发 `list-load` 事件，加载完毕后将触发 `list-loaded` 事件。

## 释放资源 `Artalk.destroy`

销毁 Artalk 实例，用于资源释放。

```js
Artalk.destroy()
```

## 事件监听 `Artalk.on`

添加 Artalk 事件监听。

```js
Artalk.on('list-loaded', () => {
  alert('评论列表加载完毕')
})
```

可监听事件类型见：[前端界面事件](./event.md)

## 解除监听 `Artalk.off`

解除 Artalk 事件监听。

```js
const handler = () => {
  alert('评论列表加载完毕')
}

Artalk.on('list-loaded', handler)
Artalk.off('list-loaded', handler)
```

## 触发事件 `Artalk.trigger`

触发 Artalk 事件。

```js
Artalk.trigger('list-loaded')
```

## 扩展插件 `Artalk.use`

该 API 用于扩展 Artalk。

```js
Artalk.use((ctx) => {
  ctx.editor.setContent("Hello World")
})
```

详情见：[插件开发](./plugs.md)

## 夜间模式 `Artalk.setDarkMode`

修改夜间模式，可以配合博客主题调用，例如当用户点击夜间模式切换按钮。

```js
Artalk.setDarkMode(true)
```

## 浏览量组件 `Artalk.loadCountWidget`

详情见：[浏览量统计](../guide/frontend/pv.md)

# Frontend API

## 安装依赖

用你喜欢的包管理器安装 Artalk：

```bash
pnpm add artalk
# 或者使用 npm
npm install artalk
```

## 创建实例 `init`

调用该函数将创建并得到一个 Artalk 的实例化对象，可用于后续操作。

```js
import Artalk from 'artalk'

const artalk = Artalk.init({
  el:        '#Comments',
  pageKey:   '/post/1',
  pageTitle: '关于引入 Artalk 的这档子事',
  server:    'http://your_domain:8080',
  site:      'Artalk 的博客',
})
```

调用该函数会异步地向后端发起请求：

1. 获取前端配置
2. 获取评论列表

配置和插件加载完毕后会触发 `mounted` 事件。

评论列表加载完毕后会触发 `list-loaded` 事件。

注：前端界面配置可能会被覆盖，详情见：[在后端控制前端](../guide/backend/fe-control.md)。

## 释放资源 `destroy`

销毁 Artalk 实例，用于资源释放。

```js
const artalk = Artalk.init({ ... })

artalk.destroy()
```

调用该函数时，将触发 `unmounted` 事件。释放资源将执行 Artalk 挂载的 DOM 节点删除，以及解除 Artalk 事件监听等操作。当该函数执行完毕后，Artalk 实例将不再可用。

在 Vue / React 等框架中，务必在组件销毁时调用该函数，否则会造成内存泄漏：

```ts
import { onUnmounted } from 'vue'

onUnmounted(() => {
  artalk.destroy()
})
```

## 更新配置 `update`

修改 Artalk 当前配置。

```js
const artalk = Artalk.init({ ... })

artalk.update({
  // 新的配置...
})
```

调用该函数将触发 `updated` 事件。

更新配置不会自动刷新评论列表，可按需继续执行 `artalk.reload()` 函数。

请注意，该函数是一个实例化对象中的成员方法，并非全局函数。

前端配置参考：[前端配置](../guide/frontend/config.md)

## 重新加载 `reload`

重新从后端加载 Artalk 评论列表数据。

```js
const artalk = Artalk.init({ ... })

artalk.reload()
```

列表加载前将触发 `list-load` 事件，加载完毕后将触发 `list-loaded` 事件。

## 监听事件 `on`

添加 Artalk 事件监听。

```js
const artalk = Artalk.init({ ... })

artalk.on('list-loaded', (comments) => {
  // comments 包含当前页面所有评论数据 (而非仅本次请求获取的部分评论)
  console.log('评论列表加载完毕: ', comments)
})
```

可参考：[前端 Event](./event.md)。

## 解除监听 `off`

解除 Artalk 事件监听。

```js
const artalk = Artalk.init({ ... })

const handler = () => {
  alert('评论列表加载完毕')
}

artalk.on('list-loaded', handler)
artalk.off('list-loaded', handler)
```

## 触发事件 `trigger`

触发 Artalk 事件。

```js
const artalk = Artalk.init({ ... })

artalk.trigger('list-loaded', [...])
```

## 加载插件 `use`

通过 `Artalk.use` 函数加载插件。请注意，该函数是一个全局函数，并非实例化对象中的方法。通过该函数加载的插件将对所有 Artalk 实例生效。

```js
import Artalk from 'artalk'

Artalk.use((ctx) => {
  ctx.editor.setContent("Hello World")
})

const artalk = Artalk.init({ ... })
```

详情见：[插件开发](./plugs.md)

## 夜间模式 `setDarkMode`

修改夜间模式，可以配合博客主题调用，例如当用户点击夜间模式切换按钮。

```js
const artalk = Artalk.init({ ... })

artalk.setDarkMode(true)
```

也可以通过调用 `Artalk.init` 时，通过 `darkMode` 参数设置初始夜间模式。

```js
const artalk = Artalk.init({
  // ...其他配置
  darkMode:  'auto',
})
```

## 浏览量组件 `loadCountWidget`

详情见：[浏览量统计](../guide/frontend/pv.md)

## 获取配置 `getConf`

获取当前 Artalk 配置。

```js
const artalk = Artalk.init({ ... })

const conf = artalk.getConf()
```

## 获取挂载元素 `getEl`

获取当前 Artalk 挂载的 DOM 元素。

```js
const artalk = Artalk.init({ ... })

const el = artalk.getEl()
```

# Event

## 基本事件

|事件名|描述|
|-|-|
| `list-load` | 评论加载事件 |
| `list-loaded` | 评论加载完成事件 |
| `list-inserted` | 新增评论插入事件 |
| `editor-submit` | 编辑器提交事件 |
| `editor-submitted` | 编辑器提交完成事件 |
| `user-changed` | 本地用户数据变更事件 |
| `conf-loaded` | 配置变更事件 |
| `sidebar-show` | 侧边栏显示事件 |
| `sidebar-hide` | 侧边栏隐藏事件 |

## 添加事件监听

```js
Artalk.use(ctx => {
  ctx.on('list-loaded', () => {
    alert('评论已加载完毕')
  })
})
```

## 解除事件监听

```js
let foo = function() { /* do something */ }

Artalk.use(ctx => {
  ctx.on('list-loaded', foo)
  ctx.off('list-loaded', foo)
})
```

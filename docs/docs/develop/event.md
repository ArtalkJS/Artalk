# Event

## 基本事件

| 事件名              | 描述                 |
| ------------------- | -------------------- |
| `created`           | 初始化后             |
| `mounted`           | 数据加载后           |
| `updated`           | 数据更新后           |
| `unmounted`         | 销毁后               |
| `list-fetch`        | 评论列表请求时       |
| `list-fetched`      | 评论列表请求后       |
| `list-load`         | 评论装载前           |
| `list-loaded`       | 评论装载后           |
| `list-failed`       | 评论加载错误时       |
| `list-goto-first`   | 评论列表归位时       |
| `list-reach-bottom` | 评论列表滚动到底部时 |
| `comment-inserted`  | 评论插入后           |
| `comment-updated`   | 评论更新后           |
| `comment-deleted`   | 评论删除后           |
| `comment-rendered`  | 评论节点渲染后       |
| `notifies-updated`  | 未读消息变更时       |
| `list-goto`         | 评论跳转时           |
| `page-loaded`       | 页面数据更新后       |
| `editor-submit`     | 编辑器提交时         |
| `editor-submitted`  | 编辑器提交后         |
| `user-changed`      | 本地用户数据变更时   |
| `sidebar-show`      | 侧边栏显示           |
| `sidebar-hide`      | 侧边栏隐藏           |

事件声明代码：[@ArtalkJS/Artalk - src/types/event.ts](https://github.com/ArtalkJS/Artalk/blob/master/ui/artalk/src/types/event.ts)

## 添加事件监听

```js
artalk.on('list-loaded', () => {
  alert('评论已加载完毕')
})
```

## 解除事件监听

```js
let foo = function () {
  /* do something */
}

artalk.on('list-loaded', foo)
artalk.off('list-loaded', foo)
```

# 划线批注

`@artalk/plugin-annotation` 为 Artalk 增加「选中文字 → 划线 → 批注」的能力。读者在正文中选中一段文字后即可创建批注，讨论始终贴着原文。

## 安装

```bash
pnpm add @artalk/plugin-annotation
```

## 启用

```ts
import Artalk from 'artalk'
import { ArtalkAnnotationPlugin } from '@artalk/plugin-annotation'

Artalk.use(ArtalkAnnotationPlugin)

const artalk = Artalk.init({
  el: '#comments',
  server: 'https://artalk.example.com',
  site: 'MySite',
  pageKey: location.pathname,
  annotation: {
    enabled: true,
    root: 'article',
    hideInList: true,
    allowGuest: true,
  },
})
```

也可以通过 script 标签引入：

```html
<script src="https://artalk.example.com/dist/Artalk.js"></script>
<script src="https://artalk.example.com/dist/plugins/artalk-plugin-annotation.js"></script>
<script>
  Artalk.init({
    el: '#comments',
    server: 'https://artalk.example.com',
    site: 'MySite',
    annotation: { enabled: true, root: 'article' },
  })
</script>
```

## 配置项

| 配置 | 类型 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `enabled` | `boolean` | `false` | 是否启用划线批注 |
| `root` | `string \| HTMLElement` | `'article'` | 允许划词的正文区域 |
| `hideInList` | `boolean` | `true` | 普通评论列表是否隐藏划线根评论 |
| `allowGuest` | `boolean` | `true` | 是否允许游客创建划线 |

插件还支持通过 `Artalk.use` 传入 `excludeSelector` 与 `panelContainer`。

## 工作方式

- 一条划线批注就是一条带 W3C `TextQuoteSelector` 锚点的根评论。
- 回复复用普通评论树。
- 普通评论列表默认排除划线根评论（`hideInList`）。
- 审核、反垃圾、通知与权限与普通评论共用。
- 正文变化导致锚点无法定位时，该划线会标记为「原文已变更」，可手动重新定位。

## SPA

路由切换后更新 pageKey：

```ts
artalk.update({ pageKey: location.pathname })
```

插件会清理旧高亮并重新加载当前页的划线。

## 限制

- 仅支持文本划线，不支持图片或区域批注。
- 划线面板只展示当前页面。

# 图片灯箱

Artalk LightBox 插件能帮助你将网站**现有的图片灯箱**功能自动集成到 Artalk 中。

::: code-group

```html [浏览器引入]
<!-- 1. 引入图片灯箱依赖，例如 LightGallery (通常博客主题已提供，无需再引入) -->
<script src="lightgallery.js"></script>

<!-- 2. 引入 Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. 引入 Artalk LightBox 插件 -->
<script src="https://unpkg.com/@artalk/plugin-lightbox/dist/artalk-plugin-lightbox.js"></script>
```

```ts [Node 引入]
// pnpm add lightgallery @artalk/plugin-lightbox
import Artalk from 'artalk'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'lightgallery/css/lightgallery.css'

Artalk.use(ArtalkLightboxPlugin, {
  // 手动配置引入灯箱库
  lightGallery: {
    lib: async () => (await import('lightgallery')).default
  }
})
```

:::

目前自动集成支持：[LightGallery](https://github.com/sachinchoolur/lightGallery) • [FancyBox](https://github.com/fancyapps/fancybox) • [lightbox2](https://github.com/lokesh/lightbox2) • [PhotoSwipe](https://photoswipe.com/)

对于暂未适配的灯箱库，我们期待你的 PR 😉：[@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-lightbox/src/main.ts)。

::: details 附：图片灯箱依赖 CDN 资源

注：一个博客主题可能包含现成的图片灯箱插件，无需重复引入。

#### LightGallery

```html
<link
  rel="stylesheet"
  href="https://unpkg.com/lightgallery@2.7.2/css/lightgallery.css"
/>
<script src="https://unpkg.com/lightgallery@2.7.2/lightgallery.min.js"></script>
```

#### FancyBox

```html
<link
  rel="stylesheet"
  href="https://unpkg.com/@fancyapps/ui@4.0.31/dist/fancybox.css"
/>
<script src="https://unpkg.com/@fancyapps/ui@4.0.31/dist/fancybox.umd.js"></script>
```

:::

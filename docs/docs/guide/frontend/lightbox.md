# 图片灯箱

Artalk 图片灯箱插件帮助你将网站**现有的图片灯箱**功能自动集成到 Artalk 中。

## 浏览器环境

如果处于浏览器环境，博客主题通常自带图片灯箱，插件将自动检测 `window` 全局对象中的灯箱库，并集成到 Artalk 中，无需额外配置。

```html
<!-- 1. 引入图片灯箱依赖，例如 LightGallery (通常博客主题已提供，无需再引入) -->
<script src="lightgallery.js"></script>

<!-- 2. 引入 Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. 引入 Artalk LightBox 插件 -->
<script src="https://unpkg.com/@artalk/plugin-lightbox/dist/artalk-plugin-lightbox.js"></script>
```

如有需要，可以通过全局变量 `ATK_LIGHTBOX_CONF` 对灯箱库进行额外的配置：

```html
<script>
// Config for LightGallery
window.ATK_LIGHTBOX_CONF = {
  plugins: [lgZoom, lgThumbnail],
  speed: 500,
  licenseKey: 'your_license_key'
}
</script>
```

## Node 环境

安装图片灯箱插件：

```bash
pnpm add @artalk/plugin-lightbox
```

配置插件，通过 `import` 动态引入灯箱库：

::: code-group

```ts [LightGallery]
import Artalk from 'artalk'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'lightgallery/css/lightgallery.css'

Artalk.use(ArtalkLightboxPlugin, {
  lightGallery: {
    lib: () => import('lightgallery'),
  },
})
```

```ts [PhotoSwipe]
import Artalk from 'artalk'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'photoswipe/style.css'

Artalk.use(ArtalkLightboxPlugin, {
  photoSwipe: {
    lib: () => import('photoswipe/lightbox'),
    pswpModule: () => import('photoswipe'),
  },
})
```

```ts [LightBox]
import Artalk from 'artalk'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'lightbox2/dist/css/lightbox.min.css'
import jQuery from 'jquery'

window.$ = jQuery

Artalk.use(ArtalkLightboxPlugin, {
  lightBox: {
    // @ts-ignore
    lib: () => import('lightbox2'),
  },
})
```

```ts [FancyBox]
import Artalk from 'artalk'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'fancybox/dist/css/jquery.fancybox.css'
import jQuery from 'jquery'

window.$ = jQuery

Artalk.use(ArtalkLightboxPlugin, {
  fancyBox: {
    // @ts-ignore
    lib: () => import('fancybox'),
  },
})
```

:::

在 Node 环境中，如有需要对灯箱库进行额外的配置，可以通过 `config` 选项配置：

```ts
Artalk.use(ArtalkLightboxPlugin, {
  // ...
  // Config for LightGallery
  config: {
    speed: 500,
  },
})
```

---

目前自动集成支持：[LightGallery](https://github.com/sachinchoolur/lightGallery) • [FancyBox](https://github.com/fancyapps/fancybox) • [lightbox2](https://github.com/lokesh/lightbox2) • [PhotoSwipe](https://photoswipe.com/)

一些灯箱库暂未适配，期待你的 PR 😉：[@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-lightbox/src/main.ts)。

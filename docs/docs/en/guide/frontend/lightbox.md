# Image Lightbox

The Artalk Image Lightbox plugin seamlessly integrates your website's **existing image lightbox** functionality into Artalk.

## Browser Environment

In a browser environment, blog themes typically come with an image lightbox. The plugin will automatically detect the lightbox library in the `window` global object and integrate it into Artalk without additional configuration.

```html
<!-- 1. Include the image lightbox dependency, such as LightGallery (usually provided by the blog theme, no need to re-include) -->
<script src="lightgallery.js"></script>

<!-- 2. Include Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. Include the Artalk LightBox plugin -->
<script src="https://unpkg.com/@artalk/plugin-lightbox/dist/artalk-plugin-lightbox.js"></script>
```

If necessary, you can configure the lightbox library with the global variable `ATK_LIGHTBOX_CONF`:

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

## Node Environment

Install the image lightbox plugin:

```bash
pnpm add @artalk/plugin-lightbox
```

Configure the plugin and dynamically import the lightbox library using `import`:

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

In a Node environment, if additional configuration for the lightbox library is required, it can be done through the `config` option:

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

Currently, the following libraries are supported for automatic integration: [LightGallery](https://github.com/sachinchoolur/lightGallery) • [FancyBox](https://github.com/fancyapps/fancybox) • [lightbox2](https://github.com/lokesh/lightbox2) • [PhotoSwipe](https://photoswipe.com/)

Some lightbox libraries are not yet supported. Contributions are welcome 😉: [@artalk/plugin-lightbox](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-lightbox/src/main.ts).

# 图片灯箱

Artalk LightBox 插件能帮助你将网站**现有的图片灯箱**功能自动集成到 Artalk 中。

```html
<!-- 1. 引入图片灯箱依赖，例如 LightGallery (通常博客主题已提供，无需再引入) -->
<script src="lightgallery.js"></script>

<!-- 2. 引入 Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. 引入 Artalk LightBox 插件 -->
<script src="https://unpkg.com/@artalk/plugin-lightbox/dist/artalk-plugin-lightbox.js"></script>
```

如上所示，额外引入一个 `artalk-plugin-lightbox.js` 文件即可。

目前自动集成支持：[LightGallery](https://github.com/sachinchoolur/lightGallery) [v2.5.0] / [FancyBox](https://github.com/fancyapps/fancybox) [v4.0.27] / [lightbox2](https://github.com/lokesh/lightbox2) [v2.11.3]

对于还未适配的图片灯箱，欢迎提交 PR -> [查看代码](https://github.com/ArtalkJS/Artalk/blob/master/ui/plugin-lightbox/main.ts)

::: details 附：图片灯箱依赖 CDN 资源

注：通常一个博客主题本来就是有图片灯箱插件的，所以无需重复引入。

#### LightGallery

```html
<link
  rel="stylesheet"
  href="https://unpkg.com/lightgallery@2.5.0/css/lightgallery.css"
/>
<script src="https://unpkg.com/lightgallery@2.5.0/lightgallery.min.js"></script>
```

#### FancyBox

```html
<link
  rel="stylesheet"
  href="https://unpkg.com/@fancyapps/ui@4.0.27/dist/fancybox.css"
/>
<script src="https://unpkg.com/@fancyapps/ui@4.0.27/dist/fancybox.umd.js"></script>
```

:::

### 配置灯箱

在引入 `artalk-plugin-lightbox.js` 之前对全局变量 `ATK_LIGHTBOX_CONF` 进行设置，如下：

```html
<script>
  window.ATK_LIGHTBOX_CONF = {
    groupAll: true,
    // ...其他配置
  }
</script>
<script src="artalk-plugin-lightbox.js"></script>
```

# 置入博客

## 通用方法

<<< @/code/quick-start/cdn.html

#### 什么时候执行 `Artalk.init({})`？

- 可以在任意位置引入 JS 和 CSS 资源，但需确保 JS 引入在执行 `Artalk.init({})` 前。
- 执行 `Artalk.init({ el: '#artalk' })` 前需要确保 `<div id="artalk"></div>` 在页面当中。

## Hugo

创建模板文件 `/主题目录/layouts/partials/comment/artalk.html`：

<!-- prettier-ignore-start -->

```html
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>

<script>
  Artalk.init({
    el: '#Comments',
    pageKey: '{{ .Permalink }}',
    pageTitle: '{{ .Title }}',
    server: '{{ $.Site.Params.artalk.server }}',
    site: '{{ $.Site.Params.artalk.site }}',
    // ...你的其他配置
  })
</script>
```

<!-- prettier-ignore-end -->

文章页模板 `/主题目录/layouts/_default/single.html` 合适的位置添加：

```html
<div class="article-comments">{{- partial "comment/artalk" . -}}</div>
```

修改 Hugo 配置文件：

::: code-group

```toml [config.toml]
[params.artalk]
server = 'https://artalk.example.org'
site = '你的站点名'
```

```yaml [config.yaml]
params:
  artalk:
    server: 'https://artalk.example.org'
    site: '你的站点名'
```

:::

## Hexo

创建 `/主题目录/layout/comment/artalk.ejs`：

```html
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<div id="Comments"></div>

<script>
  var artalk = Artalk.init({
    el: '#Comments',
    pageKey: '<%= page.permalink %>',
    pageTitle: '<%= page.title %>',
    server: '<%= theme.comment.artalk.server %>'
    site: '<%= theme.comment.artalk.site %>',
  })
</script>
```

修改文章模板文件，例如 `/主题目录/layout/post.ejs`：

```html
<div class="article-comments"><%- partial('comment/artalk') %></div>
```

编辑主题配置 `/主题目录/_config.example.yml`：

```yaml
comment:
  artalk:
    site: '你的站点名'
    server: 'https://artalk.example.org'
```

::: tip

NexT 主题可以安装 [Hexo NexT 主题的 Artalk 插件](https://github.com/leirock/hexo-next-artalk)

:::

## VitePress

可参考：[@ArtalkJS/Artalk:/docs/.vitepress](https://github.com/ArtalkJS/Artalk/tree/master/docs/.vitepress)

- `config.ts` 修改配置引入 CSS 资源
- `theme/Artalk.vue` 创建组件
- `theme/index.ts` 注册组件
- `theme/Layout.vue` 使用组件

注：SSG 应用需通过 `import()` 函数异步引入 Artalk，否则会导致构建失败。

```ts
import('artalk').then(({ default: Artalk }) => {
  artalk = Artalk.init({
    //...
  })
})
```

参考文件：[.vitepress/theme/Artalk.vue](https://github.com/ArtalkJS/Artalk/blob/master/docs/.vitepress/theme/Artalk.vue)

## VuePress

（待补充）

::: tip 提示

可参考：[置入框架](./import-framework.md)

:::

## Typecho

修改主题文章模板文件，例如 `post.php`：

```php
<!-- CSS -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet">

<!-- JS -->
<script src="/lib/artalk/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
 Artalk.init({
   el:        '#Comments',
   pageKey:   '<?php $this->permalink() ?>',
   pageTitle: '<?php $this->title() ?>',
   server:    '<后端地址>',
   site:      '<?php $this->options->title() ?>',
   // ...你的其他配置
 })
</script>
```

## WordPress

修改主题文章模板文件，例如 `single.php`：

```php
<!-- CSS -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet">

<!-- JS -->
<script src="/lib/artalk/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>
<script>
 Artalk.init({
   el:        '#Comments',
   pageKey:   '<?= addslashes(get_permalink()) ?>',
   pageTitle: '<?= addslashes(get_the_title()) ?>',
   server:    '<后端地址>',
   site:      '<?= addslashes(get_bloginfo('name')) ?>',
   // ...你的其他配置
 })
</script>
```

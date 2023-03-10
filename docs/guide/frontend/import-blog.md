# 置入博客

## 通用方法

<<< @/code/quick-start/cdn.html

## Hugo

创建模板文件 `/主题目录/layouts/partials/comment/artalk.html`：

```html
<link href="/lib/artalk/Artalk.css" rel="stylesheet">
<script src="/lib/artalk/Artalk.js"></script>

<!-- Artalk -->
<div id="Comments"></div>

<script>
  Artalk.init({
    el:        '#Comments',
    pageKey:   '{{ .Permalink }}',
    pageTitle: '{{ .Title }}',
    server:    '{{ $.Site.Params.artalk.server }}',
    site:      '{{ $.Site.Params.artalk.site }}',
    // ...你的其他配置
  })
</script>
```

文章页模板 `/主题目录/layouts/_default/single.html` 合适的位置添加：

```html
<div class="article-comments">
  {{- partial "comment/artalk" . -}}
</div>
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
<link href="/lib/artalk/Artalk.css" rel="stylesheet">
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
<div class="article-comments">
  <%- partial('comment/artalk') %>
</div>
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


## VuePress

以 [VuePress v2](https://github.com/vuepress/vuepress-next) 为例，继承默认主题 `@vuepress/theme-default`

可参考：[“/.vuepress/theme/Artalk.vue”](https://github.com/ArtalkJS/Docs/blob/eef37bca8cc0c9973bf121fdef4014dcd940f104/docs/.vuepress/theme/Artalk.vue) 创建 Artalk 评论组件。

在 `/.vuepress/theme/clientAppEnhance.ts` 文件中全局注册组件：

```ts
import { defineClientAppEnhance } from '@vuepress/client'

import Artalk from './Artalk.vue'

export default defineClientAppEnhance(({ app, router, siteData }) => {
  app.component('Artalk', Artalk)
})
```

主题布局 `/.vuepress/theme/Layout.vue`：

```vue
<template>
  <Layout>
    <template #page-bottom>
      <div class="page-meta">
        <!-- Artalk -->
        <Artalk />
      </div>
    </template>
  </Layout>
</template>

<script lang="ts">
import { defineComponent, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '@vuepress/theme-default/lib/client/layouts/Layout.vue'

export default defineComponent({
  components: { Layout },
  mounted: () => {
  }
})
</script>
```

主题配置文件 `/.vuepress/theme/index.ts`：

```ts
import { path } from '@vuepress/utils'

export default ({
  name: 'vuepress-theme-local',
  extends: '@vuepress/theme-default',
  layouts: {
    Layout: path.resolve(__dirname, 'Layout.vue'),
  },
  clientAppEnhanceFiles: path.resolve(__dirname, 'clientAppEnhance.ts'),
})
```

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

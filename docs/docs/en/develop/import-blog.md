# Import to Blog

## General Method

<<< @/code/quick-start/cdn.html

### When to execute `Artalk.init({})`?

- You can include JS and CSS resources anywhere, but ensure the JS is included before executing `Artalk.init({})`.
- Ensure `<div id="artalk"></div>` is present in the page before executing `Artalk.init({ el: '#artalk' })`.

## Hugo

Create the template file `/themes/YOUR_THEME/layouts/partials/comment/artalk.html`:

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
    // ...your additional configuration
  })
</script>
```

<!-- prettier-ignore-end -->

Add the following to the appropriate place in your article template `/themes/YOUR_THEME/layouts/_default/single.html`:

```html
<div class="article-comments">{{- partial "comment/artalk" . -}}</div>
```

Modify the Hugo configuration file:

::: code-group

```toml [config.toml]
[params.artalk]
server = 'https://artalk.example.org'
site = 'Your Site Name'
```

```yaml [config.yaml]
params:
  artalk:
    server: 'https://artalk.example.org'
    site: 'Your Site Name'
```

:::

## Hexo

Create `/themes/YOUR_THEME/layout/comment/artalk.ejs`:

```html
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<div id="Comments"></div>

<script>
  var artalk = Artalk.init({
    el: '#Comments',
    pageKey: '<%= page.permalink %>',
    pageTitle: '<%= page.title %>',
    server: '<%= theme.comment.artalk.server %>',
    site: '<%= theme.comment.artalk.site %>',
  })
</script>
```

Modify the article template file, for example `/themes/YOUR_THEME/layout/post.ejs`:

```html
<div class="article-comments"><%- partial('comment/artalk') %></div>
```

Edit the theme configuration `/themes/YOUR_THEME/_config.yml`:

```yaml
comment:
  artalk:
    site: 'Your Site Name'
    server: 'https://artalk.example.org'
```

::: tip

For NexT theme, you can install the [Artalk plugin for Hexo NexT theme](https://github.com/leirock/hexo-next-artalk)

:::

## VitePress

Refer to: [docs/.vitepress](https://github.com/ArtalkJS/Artalk/tree/master/docs/docs/.vitepress)

- Modify `config.ts` to include CSS resources
- Create `theme/Artalk.vue` component
- Register the component in `theme/index.ts`
- Use the component in `theme/Layout.vue`

Sample code: [docs/.vitepress/theme/Artalk.vue](https://github.com/ArtalkJS/Artalk/blob/master/docs/docs/.vitepress/theme/Artalk.vue)

## VuePress

(To be supplemented)

::: tip

Refer to: [Import to Frameworks](./import-framework.md)

:::

## Typecho

Modify the theme article template file, for example, `post.php`:

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
   server:    '<Server Address>',
   site:      '<?php $this->options->title() ?>',
   // ...your additional configuration
 })
</script>
```

## WordPress

Modify the theme article template file, for example, `single.php`:

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
   server:    '<Server Address>',
   site:      '<?= addslashes(get_bloginfo('name')) ?>',
   // ...your additional configuration
 })
</script>
```

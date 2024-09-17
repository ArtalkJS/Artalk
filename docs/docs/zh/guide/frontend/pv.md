# 浏览量统计

Artalk 内置页面浏览量统计和评论数统计功能，你可以在页面中显示页面的浏览量和评论数。

```html
浏览量：<span class="artalk-pv-count"></span>
评论数：<span class="artalk-comment-count"></span>
```

## 加载占位符

Artalk 加载需要时间，你可以在展示统计数的元素中加入占位符：

```html
<span class="artalk-pv-count">加载中...</span>
```

## 同时加载多个页面的统计数

例如在文章列表页，你可以显示每篇文章的浏览量和评论数。

当处于文章列表页时，仅需调用 `Artalk.loadCountWidget` 函数，而无需 `Artalk.init` (加载评论列表会使页面浏览量 PV 数 +1)。

<!-- prettier-ignore-start -->

```js
Artalk.loadCountWidget({
  server: '服务器地址',
  site: '站点名',
  pvEl: '.artalk-pv-count',
  countEl: '.artalk-comment-count',
  statPageKeyAttr: 'data-page-key',
})
```

<!-- prettier-ignore-end -->

然后放置多个 class 为 `artalk-pv-count` 的元素，通过属性 `data-page-key` 来指定需要查询的页面：

```html
<span class="artalk-pv-count" data-page-key="/test/1.html"></span>
<span class="artalk-pv-count" data-page-key="/test/2.html"></span>
<span class="artalk-pv-count" data-page-key="/test/3.html"></span>
```

评论数查询同理：

```html
<span class="artalk-comment-count" data-page-key="/test/1.html"></span>
<span class="artalk-comment-count" data-page-key="/test/2.html"></span>
<span class="artalk-comment-count" data-page-key="/test/3.html"></span>
```

## 自定义元素选择器

可通过配置项 `pvEl` 和 `countEl` 来指定元素选择器，以展示页面浏览量和评论数：

```js
Artalk.init({
  pvEl: '.your_pv_count_element',  // 页面浏览量元素选择器
  countEl: '.your_comment_count_element',  // 评论数元素选择器
})
```

## 自定义 `data-page-key` 属性名

配置项 `statPageKeyAttr` 的默认值为 `data-page-key`，Artalk 会通过该属性名来查询指定页面。为了方便博客主题的适配，可通过它自定义属性名，例如将其替换为 `data-path`：

```js
Artalk.loadCountWidget({
  statPageKeyAttr: 'data-path',
})
```

此时，对应的 HTML 代码应如下所示：

```html
<span class="artalk-pv-count" data-path="/test/1.html"></span>
```

这样，`data-path` 属性值将被用于查询指定的页面。

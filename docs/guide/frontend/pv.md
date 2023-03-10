# 浏览量统计

Artalk 内置页面浏览量统计功能，你可以在你的页面任意位置，放置 HTML 标签：

```html
<span id="ArtalkPV"></span>
```

当 Artalk 加载完毕时，该标签内容将修改为页面的浏览量计数。

Artalk 加载需要时间，所以你可以给它一个占位字符：

```html
<span id="ArtalkPV">加载中...</span>
```

你也可以绑定使用其他的元素，修改 Artalk 配置项，例如：

```js
Artalk.init({
  pvEl: '.your_element',
})
```

### 显示多个页面的浏览量

你能在除评论页面之外的任何页面，例如「文章列表」页，显示页面浏览量或评论数：

```js
Artalk.loadCountWidget({
  server:  '服务器地址',
  site:    '站点名',
  pvEl:    '#ArtalkPV',
  countEl: '#ArtalkCount',
});
```

在非评论页时，就无需再 Artalk.init 实例了 (否则页面浏览量 PV 数会无故增加)，仅调用 `loadCountWidget` 静态方法即可。

然后你可以放置多个 `#ArtalkPV` 元素，通过属性 `data-page-key` 来指定需要查询的页面：

```html
<span id="ArtalkPV" data-page-key="/test/1.html"></span>
<span id="ArtalkPV" data-page-key="/test/2.html"></span>
<span id="ArtalkPV" data-page-key="/test/3.html"></span>
<span id="ArtalkCount" data-page-key="/test/1.html"></span>
<!-- ... -->
```

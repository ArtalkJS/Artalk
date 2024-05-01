# 图片懒加载

Artalk 支持为评论中的图片添加懒加载功能，以减少页面加载时间。该功能默认禁用，可在控制中心「设置 - 界面配置」找到「图片懒加载」配置项 (frontend.imgLazyLoad)。该配置项有三个选项，分别是 `native`、`data-src` 和 `false`。

## 浏览器原生

将图片懒加载设置为 `native` 时，

Artalk 将使用浏览器原生的图片懒加载功能，即为图片标签添加 `loading="lazy"` 属性。例如：

```html
<img src="1.png" loading="lazy" />
```

请注意部分浏览器可能不支持该属性，或者浏览器实现的懒加载策略不同，详情请参考：[MDN - 懒加载](https://developer.mozilla.org/zh-CN/docs/Web/Performance/Lazy_loading)。

## 懒加载库

将图片懒加载设置为 `data-src` 时，

Artalk 将为图片标签添加 `data-src` 和 `class="lazyload"` 属性。例如：

```html
<img data-src="1.png" class="lazyload" />
```

这时，你需要在页面中引入一个图片懒加载库：[vanilla-lazyload](https://github.com/verlok/vanilla-lazyload)

在 HTML 的 `<head>` 中添加以下代码：

```html
<script src="https://cdn.jsdelivr.net/npm/vanilla-lazyload/dist/lazyload.min.js"></script>
```

然后编写 JS 代码，在 Artalk 初始化之前，添加监听事件：当评论渲染完成时，调用图片懒加载库的 `update` 方法。

```js
const lazyLoadInstance = new LazyLoad({
  elements_selector: ".lazyload",
})

Artalk.use((ctx) => {
  ctx.on('comment-rendered', () => lazyLoadInstance.update())
})

const artalk = Artalk.init({ /* ... */ })
```

推荐使用该方式，因为它兼容性更好，在不同浏览器下有一致的表现。


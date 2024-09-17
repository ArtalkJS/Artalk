# 图片懒加载

Artalk 支持为评论中的图片添加懒加载功能，以减少页面加载时间。该功能默认禁用，可在控制中心「设置 - 界面配置」找到「图片懒加载」配置项 (frontend.imgLazyLoad)：

| 配置值 | 说明 |
| --- | --- |
| `native` | 使用浏览器原生的图片懒加载 |
| `data-src` | 使用懒加载库实现图片懒加载 |
| `false` | 禁用图片懒加载 |

## 浏览器原生

当图片懒加载设置为 `native` 时，Artalk 将使用浏览器原生的图片懒加载功能，即为图片标签添加 `loading="lazy"` 属性。例如：

```html
<img src="1.png" loading="lazy" />
```

这不需要额外引入库和编写代码，是最简单的方式，但部分浏览器可能不支持该属性，或者实现的懒加载策略不同，详情参考：[MDN - 懒加载](https://developer.mozilla.org/zh-CN/docs/Web/Performance/Lazy_loading)。

## 懒加载库

使用该方式能确保在不同的浏览器下有一致的表现，并且能定制更多功能，如图片加载动画等。一个博客主题可能自带懒加载库，借助其附带的库，让 Artalk 也同时支持懒加载。

当图片懒加载设置为 `data-src` 时，Artalk 将为图片标签添加 `class="lazyload"` 和 `data-src` 属性。例如：

```html
<img data-src="1.png" class="lazyload" />
```

这时，你需要引入一个额外的图片懒加载库：[vanilla-lazyload](https://github.com/verlok/vanilla-lazyload)；在页面的 `<head>` 中添加以下代码：

```html
<script src="https://unpkg.com/vanilla-lazyload/dist/lazyload.iife.min.js"></script>
```

编写代码，在 Artalk 初始化之前，添加事件监听：当评论列表发生更新后，调用图片懒加载库的 `update` 方法：

```js
// 初始化图片懒加载库
const lazyLoadInstance = new LazyLoad({
  elements_selector: '.lazyload',
  threshold: 0,
})

// 监听 Artalk 事件
Artalk.use((ctx) => {
  ctx.on('list-loaded', () => lazyLoadInstance.update())
})

// 初始化 Artalk
const artalk = Artalk.init({ /* ... */ })
```

# Image Lazy Loading

Artalk supports adding lazy loading functionality to images within comments to reduce page load times. This feature is disabled by default. You can find the "Image Lazy Load" configuration item (frontend.imgLazyLoad) in the Dashboard under "Settings - UI Configuration":

| Configuration Value | Description |
| --- | --- |
| `native` | Use the browser's native image lazy loading |
| `data-src` | Use a lazy loading library for image lazy loading |
| `false` | Disable image lazy loading |

## Native Browser Support

When the image lazy loading is set to `native`, Artalk will use the browser's native image lazy loading feature, which adds the `loading="lazy"` attribute to the image tag. For example:

```html
<img src="1.png" loading="lazy" />
```

This method does not require additional libraries or code and is the simplest approach. However, some browsers may not support this attribute, or their lazy loading strategies may vary. For more details, refer to: [MDN - Lazy Loading](https://developer.mozilla.org/zh-CN/docs/Web/Performance/Lazy_loading).

## Lazy Loading Library

Using a lazy loading library ensures consistent behavior across different browsers and can offer additional features, such as image loading animations. If your blog theme comes with a lazy loading library, you can leverage it to enable lazy loading in Artalk as well.

When the image lazy loading is set to `data-src`, Artalk will add the `class="lazyload"` and `data-src` attributes to the image tag. For example:

```html
<img data-src="1.png" class="lazyload" />
```

In this case, you need to include an additional image lazy loading library, such as [vanilla-lazyload](https://github.com/verlok/vanilla-lazyload). Add the following code to the `<head>` of your page:

```html
<script src="https://unpkg.com/vanilla-lazyload/dist/lazyload.iife.min.js"></script>
```

Then, write the code to add an event listener before initializing Artalk. This listener will call the `update` method of the lazy loading library when the comment list is updated:

```js
// Initialize the image lazy loading library
const lazyLoadInstance = new LazyLoad({
  elements_selector: '.lazyload',
  threshold: 0,
})

// Listen to Artalk events
Artalk.use((ctx) => {
  ctx.on('list-loaded', () => lazyLoadInstance.update())
})

// Initialize Artalk
const artalk = Artalk.init({ /* ... */ })
```

By following these steps, you can ensure that images within Artalk comments are efficiently lazy-loaded, improving the overall performance and user experience of your site.

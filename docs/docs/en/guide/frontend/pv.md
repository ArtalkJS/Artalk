# Page View Statistics

Artalk integrates built-in page view and comment count statistics, allowing you to display the page's view count and comment count within the page.

```html
Views: <span class="artalk-pv-count"></span>
Comments: <span class="artalk-comment-count"></span>
```

## Loading Placeholder

Artalk requires some time to load. You can add a placeholder within the elements displaying the statistics:

```html
<span class="artalk-pv-count">Loading...</span>
```

## Loading Statistics for Multiple Pages Simultaneously

For instance, on an article list page, you can display the view count and comment count for each article.

When on the article list page, simply call the `Artalk.loadCountWidget` function without `Artalk.init` (loading the comment list will increment the page view count by 1).

<!-- prettier-ignore-start -->

```js
Artalk.loadCountWidget({
  server: 'server address',
  site: 'site name',
  pvEl: '.artalk-pv-count',
  countEl: '.artalk-comment-count',
  statPageKeyAttr: 'data-page-key',
})
```

<!-- prettier-ignore-end -->

Then place multiple elements with the class `artalk-pv-count`, using the attribute `data-page-key` to specify the page to query:

```html
<span class="artalk-pv-count" data-page-key="/test/1.html"></span>
<span class="artalk-pv-count" data-page-key="/test/2.html"></span>
<span class="artalk-pv-count" data-page-key="/test/3.html"></span>
```

The same applies for comment count queries:

```html
<span class="artalk-comment-count" data-page-key="/test/1.html"></span>
<span class="artalk-comment-count" data-page-key="/test/2.html"></span>
<span class="artalk-comment-count" data-page-key="/test/3.html"></span>
```

## Custom Element Selectors

You can specify element selectors for displaying page views and comment counts through the `pvEl` and `countEl` configuration options:

```js
Artalk.init({
  pvEl: '.your_pv_count_element',  // Selector for page view elements
  countEl: '.your_comment_count_element',  // Selector for comment count elements
})
```

## Custom `data-page-key` Attribute Name

The default value for the `statPageKeyAttr` configuration option is `data-page-key`. Artalk uses this attribute name to query the specified page. To facilitate adaptation with blog themes, you can customize the attribute name, for example, replacing it with `data-path`:

```js
Artalk.loadCountWidget({
  statPageKeyAttr: 'data-path',
})
```

In this case, the corresponding HTML code should be as follows:

```html
<span class="artalk-pv-count" data-path="/test/1.html"></span>
```

Thus, the value of the `data-path` attribute will be used to query the specified page.

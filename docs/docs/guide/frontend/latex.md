# Latex

对于学术性站点，支持 Latex 是一个刚需，可以使用 [Katex](https://katex.org/) 这款被前端广泛使用的 Latex 语法解析器。Artalk 为保证简洁性，并未内置 Latex 语法解析功能，但你可以在页面中引入 Latex 和 Artalk Katex 插件，让评论系统获得 Latex 语法支持。

```html
<!-- 1. 引入 Katex -->
<link
  rel="stylesheet"
  href="https://unpkg.com/katex@0.15.3/dist/katex.min.css"
/>
<script defer src="https://unpkg.com/katex@0.15.3/dist/katex.min.js"></script>

<!-- 2. 引入 Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. 引入 Artalk Katex 插件 -->
<script src="https://unpkg.com/@artalk/plugin-katex/dist/artalk-plugin-katex.js"></script>
```

之后，你可以回复：

```md
$$ P(A) = \sum P(\{ (e_1,...,e_N) \}) = {{N}\choose{k}} \cdot p^kq^{N-k} $$
```

查看效果：

![](/images/latex-support/1.png)

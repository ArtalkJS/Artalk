# LaTeX

对于学术网站来说，支持 LaTeX 是一个刚需，可以使用 [Katex](https://katex.org/) 这款被前端广泛使用的 LaTeX 语法解析器来继承该功能。Artalk 为保持简洁，并未内置 LaTeX 语法解析功能，你可以手动在页面中引入 Katex 和 Artalk Katex 插件，让评论系统获得 LaTeX 语法支持。

```html
<!-- 1. 引入 Katex -->
<link rel="stylesheet" href="https://unpkg.com/katex@0.16.7/dist/katex.min.css" />
<script defer src="https://unpkg.com/katex@0.16.7/dist/katex.min.js"></script>

<!-- 2. 引入 Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. 引入 Artalk Katex 插件 -->
<script src="https://unpkg.com/@artalk/plugin-katex/dist/artalk-plugin-katex.js"></script>
```

你可以发表评论：

```md
$$
P(A|B_1, B_2, \ldots, B_n) = \frac{P(B_1, B_2, \ldots, B_n|A) \cdot P(A)}{P(B_1, B_2, \ldots, B_n)}
$$
```

检验效果：

![](/images/latex-support/1.png)

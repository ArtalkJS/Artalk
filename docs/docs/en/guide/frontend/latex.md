# LaTeX

For academic websites, LaTeX support is essential. You can use [Katex](https://katex.org/), a widely-used LaTeX syntax parser in the frontend, to implement this functionality. To maintain simplicity, Artalk does not natively include LaTeX syntax parsing capabilities. However, you can manually incorporate Katex and the Artalk Katex plugin into your page to enable LaTeX syntax support in the comment system.

::: code-group

```html [Browser Inclusion]
<!-- 1. Include Katex -->
<link rel="stylesheet" href="https://unpkg.com/katex@0.16.7/dist/katex.min.css" />
<script src="https://unpkg.com/katex@0.16.7/dist/katex.min.js"></script>

<!-- 2. Include Artalk -->
<link href="/lib/artalk/Artalk.css" rel="stylesheet" />
<script src="/lib/artalk/Artalk.js"></script>

<!-- 3. Include Artalk Katex Plugin -->
<script src="https://unpkg.com/@artalk/plugin-katex/dist/artalk-plugin-katex.js"></script>
```

```ts [Node Inclusion]
// pnpm add katex @artalk/plugin-katex
import Artalk from 'artalk'
import { ArtalkKatexPlugin } from '@artalk/plugin-katex'
import 'katex/dist/katex.min.css'

Artalk.use(ArtalkKatexPlugin)
```

:::

You can post a comment:

```md
$$
P(A|B_1, B_2, \ldots, B_n) = \frac{P(B_1, B_2, \ldots, B_n|A) \cdot P(A)}{P(B_1, B_2, \ldots, B_n)}
$$
```

Check the result:

![](/images/latex-support/1.png)

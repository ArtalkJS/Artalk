# 置入框架

## 引入 Artalk

通过包管理工具引入 Artalk，推荐使用 [pnpm](https://pnpm.io/zh/)

```bash
pnpm add artalk
```

## Vue

Vue 3 + TypeScript 例：

```vue
<template>
  <div class="artalk-comments"></div>
</template>

<script lang="ts">
import 'artalk/dist/Artalk.css'

import { defineComponent } from 'vue'
import Artalk from 'artalk'

export default defineComponent({
  mounted: () => {
    const artalk = Artalk.init({
      el:        this.$el,
      pageKey:   `${location.pathname}`,
      pageTitle: `${document.title}`,
      server:    'http://localhost:8080',
      site:      'Artalk 的博客',
      // ...
    })
  }
})
</script>
```

::: tip 提示

VuePress 可参考：[“VuePress 引入”](./import-blog.md#vuepress)

:::

## React

```jsx
import 'artalk/dist/Artalk.css'

import React, { createRef } from 'react'
import Artalk from 'artalk'

export default class Artalk extends React.Component {
  el = createRef()

  componentDidMount () {
    const artalk = Artalk.init({
      el: this.el.current,
      pageKey:   `${location.pathname}`,
      pageTitle: `${document.title}`,
      server:    'http://localhost:8080',
      site:      'Artalk 的博客',
      // ...
    })
  }

  render () {
    return (
      <div ref={this.el} class="artalk-comments" />
    )
  }
}
```

## Svelte

```html
<script>
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

let comments;

onMount(() => {
  Artalk.init({
    el:        comments,
    pageKey:   `${location.pathname}`,
    pageTitle: `${document.title}`,
    server:    'http://localhost:8080',
    site:      'Artalk 的博客',
    // ...
  })
})
</script>

<div bind:this={comments} class="artalk-comments"></div>
```

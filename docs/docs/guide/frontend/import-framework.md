# 置入框架

## 安装 Artalk

通过包管理工具引入 Artalk，推荐使用 [pnpm](https://pnpm.io/zh/)

```bash
pnpm add artalk
```

## Vue

Vue 3 + TypeScript 例：

```vue
<script lang="ts" setup>
import Artalk from 'artalk'
import { onBeforeUnmounted, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import 'artalk/dist/Artalk.css'

const el = ref<HTMLElement>()
const route = useRoute()

let artalk: Artalk

onMounted(() => {
  artalk = Artalk.init({
    el: el.value,
    pageKey: route.path,
    pageTitle: `${document.title}`,
    server: 'http://localhost:8080',
    site: 'Artalk 的博客',
    // ...
  })
})

onBeforeUnmounted(() => {
  artalk.destroy()
})
</script>

<template>
  <view ref="el"></view>
</template>
```

## React

```tsx
import React, { useEffect, useRef } from 'react'
import { useLocation } from 'react-router-dom'
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

const ArtalkComment = () => {
  const el = useRef(null)
  const location = useLocation()
  let artalk = null

  useEffect(() => {
    artalk = new Artalk({
      el: el.current,
      pageKey: location.pathname,
      pageTitle: document.title,
      server: 'http://localhost:8080',
      site: 'Artalk 的博客',
      // ...
    })

    return () => {
      if (artalk) {
        artalk.destroy()
      }
    }
  }, [location.pathname])

  return <div ref={el}></div>
}

export default ArtalkComment
```

```jsx
import React, { createRef } from 'react'
import 'artalk/dist/Artalk.css'
import Artalk from 'artalk'

export default class Artalk extends React.Component {
  el = createRef()
  artalk = null

  componentDidMount() {
    this.artalk = Artalk.init({
      el: this.el.current,
      pageKey: `${location.pathname}`,
      pageTitle: `${document.title}`,
      server: 'http://localhost:8080',
      site: 'Artalk 的博客',
      // ...
    })
  }

  componentWillUnmount() {
    this.artalk?.destroy()
  }

  render() {
    return <div ref={this.el} />
  }
}
```

## Svelte

```html
<script>
  import Artalk from 'artalk'
  import { onMount, onDestroy } from 'svelte'

  import 'artalk/dist/Artalk.css'

  let el
  let artalk

  onMount(() => {
    artalk = Artalk.init({
      el: el,
      pageKey: location.pathname,
      pageTitle: document.title,
      server: 'http://localhost:8080',
      site: 'Artalk 的博客',
      // ...
    })

    onDestroy(() => {
      if (artalk) {
        artalk.destroy()
      }
    })
  })
</script>

<div bind:this="{el}"></div>
```

# 置入框架

## 安装 Artalk

通过包管理工具引入 Artalk，推荐使用 [pnpm](https://pnpm.io/zh/)

```bash
pnpm add artalk
```

## Vue

::: code-group

```vue [Vue 3]
<script lang="ts" setup>
import Artalk from 'artalk'
import { onMounted, onBeforeUnmount, ref } from 'vue'

import 'artalk/Artalk.css'

const el = ref<HTMLElement>()

let artalk: Artalk

onMounted(() => {
  artalk = Artalk.init({
    el: el.value,
    pageKey: location.pathname,
    pageTitle: `${document.title}`,
    server: 'http://localhost:8080',
    site: 'Artalk 的博客',
    // ...
  })
})

onBeforeUnmount(() => {
  artalk.destroy()
})
</script>

<template>
  <view ref="el"></view>
</template>
```

```vue [Vue 3 + Vue Router]
<script lang="ts" setup>
import Artalk from 'artalk'
import { onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'

import 'artalk/Artalk.css'

const el = ref<HTMLElement>()
const router = useRouter()

let artalk: Artalk

onMounted(() => {
  artalk = Artalk.init({
    el: el.value,
    pageKey: router.route.path,
    pageTitle: document.title,
    server: 'http://localhost:8080',
    site: 'Artalk 的博客',
    // ...
  })
})

watch(
  () => router.route.path,
  (path) => {
    nextTick(() => {
      artalk.update({
        pageKey: path,
        pageTitle: document.title,
      })
      artalk.reload()
    })
  }
)

onBeforeUnmount(() => {
  artalk.destroy()
})
</script>

<template>
  <view ref="el"></view>
</template>
```

:::

## React

::: code-group

```tsx [React Hooks]
import { useCallback, useRef } from 'react'
import { useLocation } from 'react-router-dom'
import 'artalk/Artalk.css'
import Artalk from 'artalk'

const ArtalkComment = () => {
  const { pathname } = useLocation()
  const artalk = useRef<Artalk>()

  const handleContainerInit = useCallback(
    (node: HTMLDivElement | null) => {
      if (!node) {
        return
      }
      if (artalk.current) {
        artalk.current.destroy()
        artalk.current = undefined
      }
      artalk.current = Artalk.init({
        el: node,
        pageKey: pathname,
        pageTitle: document.title,
        server: 'http://localhost:8080',
        site: 'Artalk 的博客',
        // ...
      })
    },
    [pathname],
  )

  return <div ref={handleContainerInit}></div>
}

export default ArtalkComment
```



```jsx [React Class]
import React, { createRef } from 'react'
import 'artalk/Artalk.css'
import Artalk from 'artalk'

export default class Artalk extends React.Component {
  el = createRef()
  artalk = null

  componentDidMount() {
    this.artalk = Artalk.init({
      el: this.el.current,
      pageKey: location.pathname,
      pageTitle: document.title,
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

:::

## SolidJS

```tsx
import { onCleanup, onMount } from 'solid-js'
import Artalk from 'artalk'
import 'artalk/Artalk.css'

const ArtalkComment = () => {
  let el: HTMLDivElement
  let artalk: Artalk

  onMount(() => {
    artalk = Artalk.init({
      el: el,
      pageKey: location.pathname,
      pageTitle: document.title,
      server: 'http://localhost:8080',
      site: 'Artalk 的博客',
      // ...
    })

    onCleanup(() => {
      artalk.destroy()
    })
  })

  return <div ref={el} />
}
```

## Svelte

```html
<script>
import Artalk from 'artalk'
import { onMount, onDestroy } from 'svelte'

import 'artalk/Artalk.css'

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

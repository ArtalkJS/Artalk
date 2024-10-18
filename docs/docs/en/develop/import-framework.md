# Import to Frameworks

## Install Artalk

Install Artalk using a package manager. It is recommended to use [pnpm](https://pnpm.io/zh/).

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
    site: 'Artalk Blog',
    // ...
  })
})

onBeforeUnmount(() => {
  artalk.destroy()
})
</script>

<template>
  <div ref="el"></div>
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
    pageKey: router.currentRoute.value.path,
    pageTitle: document.title,
    server: 'http://localhost:8080',
    site: 'Artalk Blog',
    // ...
  })
})

watch(
  () => router.currentRoute.value.path,
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
  <div ref="el"></div>
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
        site: 'Artalk Blog',
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

export default class ArtalkComponent extends React.Component {
  el = createRef()
  artalk = null

  componentDidMount() {
    this.artalk = Artalk.init({
      el: this.el.current,
      pageKey: location.pathname,
      pageTitle: document.title,
      server: 'http://localhost:8080',
      site: 'Artalk Blog',
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
      site: 'Artalk Blog',
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
    site: 'Artalk Blog',
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

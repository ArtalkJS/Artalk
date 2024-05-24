<template>
  <div ref="el" style="margin-top: 20px"></div>
</template>

<script setup lang="ts">
import { watch, nextTick, ref, onMounted, onUnmounted } from 'vue'
import { useData, useRouter } from 'vitepress'
import Artalk from 'artalk'
import 'artalk/dist/artalk.css'

Artalk.use((ctx) => {
  console.log(ctx)
})

const el = ref<HTMLElement>()

const router = useRouter()
const page = useData().page

let artalk: Artalk

onMounted(() => {
  nextTick(() => {
    initArtalk(getConfByPage())
  })
})

watch(
  () => router.route.path,
  () => {
    nextTick(() => {
      artalk.update(getConfByPage())
      artalk.reload()
    })
  },
)

onUnmounted(() => {
  artalk.destroy()
})

function initArtalk(conf: any) {
  artalk = Artalk.init({
    el: el.value,
    emoticons: '/assets/emoticons/default.json',
    gravatar: {
      mirror: 'https://weavatar.com/avatar/',
    },
    ...conf,
  })

  loadExtraFuncs()
}

function getConfByPage() {
  return {
    pageKey: 'https://artalk.js.org' + router.route.path,
    pageTitle: page.value.title,
    server: 'http://localhost:23366',
    site: 'ArtalkDocs',
  }
}

function loadExtraFuncs() {
  // 夜间模式
  const darkMode = document.querySelector('html')!.classList.contains('dark')
  artalk.setDarkMode(darkMode)

  new MutationObserver((mList) => {
    mList.forEach((m) => {
      if (m.attributeName !== 'class') return

      // @ts-ignore
      const darkMode = m.target.classList.contains('dark')
      artalk.setDarkMode(darkMode)
    })
  }).observe(document.querySelector('html')!, { attributes: true })
}
</script>

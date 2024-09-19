<template>
  <div ref="el" style="margin-top: 20px"></div>
</template>

<script setup lang="ts">
import { watch, nextTick, ref, onMounted, onUnmounted } from 'vue'
import { useData, useRouter } from 'vitepress'
import Artalk from 'artalk'
import { ArtalkKatexPlugin } from '@artalk/plugin-katex'
import { ArtalkLightboxPlugin } from '@artalk/plugin-lightbox'
import 'lightgallery/css/lightgallery.css'
import 'katex/dist/katex.min.css'
import 'artalk/Artalk.css'

const el = ref<HTMLElement | null>(null)

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
  Artalk.use(ArtalkKatexPlugin)
  Artalk.use(ArtalkLightboxPlugin, {
    lightGallery: {
      lib: async () => (await import('lightgallery')).default,
    },
  })

  Artalk.use((ctx) => {
    ctx.on('mounted', () => {
      // TODO: Optimize the locale setting, which is overwritten by the remote server config
      const locale = document.documentElement.lang.includes('zh') ? 'zh-CN' : 'en'
      if (locale !== ctx.getConf().locale) ctx.updateConf({ locale })
    })
  })

  artalk = Artalk.init({
    el: el.value,
    emoticons: '/assets/emoticons/default.json',
    ...conf,
  })

  loadExtraFuncs()
}

function getConfByPage() {
  return {
    pageKey: `https://artalk.js.org/${router.route.path
      .replace(/^\/+/, '')
      .replace(/^zh\//, '')
      .replace(/\.html$/, '')}.html`,
    pageTitle: page.value.title,
    server: 'https://artalk.qwqaq.com',
    site: 'ArtalkDocs',
  }
}

function loadExtraFuncs() {
  // 夜间模式
  const darkMode = document.querySelector('html').classList.contains('dark')
  artalk.setDarkMode(darkMode)

  new MutationObserver((mList) => {
    mList.forEach((m) => {
      if (m.attributeName !== 'class') return

      const darkMode = (<HTMLElement>m.target).classList.contains('dark')
      artalk.setDarkMode(darkMode)
    })
  }).observe(document.querySelector('html'), { attributes: true })
}
</script>

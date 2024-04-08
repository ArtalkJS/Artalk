<template>
  <div ref="el" style="margin-top: 20px"></div>
</template>

<script setup lang="ts">
import { watch, nextTick, ref, onMounted, onUnmounted } from 'vue'
import { useData, useRouter } from 'vitepress'
import Artalk from 'artalk'
import 'artalk/dist/Artalk.css'

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
  artalk = Artalk.init({
    el: el.value,
    emoticons: '/assets/emoticons/default.json',
    gravatar: {
      mirror: 'https://cravatar.cn/avatar/',
    },
    ...conf,
  })

  loadExtraFuncs()
}

function getConfByPage() {
  return {
    pageKey: 'https://artalk.js.org' + router.route.path,
    pageTitle: page.value.title,
    server: 'https://artalk.qwqaq.com',
    site: 'ArtalkDocs',
  }
}

function loadExtraFuncs() {
  // 图片灯箱插件
  artalk.on('list-loaded', () => {
    document
      .querySelectorAll('.atk-comment .atk-content')
      .forEach(($content) => {
        const imgEls = $content.querySelectorAll<HTMLImageElement>(
          'img:not([atk-emoticon]):not([atk-lightbox])',
        )
        imgEls.forEach((imgEl) => {
          imgEl.setAttribute('atk-lightbox', '')
          const linkEl = document.createElement('a')
          linkEl.setAttribute('class', 'atk-img-link')
          linkEl.setAttribute('href', imgEl.src)
          linkEl.setAttribute('data-src', imgEl.src)
          linkEl.append(imgEl.cloneNode())
          imgEl.replaceWith(linkEl)
        })
        // @ts-ignore
        if (imgEls.length) lightGallery($content, { selector: '.atk-img-link' })
      })
  })

  // 夜间模式
  const darkMode = document.querySelector('html').classList.contains('dark')
  artalk.setDarkMode(darkMode)

  new MutationObserver((mList) => {
    mList.forEach((m) => {
      if (m.attributeName !== 'class') return

      // @ts-ignore
      const darkMode = m.target.classList.contains('dark')
      artalk.setDarkMode(darkMode)
    })
  }).observe(document.querySelector('html'), { attributes: true })
}
</script>

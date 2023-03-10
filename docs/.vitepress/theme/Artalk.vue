<template>
  <div ref="artalkEl" style="margin-top: 20px;"></div>
</template>

<script setup lang="ts">
import { watch, nextTick, ref, onMounted } from 'vue'
import { useData, useRouter } from 'vitepress'

const artalkEl = ref<HTMLElement | null>(null)

const router = useRouter()
const page = useData().page

onMounted(() => {
  const script = document.createElement('script')
  script.src = `https://artalk.qwqaq.com/dist/Artalk.js`
  document.getElementsByTagName('head')[0].appendChild(script)
  script.onload = () => {
    initArtalk(page.value)
  }
})

watch(() => router.route.data.relativePath, (path) => {
  nextTick(() => {
    Artalk.update(getArtalkConfByPage(page.value))
    Artalk.reload()
  })
})

function getArtalkConfByPage(page: any) {
  return {
    pageKey:   'https://artalk.js.org'+location.pathname,
    pageTitle:  page.title,
    server:    'https://artalk.qwqaq.com',
    site:      'ArtalkDocs',
  }
}

function initArtalk(page: any) {
  Artalk.init({
    el:        artalkEl.value,
    emoticons: '/assets/emoticons/default.json',
    gravatar:   {
      mirror: 'https://cravatar.cn/avatar/'
    },
    ...getArtalkConfByPage(page)
  })

  // 图片灯箱插件
  Artalk.on('list-loaded', () => {
    document.querySelectorAll('.atk-comment .atk-content').forEach(($content) => {
      const imgEls = $content.querySelectorAll<HTMLImageElement>('img:not([atk-emoticon]):not([atk-lightbox])');
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
  Artalk.setDarkMode(darkMode)

  new MutationObserver((mList) => {
    mList.forEach((m) => {
      if (m.attributeName !== 'class') return

      // @ts-ignore
      const darkMode = m.target.classList.contains('dark')
      Artalk.setDarkMode(darkMode)
    })
  }).observe(document.querySelector('html'), { attributes: true })
}
</script>

<style>

</style>

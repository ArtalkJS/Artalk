<template>
  <div></div>
</template>

<script lang="ts">
import { defineComponent, watch, nextTick } from 'vue'
import { useData, useRouter } from 'vitepress'
import * as ArtalkCDN from '../../code/ArtalkCDN.json'

let isFirstLoad = true

function loadArtalk() {
  const script = document.createElement('script')
  script.async = true
  script.defer = true
  script.src = ArtalkCDN.JS
  document.getElementsByTagName('head')[0].appendChild(script)
  return script
}

function initArtalk(pageData: any) {
  const pEl = document.querySelector<HTMLElement>('.VPDoc .container > .content > .content-container')
  if (!pEl) return
  let artalkEl = pEl.querySelector<HTMLElement>('#ArtalkComment')
  if (artalkEl) artalkEl.remove()
  artalkEl = document.createElement('div')
  artalkEl.id = 'ArtalkComment'
  artalkEl.style.marginTop = '20px'
  pEl.appendChild(artalkEl)

  const conf = {
    el:        '#ArtalkComment',
    pageKey:   `https://artalk.js.org${location.pathname}`,
    pageTitle:  pageData.title,
    server:    'https://artalk.qwqaq.com',
    site:      'ArtalkDocs',
    emoticons: '/assets/emoticons/default.json',
    gravatar:   {
      mirror: 'https://cravatar.cn/avatar/'
    }
  }

  const artalk = new (window as any).Artalk(conf);

  confArtalk(artalk)
}

function confArtalk(artalk) {
  // lightGallery
  artalk.on('list-loaded', () => {
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

  // dark_mode
  const darkMode = document.querySelector('html').classList.contains('dark')
  artalk.setDarkMode(darkMode)

  // dark_mode 监听
  new MutationObserver((mList) => {
    mList.forEach((m) => {
      if (m.attributeName !== 'class') return

      // @ts-ignore
      const darkMode = m.target.classList.contains('dark')
      artalk.setDarkMode(darkMode)
    })
  }).observe(document.querySelector('html'), { attributes: true })
}

export default defineComponent({
  mounted: () => {
    const pageData = useData().page.value
    const router = useRouter()

    loadArtalk().onload = () => {
      initArtalk(pageData)
      isFirstLoad = false
    }

    watch(() => router.route.data.relativePath, (path) => {
      if (isFirstLoad) return
      nextTick(() => {
        initArtalk(pageData)
      })
    }, { immediate: false });
  }
})
</script>

<style>

</style>

<template>
  <iframe
    src="https://artransfer.netlify.app/?iframe=1"
    style="width: 100%; height: 520px; border: 0"
    id="artransferIframe"
  ></iframe>
</template>

<script lang="ts">
import { defineComponent, nextTick } from 'vue'

function setArtransferIframeDarkMode(value: boolean) {
  const iframe = document.querySelector<HTMLIFrameElement>('#artransferIframe')!
  if (value) {
    iframe.src = iframe.src + '&dark=1'
  } else {
    iframe.src = iframe.src.replace(/&dark=1$/, '')
  }
}

export default defineComponent({
  mounted: () => {
    const darkMode = document.querySelector('html').classList.contains('dark')
    setArtransferIframeDarkMode(darkMode)

    // dark_mode 监听
    new MutationObserver((mList) => {
      mList.forEach((m) => {
        if (m.attributeName !== 'class') return

        // @ts-ignore
        const darkMode = m.target.classList.contains('dark')
        setArtransferIframeDarkMode(darkMode)
      })
    }).observe(document.querySelector('html'), { attributes: true })
  },
})
</script>

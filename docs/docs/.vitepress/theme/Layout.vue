<template>
  <Layout>
    <template #doc-after>
      <Artalk v-if="frontmatter.comments !== false" />
    </template>
    <template #aside-outline-after>
      <div
        v-if="desktopAside && isArticle"
        :key="page.relativePath"
        v-display-ad
        class="site-ad aside-ad"
        :data-label="adLabel"
      />
    </template>
    <template #doc-footer-before>
      <div
        v-if="isArticle"
        :key="page.relativePath"
        v-display-ad
        class="site-ad"
        :data-label="adLabel"
      />
    </template>
  </Layout>
</template>

<script setup>
import DefaultTheme from 'vitepress/theme'
import { inBrowser, useData } from 'vitepress'
import { computed, onMounted, ref, watchEffect, watchPostEffect } from 'vue'
import './style.scss'
import { loadAdSense, showDisplayAd } from '../../../shared/adsense'
import { loadClarity } from '../../../shared/clarity'

const { Layout } = DefaultTheme
const { lang, page, frontmatter } = useData()

watchEffect(() => {
  if (inBrowser)
    document.cookie = `nf_lang=${lang.value}; expires=Mon, 1 Jan 2030 00:00:00 UTC; path=/`
})

const isArticle = computed(() => /^(en|zh)\/(guide|develop)\//.test(page.value.relativePath))
const adLabel = computed(() => (lang.value.startsWith('zh') ? '广告' : 'Advertisements'))
const desktopAside = ref(false)
const vDisplayAd = { mounted: showDisplayAd }

watchPostEffect((onCleanup) => {
  const query = window.matchMedia('(min-width: 1280px) and (hover: hover) and (pointer: fine)')
  const update = () => {
    desktopAside.value = query.matches
  }
  update()
  query.addEventListener('change', update)
  onCleanup(() => query.removeEventListener('change', update))
})

onMounted(loadAdSense)
onMounted(loadClarity)
</script>

<style scoped>
.aside-ad {
  max-width: 200px;
  margin-inline: 0;
}

.aside-ad :deep(ins) {
  height: 200px;
}
</style>

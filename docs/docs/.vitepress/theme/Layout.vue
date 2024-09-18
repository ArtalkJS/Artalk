<script setup>
import DefaultTheme from 'vitepress/theme'
import { inBrowser, useData } from 'vitepress'
import { watchEffect } from 'vue'
import './style.scss'

const { Layout } = DefaultTheme

const { lang } = useData()
watchEffect(() => {
  if (inBrowser)
    document.cookie = `nf_lang=${lang.value}; expires=Mon, 1 Jan 2030 00:00:00 UTC; path=/`
})
</script>

<template>
  <Layout>
    <template #doc-after>
      <Artalk />
    </template>
  </Layout>
</template>

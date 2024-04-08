<script setup>
import DefaultTheme from 'vitepress/theme'
import { onMounted } from 'vue'
import './style.scss'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #doc-after>
      <Artalk />
    </template>
  </Layout>
</template>

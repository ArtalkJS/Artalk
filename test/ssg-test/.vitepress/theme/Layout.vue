<script setup>
import DefaultTheme from 'vitepress/theme'
import Artalk from './Artalk.vue'
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

<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import Artalk from 'artalk'
import 'artalk/dist/Artalk.css'

const el = ref<HTMLElement>()
const route = useRoute()

let artalk: Artalk

onMounted(() => {
  artalk = Artalk.init({
    el: el.value,
    server: 'http://localhost:23366',
    site: 'ArtalkDocs',
    pageKey: route.path,
  })
})

onUnmounted(() => {
  artalk.destroy()
})
</script>

<template>
  <view ref="el"></view>
</template>

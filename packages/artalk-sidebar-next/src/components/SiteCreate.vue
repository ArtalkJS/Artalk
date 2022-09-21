<script setup lang="ts">
import type { SiteData } from 'artalk/types/artalk-data'

const emit = defineEmits(['close', 'done'])

const site = ref<{name: string, urls: string}>({
  name: '',
  urls: ''
})

function createSite() {
  emit('done')
}

function close() {
  emit('close')
}
</script>

<template>
  <div class="atk-site-add">
    <div class="atk-header">
      <div class="atk-title">新增站点</div>
      <div class="atk-close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-form">
      <input v-model="site.name" type="text" name="AtkSiteName" placeholder="站点名称" autocomplete="off">
      <input v-model="site.urls" type="text" name="AtkSiteUrls" placeholder="站点 URL（多个用逗号隔开）" autocomplete="off">
      <button class="atk-btn" name="AtkSubmit" @submit="createSite()">创建</button>
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>

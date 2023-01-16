<script setup lang="ts">
import { artalk } from '../global'
import type { SiteData } from 'artalk/types/artalk-data'

const props = defineProps<{
  initVal?: { name: string, urls: string }
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'done', siteNew: SiteData): void
}>()

const isLoading = ref(false)
const site = ref<{name: string, urls: string}>({
  name: '',
  urls: ''
})

onMounted(() => {
  site.value.name = props.initVal?.name || ''
  site.value.urls = props.initVal?.urls || ''
})

async function submit() {
  const siteName = site.value.name.trim()
  const siteUrls = site.value.urls.trim()

  if (siteName === '') { alert('请输入站点名称'); return }

  isLoading.value = true
  let s: SiteData
  try {
    s = await artalk!.ctx.getApi().site.siteAdd(siteName, siteUrls)
  } catch (err: any) {
    window.alert(`创建失败：${err.msg || ''}`)
    console.error(err)
    return
  } finally { isLoading.value = false }

  emit('done', s)
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
    <form class="atk-form" @submit.prevent="submit()">
      <input v-model="site.name" type="text" name="AtkSiteName" placeholder="站点名称" autocomplete="off">
      <input v-model="site.urls" type="text" name="AtkSiteUrls" placeholder="站点 URL（多个用逗号隔开）" autocomplete="off">
      <button type="submit" class="atk-btn" name="AtkSubmit">创建</button>
    </form>
    <LoadingLayer v-if="isLoading" />
  </div>
</template>

<style scoped lang="scss">

</style>

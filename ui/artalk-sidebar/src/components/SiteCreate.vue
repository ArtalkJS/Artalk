<script setup lang="ts">
import { artalk } from '../global'
import type { ArtalkType } from 'artalk'

const { t } = useI18n()

const props = defineProps<{
  initVal?: { name: string, urls: string }
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'done', siteNew: ArtalkType.SiteData): void
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
  const siteUrls = site.value.urls.trim().split(',').map((v) => v.trim()).filter((v) => !!v)

  if (siteName === '') { alert('请输入站点名称'); return }

  isLoading.value = true
  let s: ArtalkType.SiteData
  try {
    s = (await artalk!.ctx.getApi().sites.createSite({
      name: siteName,
      urls: siteUrls
    })).data
  } catch (err: any) {
    window.alert(`创建失败：${err.message || ''}`)
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
      <div class="atk-title">{{ t('createSite') }}</div>
      <div class="atk-close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <form class="atk-form" @submit.prevent="submit()">
      <input v-model="site.name" type="text" name="AtkSiteName" :placeholder="t('siteName')" autocomplete="off">
      <input v-model="site.urls" type="text" name="AtkSiteUrls" :placeholder="`${t('siteUrls')} (${t('multiSepHint')})`" autocomplete="off">
      <button type="submit" class="atk-btn" name="AtkSubmit">{{ t('add') }}</button>
    </form>
    <LoadingLayer v-if="isLoading" />
  </div>
</template>

<style scoped lang="scss">

</style>

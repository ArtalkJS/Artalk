<script setup lang="ts">
import type { SiteData } from 'artalk/types/artalk-data'
import { artalk } from '../global'

const props = defineProps<{
  site: SiteData
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'update', page: SiteData): void
  (evt: 'remove', id: number): void
}>()

const { site } = toRefs(props)
const isLoading = ref(false)
const editFieldKey = ref<keyof SiteData|null>(null)
const editFieldVal = computed(() => {
  if (editFieldKey.value === 'urls') return site.value.urls_raw || ''
  return String(editFieldKey ? site.value[editFieldKey.value!] || '' : '')
})

function close() {
  emit('close')
}

function openURL(url: string) {
  window.open(url)
}

function rename() {
  editFieldKey.value = 'name'
}

function editURL() {
  editFieldKey.value = 'urls'
}

function del() {
  const del = async () => {
    isLoading.value = true
    try {
      await artalk!.ctx.getApi().site.siteDel(site.value.id, true)
    } catch (err: any) {
      console.log(err)
      alert(`删除失败 ${String(err)}`)
      return
    } finally { isLoading.value = false }
    emit('remove', site.value.id)
  }
  if (window.confirm(
    `确认删除站点 "${site.value.name}"？将会删除所有相关数据`
  )) del()
}

async function onFieldEditorYes(val: string) {
  if (editFieldVal.value !== val) {
    isLoading.value = true
    let s: SiteData
    try {
      s = await artalk!.ctx.getApi().site.siteEdit({ ...site.value, [editFieldKey.value as any]: val })
    } catch (err: any) {
      alert(`修改失败：${err.msg || '未知错误'}`)
      console.error(err)
      return false
    } finally { isLoading.value = false }
    emit('update', s)
  }

  editFieldKey.value = null
  return true
}

function onFiledEditorNo() {
  editFieldKey.value = null
}
</script>

<template>
  <div class="atk-site-edit">
    <div class="atk-header">
      <div class="atk-site-info">
        <span
          class="atk-site-name"
          @click="!site.first_url || openURL(site.first_url)"
        >{{ site.name }}</span>
        <span class="atk-site-urls">
          <div
            v-for="(url) in site.urls"
            class="atk-url-item"
            @click="openURL(url)"
          >{{ url }}</div>
        </span>
      </div>
      <div class="atk-close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-main">
      <div class="atk-site-text-actions">
        <div class="atk-item atk-rename-btn" @click="rename()">重命名</div>
        <div class="atk-item atk-edit-url-btn" @click="editURL()">修改 URL</div>
        <!--<div class="atk-item atk-export-btn">导出</div>
        <div class="atk-item atk-import-btn">导入</div>-->
      </div>
      <div class="atk-site-btn-actions">
        <div class="atk-item atk-del-btn" @click="del()">
          <i class="atk-icon atk-icon-del"></i>
        </div>
      </div>
      <LoadingLayer v-if="isLoading" style="z-index: 1000" />
      <ItemTextEditor
        v-if="!!editFieldKey"
        :init-value="editFieldVal"
        @yes="onFieldEditorYes"
        @no="onFiledEditorNo"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>

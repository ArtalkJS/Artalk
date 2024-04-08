<script setup lang="ts">
import type { ArtalkType } from 'artalk'
import { artalk } from '../global'

const props = defineProps<{
  site: ArtalkType.SiteData
}>()

const emit = defineEmits<{
  (evt: 'close'): void
  (evt: 'update', page: ArtalkType.SiteData): void
  (evt: 'remove', id: number): void
}>()

const { site } = toRefs(props)
const isLoading = ref(false)
const editFieldKey = ref<keyof ArtalkType.SiteData | null>(null)
const editFieldVal = computed(() => {
  if (editFieldKey.value === 'urls') return site.value.urls_raw || ''
  return String(editFieldKey ? site.value[editFieldKey.value!] || '' : '')
})

const { t } = useI18n()

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
      await artalk!.ctx.getApi().sites.deleteSite(site.value.id)
    } catch (err: any) {
      console.log(err)
      alert(`删除失败 ${String(err)}`)
      return
    } finally {
      isLoading.value = false
    }
    emit('remove', site.value.id)
  }
  if (window.confirm(`确认删除站点 "${site.value.name}"？将会删除所有相关数据`)) del()
}

async function onFieldEditorYes(val: string) {
  if (!editFieldKey.value) return

  if (editFieldVal.value !== val) {
    isLoading.value = true
    let s: ArtalkType.SiteData
    try {
      let finalVal: string | string[] = val
      if (Array.isArray(site.value[editFieldKey.value]))
        finalVal = val
          .split(',')
          .map((v) => v.trim())
          .filter((v) => !!v)
      s = (
        await artalk!.ctx.getApi().sites.updateSite(site.value.id, {
          ...site.value,
          [editFieldKey.value]: finalVal,
        })
      ).data
    } catch (err: any) {
      alert(`修改失败：${err.message || '未知错误'}`)
      console.error(err)
      return false
    } finally {
      isLoading.value = false
    }
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
        <span class="atk-site-name" @click="!site.first_url || openURL(site.first_url)">
          {{ site.name }}
        </span>
        <span class="atk-site-urls">
          <div v-for="(url, i) in site.urls" :key="i" class="atk-url-item" @click="openURL(url)">
            {{ url }}
          </div>
        </span>
      </div>
      <div class="atk-close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
    <div class="atk-main">
      <div class="atk-site-text-actions">
        <div class="atk-item atk-rename-btn" @click="rename()">
          {{ t('rename') }}
        </div>
        <div class="atk-item atk-edit-url-btn" @click="editURL()">{{ t('edit') }} URL</div>
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

<style scoped lang="scss"></style>

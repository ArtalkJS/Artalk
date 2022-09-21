<script setup lang="ts">
import type { SiteData } from 'artalk/types/artalk-data'

const props = defineProps<{
  site: SiteData
}>()

const emit = defineEmits(['close'])

const { site } = toRefs(props)

function close() {
  emit('close')
}

function openURL(url: string) {
  window.open(url)
}

function rename() {

}

function editURL() {

}

function del() {

}
</script>

<template>
  <div class="atk-site-edit">
    <div class="atk-header">
      <div class="atk-site-info">
        <span class="atk-site-name">{{ site.name }}</span>
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
    </div>
  </div>
</template>

<style scoped lang="scss">

</style>

<script setup lang="ts">
import type { PageData } from 'artalk/types/artalk-data'

const props = defineProps<{
  page: PageData
}>()

const emit = defineEmits(['close'])

const { page } = toRefs(props)

function editTitle() {

}

function editKey() {

}

function editAdminOnly() {

}

function sync() {

}

function del() {

}

function close() {
  emit('close')
}
</script>

<template>
  <div class="atk-page-edit-layer">
    <div class="atk-page-main-actions">
      <div class="atk-item atk-title-edit-btn" @click="editTitle()">标题修改</div>
      <div class="atk-item atk-key-edit-btn" @click="editKey()">KEY 变更</div>
      <div
        class="atk-item atk-admin-only-btn"
        :class="!page.admin_only ? 'atk-green' : 'atk-yellow'"
        @click="editAdminOnly()"
      >{{ !page.admin_only ? '所有人可评' : '管理员可评' }}</div>
    </div>
    <div class="atk-page-actions">
      <div class="atk-item atk-sync-btn" @click="sync()">
        <i class="atk-icon atk-icon-sync"></i>
      </div>
      <div class="atk-item atk-del-btn" @click="del()">
        <i class="atk-icon atk-icon-del"></i>
      </div>
      <div class="atk-item atk-close-btn" @click="close()">
        <i class="atk-icon atk-icon-close"></i>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.atk-page-edit-layer {
  z-index: 999;
  background: var(--at-color-bg);
  position: absolute;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.atk-page-main-actions {
  @extend .atk-list-text-actions;
}
</style>

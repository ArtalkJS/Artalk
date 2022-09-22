<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNavStore } from './stores/nav'

const nav = useNavStore()
const { scrollableArea } = storeToRefs(nav)
</script>

<template>
  <Header />
  <Tab />

  <div ref="scrollableArea" class="main artalk atk-sidebar">
    <div class="atk-sidebar-inner">
      <router-view />
    </div>
  </div>
</template>

<style scoped lang="scss">
.main {
  overflow-y: auto;
  height: calc(100vh - 61px - 41px);
  padding-bottom: 50px;

  // 分页条占位
  :deep(.atk-pagination-wrap) {
    z-index: 200;
    position: fixed;
    width: 100%;
    bottom: 0;
    left: 0;
    background: var(--at-color-bg);
    border-top: 1px solid var(--at-color-border);
  }
}
</style>

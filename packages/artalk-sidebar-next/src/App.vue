<script setup lang="ts">

</script>

<template>
  <Header />
  <Tab />

  <div class="main artalk">
    <router-view />
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

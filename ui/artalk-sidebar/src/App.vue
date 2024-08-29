<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNavStore } from './stores/nav'
import { bootParams } from './global'

const nav = useNavStore()
const { scrollableArea } = storeToRefs(nav)

const darkMode = ref(bootParams.darkMode)

;(function initDarkModeWatchMedia() {
  if (!window.matchMedia) return
  const query = window.matchMedia('(prefers-color-scheme: dark)')
  query.addEventListener('change', (e) => {
    darkMode.value = e.matches
  })
})()
</script>

<template>
  <div class="app-wrap artalk atk-sidebar" :class="{ 'atk-dark-mode': darkMode }">
    <AppHeader />
    <AppNavigation />

    <div class="main">
      <div ref="scrollableArea" class="atk-sidebar-inner">
        <div class="atk-sidebar-container">
          <router-view />
        </div>
      </div>
      <LoadingLayer v-if="nav.isPageLoading" />
    </div>
  </div>
</template>

<style scoped lang="scss">
$headerHeight: 61px;
$subHeaderHeight: 41px;
$maxWidth: 1100px;
$sidebarWidth: 280px;

.app-wrap {
  background: var(--at-color-bg);
  color: var(--at-color-font);
}

.main {
  position: relative;

  .atk-sidebar-inner {
    overflow-y: auto;
    height: calc(100vh - $headerHeight - $subHeaderHeight);
    padding-bottom: 50px;

    @media (min-width: 1024px) {
      height: calc(100vh - $headerHeight - 51px);
    }
  }

  // The placeholder area for the pagination bar
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

:deep(.atk-sidebar-container) {
  position: relative;
}

@media (min-width: 1024px) {
  :deep(.atk-sidebar-container) {
    max-width: $maxWidth;
    margin: 0 auto;
    width: 100%;
  }
}

@media (min-width: 1024px) and (max-width: calc($maxWidth + $sidebarWidth * 2)) {
  :deep(.atk-sidebar-container) {
    margin-left: $sidebarWidth;
    max-width: calc(100% - $sidebarWidth);
  }
}
</style>

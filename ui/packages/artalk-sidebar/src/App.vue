<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNavStore } from './stores/nav'
import { useUserStore } from './stores/user'
import global, { bootParams, createArtalkInstance } from './global'

const nav = useNavStore()
const user = useUserStore()
const route = useRoute()
const router = useRouter()
const { scrollableArea } = storeToRefs(nav)
const i18n = useI18n()

onBeforeMount(() => {
  // 获取语言
  if (!global.getBootParams().locale) {
    global.getArtalk().ctx.getApi().system.conf().then(resp => {
      if (resp.locale && typeof resp.locale == 'string') {
        i18n.locale.value = resp.locale
      }
    })
  }

  if (bootParams.user?.email) {
    global.getArtalk().ctx.get('user').update(bootParams.user)
  } else {
    try {
      global.importUserDataFromArtalkInstance()
    } catch (e) {
      // console.error(e)
      router.replace('/login')
      return
    }
  }

  // 验证登陆身份有效性
  global.getArtalk().ctx.getApi().user.loginStatus().then(resp => {
    if (resp.is_admin && !resp.is_login) {
      global.getArtalk().ctx.get('user').logout()
      user.logout()
      router.replace('/login')
    }
  })
})

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
    <Header />
    <Tab />

    <div class="main">
      <div ref="scrollableArea" class="atk-sidebar-inner">
        <router-view />
      </div>
      <LoadingLayer v-if="nav.isPageLoading" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.app-wrap {
  background: var(--at-color-bg);
  color: var(--at-color-font);
}

.main {
  position: relative;

  .atk-sidebar-inner {
    overflow-y: auto;
    height: calc(100vh - 61px - 41px);
    padding-bottom: 50px;
  }

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

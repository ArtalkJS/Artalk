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

onBeforeMount(() => {
  if (bootParams.user?.email) {
    global.getArtalk().ctx.user.update(bootParams.user)
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
      global.getArtalk().ctx.user.logout()
      user.logout()
      router.replace('/login')
    }
  })
})
</script>

<template>
  <Header />
  <Tab />

  <div class="main artalk atk-sidebar">
    <div ref="scrollableArea" class="atk-sidebar-inner">
      <router-view />
    </div>
    <LoadingLayer v-if="nav.isPageLoading" />
  </div>
</template>

<style scoped lang="scss">
.main {
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

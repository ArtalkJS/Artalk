<script setup lang="ts">
import { storeToRefs } from 'pinia'
import type Artalk from 'artalk'
import { useNavStore } from './stores/nav'
import { useUserStore } from './stores/user'
import { getArtalk, bootParams } from './global'

const nav = useNavStore()
const user = useUserStore()
const route = useRoute()
const router = useRouter()
const { scrollableArea } = storeToRefs(nav)

const artalkLoaded = ref(false)

onBeforeMount(() => {
  const artalk = getArtalk()
  if (!artalk) {
    throw new Error('Artalk instance not initialized')
  }

  artalk.on('mounted', () => {
    if (artalkLoaded.value) return
    artalkLoaded.value = true

    syncArtalk(artalk)
  })
})

function syncArtalk(artalk: Artalk) {
  // access from open sidebar or directly by url
  if (bootParams.user?.email) {
    // sync user from sidebar to artalk
    artalk.ctx.get('user').update(bootParams.user)
  } else {
    // sync user from artalk to sidebar
    try {
      useUserStore().sync()
    } catch {
      nextTick(() => {
        router.replace('/login')
      })
      return
    }
  }

  // 验证登录身份有效性
  const artalkUser = artalk.ctx.get('user')
  const artalkUserData = artalkUser.getData()

  const logout = () => {
    user.logout()
    nextTick(() => {
      router.replace('/login')
    })
  }

  // Remove login failed dialog if sidebar
  artalk.ctx.getApiHandlers().remove('need_login')
  artalk.ctx.getApiHandlers().add('need_login', async () => {
    logout()
    throw new Error('Need login')
  })

  // Check user status
  artalk.ctx
    .getApi()
    .user.getUserStatus({
      email: artalkUserData.email,
      name: artalkUserData.nick,
    })
    .then((res) => {
      if (res.data.is_admin && !res.data.is_login) {
        logout()
      } else {
        // 将全部通知标记为已读
        artalk.ctx.getApi().notifies.markAllNotifyRead({
          email: artalkUserData.email,
          name: artalkUserData.nick,
        })
      }
    })
}

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
  <div
    v-if="artalkLoaded"
    class="app-wrap artalk atk-sidebar"
    :class="{ 'atk-dark-mode': darkMode }"
  >
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

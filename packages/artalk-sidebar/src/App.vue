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
const artalkLoaded = ref(false)

const LinkMap: {[key:string]:string} = {
  comments: '/comments',
  pages: '/pages',
  sites: '/sites',
  settings: '/settings'
}

onBeforeMount(() => {
  createArtalkInstance().then(artalkInstance => {
    // 初始化 Artalk
    global.setArtalk(artalkInstance)

    artalkLoaded.value = true

    // 更新用户资料
    if (bootParams.user?.email) {
      artalkInstance.ctx.user.update(bootParams.user)
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
    artalkInstance.ctx.getApi().user.loginStatus().then(resp => {
      if (resp.is_admin && !resp.is_login) {
        router.replace('/login')
        return
      }
    })

    // 首页跳转
    if (route.path === '/') {
      if (bootParams.view) {
        const splitted = bootParams.view.split('|')
        if (splitted[0]) bootParams.view = splitted[0]
        if (splitted[1]) bootParams.viewParams = JSON.parse(splitted[1])
      }

      router.replace(LinkMap[bootParams.view] || '/comments')
    }
  })
})
</script>

<template>
  <div v-if="artalkLoaded">
    <Header />
    <Tab />

    <div class="main artalk atk-sidebar">
      <div ref="scrollableArea" class="atk-sidebar-inner">
        <router-view />
      </div>
      <LoadingLayer v-if="nav.isPageLoading" />
    </div>
  </div>
  <LoadingLayer v-if="!artalkLoaded" />
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

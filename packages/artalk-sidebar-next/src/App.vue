<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNavStore } from './stores/nav'
import global, { bootParams, createArtalkInstance } from './global'

const nav = useNavStore()
const router = useRouter()
const { scrollableArea } = storeToRefs(nav)
const artalkLoaded = ref(false)

onMounted(() => {
  createArtalkInstance().then(artalkInstance => {
    // 初始化 Artalk
    global.setArtalk(artalkInstance)

    // 更新用户资料
    global.getArtalk()!.ctx.user.update(bootParams.user)

    // 验证登陆身份有效性
    // artalkInstance.ctx.getApi().user.loginStatus()
    //   .then(resp => {
    //     if (resp.is_admin && !resp.is_login) {
    //       router.replace('/login')
    //     }
    //   })

    artalkLoaded.value = true
  })
})
</script>

<template>
  <div v-if="artalkLoaded">
    <Header />
    <Tab />

    <div ref="scrollableArea" class="main artalk atk-sidebar">
      <div class="atk-sidebar-inner">
        <router-view />
      </div>
    </div>
  </div>
  <LoadingLayer v-if="!artalkLoaded" />
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

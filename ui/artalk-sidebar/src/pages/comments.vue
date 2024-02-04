<script setup lang="ts">
import { artalk } from '../global'
import Artalk from 'artalk'
import { useNavStore } from '../stores/nav'
import { useUserStore } from '../stores/user'
import { storeToRefs } from 'pinia'

const wrapEl = ref<HTMLElement>()
const listEl = ref<HTMLElement>()
const user = useUserStore()
const nav = useNavStore()
const { curtTab } = storeToRefs(nav)
const { site: curtSite } = storeToRefs(user)
const { t } = useI18n()

const search = ref('')

onMounted(() => {
  // 初始化导航条
  if (user.isAdmin) {
    nav.updateTabs({
      all: 'all',
      pending: 'pending',
      personal_all: 'personal',
    }, 'all')
  } else {
    nav.updateTabs({
      all: 'all',
      mentions: 'mentions',
      mine: 'mine',
      pending: 'pending',
    }, 'all')
  }

  watch(curtTab, (curtTab) => {
    artalk!.ctx.fetch({
      offset: 0
    })
  })

  watch(curtSite, (value) => {
    artalk!.ctx.reload()
  })

  artalk!.ctx.on('comment-rendered', (comment) => {
    const pageURL = comment.getData().page_url
    comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
  })

  artalk!.ctx.updateConf({
    listFetchParamsModifier: (params) => {
      params.site_name = curtSite.value // 站点名

      let scope = user.isAdmin ? 'site' : 'user'
      let type = curtTab.value

      if (curtTab.value === 'personal_all') {
        scope = 'user'
        type = 'all'
      }

      params.scope = scope
      params.type = type

      if (search.value) params.search = search.value
    },
    scrollRelativeTo: () => wrapEl.value!,
  })

  artalk!.reload()

  const $el = artalk!.ctx.get('list')!.$el

  // @ts-ignore
  $el.querySelector<HTMLElement>('.atk-list-header').style.display = 'none'
  // @ts-ignore
  $el.querySelector<HTMLElement>('.atk-list-footer').style.display = 'none'

  listEl.value?.append($el)

  // 搜索功能
  nav.enableSearch((value: string) => {
    search.value = value
    artalk!.reload()
  }, () => {
    if (search.value === '') return
    search.value = ''
    artalk!.reload()
  })
})
</script>

<template>
  <div ref="wrapEl" class="comments-wrap">
    <div ref="listEl" />
  </div>
</template>

<style scoped lang="scss">
.comments-wrap {
  overflow-y: auto;
  height: 100%;

  :deep(.atk-comment-wrap) {
    border-bottom: 1px solid var(--at-color-border);
  }
}
</style>

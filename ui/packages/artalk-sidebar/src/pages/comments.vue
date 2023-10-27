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
      admin_all: 'all',
      admin_pending: 'pending',
      all: 'personal',
    }, 'admin_all')
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

  artalk!.ctx.updateConf({
    noComment: `<div class="atk-sidebar-no-content">${t('noContent')}</div>`,
    pagination: {
      pageSize: 20,
      readMore: false,
      autoLoad: false,
    },
    listUnreadHighlight: true,
    listFetchParamsModifier: (params) => {
      params.type = curtTab.value // 列表数据类型
      params.site_name = curtSite.value // 站点名
      if (search.value) params.search = search.value
    },
    listScrollListenerAt: wrapEl.value,
  })

  artalk!.ctx.on('comment-rendered', (comment) => {
    const pageURL = comment.getData().page_url
    comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
    comment.getConf().onReplyBtnClick = () => {
      artalk!.ctx.replyComment(comment.getData(), comment.getEl())
    }
  })

  artalk!.reload()

  listEl.value?.append(artalk!.ctx.get('list')!.$el)

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

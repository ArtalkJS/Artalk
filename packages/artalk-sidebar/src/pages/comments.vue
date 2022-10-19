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

const search = ref('')

onMounted(() => {
  // 初始化导航条
  if (user.isAdmin) {
    nav.updateTabs({
      admin_all: '全部',
      admin_pending: '待审',
      all: '个人',
    }, 'admin_all')
  } else {
    nav.updateTabs({
      mentions: '提及',
      all: '全部',
      mine: '我的',
      pending: '待审',
    }, 'mentions')
  }

  watch(curtTab, (curtTab) => {
    list.fetchComments(0)
  })

  watch(curtSite, (value) => {
    list.reload()
  })

  // 初始化评论列表
  const list = new Artalk.ListLite(artalk!.ctx)
  artalk!.ctx.setList(list)

  list.flatMode = true
  list.unreadHighlight = true
  list.scrollListenerAt = wrapEl.value
  list.pageMode = 'pagination'
  list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
  list.renderComment = (comment) => {
    const pageURL = comment.getData().page_url
    comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
    comment.getConf().onReplyBtnClick = () => {
      artalk!.ctx.replyComment(comment.getData(), comment.getEl(), true)
    }
  }
  list.paramsEditor = (params) => {
    params.type = curtTab.value // 列表数据类型
    params.site_name = curtSite.value // 站点名
    if (search.value) params.search = search.value
  }
  artalk!.on('list-inserted', (data) => {
    wrapEl.value!.scrollTo(0, 0)
  })

  list.reload()

  listEl.value?.append(list.$el)

  // 搜索功能
  nav.enableSearch((value: string) => {
    search.value = value
    list.reload()
  }, () => {
    if (search.value === '') return
    search.value = ''
    list.reload()
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

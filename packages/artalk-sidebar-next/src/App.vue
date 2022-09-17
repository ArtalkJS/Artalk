<script setup lang="ts">
import global from './global'
import ListLite from 'artalk/src/list/list-lite'

const commentListEl = ref<HTMLElement>()

onMounted(() => {
  const list = new ListLite(global.artalk!.ctx)
  global.artalk!.ctx.setList(list)
  list.flatMode = true
  list.unreadHighlight = true
  // list.scrollListenerAt = this.$parent
  list.pageMode = 'pagination'
  list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
  list.renderComment = (comment) => {
    const pageURL = comment.getData().page_url
    comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
    comment.getConf().onReplyBtnClick = () => {
      global.artalk!.ctx.replyComment(comment.getData(), comment.getEl(), true)
    }
  }

  list.paramsEditor = (params) => {
    params.site_name = 'ArtalkDocs'
  }

  global.artalk!.ctx.on('list-inserted', (data) => {
    // ;(this.$el.parentNode as any)?.scrollTo(0, 0)
  })

  list.reload()

  commentListEl.value?.append(list.$el)
})
</script>

<template>
  <Header />
  <Tab />

  <div class="artalk">
    <div ref="commentListEl"></div>
  </div>

  <h2>{{ $route.name }}</h2>
  <RouterLink to="/">Home</RouterLink>
  <RouterLink to="/about">About</RouterLink>
  <main>
    <router-view />
  </main>
</template>

<style scoped lang="scss">

</style>

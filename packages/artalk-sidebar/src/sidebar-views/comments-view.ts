import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import ListLite from 'artalk/src/list/list-lite'

import SidebarView from '../sidebar-view'

export default class MessageView extends SidebarView {
  static viewName = 'comments'
  static viewTitle = '评论'

  viewTabs: any = {}
  viewActiveTab = ''

  list!: ListLite

  mount() {
    // tabs
    if (this.ctx.user.data.isAdmin) {
      this.viewTabs = {
        admin_all: '全部',
        admin_pending: '待审',
        all: '个人',
      }
      this.viewActiveTab = 'admin_all'
    } else {
      this.viewTabs = {
        mentions: '提及',
        all: '全部',
        mine: '我的',
        pending: '待审',
      }
      this.viewActiveTab = 'mentions'
    }

    this.list = new ListLite(this.ctx)
    this.list.flatMode = true
    this.list.unreadHighlight = true
    this.list.scrollListenerAt = this.$parent
    this.list.pageMode = 'pagination'
    this.list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
    this.list.renderComment = (comment) => {
      const pageURL = comment.getData().page_url
      comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
      comment.getConf().onReplyBtnClick = () => {
        this.ctx.trigger('editor-reply', {data: comment.getData(), $el: comment.getEl(), scroll: true})
      }
    }
    this.list.paramsEditor = (params) => {
      params.site_name = this.sidebar.curtSite
    }
    this.ctx.on('list-insert', (data) => {
      this.list.insertComment(data)
      ;(this.$el.parentNode as any)?.scrollTo(0, 0)
    })

    this.$el.innerHTML = ''
    this.$el.append(this.list.$el)

    this.switchTab(this.viewActiveTab)
  }

  switchTab(tab: string): boolean|void {
    //console.log(tab, siteName)
    this.viewActiveTab = tab
    this.list.paramsEditor = (params) => {
      params.type = tab as any
      params.site_name = this.sidebar.curtSite
    }
    this.list.fetchComments(0)

    return true
  }
}


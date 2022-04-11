import Context from '../../context'
import SidebarView from '../sidebar-view'
import * as Utils from '../../lib/utils'
import ListLite from '../list-lite'

export default class MessageView extends SidebarView {
  static viewName = 'comments'
  static viewTitle = '评论'

  viewTabs: any = {}
  viewActiveTab = ''

  list!: ListLite

  mount(siteName: string) {
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
      comment.setOpenURL(`${comment.data.page_key}#atk-comment-${comment.data.id}`)
      comment.onReplyBtnClick = () => {
        if (this.ctx.conf.editorTravel === true) {
          this.ctx.trigger('editor-reply', {data: comment.data, $el: comment.$el, scroll: false})
        } else {
          this.ctx.trigger('sidebar-hide')
          this.ctx.trigger('editor-reply', {data: comment.data, $el: comment.$el, scroll: true})
        }
      }
    }
    this.list.paramsEditor = (params) => {
      params.site_name = siteName
    }

    this.$el.innerHTML = ''
    this.$el.append(this.list.$el)

    this.switchTab(this.viewActiveTab, siteName)
  }

  switchTab(tab: string, siteName: string): boolean|void {
    //console.log(tab, siteName)
    this.viewActiveTab = tab
    this.list.paramsEditor = (params) => {
      params.type = tab as any
      params.site_name = siteName
    }
    this.list.fetchComments(0)

    return true
  }
}


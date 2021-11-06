import Context from '../../context'
import SidebarView from '../sidebar-view'
import * as Utils from '../../lib/utils'
import ListLite from '../list-lite'

export default class MessageView extends SidebarView {
  static viewName = 'comments'
  static viewTitle = '评论'

  viewTabs = {
    mentions: '提及',
    all: '全部',
    mine: '我的',
    pending: '待审',
  }
  viewActiveTab = 'mentions'

  list?: ListLite

  mount() {
    this.list = new ListLite(this.ctx)
    this.list.flatMode = true
    this.list.unreadHighlight = true
    this.list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
    this.list.renderComment = (comment) => {
      comment.setOpenURL(`${comment.data.page_key}#atk-comment-${comment.data.id}`)
    }

    this.$el.innerHTML = ''
    this.$el.append(this.list.$el)

    this.switch(this.viewActiveTab)
  }

  switch(tab: string): boolean|void {
    if (!this.list) return false
    this.viewActiveTab = tab
    this.list.type = tab as any;
    this.list.isFirstLoad = true
    this.list.reqComments()

    return true
  }
}


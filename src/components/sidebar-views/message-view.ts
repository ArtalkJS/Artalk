import Context from '../../context'
import SidebarView from './sidebar-view'
import * as Utils from '../../lib/utils'
import ListLite from '../list-lite'

export default class MessageView extends SidebarView {
  name = 'message'
  title = '通知中心'
  actions = {
    mentions: '提及',
    all: '全部',
    mine: '我的',
    pending: '待审',
  }
  activeAction = ''

  list?: ListLite
  type: string = 'mentions'

  init () {
    this.list = CreateCommentList(this.ctx)
    this.el.innerHTML = ''
    this.el.append(this.list.el)

    this.activeAction = this.type
    this.switch(this.type)
  }

  switch (action: string) {
    if (!this.list) return
    this.type = action
    this.list.type = action as any;
    this.list.isFirstLoad = true
    this.list.reqComments()
  }
}

export function CreateCommentList (ctx: Context) {
  const list = new ListLite(ctx)
  list.flatMode = true
  list.unreadHighlight = true
  list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
  list.renderComment = (comment) => {
    comment.setOpenURL(`${comment.data.page_key}#atk-comment-${comment.data.id}`)
  }

  return list
}

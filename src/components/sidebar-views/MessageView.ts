import Context from '../../Context'
import SidebarView from './SidebarView'
import * as Utils from '../../lib/utils'
import Comment from '../Comment'
import ListLite from '../ListLite'

export default class MessageView extends SidebarView {
  list?: ListLite
  type: string = 'mentions'

  name = 'message'
  title = '通知中心'
  actions = {
    mentions: '提及',
    all: '全部',
    mine: '我的',
    pending: '待审',
  }
  activeAction = this.type

  render () {
    this.list = new ListLite(this.ctx)
    this.list.flatMode = true
    this.list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
    this.list.renderComment = this.renderComment
    this.el.append(this.list.el)

    this.list.type = this.type as any
    this.list.isFirstLoad = true
    this.list.reqComments()

    return this.el
  }

  switch (action: string) {
    if (!this.list) return
    this.type = action
    this.list.type = action as any;
    this.list.isFirstLoad = true
    this.list.reqComments()
  }

  renderComment (comment: Comment) {
    // comment.el.querySelector('[data-atk-action="comment-reply"]')!.remove()

    comment.el.style.cursor = 'pointer'
    comment.el.addEventListener('mouseover', () => {
      comment.el.style.backgroundColor = 'var(--at-color-bg-grey)'
    })

    comment.el.addEventListener('mouseout', () => {
      comment.el.style.backgroundColor = ''
    })

    comment.el.addEventListener('click', (evt) => {
      evt.preventDefault()
      window.location.href = `${comment.data.page_key}#artalk-comment-${comment.data.id}`
    })

    // this.contentEl.appendChild(comment.getEl())
  }
}

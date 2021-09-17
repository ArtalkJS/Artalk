import '../style/sidebar.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment  from './Comment'
import SidebarHTML from './html/sidebar.html?raw'
import { ListData, CommentData } from '~/types/artalk-data'
import BuildLayer, { Layer } from './Layer'
import Api from '../lib/api'
import ListLite from './ListLite'

export default class Sidebar extends Component {
  public el: HTMLElement
  public layer?: Layer
  public contentEl: HTMLElement
  public list: ListLite
  public type: string = 'mentions'

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(SidebarHTML)
    this.contentEl = this.el.querySelector('.atk-sidebar-content')!

    this.el.querySelector('.atk-sidebar-close')!.addEventListener('click', () => {
      this.hide()
    })

    this.initActionBar()

    this.ctx.addEventListener('sidebar-show', () => (this.show()))
    this.ctx.addEventListener('sidebar-hide', () => (this.hide()))

    this.list = new ListLite(this.ctx)
    this.list.flatMode = true
    this.list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
    this.list.renderComment = this.renderComment
    this.contentEl.append(this.list.el)
  }

  initActionBar() {
    const actionsEl = this.el.querySelector('.atk-sidebar-actions')!
    actionsEl.addEventListener('click', (evt) => {
      const el = evt.target as HTMLElement
      const type = el.getAttribute('data-atk-action')
      if (!type) return

      actionsEl.querySelectorAll('.atk-active').forEach((item) => {
        item.classList.remove('atk-active')
      })
      el.classList.add('atk-active')
      this.type = type
      this.list.type = (this.type as any);
      this.list.isFirstLoad = true
      this.list.reqComments()
    })
  }

  show () {
    this.el.style.transform = '' // 动画清除，防止二次打开失效

    this.layer = BuildLayer(this.ctx, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)

    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)' // 执行动画
    }, 20)

    this.list.type = this.type as any
    this.list.isFirstLoad = true
    this.list.reqComments()
  }

  hide () {
    this.el.style.transform = ''
    this.layer?.dispose() // 用完即销毁
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

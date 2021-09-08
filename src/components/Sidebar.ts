import '@/style/sidebar.less'

import Context from '@/Context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import Comment  from './Comment'
import SidebarHTML from './html/sidebar.html?raw'
import { ListData, CommentData } from '~/types/artalk-data'
import BuildLayer, { Layer } from './Layer'

export default class Sidebar extends Component {
  public el: HTMLElement
  public layer?: Layer
  public contentEl: HTMLElement

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(SidebarHTML)
    this.contentEl = this.el.querySelector('.artalk-sidebar-content')!

    this.el.querySelector('.artalk-sidebar-close')!.addEventListener('click', () => {
      this.hide()
    })

    this.initActionBar()

    this.ctx.addEventListener('sidebar-show', () => (this.show()))
    this.ctx.addEventListener('sidebar-hide', () => (this.hide()))
  }

  initActionBar() {
    const actionsEl = this.el.querySelector('.artalk-sidebar-actions')!
    actionsEl.addEventListener('click', (evt) => {
      const el = evt.target as HTMLElement
      const type = el.getAttribute('data-artalk-action')
      if (!type) return

      actionsEl.querySelectorAll('.artalk-active').forEach((item) => {
        item.classList.remove('artalk-active')
      })
      el.classList.add('artalk-active')
      this.reqComment(type)
    })
  }

  show () {
    this.layer = BuildLayer(this.ctx, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)
    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)'
    }, 20)

    this.reqComment('mentions')
  }

  hide () {
    this.el.style.transform = ''
    this.layer?.dispose() // 用完即销毁
  }

  reqComment (type: string) {
    this.contentEl.innerHTML = ''

    let reqData: any = {
      nick: this.ctx.user.data.nick,
      email: this.ctx.user.data.email,
      type,
      limit: 999,
    }

    if (this.ctx.user.data.isAdmin) {
      reqData = { token: this.ctx.user.data.token, ...reqData }
    }

    // TODO: sidebar Req
    // this.artalk.request('CommentGetV2', reqData, () => {
    //   Ui.showLoading(this.contentEl)
    // }, () => {
    //   Ui.hideLoading(this.contentEl)
    // }, (msg, data) => {
    //   this.contentEl.innerHTML = ''
    //   if (Array.isArray(data.comments)) {
    //     (data.comments as CommentData[]).forEach((item) => {
    //       this.putComment(item)
    //     });
    //   }
    //   if (!data.comments || !Array.isArray(data.comments) || data.comments.length <= 0) {
    //     this.showNoComment()
    //   }
    // }, (msg, data) => {

    // })
  }

  putComment (data: CommentData) {
    const comment = new Comment(this.ctx, data)

    comment.el.querySelector('[data-comment-action="reply"]')!.remove()
    comment.el.style.cursor = 'pointer'
    comment.el.addEventListener('mouseover', () => {
      comment.el.style.backgroundColor = 'var(--at-color-bg-grey)'
    })

    comment.el.addEventListener('mouseout', () => {
      comment.el.style.backgroundColor = ''
    })

    comment.el.addEventListener('click', (evt) => {
      evt.preventDefault()
      window.location.href = `${(data as any).page_key}#artalk-comment-${comment.data.id}`
    })

    this.contentEl.appendChild(comment.getEl())
  }

  showNoComment() {
    this.contentEl.innerHTML = '<div class="artalk-sidebar-no-content">无内容</div>'
  }
}

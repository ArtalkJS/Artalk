import Comment from './Comment'
import '../css/sidebar.less'
import Artalk from '../Artalk'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'
import { CommentData } from '~/types/artalk-data'
import BuildLayer, { Layer } from './Layer'

export default class Sidebar extends ArtalkContext {
  public el: HTMLElement
  public layer: Layer
  public contentEl: HTMLElement

  constructor (artalk: Artalk) {
    super(artalk)

    this.el = Utils.createElement(require('../templates/Sidebar.ejs')(this))
    this.contentEl = this.el.querySelector('.artalk-sidebar-content')

    this.el.querySelector('.artalk-sidebar-close').addEventListener('click', () => {
      this.hide()
    })

    this.initActionBar()
  }

  initActionBar() {
    const actionsEl = this.el.querySelector('.artalk-sidebar-actions')
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
    this.layer = BuildLayer(this.artalk, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)
    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)'
    }, 20)

    this.reqComment('mentions')
  }

  hide () {
    this.el.style.transform = ''
    this.layer.dispose() // 用完即销毁
  }

  reqComment (type: string) {
    this.contentEl.innerHTML = ''

    let reqObj: any = {
      nick: this.artalk.user.data.nick,
      email: this.artalk.user.data.email,
      type,
      limit: 999,
    }

    if (this.artalk.user.data.isAdmin) {
      reqObj = { password: this.artalk.user.data.password, ...reqObj }
    }

    this.artalk.request('CommentGetV2', reqObj, () => {
      this.artalk.ui.showLoading(this.contentEl)
    }, () => {
      this.artalk.ui.hideLoading(this.contentEl)
    }, (msg, data) => {
      this.contentEl.innerHTML = ''
      if (Array.isArray(data.comments)) {
        (data.comments as CommentData[]).forEach((item) => {
          this.putComment(item)
        });
      }
      if (!data.comments || !Array.isArray(data.comments) || data.comments.length <= 0) {
        this.showNoComment()
      }
    }, (msg, data) => {

    })
  }

  putComment (data: CommentData) {
    const comment = new Comment(this.artalk, null, data)

    comment.elem.querySelector('[data-comment-action="reply"]').remove()
    comment.elem.style.cursor = 'pointer'
    comment.elem.addEventListener('mouseover', () => {
      comment.elem.style.backgroundColor = 'var(--at-color-bg-grey)'
    })

    comment.elem.addEventListener('mouseout', () => {
      comment.elem.style.backgroundColor = ''
    })

    comment.elem.addEventListener('click', (evt) => {
      evt.preventDefault()
      window.location.href = `${(data as any).page_key}#artalk-comment-${comment.data.id}`
    })

    this.contentEl.appendChild(comment.getElem())
  }

  showNoComment() {
    this.contentEl.innerHTML = '<div class="artalk-sidebar-no-content">无内容</div>'
  }
}

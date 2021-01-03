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
  }

  show () {
    this.layer = BuildLayer(this.artalk, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)
    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)'
    }, 20)

    this.artalk.request('CommentReplyGet', {
      nick: this.artalk.user.data.nick,
      email: this.artalk.user.data.email
    }, () => {
      this.artalk.ui.showLoading(this.contentEl)
    }, () => {
      this.artalk.ui.hideLoading(this.contentEl)
    }, (msg, data) => {
      this.contentEl.innerHTML = ''
      if (Array.isArray(data.reply_comments)) {
        (data.reply_comments as CommentData[]).forEach((item) => {
          this.putComment(item)
        });
      }
    }, (msg, data) => {

    })
  }

  hide () {
    this.el.style.transform = ''
    this.layer.dispose() // 用完即销毁
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
}

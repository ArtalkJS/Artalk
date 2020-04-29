import Comment from './Comment'
import '../css/sidebar.less'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'
import { CommentData } from '~/types/artalk-data'
import Layer from './Layer'

export default class Sidebar extends ArtalkContext {
  public el: HTMLElement
  public layer: Layer
  public contentEl: HTMLElement

  constructor () {
    super()

    this.el = Utils.createElement(require('../templates/Sidebar.ejs')(this))
    this.layer = new Layer('sidebar', this.el)
    this.contentEl = this.el.querySelector('.artalk-sidebar-content')

    this.el.querySelector('.artalk-sidebar-close').addEventListener('click', () => {
      this.hide()
    })
  }

  show () {
    this.layer.show()
    this.contentEl.scrollTo(0, 0)
    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)'
    }, 20)

    this.artalk.request('CommentReplyGet', {
      nick: this.artalk.user.nick,
      email: this.artalk.user.email
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
    this.layer.hide()
  }

  putComment (data: CommentData) {
    const comment = new Comment(null, data)

    comment.elem.querySelector('[data-comment-action="reply"]').remove()
    comment.elem.style.cursor = 'pointer'
    comment.elem.addEventListener('mouseover', () => {
      comment.elem.style.backgroundColor = '#F4F4F4'
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

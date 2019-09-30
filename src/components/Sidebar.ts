import Comment from './Comment'
import '../css/sidebar.scss'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'
import { CommentData } from '~/types/artalk-data'

export default class Sidebar extends ArtalkContext {
  public el: HTMLElement
  public wrapEl: HTMLElement
  public contentEl: HTMLElement
  public maskEl: HTMLElement

  constructor () {
    super()

    this.wrapEl = Utils.createElement(require('../templates/Sidebar.ejs')(this))
    document.querySelector('body').appendChild(this.wrapEl)

    this.el = this.wrapEl.querySelector('.artalk-sidebar')
    this.contentEl = this.el.querySelector('.artalk-sidebar-content')
    this.maskEl = this.wrapEl.querySelector('.artalk-layer-mask')

    this.maskEl.addEventListener('click', () => {
      this.hide()
    })
    this.el.querySelector('.artalk-sidebar-close').addEventListener('click', () => {
      this.hide()
    })
  }

  show () {
    this.wrapEl.style.display = 'block'
    this.maskEl.style.display = 'block'
    this.maskEl.classList.add('artalk-fade-in')
    this.el.style.transform = 'translate(0, 0)'

    this.artalk.request('CommentReplyGet', {
      nick: this.artalk.editor.user.nick,
      email: this.artalk.editor.user.email
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
    setTimeout(() => {
      this.wrapEl.style.display = 'none'
    }, 450)

    this.maskEl.classList.add('artalk-fade-out')
    setTimeout(() => {
      this.maskEl.style.display = 'none'
      this.maskEl.classList.remove('artalk-fade-out')
    }, 200)
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
      window.location.href = (data as any).page_key
    })

    this.contentEl.appendChild(comment.getElem())
  }
}

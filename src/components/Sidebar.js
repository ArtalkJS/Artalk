import $ from 'jquery'
import Comment from './Comment.js'
import '../css/sidebar.scss'

export default class Sidebar {
  constructor (artalk) {
    this.artalk = artalk

    this.wrapEl = $(require('../templates/Sidebar.ejs')(this)).appendTo('body')
    this.el = this.wrapEl.find('.artalk-sidebar')
    this.contentEl = this.el.find('.artalk-sidebar-content')
    this.maskEl = this.wrapEl.find('.artalk-layer-mask')

    this.maskEl.click(() => {
      this.hide()
    })
    this.el.find('.artalk-sidebar-close').click(() => {
      this.hide()
    })
  }

  show () {
    this.wrapEl.show()
    this.maskEl.show()
    this.maskEl.addClass('artalk-fade-in')
    this.el.css('transform', 'translate(0, 0)')

    this.artalk.request('CommentReplyGet', {
      nick: this.artalk.editor.user.nick,
      email: this.artalk.editor.user.email
    }, () => {
      this.artalk.ui.showLoading(this.contentEl)
    }, () => {
      this.artalk.ui.hideLoading(this.contentEl)
    }, (msg, data) => {
      this.contentEl.html('')
      for (let i in data.reply_comments) {
        let item = data.reply_comments[i]
        this.putComment(item)
      }
    }, (msg, data) => {

    })
  }

  hide () {
    this.el.css('transform', '')
    setTimeout(() => {
      this.wrapEl.hide()
    }, 450)

    this.maskEl.addClass('artalk-fade-out')
    setTimeout(() => {
      this.maskEl.hide()
      this.maskEl.removeClass('artalk-fade-out')
    }, 200)
  }

  putComment (data) {
    let comment = new Comment(this, data)

    comment.elem.find('[data-comment-action="reply"]').remove()
    comment.elem.css('cursor', 'pointer')
    comment.elem.hover(() => {
      comment.elem.css('background-color', '#F4F4F4')
    }, () => {
      comment.elem.css('background-color', '')
    })
    comment.elem.click((evt) => {
      evt.preventDefault()
      window.location.href = data.page_key
    })

    $(comment.getElem()).appendTo(this.contentEl)
  }
}

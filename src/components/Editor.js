import $ from 'jquery'
import '../css/editor.scss'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug.js'
import PreviewPlug from './editor-plugs/PreviewPlug.js'

export default class Editor {
  constructor (artalk) {
    this.artalk = artalk

    this.plugs = [new EmoticonsPlug(this), new PreviewPlug(this)]
    this.el = $(require('./Editor.ejs')(this)).appendTo(this.artalk.el)

    this.initUser()
    this.initTextarea()
    this.initEditorPlug()
    this.initSubmit()
  }

  initUser () {
    const storeUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: storeUser.nick || null,
      email: storeUser.email || null,
      link: storeUser.link || null
    }

    this.headerEl = this.el.find('.artalk-editor-header')
    let userFields = Object.keys(this.user)
    for (let i in userFields) {
      let field = userFields[i]
      let inputEl = this.headerEl.find(`[name="${field}"]`)
      inputEl.val(this.user[field] || '')
      inputEl.bind('input propertychange', (evt) => {
        let inputEl = $(evt.currentTarget)
        this.user[field] = $.trim(inputEl.val())
        window.localStorage.setItem('ArtalkUser', JSON.stringify(this.user))
      })
    }
  }

  initTextarea () {
    this.textareaEl = this.el.find('.artalk-editor-textarea')

    /* Textarea Ui Fix */
    // Tab
    this.textareaEl.on('keydown', (e) => {
      var keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.addContent('\t')
      }
    })

    // Auto Adjust Height
    let adjustHeight = () => {
      var diff = this.textareaEl.outerHeight() - this.textareaEl[0].clientHeight
      this.textareaEl[0].style.height = 0
      this.textareaEl[0].style.height = this.textareaEl[0].scrollHeight + diff + 'px'
    }

    this.textareaEl.bind('input propertychange', (evt) => {
      adjustHeight()
    })
  }

  initEditorPlug () {
    this.plugWrapEl = this.el.find('.artalk-editor-plug-wrap')

    let openedPlugName = null
    this.el.find('.artalk-editor-plug-switcher').click((evt) => {
      let btnElem = $(evt.currentTarget)
      let plugIndex = btnElem.attr('data-plug-index')
      let plug = this.plugs[plugIndex]
      if (!plug) {
        throw Error(`Plug index=${plugIndex} was not found`)
      }

      this.el.find('.artalk-editor-plug-switcher.active').removeClass('active')

      if (openedPlugName === plug.getName()) {
        this.plugWrapEl.css('display', 'none')
        openedPlugName = null
        return
      }

      if (!this.plugWrapEl.find(`[data-plug-name="${plug.getName()}"]`).length) {
        // 需要初始化
        plug.getElem()
          .attr('data-plug-name', plug.getName())
          .css('display', 'none')
          .appendTo(this.plugWrapEl)
      }

      this.plugWrapEl.children().each((i, plugElem) => {
        plugElem = $(plugElem)
        if (plugElem.attr('data-plug-name') !== plug.getName()) {
          plugElem.css('display', 'none')
        } else {
          plugElem.css('display', '')
          plug.onShow()
        }
      })

      this.plugWrapEl.css('display', '')
      openedPlugName = plug.getName()

      btnElem.addClass('active')
    })
  }

  addContent (val) {
    if (document.selection) {
      this.textareaEl.focus()
      document.selection.createRange().text = val
      this.textareaEl.focus()
    } else if (this.textareaEl[0].selectionStart || this.textareaEl[0].selectionStart === 0) {
      let sStart = this.textareaEl[0].selectionStart
      let sEnd = this.textareaEl[0].selectionEnd
      let sT = this.textareaEl[0].scrollTop
      this.textareaEl.val(this.textareaEl.val().substring(0, sStart) + val + this.textareaEl.val().substring(sEnd, this.textareaEl.val().length))
      this.textareaEl.focus()
      this.textareaEl[0].selectionStart = sStart + val.length
      this.textareaEl[0].selectionEnd = sStart + val.length
      this.textareaEl[0].scrollTop = sT
    } else {
      this.textareaEl.focus()
      this.textareaEl[0].value += val
    }
  }

  setContent (val) {
    this.textareaEl.val(val)
  }

  clearContent () {
    this.textareaEl.val('')
  }

  getContent () {
    return this.textareaEl.val()
  }

  getContentMarked () {
    return this.artalk.marked(this.textareaEl.val())
  }

  initSubmit () {
    this.submitBtn = this.el.find('.artalk-send-btn')
    this.submitBtn.click((evt) => {
      let btnEl = evt.currentTarget
      this.onSubmit(btnEl)
    })
  }

  onSubmit (btnEl) {
    $.ajax({
      type: 'POST',
      url: this.artalk.opts.serverUrl,
      data: {
        action: 'CommentAdd',
        content: this.getContent(),
        nick: this.user.nick,
        email: this.user.email,
        link: this.user.link,
        rid: 0,
        page_key: this.artalk.opts.pageKey
      },
      success: (obj) => {
        console.log(obj)
      },
      dataType: 'json'
    })
  }
}

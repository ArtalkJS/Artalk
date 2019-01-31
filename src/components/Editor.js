import $ from 'jquery'
import '../css/editor.scss'
import Comment from './Comment.js'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug.js'
import PreviewPlug from './editor-plugs/PreviewPlug.js'

export default class Editor {
  constructor (artalk) {
    this.artalk = artalk

    this.plugs = [new EmoticonsPlug(this), new PreviewPlug(this)]
    this.el = $(require('../templates/Editor.ejs')(this)).appendTo(this.artalk.el)

    this.initUser()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()
  }

  initUser () {
    const storeUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: storeUser.nick || null,
      email: storeUser.email || null,
      link: storeUser.link || null,
      password: storeUser.password || null
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
        if (field === 'nick' || field === 'email') {
          this.user.password = null
        }
        this.userLocalStorageSave()
      })
    }
  }

  userLocalStorageSave () {
    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.user))
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

  clearEditor () {
    this.textareaEl.val('')
    this.cancelReply()
  }

  getContent () {
    return this.textareaEl.val()
  }

  getContentMarked () {
    return this.artalk.marked(this.textareaEl.val())
  }

  initBottomPart () {
    this.bottomPartLeftEl = this.el.find('.artalk-editor-bottom-part.artalk-left')
    this.bottomPartRightEl = this.el.find('.artalk-editor-bottom-part.artalk-right')

    this.initReply()
    this.initSubmit()
  }

  initReply () {
    this.replyComment = null
    this.sendReplyEl = null
  }

  setReply (comment) {
    if (this.replyComment !== null) {
      this.cancelReply()
    }

    if (this.sendReplyEl === null) {
      this.sendReplyEl = $('<div class="artalk-send-reply"><span class="artalk-text"></span><span class="artalk-cancel" title="取消 AT">×</span></div>')
      this.sendReplyEl.find('.artalk-text').text(`回复 -> ${comment.data.nick}`)
      this.sendReplyEl.find('.artalk-cancel').click(() => {
        this.cancelReply()
      })
      this.sendReplyEl.appendTo(this.bottomPartRightEl)
    }
    this.replyComment = comment
    this.artalk.scrollToView(this.el)
    this.textareaEl.focus()
  }

  cancelReply () {
    if (this.sendReplyEl !== null) {
      this.sendReplyEl.remove()
      this.sendReplyEl = null
    }
    this.replyComment = null
  }

  getReplyComment () {
    return this.replyComment
  }

  initSubmit () {
    this.submitBtn = this.el.find('.artalk-send-btn')
    this.submitBtn.click((evt) => {
      let btnEl = evt.currentTarget
      this.submit(btnEl)
    })
  }

  submit () {
    if ($.trim(this.getContent()) === '') {
      this.textareaEl.focus()
      return
    }

    $.ajax({
      type: 'POST',
      url: this.artalk.opts.serverUrl,
      data: {
        action: 'CommentAdd',
        content: this.getContent(),
        nick: this.user.nick,
        email: this.user.email,
        link: this.user.link,
        rid: this.getReplyComment() === null ? 0 : this.getReplyComment().data.id,
        page_key: this.artalk.opts.pageKey,
        password: this.user.password || null
      },
      dataType: 'json',
      beforeSend: () => {
        this.artalk.showLoading(this.el)
      },
      success: (obj) => {
        this.artalk.hideLoading(this.el)
        if (obj.success) {
          let newComment = new Comment(this.artalk.list, obj.data.comment)
          if (this.getReplyComment() === null) {
            this.artalk.list.putOneComment(newComment)
          } else {
            this.getReplyComment().setChild(newComment)
          }
          this.artalk.scrollToView(newComment.getElem())
          this.clearEditor()
        } else {
          if (typeof obj.data.need_password === 'boolean' && obj.data.need_password === true) {
            // 管理员密码验证
            this.showAdminCheck()
          } else {
            this.showNotify('评论失败，' + obj.msg, 'e')
          }
        }
      },
      error: () => {
        this.artalk.hideLoading(this.el)
        this.showNotify('评论失败，网络错误', 'e')
      }
    })
  }

  showAdminCheck () {
    let formElem = $(`<div>输入密码来验证管理员身份：<input type="password" required placeholder="输入密码..."></div>`)
    let input = formElem.find('[type="password"]')
    setTimeout(() => {
      input.focus() // 延迟保证有效
    }, 80)
    this.artalk.showLayerDialog(this.el, formElem, (dialogElem, btnElem) => {
      let inputVal = $.trim(input.val())
      let btnRawText = btnElem.text()
      let btnTextRestore = () => {
        btnElem.text(btnRawText)
      }
      $.ajax({
        type: 'POST',
        url: this.artalk.opts.serverUrl,
        dataType: 'json',
        data: {
          action: 'AdminCheck',
          nick: this.user.nick,
          email: this.user.email,
          password: inputVal
        },
        beforeSend: () => {
          btnElem.text('加载中...')
        },
        success: (obj) => {
          if (obj.success) {
            // 密码验证成功
            this.user.password = inputVal
            this.userLocalStorageSave()
            dialogElem.remove()
            this.submit()
          } else {
            btnElem.text('密码错误')
            input.focus(() => {
              btnTextRestore()
            })
            setTimeout(() => {
              btnTextRestore()
            }, 3000)
          }
        },
        error: () => {
          btnElem.text('网络错误')
        }
      })

      return false
    }, () => true)
  }

  showNotify (msg, type) {
    let colors = { s: '#57d59f', e: '#ff6f6c', w: '#ffc721', i: '#2ebcfc' }
    if (!colors[type]) {
      throw Error('showNotify 的 type 有问题！仅支持：' + Object.keys(colors).join(', '))
    }

    let timeout = 3000 // 持续显示时间 ms
    let wrapElem = this.el.find('.artalk-editor-notify-wrap')
    if (!wrapElem.length) {
      wrapElem = $('<div class="artalk-editor-notify-wrap"></div>').appendTo(this.el)
    }

    let notifyElem = $(`<div class="artalk-editor-notify-item artalk-fade-in" style="background-color: ${colors[type]}"><span class="artalk-editor-notify-content"></span></div>`)
    notifyElem.find('.artalk-editor-notify-content').html($('<div/>').text(msg).html().replace('\n', '<br/>'))
    notifyElem.appendTo(wrapElem)

    let notifyRemove = () => {
      notifyElem.addClass('artalk-fade-out')
      setTimeout(() => {
        notifyElem.remove()
      }, 200)
    }

    let timeoutFn
    if (timeout > 0) {
      timeoutFn = setTimeout(() => {
        notifyRemove()
      }, timeout)
    }

    notifyElem.click(() => {
      notifyRemove()
      clearTimeout(timeoutFn)
    })
  }
}

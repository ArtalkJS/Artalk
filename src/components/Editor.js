import $ from 'jquery'
import '../css/editor.scss'
import Comment from './Comment.js'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug.js'
import PreviewPlug from './editor-plugs/PreviewPlug.js'

export default class Editor {
  constructor (artalk) {
    this.artalk = artalk

    this.plugList = [EmoticonsPlug, PreviewPlug]
    this.el = $(require('../templates/Editor.ejs')(this)).appendTo(this.artalk.el)
    this.headerEl = this.el.find('.artalk-editor-header')
    this.textareaEl = this.el.find('.artalk-editor-textarea')
    this.plugWrapEl = this.el.find('.artalk-editor-plug-wrap')
    this.bottomPartLeftEl = this.el.find('.artalk-editor-bottom-part.artalk-left')
    this.plugSwitherWrapEl = this.el.find('.artalk-editor-plug-switcher-wrap')
    this.bottomPartRightEl = this.el.find('.artalk-editor-bottom-part.artalk-right')
    this.submitBtn = this.el.find('.artalk-send-btn')
    this.notifyWrapEl = this.el.find('.artalk-editor-notify-wrap')

    this.initUser()
    this.initHeader()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()
  }

  initUser () {
    let localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: localUser.nick || null,
      email: localUser.email || null,
      link: localUser.link || null,
      password: localUser.password || null
    }
  }

  initHeader () {
    for (let field in this.user) {
      let inputEl = this.headerEl.find(`[name="${field}"]`)
      if (inputEl.length) {
        inputEl.val(this.user[field] || '')
        // 输入框内容变化事件
        inputEl.bind('input propertychange', (evt) => {
          this.user[field] = $.trim(inputEl.val())
          this.user.password = null
          this.saveUser()
        })
      }
    }
  }

  /**
   * 保存用户到 localStorage 中
   */
  saveUser () {
    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.user))
  }

  initTextarea () {
    // 修复按下 Tab 输入的内容
    this.textareaEl.on('keydown', (e) => {
      var keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.insertContent('\t')
      }
    })

    // 输入框高度随内容而变化
    this.textareaEl.bind('input propertychange', (evt) => {
      var diff = this.textareaEl.outerHeight() - this.textareaEl[0].clientHeight
      this.textareaEl[0].style.height = 0 // 若不加此行，内容减少，高度回不去
      this.textareaEl[0].style.height = this.textareaEl[0].scrollHeight + diff + 'px'
    })
  }

  initEditorPlug () {
    this.plugs = {}
    let openedPlugName = null

    // 依次实例化 plug
    for (let i in this.plugList) {
      let plug = new (this.plugList[i])(this)
      this.plugs[plug.getName()] = plug

      // 切换按钮
      let btnElem = $(`<span class="artalk-editor-action artalk-editor-plug-switcher" data-plug-index="${i}">${plug.getBtnHtml()}</span>`)
      btnElem.appendTo(this.plugSwitherWrapEl)
      btnElem.click(() => {
        this.plugSwitherWrapEl.find('.active').removeClass('active')

        // 若点击已打开的，则收起
        if (openedPlugName === plug.getName()) {
          plug.onHide()
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

        this.plugWrapEl.children().each((i, plugItemEl) => {
          plugItemEl = $(plugItemEl)
          let plugItemName = plugItemEl.attr('data-plug-name')
          if (plugItemName === plug.getName()) {
            plugItemEl.css('display', '')
            this.plugs[plugItemName].onShow()
          } else {
            plugItemEl.css('display', 'none')
            this.plugs[plugItemName].onHide()
          }
        })

        this.plugWrapEl.css('display', '')
        openedPlugName = plug.getName()

        btnElem.addClass('active')
      })
    }
  }

  insertContent (val) {
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
    this.artalk.ui.scrollIntoView(this.el)
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
        this.artalk.ui.showLoading(this.el)
      },
      success: (obj) => {
        this.artalk.ui.hideLoading(this.el)
        if (obj.success) {
          let newComment = new Comment(this.artalk.list, obj.data.comment)
          if (this.getReplyComment() === null) {
            this.artalk.list.putComment(newComment)
          } else {
            this.getReplyComment().setChild(newComment)
          }
          this.artalk.ui.scrollIntoView(newComment.getElem())
          this.clearEditor()
        } else {
          if (obj.data !== null && typeof obj.data.need_password === 'boolean' && obj.data.need_password === true) {
            // 管理员密码验证
            this.showAdminCheck()
          } else {
            this.showNotify('评论失败，' + obj.msg, 'e')
          }
        }
      },
      error: () => {
        this.artalk.ui.hideLoading(this.el)
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
    this.artalk.ui.showDialog(this.el, formElem, (dialogElem, btnElem) => {
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
            this.saveUser()
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
    this.artalk.ui.showNotify(msg, type, this.notifyWrapEl)
  }
}

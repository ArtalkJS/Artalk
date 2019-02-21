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
    this.textareaWrapEl = this.el.find('.artalk-editor-textarea-wrap')
    this.textareaEl = this.el.find('.artalk-editor-textarea')
    this.plugWrapEl = this.el.find('.artalk-editor-plug-wrap')
    this.bottomPartLeftEl = this.el.find('.artalk-editor-bottom-part.artalk-left')
    this.plugSwitherWrapEl = this.el.find('.artalk-editor-plug-switcher-wrap')
    this.bottomPartRightEl = this.el.find('.artalk-editor-bottom-part.artalk-right')
    this.submitBtn = this.el.find('.artalk-send-btn')
    this.notifyWrapEl = this.el.find('.artalk-editor-notify-wrap')

    this.initLocalStorge()
    this.initHeader()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()
  }

  initLocalStorge () {
    let localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: localUser.nick || null,
      email: localUser.email || null,
      link: localUser.link || null,
      password: localUser.password || null
    }

    let localContent = window.localStorage.getItem('ArtalkContent') || ''
    if ($.trim(localContent) !== '') {
      this.showNotify('已自动恢复', 'i')
      this.setContent(localContent)
    }
    this.textareaEl.bind('input propertychange', (evt) => {
      this.saveContent()
    })
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

  saveContent () {
    window.localStorage.setItem('ArtalkContent', $.trim(this.getContent()))
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
      this.setContent(this.textareaEl.val().substring(0, sStart) + val + this.textareaEl.val().substring(sEnd, this.textareaEl.val().length))
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
    this.saveContent()
    if (!!this.plugs && !!this.plugs['preview']) {
      this.plugs['preview'].updateContent()
    }
  }

  clearEditor () {
    this.setContent('')
    this.cancelReply()
  }

  getContent () {
    return this.textareaEl.val()
  }

  getContentMarked () {
    return this.artalk.marked(this.getContent())
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
      this.sendReplyEl.find('.artalk-text').text(`@${comment.data.nick}`)
      this.sendReplyEl.find('.artalk-cancel').click(() => {
        this.cancelReply()
      })
      this.sendReplyEl.appendTo(this.textareaWrapEl)
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

    this.artalk.request('CommentAdd', {
      content: this.getContent(),
      nick: this.user.nick,
      email: this.user.email,
      link: this.user.link,
      rid: this.getReplyComment() === null ? 0 : this.getReplyComment().data.id,
      page_key: this.artalk.opts.pageKey,
      password: this.user.password,
      captcha: this.submitCaptchaVal || null
    }, () => {
      this.artalk.ui.showLoading(this.el)
    }, () => {
      this.artalk.ui.hideLoading(this.el)
    }, (msg, data) => {
      let newComment = new Comment(this.artalk.list, data.comment)
      if (this.getReplyComment() !== null) {
        this.getReplyComment().setChild(newComment)
      } else {
        this.artalk.list.putComment(newComment)
      }
      this.artalk.ui.scrollIntoView(newComment.getElem())
      this.clearEditor()
    }, (msg, data) => {
      if (data !== null && !$.isEmptyObject(data) && typeof data.need_password === 'boolean' && data.need_password === true) {
        // 管理员密码验证
        this.showCheck('密码')
      } else if (data !== null && !$.isEmptyObject(data) && typeof data.need_captcha === 'boolean' && data.need_captcha === true) {
        // 验证码验证
        this.submitCaptchaImgData = data.img_data
        this.showCheck('验证码')
      } else {
        this.showNotify(`评论失败，${msg}`, 'e')
      }
    })
  }

  showCheck (typeKey) {
    let typeList = {
      '密码': {
        body: () => $('<span>输入密码来验证管理员身份：</span>'),
        reqAct: 'AdminCheck',
        reqObj: (inputVal) => {
          return {
            nick: this.user.nick,
            email: this.user.email,
            password: inputVal
          }
        },
        onSuccess: (msg, data, inputVal) => {
          this.user.password = inputVal
          this.saveUser()
        }
      },
      '验证码': {
        body: () => {
          let elem = typeList['验证码'].elem = $(`<span><img class="artalk-captcha-img" src="${this.submitCaptchaImgData || ''}" alt="验证码">输入验证码继续：</span>`)
          elem.find('.artalk-captcha-img').click(() => {
            typeList['验证码'].refresh()
          })
          return elem
        },
        reqAct: 'CaptchaCheck',
        reqObj: (inputVal) => {
          return {
            captcha: inputVal
          }
        },
        onSuccess: (msg, data, inputVal) => {
          this.submitCaptchaVal = inputVal
        },
        refresh: (imgData) => {
          let elem = typeList['验证码'].elem
          let imgEl = elem.find('.artalk-captcha-img')
          if (!imgData) {
            this.artalk.request('CaptchaCheck', { refresh: true }, () => {}, () => {}, (msg, data) => {
              imgEl.attr('src', data.img_data)
            }, () => {})
          } else {
            imgEl.attr('src', imgData)
          }
        }
      }
    }

    let type = typeList[typeKey]
    if (!type) {
      throw Error('type 有问题，仅支持：' + Object.keys(type).join(', '))
    }

    let formElem = $(`<div></div>`)
    $(type.body()).appendTo(formElem)
    let input = $(`<input id="check" type="${(typeKey === '密码' ? 'password' : 'text')}" required placeholder="输入${typeKey}...">`).appendTo(formElem)
    setTimeout(() => {
      input.focus() // 延迟保证有效
    }, 80)
    this.artalk.ui.showDialog(this.el, formElem, (dialogElem, btnElem) => {
      let inputVal = $.trim(input.val())
      let btnRawText = btnElem.text()
      let btnTextSet = (btnText) => {
        btnElem.text(btnText)
        btnElem.addClass('error')
      }
      let btnTextRestore = () => {
        btnElem.text(btnRawText)
        btnElem.removeClass('error')
      }

      this.artalk.request(type.reqAct, type.reqObj(inputVal), () => {
        btnElem.text('加载中...')
      }, () => {

      }, (msg, data) => {
        type.onSuccess(msg, data, inputVal)
        dialogElem.remove()
        this.submit()
      }, (msg, data) => {
        btnTextSet(msg)
        if (typeKey === '验证码') {
          type.refresh(data.img_data)
        }
        let tf = setTimeout(() => {
          btnTextRestore()
        }, 3000)
        input.focus(() => {
          btnTextRestore()
          clearTimeout(tf)
        })
      })

      return false
    }, () => true)
  }

  showNotify (msg, type) {
    this.artalk.ui.showNotify(msg, type, this.notifyWrapEl)
  }
}

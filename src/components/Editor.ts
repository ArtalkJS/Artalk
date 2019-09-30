import '../css/editor.scss'
import Comment from './Comment'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug'
import PreviewPlug from './editor-plugs/PreviewPlug'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

export default class Editor extends ArtalkContext {
  private readonly LOADABLE_PLUG_LIST = [EmoticonsPlug, PreviewPlug]
  public plugList: { [name: string]: any }

  public el: HTMLElement

  public headerEl: HTMLElement
  public textareaWrapEl: HTMLElement
  public textareaEl: HTMLTextAreaElement
  public plugWrapEl: HTMLElement
  public bottomPartLeftEl: HTMLElement
  public plugSwitcherWrapEl: HTMLElement
  public bottomPartRightEl: HTMLElement
  public submitBtn: HTMLButtonElement
  public notifyWrapEl: HTMLElement

  public user: {
    nick: string|null,
    email: string|null,
    link: string|null,
    password: string|null
  }

  private replyComment: Comment
  private sendReplyEl: HTMLElement

  private submitCaptchaVal: string
  private submitCaptchaImgData: string

  constructor () {
    super()

    this.el = Utils.createElement(require('../templates/Editor.ejs')(this))
    this.el = Utils.createElement(require('../templates/Editor.ejs')(this))
    this.artalk.el.appendChild(this.el)

    this.headerEl = this.el.querySelector('.artalk-editor-header')
    this.textareaWrapEl = this.el.querySelector('.artalk-editor-textarea-wrap')
    this.textareaEl = this.el.querySelector('.artalk-editor-textarea')
    this.plugWrapEl = this.el.querySelector('.artalk-editor-plug-wrap')
    this.bottomPartLeftEl = this.el.querySelector('.artalk-editor-bottom-part.artalk-left')
    this.plugSwitcherWrapEl = this.el.querySelector('.artalk-editor-plug-switcher-wrap')
    this.bottomPartRightEl = this.el.querySelector('.artalk-editor-bottom-part.artalk-right')
    this.submitBtn = this.el.querySelector('.artalk-send-btn')
    this.notifyWrapEl = this.el.querySelector('.artalk-editor-notify-wrap')

    this.initLocalStorage()
    this.initHeader()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()
  }

  initLocalStorage () {
    const localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: localUser.nick || null,
      email: localUser.email || null,
      link: localUser.link || null,
      password: localUser.password || null
    }

    const localContent = window.localStorage.getItem('ArtalkContent') || ''
    if (localContent.trim() !== '') {
      this.showNotify('已自动恢复', 'i')
      this.setContent(localContent)
    }
    this.textareaEl.addEventListener('input', () => {
      this.saveContent()
    })
  }

  initHeader () {
    Object.keys(this.user).forEach((field) => {
      const inputEl: HTMLInputElement = this.headerEl.querySelector(`[name="${field}"]`)
      if (inputEl !== null && inputEl instanceof HTMLInputElement) {
        inputEl.value = this.user[field] || ''
        // 输入框内容变化事件
        inputEl.addEventListener('input', (evt) => {
          this.user[field] = inputEl.value.trim()
          this.user.password = null
          this.saveUser()
        })
      }
    })
  }

  /**
   * 保存用户到 localStorage 中
   */
  saveUser () {
    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.user))
  }

  saveContent () {
    window.localStorage.setItem('ArtalkContent', this.getContent().trim())
  }

  initTextarea () {
    // 修复按下 Tab 输入的内容
    this.textareaEl.addEventListener('keydown', (e) => {
      const keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.insertContent('\t')
      }
    })

    // 输入框高度随内容而变化
    this.textareaEl.addEventListener('input', (evt) => {
      const diff = this.textareaEl.offsetHeight - this.textareaEl.clientHeight
      this.textareaEl.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
      this.textareaEl.style.height = `${this.textareaEl.scrollHeight + diff}px`
    })
  }

  initEditorPlug () {
    this.plugList = {}
    let openedPlugName = null

    // 依次实例化 plug
    this.LOADABLE_PLUG_LIST.forEach((PlugObj) => {
      const plug = new PlugObj(this)
      this.plugList[plug.getName()] = plug

      // 切换按钮
      const btnElem = Utils.createElement(`<span class="artalk-editor-action artalk-editor-plug-switcher">${plug.getBtnHtml()}</span>`)
      this.plugSwitcherWrapEl.appendChild(btnElem)
      btnElem.addEventListener('click', () => {
        this.plugSwitcherWrapEl.querySelectorAll('.active').forEach(item => item.classList.remove('active'))

        // 若点击已打开的，则收起
        if (plug.getName() === openedPlugName) {
          plug.onHide()
          this.plugWrapEl.style.display = 'none'
          openedPlugName = null
          return
        }

        if (this.plugWrapEl.querySelector(`[data-plug-name="${plug.getName()}"]`) === null) {
          // 需要初始化
          const plugEl = plug.getElem()
          plugEl.setAttribute('data-plug-name', plug.getName())
          plugEl.style.display = 'none'
          this.plugWrapEl.appendChild(plugEl)
        }

        Array.from(this.plugWrapEl.children).forEach((plugItemEl: HTMLElement) => {
          const plugItemName = plugItemEl.getAttribute('data-plug-name')
          if (plugItemName === plug.getName()) {
            plugItemEl.style.display = ''
            this.plugList[plugItemName].onShow()
          } else {
            plugItemEl.style.display = 'none'
            this.plugList[plugItemName].onHide()
          }
        })

        this.plugWrapEl.style.display = ''
        openedPlugName = plug.getName()

        btnElem.classList.add('active')
      })
    })
  }

  insertContent (val: string) {
    if ((document as any).selection) {
      this.textareaEl.focus();
      (document as any).selection.createRange().text = val
      this.textareaEl.focus()
    } else if (this.textareaEl.selectionStart || this.textareaEl.selectionStart === 0) {
      const sStart = this.textareaEl.selectionStart
      const sEnd = this.textareaEl.selectionEnd
      const sT = this.textareaEl.scrollTop
      this.setContent(this.textareaEl.value.substring(0, sStart) + val + this.textareaEl.value.substring(sEnd, this.textareaEl.value.length))
      this.textareaEl.focus()
      this.textareaEl.selectionStart = sStart + val.length
      this.textareaEl.selectionEnd = sStart + val.length
      this.textareaEl.scrollTop = sT
    } else {
      this.textareaEl.focus()
      this.textareaEl.value += val
    }
  }

  setContent (val: string) {
    this.textareaEl.value = val
    this.saveContent()
    if (!!this.plugList && !!this.plugList.preview) {
      this.plugList.preview.updateContent()
    }
  }

  clearEditor () {
    this.setContent('')
    this.cancelReply()
  }

  getContent () {
    return this.textareaEl.value || ''
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

  setReply (comment: Comment) {
    if (this.replyComment !== null) {
      this.cancelReply()
    }

    if (this.sendReplyEl === null) {
      this.sendReplyEl = Utils.createElement('<div class="artalk-send-reply"><span class="artalk-text"></span><span class="artalk-cancel" title="取消 AT">×</span></div>');
      (this.sendReplyEl.querySelector('.artalk-text') as HTMLElement).innerText = `@${comment.data.nick}`
      this.sendReplyEl.querySelector('.artalk-cancel').addEventListener('click', () => {
        this.cancelReply()
      })
      this.textareaWrapEl.appendChild(this.sendReplyEl)
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
    this.submitBtn.addEventListener('click', (evt) => {
      const btnEl = evt.currentTarget
      this.submit()
    })
  }

  submit () {
    if (this.getContent().trim() === '') {
      this.textareaEl.focus()
      return
    }

    this.artalk.request('CommentAdd', {
      content: this.getContent(),
      nick: this.user.nick,
      email: this.user.email,
      link: this.user.link,
      rid: this.getReplyComment() === null ? 0 : this.getReplyComment().data.id,
      page_key: this.artalk.conf.pageKey,
      password: this.user.password,
      captcha: this.submitCaptchaVal || null
    }, () => {
      this.artalk.ui.showLoading(this.el)
    }, () => {
      this.artalk.ui.hideLoading(this.el)
    }, (msg, data) => {
      const newComment = new Comment(this.artalk.list, data.comment)
      if (this.getReplyComment() !== null) {
        this.getReplyComment().setChild(newComment)
      } else {
        this.artalk.list.putComment(newComment)
      }
      this.artalk.ui.scrollIntoView(newComment.getElem())
      this.clearEditor()
    }, (msg, data) => {
      if ((typeof data === 'object') && data !== null && typeof data.need_password === 'boolean' && data.need_password === true) {
        // 管理员密码验证
        this.showCheck('密码')
      } else if ((typeof data === 'object') && data !== null && typeof data.need_captcha === 'boolean' && data.need_captcha === true) {
        // 验证码验证
        this.submitCaptchaImgData = data.img_data
        this.showCheck('验证码')
      } else {
        this.showNotify(`评论失败，${msg}`, 'e')
      }
    })
  }

  private readonly CHECKER_TYPE_LIST: {[key: string]: Checker} = {
    '密码': {
      body: () => Utils.createElement('<span>输入密码来验证管理员身份：</span>'),
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
        const elem = Utils.createElement(`<span><img class="artalk-captcha-img" src="${this.submitCaptchaImgData || ''}" alt="验证码">输入验证码继续：</span>`)
        this.CHECKER_TYPE_LIST['验证码'].elem = elem;
        (elem.querySelector('.artalk-captcha-img') as HTMLElement).onclick = () => {
          this.CHECKER_TYPE_LIST['验证码'].refresh()
        }
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
      refresh: (imgData?: string) => {
        const { elem } = this.CHECKER_TYPE_LIST['验证码']
        const imgEl = elem.querySelector('.artalk-captcha-img')
        if (!imgData) {
          this.artalk.request('CaptchaCheck', { refresh: true }, () => {}, () => {}, (msg, data) => {
            imgEl.setAttribute('src', data.img_data)
          }, () => {})
        } else {
          imgEl.setAttribute('src', imgData)
        }
      }
    }
  }

  showCheck (typeKey: string) {
    const type = this.CHECKER_TYPE_LIST[typeKey]
    if (!type) {
      throw Error(`type 有问题，仅支持：${Object.keys(type).join(', ')}`)
    }

    const formElem = Utils.createElement()
    formElem.appendChild(type.body())

    const input = Utils.createElement(`<input id="check" type="${(typeKey === '密码' ? 'password' : 'text')}" required placeholder="输入${typeKey}...">`) as HTMLInputElement
    formElem.appendChild(input)
    setTimeout(() => {
      input.focus() // 延迟保证有效
    }, 80)
    this.artalk.ui.showDialog(this.el, formElem, (dialogElem, btnElem: HTMLElement) => {
      const inputVal = input.value.trim()
      const btnRawText = btnElem.innerText
      const btnTextSet = (btnText: string) => {
        btnElem.innerText = btnText
        btnElem.classList.add('error')
      }
      const btnTextRestore = () => {
        btnElem.innerText = btnRawText
        btnElem.classList.remove('error')
      }

      this.artalk.request(type.reqAct, type.reqObj(inputVal), () => {
        btnElem.innerText = '加载中...'
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
        const tf = setTimeout(() => {
          btnTextRestore()
        }, 3000)
        input.onfocus = () => {
          btnTextRestore()
          clearTimeout(tf)
        }
      })

      return false
    }, () => true)
  }

  showNotify (msg: string, type) {
    this.artalk.ui.showNotify(msg, type, this.notifyWrapEl)
  }
}

interface Checker {
  elem?: HTMLElement
  body: () => HTMLElement
  reqAct: string
  reqObj: (inputVal: string) => any
  onSuccess: (msg: string, data: any, inputVal: string) => void
  refresh?: Function
}

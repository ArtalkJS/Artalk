import '../css/editor.less'
import Comment from './Comment'
import EmoticonsPlug from './editor-plugs/EmoticonsPlug'
import PreviewPlug from './editor-plugs/PreviewPlug'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'
import Layer from './Layer'
import Checker from './Checker'

export default class Editor extends ArtalkContext {
  private readonly LOADABLE_PLUG_LIST = [EmoticonsPlug, PreviewPlug]
  public plugList: { [name: string]: any } = {}

  public el: HTMLElement

  public headerEl: HTMLElement
  public textareaWrapEl: HTMLElement
  public textareaEl: HTMLTextAreaElement
  public closeCommentEl: HTMLTextAreaElement
  public plugWrapEl: HTMLElement
  public bottomEl: HTMLElement
  public bottomPartLeftEl: HTMLElement
  public plugSwitcherWrapEl: HTMLElement
  public bottomPartRightEl: HTMLElement
  public submitBtn: HTMLButtonElement
  public notifyWrapEl: HTMLElement

  private replyComment: Comment
  private sendReplyEl: HTMLElement

  private get user () {
    return this.artalk.user
  }

  constructor () {
    super()

    this.el = Utils.createElement(require('../templates/Editor.ejs')(this))
    this.artalk.el.appendChild(this.el)

    this.headerEl = this.el.querySelector('.artalk-editor-header')
    this.textareaWrapEl = this.el.querySelector('.artalk-editor-textarea-wrap')
    this.textareaEl = this.el.querySelector('.artalk-editor-textarea')
    this.closeCommentEl = this.el.querySelector('.artalk-close-comment')
    this.plugWrapEl = this.el.querySelector('.artalk-editor-plug-wrap')
    this.bottomEl = this.el.querySelector('.artalk-editor-bottom')
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
    Object.keys(this.user.data).forEach((field) => {
      const inputEl: HTMLInputElement = this.headerEl.querySelector(`[name="${field}"]`)
      if (inputEl !== null && inputEl instanceof HTMLInputElement) {
        inputEl.value = this.user.data[field] || ''
        // 输入框内容变化事件
        inputEl.addEventListener('input', (evt) => {
          this.user.data[field] = inputEl.value.trim()

          // 若修改的是 nick or email
          if (field !== 'link') {
            this.user.data.password = ''
            this.user.data.isAdmin = false
            if (this.user.checkHasBasicUserInfo()
              && this.artalk.list.checkNickEmailIsAdmin(this.user.data.nick, this.user.data.email)) {
              // 昵称为管理员，显示管理员密码验证 dialog
              Checker.checkAction('管理员', () => {
                this.artalk.list.refreshUI()
              })
            }
          }

          this.saveUser()
          this.artalk.list.refreshUI()
        })
      }
    })
  }

  saveUser () {
    this.user.save()
  }

  saveContent () {
    window.localStorage.setItem('ArtalkContent', this.getContentOriginal().trim())
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
      this.adjustTextareaHeight()
    })
  }

  adjustTextareaHeight () {
    const diff = this.textareaEl.offsetHeight - this.textareaEl.clientHeight
    this.textareaEl.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.textareaEl.style.height = `${this.textareaEl.scrollHeight + diff}px`
  }

  openedPlugName: string = null

  initEditorPlug () {
    this.plugList = {}
    this.plugWrapEl.innerHTML = ''
    this.plugWrapEl.style.display = 'none'
    this.openedPlugName = null
    this.plugSwitcherWrapEl.innerHTML = ''

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
        if (plug.getName() === this.openedPlugName) {
          plug.onHide()
          this.plugWrapEl.style.display = 'none'
          this.openedPlugName = null
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
        this.openedPlugName = plug.getName()

        btnElem.classList.add('active')
      })
    })
  }

  /** 关闭编辑器插件 */
  closePlug () {
    this.plugWrapEl.innerHTML = ''
    this.plugWrapEl.style.display = 'none'
    this.openedPlugName = null
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
    this.adjustTextareaHeight()
  }

  clearEditor () {
    this.setContent('')
    this.cancelReply()
  }

  getContent () {
    let content = this.getContentOriginal()

    // 表情包处理
    if (this.plugList && this.plugList.emoticons) {
      const emoticonsPlug = this.plugList.emoticons as EmoticonsPlug
      content = emoticonsPlug.transEmoticonImageText(content)
    }

    return content
  }

  getContentOriginal () {
    return this.textareaEl.value || '' // Tip: !!"0" === true
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
      nick: this.user.data.nick,
      email: this.user.data.email,
      link: this.user.data.link,
      rid: this.getReplyComment() === null ? 0 : this.getReplyComment().data.id,
      page_key: this.artalk.conf.pageKey,
      password: this.user.data.password,
      captcha: Checker.submitCaptchaVal || ''
    }, () => {
      this.artalk.ui.showLoading(this.el)
    }, () => {
      this.artalk.ui.hideLoading(this.el)
    }, (msg, data) => {
      const newComment = new Comment(this.artalk.list, data.comment)
      if (this.getReplyComment() !== null) {
        this.getReplyComment().putChild(newComment)
      } else {
        this.artalk.list.putRootComment(newComment)
      }
      this.clearEditor() // 清空编辑器

      this.artalk.ui.scrollIntoView(newComment.getElem())
      newComment.playFadeInAnim() // 播放评论渐出动画
      this.artalk.list.data.total += 1 // 评论数增加 1
      this.artalk.list.refreshUI() // 更新 list 界面
    }, (msg, data) => {
      if ((typeof data === 'object') && data !== null && typeof data.need_password === 'boolean' && data.need_password === true) {
        // 管理员密码验证
        Checker.checkAction('管理员', () => {
          this.submit()
        })
      } else if ((typeof data === 'object') && data !== null && typeof data.need_captcha === 'boolean' && data.need_captcha === true) {
        // 验证码验证
        Checker.submitCaptchaImgData = data.img_data
        Checker.checkAction('验证码', () => {
          this.submit()
        })
      } else {
        this.showNotify(`评论失败，${msg}`, 'e')
      }
    })
  }

  showNotify (msg: string, type) {
    this.artalk.ui.showNotify(msg, type, this.notifyWrapEl)
  }

  /** 关闭评论 */
  closeComment () {
    this.closeCommentEl.style.display = ''

    if (!this.artalk.user.data.isAdmin) {
      this.textareaEl.style.display = 'none'
      this.closePlug()
      this.bottomEl.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.textareaEl.style.display = ''
      this.bottomEl.style.display = ''
    }
  }

  /** 打开评论 */
  openComment () {
    this.closeCommentEl.style.display = 'none'
    this.textareaEl.style.display = ''
    this.bottomEl.style.display = ''
  }
}


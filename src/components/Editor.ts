import '../style/editor.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import EditorHTML from './html/editor.html?raw'

import EmoticonsPlug from './editor-plugs/EmoticonsPlug'
import PreviewPlug from './editor-plugs/PreviewPlug'
import { CommentData } from '~/types/artalk-data'
import Api from '../lib/api'

export default class Editor extends Component {
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

  private replyComment: CommentData|null = null
  private sendReplyEl: HTMLElement|null = null

  private get user () {
    return this.ctx.user
  }

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(EditorHTML)

    this.headerEl = this.el.querySelector('.atk-editor-header')!
    this.textareaWrapEl = this.el.querySelector('.atk-editor-textarea-wrap')!
    this.textareaEl = this.el.querySelector('.atk-editor-textarea')!
    this.closeCommentEl = this.el.querySelector('.atk-close-comment')!
    this.plugWrapEl = this.el.querySelector('.atk-editor-plug-wrap')!
    this.bottomEl = this.el.querySelector('.atk-editor-bottom')!
    this.bottomPartLeftEl = this.el.querySelector('.atk-editor-bottom-part.atk-left')!
    this.plugSwitcherWrapEl = this.el.querySelector('.atk-editor-plug-switcher-wrap')!
    this.bottomPartRightEl = this.el.querySelector('.atk-editor-bottom-part.atk-right')!
    this.submitBtn = this.el.querySelector('.atk-send-btn')!
    this.notifyWrapEl = this.el.querySelector('.atk-editor-notify-wrap')!

    this.initLocalStorage()
    this.initHeader()
    this.initTextarea()
    this.initEditorPlug()
    this.initBottomPart()

    // 监听事件
    this.ctx.addEventListener('editor-open-comment', () => (this.openComment()))
    this.ctx.addEventListener('editor-close-comment', () => (this.closeComment()))
    this.ctx.addEventListener('editor-reply', (commentData) => (this.setReply(commentData)))
    this.ctx.addEventListener('editor-show-loading', () => (Ui.showLoading(this.el)))
    this.ctx.addEventListener('editor-hide-loading', () => (Ui.hideLoading(this.el)))
    this.ctx.addEventListener('editor-notify', (f) => (this.showNotify(f.msg, f.type)))
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
      const inputEl = this.getInputEl(field)
      if (inputEl && inputEl instanceof HTMLInputElement) {
        inputEl.value = this.user.data[field] || ''
        // 绑定事件
        inputEl.addEventListener('input', () => this.onHeaderInputChanged(field, inputEl))
      }
    })
  }

  getInputEl (field: string) {
    const inputEl = this.headerEl.querySelector<HTMLInputElement>(`[name="${field}"]`)
    return inputEl
  }

  queryUserInfo = {
    timeout: <number|null>null,
    abortFunc: <(() => void)|null>null
  }

  /** header 输入框内容变化事件 */
  onHeaderInputChanged (field: string, inputEl: HTMLInputElement) {
    this.user.data[field] = inputEl.value.trim()

    // 若修改的是 nick or email
    if (field === 'nick' || field === 'email') {
      this.user.data.token = '' // 清除 token 登陆状态
      this.user.data.isAdmin = false

      // 获取用户信息
      if (this.queryUserInfo.timeout !== null) window.clearTimeout(this.queryUserInfo.timeout) // 清除待发出的请求
      if (this.queryUserInfo.abortFunc !== null) this.queryUserInfo.abortFunc() // 之前发出未完成的请求立刻中止

      this.queryUserInfo.timeout = window.setTimeout(() => {
        this.queryUserInfo.timeout = null // 清理

        const {req, abort} = new Api(this.ctx).userGet(
          this.user.data.nick, this.user.data.email
        )
        this.queryUserInfo.abortFunc = abort
        req.then(data => {
          if (!data.is_login) {
            this.user.data.token = ''
            this.user.data.isAdmin = false
          }

          // 若用户为管理员，执行登陆操作
          if (this.user.checkHasBasicUserInfo() && !data.is_login && data.user && data.user.is_admin) {
            this.showLoginDialog()
          }

          // 自动填入 link
          if (data.user && data.user.link) {
            this.user.data.link = data.user.link
            this.getInputEl('link')!.value = data.user.link
          }
        })
        .finally(() => {
          this.queryUserInfo.abortFunc = null // 清理
        })
      }, 400) // 延迟执行，减少请求次数
    }

    this.saveUser()
  }

  showLoginDialog () {
    this.ctx.dispatchEvent('checker-admin', {
      onSuccess: () => {
      }
    })
  }

  saveUser () {
    this.user.save()
    this.ctx.dispatchEvent('user-changed')
  }

  saveContent () {
    window.localStorage.setItem('ArtalkContent', this.getContentOriginal().trim())
  }

  initTextarea () {
    // 占位符
    this.textareaEl.placeholder = this.ctx.conf.placeholder || ''

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

  openedPlugName: string|null = null

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
      const btnElem = Utils.createElement(`<span class="atk-editor-action atk-editor-plug-switcher">${plug.getBtnHtml()}</span>`)
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
          const plugEl = plug.getEl()
          plugEl.setAttribute('data-plug-name', plug.getName())
          plugEl.style.display = 'none'
          this.plugWrapEl.appendChild(plugEl)
        }

        (Array.from(this.plugWrapEl.children) as HTMLElement[]).forEach((plugItemEl: HTMLElement) => {
          const plugItemName = plugItemEl.getAttribute('data-plug-name')!
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
    return Utils.marked(this.ctx, this.getContent())
  }

  initBottomPart () {
    this.initReply()
    this.initSubmit()
  }

  initReply () {
    this.replyComment = null
    this.sendReplyEl = null
  }

  setReply (commentData: CommentData) {
    if (this.replyComment !== null) {
      this.cancelReply()
    }

    if (this.sendReplyEl === null) {
      this.sendReplyEl = Utils.createElement('<div class="atk-send-reply-wrap"><div class="atk-send-reply">回复 <span class="atk-text"></span><span class="atk-cancel" title="取消 AT">×</span></div></div>');
      this.sendReplyEl.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      this.sendReplyEl.addEventListener('click', () => {
        this.cancelReply()
      })
      this.textareaWrapEl.prepend(this.sendReplyEl)
    }
    this.replyComment = commentData
    Ui.scrollIntoView(this.el)
    this.textareaEl.focus()
  }

  cancelReply () {
    if (this.sendReplyEl !== null) {
      this.sendReplyEl.remove()
      this.sendReplyEl = null
    }
    this.replyComment = null
  }

  initSubmit () {
    this.submitBtn.innerText = this.ctx.conf.sendBtn || 'Send'

    this.submitBtn.addEventListener('click', (evt) => {
      const btnEl = evt.currentTarget
      this.submit()
    })
  }

  async submit () {
    if (this.getContent().trim() === '') {
      this.textareaEl.focus()
      return
    }

    Ui.showLoading(this.el)

    try {
      const nComment = await new Api(this.ctx).add({
        content: this.getContent(),
        nick: this.user.data.nick,
        email: this.user.data.email,
        link: this.user.data.link,
        rid: this.replyComment === null ? 0 : this.replyComment.id
      })

      this.ctx.dispatchEvent('list-insert', nComment)
      this.clearEditor() // 清空编辑器
    } catch (err: any) {
      this.showNotify(`评论失败，${err.msg || String(err)}`, 'e')
    } finally {
      Ui.hideLoading(this.el)
    }
  }

  showNotify (msg: string, type) {
    Ui.showNotify(this.notifyWrapEl, msg, type)
  }

  /** 关闭评论 */
  closeComment () {
    this.closeCommentEl.style.display = ''

    if (!this.user.data.isAdmin) {
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


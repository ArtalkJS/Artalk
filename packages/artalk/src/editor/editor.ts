import '@/style/editor.less'

import { CommentData } from '~/types/artalk-data'
import Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Api from '../api'

import EditorHTML from './editor.html?raw'

import EmoticonsPlug from './plugs/emoticons-plug'
import UploadPlug from './plugs/upload-plug'
import PreviewPlug from './plugs/preview-plug'
import EditorPlug from './plugs/editor-plug'
import HeaderInputPlug from './plugs/header-input-plug'

export default class Editor extends Component {
  private get user() { return this.ctx.user }

  public $header: HTMLElement
  public $nick: HTMLInputElement
  public $email: HTMLInputElement
  public $link: HTMLInputElement
  public get $inputs() {
    return { nick: this.$nick, email: this.$email, link: this.$link }
  }

  public $textareaWrap: HTMLElement
  public $textarea: HTMLTextAreaElement
  public $bottom: HTMLElement
  public $submitBtn: HTMLButtonElement
  public $notifyWrap: HTMLElement

  private replyComment: CommentData|null = null
  private $sendReply: HTMLElement|null = null

  private isTraveling = false

  /** 启用的插件 */
  private readonly ENABLED_PLUGS = [ EmoticonsPlug, UploadPlug, PreviewPlug, HeaderInputPlug ]
  public plugList: { [name: string]: EditorPlug } = {}
  private openedPlugName: string|null = null
  public $plugPanelWrap: HTMLElement
  public $plugBtnWrap: HTMLElement

  public constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(EditorHTML)

    this.$header = this.$el.querySelector('.atk-header')!
    this.$nick = this.$header.querySelector('[name="nick"]')!
    this.$email = this.$header.querySelector('[name="email"]')!
    this.$link = this.$header.querySelector('[name="link"]')!

    this.$textareaWrap = this.$el.querySelector('.atk-textarea-wrap')!
    this.$textarea = this.$el.querySelector('.atk-textarea')!
    this.$bottom = this.$el.querySelector('.atk-bottom')!
    this.$submitBtn = this.$el.querySelector('.atk-send-btn')!
    this.$notifyWrap = this.$el.querySelector('.atk-notify-wrap')!

    this.$plugBtnWrap = this.$el.querySelector('.atk-plug-btn-wrap')!
    this.$plugPanelWrap = this.$el.querySelector('.atk-plug-panel-wrap')!

    // 执行初始化
    this.initLocalStorage()
    this.initHeader()
    this.initTextarea()
    this.initPlugs()
    this.initSubmitBtn()

    // 监听事件
    this.ctx.on('editor-open', () => (this.open()))
    this.ctx.on('editor-close', () => (this.close()))
    this.ctx.on('editor-reply', (p) => (this.setReply(p.data, p.$el, p.scroll)))
    this.ctx.on('editor-reply-cancel', () => (this.cancelReply()))
    this.ctx.on('editor-show-loading', () => (Ui.showLoading(this.$el)))
    this.ctx.on('editor-hide-loading', () => (Ui.hideLoading(this.$el)))
    this.ctx.on('editor-notify', (f) => (this.showNotify(f.msg, f.type)))
    this.ctx.on('editor-travel', ($el) => (this.travel($el)))
    this.ctx.on('editor-travel-back', () => (this.travelBack()))
    this.ctx.on('conf-updated', () => {})
  }

  private initLocalStorage() {
    const localContent = window.localStorage.getItem('ArtalkContent') || ''
    if (localContent.trim() !== '') {
      this.showNotify(this.$t('restoredMsg'), 'i')
      this.setContent(localContent)
    }

    this.$textarea.addEventListener('input', () => (this.saveToLocalStorage()))
  }

  private saveToLocalStorage() {
    window.localStorage.setItem('ArtalkContent', this.getContentOriginal().trim())
  }

  private initHeader() {
    Object.entries(this.$inputs).forEach(([key, $input]) => {
      $input.value = this.user.data[key] || ''
      $input.addEventListener('input', () => this.onHeaderInput(key, $input))

      // 设置 Placeholder
      $input.placeholder = `${this.$t(key as any)}`
    })
  }

  private onHeaderInput(key: string, $input: HTMLInputElement) {
    this.user.update({
      [key]: $input.value.trim()
    })

    // 插件监听事件响应
    Object.entries(this.plugList).forEach(([plugName, plug]) => {
      if (plug.onHeaderInput) plug.onHeaderInput(key, $input)
    })
  }

  private initTextarea() {
    // 占位符
    this.$textarea.placeholder = this.ctx.conf.placeholder || this.$t('placeholder')

    // 修复按下 Tab 输入的内容
    this.$textarea.addEventListener('keydown', (e) => {
      const keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.insertContent('\t')
      }
    })

    // 输入框高度随内容而变化
    this.$textarea.addEventListener('input', () => {
      this.adjustTextareaHeight()
    })
  }

  private initSubmitBtn() {
    this.$submitBtn.innerText = this.ctx.conf.sendBtn || this.$t('send')
    this.$submitBtn.addEventListener('click', () => (this.submit()))
  }

  /** 最终用于 submit 的数据 */
  public getFinalContent() {
    let content = this.getContentOriginal()

    // 表情包处理
    if (this.plugList.emoticons) {
      content = (this.plugList.emoticons as EmoticonsPlug).transEmoticonImageText(content)
    }

    return content
  }

  public getContentOriginal() {
    return this.$textarea.value || '' // Tip: !!"0" === true
  }

  public getContentMarked() {
    return Utils.marked(this.ctx, this.getFinalContent())
  }

  public setContent(val: string) {
    this.$textarea.value = val
    this.saveToLocalStorage()
    if (this.plugList.preview) {
      ;(this.plugList.preview as PreviewPlug).updateContent()
    }

    // 延迟执行防止无效
    window.setTimeout(() => {
      this.adjustTextareaHeight()
    }, 80)
  }

  public insertContent(val: string) {
    if ((document as any).selection) {
      this.$textarea.focus();
      (document as any).selection.createRange().text = val
      this.$textarea.focus()
    } else if (this.$textarea.selectionStart || this.$textarea.selectionStart === 0) {
      const sStart = this.$textarea.selectionStart
      const sEnd = this.$textarea.selectionEnd
      const sT = this.$textarea.scrollTop
      this.setContent(this.$textarea.value.substring(0, sStart) + val + this.$textarea.value.substring(sEnd, this.$textarea.value.length))
      this.$textarea.focus()
      this.$textarea.selectionStart = sStart + val.length
      this.$textarea.selectionEnd = sStart + val.length
      this.$textarea.scrollTop = sT
    } else {
      this.$textarea.focus()
      this.$textarea.value += val
    }
  }

  public clearEditor() {
    this.setContent('')
    this.cancelReply()
  }

  private adjustTextareaHeight() {
    const diff = this.$textarea.offsetHeight - this.$textarea.clientHeight
    this.$textarea.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.$textarea.style.height = `${this.$textarea.scrollHeight + diff}px`
  }

  public setReply(commentData: CommentData, $comment: HTMLElement, scroll = true) {
    if (this.replyComment !== null) {
      this.cancelReply()
    }

    if (this.$sendReply === null) {
      this.$sendReply = Utils.createElement(`<div class="atk-send-reply">${this.$t('reply')} <span class="atk-text"></span><span class="atk-cancel">×</span></div>`);
      this.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      this.$sendReply.addEventListener('click', () => {
        this.cancelReply()
      })
      this.$textareaWrap.append(this.$sendReply)
    }
    this.replyComment = commentData

    if (this.ctx.conf.editorTravel === true) {
      this.travel($comment)
    }

    if (scroll) Ui.scrollIntoView(this.$el)

    this.$textarea.focus()
  }

  public cancelReply() {
    if (this.$sendReply !== null) {
      this.$sendReply.remove()
      this.$sendReply = null
    }
    this.replyComment = null

    if (this.ctx.conf.editorTravel === true) {
      this.travelBack()
    }
  }

  public showNotify(msg: string, type: "i"|"s"|"w"|"e") {
    Ui.showNotify(this.$notifyWrap, msg, type)
  }

  /** 提交评论 */
  public async submit() {
    if (this.getFinalContent().trim() === '') {
      this.$textarea.focus()
      return
    }

    this.ctx.trigger('editor-submit')

    Ui.showLoading(this.$el)

    try {
      const nComment = await new Api(this.ctx).add({
        content: this.getFinalContent(),
        nick: this.user.data.nick,
        email: this.user.data.email,
        link: this.user.data.link,
        rid: (this.replyComment === null) ? 0 : this.replyComment.id,
        page_key: (this.replyComment === null) ? this.ctx.conf.pageKey : this.replyComment.page_key,
        page_title: (this.replyComment === null) ? this.ctx.conf.pageTitle : undefined,
        site_name: (this.replyComment === null) ? this.ctx.conf.site : this.replyComment.site_name
      })

      // 回复不同页面的评论
      if (this.replyComment !== null && this.replyComment.page_key !== this.ctx.conf.pageKey) {
        window.open(`${this.replyComment.page_url}#atk-comment-${nComment.id}`)
      }

      this.ctx.trigger('list-insert', nComment)
      this.clearEditor() // 清空编辑器
      this.ctx.trigger('editor-submitted')
    } catch (err: any) {
      console.error(err)
      this.showNotify(`${this.$t('commentFail')}，${err.msg || String(err)}`, 'e')
      return
    } finally {
      Ui.hideLoading(this.$el)
    }
  }

  /** 关闭评论 */
  public close() {
    if (!this.$textareaWrap.querySelector('.atk-comment-closed'))
      this.$textareaWrap.prepend(Utils.createElement('<div class="atk-comment-closed">仅管理员可评论</div>'))

    if (!this.user.data.isAdmin) {
      this.$textarea.style.display = 'none'
      this.closePlugPanel()
      this.$bottom.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.$textarea.style.display = ''
      this.$bottom.style.display = ''
    }
  }

  /** 打开评论 */
  public open() {
    this.$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.$textarea.style.display = ''
    this.$bottom.style.display = ''
  }

  /** 移动评论框到置顶元素之后 */
  public travel($afterEl: HTMLElement) {
    if (this.isTraveling) return
    this.isTraveling = true
    this.$el.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

    const $travelPlace = Utils.createElement('<div></div>')
    $afterEl.after($travelPlace)
    $travelPlace.replaceWith(this.$el)

    this.$el.classList.add('atk-fade-in') // 添加渐入动画
  }

  /** 评论框归位 */
  public travelBack() {
    if (!this.isTraveling) return
    this.isTraveling = false
    this.ctx.$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(this.$el)

    // 取消回复
    if (this.replyComment !== null) this.cancelReply()
  }

  /** 插件初始化 */
  private initPlugs() {
    this.plugList = {}
    this.$plugPanelWrap.innerHTML = ''
    this.$plugPanelWrap.style.display = 'none'
    this.openedPlugName = null
    this.$plugBtnWrap.innerHTML = ''

    const disabledPlugs: string[] = []
    if (!this.conf.emoticons) disabledPlugs.push('emoticons')

    // 初始化 Editor 插件
    this.ENABLED_PLUGS.forEach((Plug) => {
      const plugName = Plug.Name

      // 禁用的插件
      if (disabledPlugs.includes(plugName)) return

      // 插件对象实例化
      const plug = new Plug(this)
      this.plugList[plugName] = plug

      // 插件按钮
      const $btn = plug.getBtn()
      if ($btn) {
        this.$plugBtnWrap.appendChild($btn)
        $btn.onclick = $btn.onclick || (() => {
          // 其他按钮去除 Active
          this.$plugBtnWrap.querySelectorAll('.active').forEach(item => item.classList.remove('active'))

          // 若点击已打开的，则关闭打开的面板
          if (plugName === this.openedPlugName) {
            this.closePlugPanel()
            return
          }

          this.openPlugPanel(plugName)

          // 当前按钮添加 Active
          $btn.classList.add('active')
        })

        // 插件面板初始化
        const $panel = plug.getPanel()
        if ($panel) {
          $panel.setAttribute('data-plug-name', plugName)
          $panel.style.display = 'none'
          this.$plugPanelWrap.appendChild($panel)
        }
      }
    })
  }

  /** 展开插件面板 */
  public openPlugPanel(plugName: string) {
    Object.entries(this.plugList).forEach(([aPlugName, plug]) => {
      const plugPanel = plug.getPanel()
      if (!plugPanel) return

      if (aPlugName === plugName) {
        plugPanel.style.display = ''
        if (plug.onPanelShow) plug.onPanelShow()
      } else {
        plugPanel.style.display = 'none'
        if (plug.onPanelHide) plug.onPanelHide()
      }
    })

    this.$plugPanelWrap.style.display = ''
    this.openedPlugName = plugName
  }

  /** 收起插件面板 */
  public closePlugPanel() {
    if (!this.openedPlugName) return

    const plug = this.plugList[this.openedPlugName]
    if (!plug) return

    if (plug.onPanelHide) plug.onPanelHide()

    this.$plugPanelWrap.style.display = 'none'
    this.openedPlugName = null
  }
}

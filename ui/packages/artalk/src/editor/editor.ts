import { CommentData } from '~/types/artalk-data'
import Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import { render, EditorUI } from './editor-ui'
import marked from '../lib/marked'

import EmoticonsPlug from './plugs/emoticons-plug'
import UploadPlug from './plugs/upload-plug'
import PreviewPlug from './plugs/preview-plug'
import EditorPlug from './plugs/editor-plug'
import HeaderInputPlug from './plugs/header-input-plug'

export default class Editor extends Component {
  private get user() { return this.ctx.user }

  private ui: EditorUI

  public get $inputs() {
    return { nick: this.ui.$nick, email: this.ui.$email, link: this.ui.$link }
  }

  // Editor 的不同状态
  private isTraveling = false

  /** 回复评论 */
  private replyComment: CommentData|null = null
  public get isReplyMode() { return this.replyComment !== null }

  /** 编辑评论 */
  private editComment: CommentData|null = null
  public get isEditMode() { return this.editComment !== null }

  /** 启用的插件 */
  private readonly ENABLED_PLUGS = [ EmoticonsPlug, UploadPlug, PreviewPlug, HeaderInputPlug ]
  public plugList: { [name: string]: EditorPlug } = {}
  private openedPlugName: string|null = null

  public constructor(ctx: Context) {
    super(ctx)

    this.ui = render()
    this.$el = this.ui.$el

    // 监听事件
    this.ctx.on('conf-loaded', () => {
      // 执行初始化
      this.initLocalStorage()
      this.initHeader()
      this.initTextarea()
      this.initSubmitBtn()
      this.initPlugs()
    })
  }

  public getUI() {
    return this.ui
  }

  private initLocalStorage() {
    const localContent = window.localStorage.getItem('ArtalkContent') || ''
    if (localContent.trim() !== '') {
      this.showNotify(this.$t('restoredMsg'), 'i')
      this.setContent(localContent)
    }

    this.ui.$textarea.addEventListener('input', () => (this.saveToLocalStorage()))
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
    if (this.isEditMode) return // 评论编辑模式，不修改个人信息

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
    this.ui.$textarea.placeholder = this.ctx.conf.placeholder || this.$t('placeholder')

    // 修复按下 Tab 输入的内容
    this.ui.$textarea.addEventListener('keydown', (e) => {
      const keyCode = e.keyCode || e.which

      if (keyCode === 9) {
        e.preventDefault()
        this.insertContent('\t')
      }
    })

    // 输入框高度随内容而变化
    this.ui.$textarea.addEventListener('input', () => {
      this.adjustTextareaHeight()
    })
  }

  private refreshSendBtnText() {
    if (this.isEditMode) this.ui.$submitBtn.innerText = this.$t('save')
    else this.ui.$submitBtn.innerText = this.ctx.conf.sendBtn || this.$t('send')
  }

  private initSubmitBtn() {
    this.refreshSendBtnText()
    this.ui.$submitBtn.addEventListener('click', () => (this.submit()))
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
    return this.ui.$textarea.value || '' // Tip: !!"0" === true
  }

  public getContentMarked() {
    return marked(this.ctx, this.getFinalContent())
  }

  public setContent(val: string) {
    this.ui.$textarea.value = val
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
      this.ui.$textarea.focus();
      (document as any).selection.createRange().text = val
      this.ui.$textarea.focus()
    } else if (this.ui.$textarea.selectionStart || this.ui.$textarea.selectionStart === 0) {
      const sStart = this.ui.$textarea.selectionStart
      const sEnd = this.ui.$textarea.selectionEnd
      const sT = this.ui.$textarea.scrollTop
      this.setContent(this.ui.$textarea.value.substring(0, sStart) + val + this.ui.$textarea.value.substring(sEnd, this.ui.$textarea.value.length))
      this.ui.$textarea.focus()
      this.ui.$textarea.selectionStart = sStart + val.length
      this.ui.$textarea.selectionEnd = sStart + val.length
      this.ui.$textarea.scrollTop = sT
    } else {
      this.ui.$textarea.focus()
      this.ui.$textarea.value += val
    }
  }

  private adjustTextareaHeight() {
    const diff = this.ui.$textarea.offsetHeight - this.ui.$textarea.clientHeight
    this.ui.$textarea.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.ui.$textarea.style.height = `${this.ui.$textarea.scrollHeight + diff}px`
  }

  /** 设置回复评论 */
  public setReply(commentData: CommentData, $comment: HTMLElement, scroll = true) {
    if (this.isEditMode) this.cancelEditComment()
    if (this.isReplyMode) this.cancelReply()

    if (!this.ui.$sendReply) {
      this.ui.$sendReply = Utils.createElement(`<div class="atk-send-reply">${this.$t('reply')} <span class="atk-text"></span><span class="atk-cancel">×</span></div>`);
      this.ui.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      this.ui.$sendReply.addEventListener('click', () => {
        this.cancelReply()
      })
      this.ui.$textareaWrap.append(this.ui.$sendReply)
    }
    this.replyComment = commentData

    if (this.ctx.conf.editorTravel === true) {
      this.travel($comment)
    }

    if (scroll) Ui.scrollIntoView(this.ui.$el)

    this.ui.$textarea.focus()
  }

  /** 取消回复评论 */
  public cancelReply() {
    if (this.ui.$sendReply) {
      this.ui.$sendReply.remove()
      this.ui.$sendReply = undefined
    }
    this.replyComment = null

    if (this.ctx.conf.editorTravel === true) {
      this.travelBack()
    }
  }

  /** 设置编辑评论 */
  public setEditComment(commentData: CommentData, $comment: HTMLElement) {
    if (this.isEditMode) this.cancelEditComment()
    if (this.isReplyMode) this.cancelReply()

    if (this.ui.$editCancelBtn === null) {
      this.ui.$editCancelBtn = Utils.createElement(`<div class="atk-send-reply">${this.$t('editCancel')}<span class="atk-cancel">×</span></div>`);
      this.ui.$editCancelBtn.onclick = () => {
        this.cancelEditComment()
      }
      this.ui.$textareaWrap.append(this.ui.$editCancelBtn)
    }
    this.editComment = commentData

    this.ui.$header.style.display = 'none' // TODO 暂时隐藏

    this.travel($comment)

    this.ui.$nick.value = commentData.nick
    this.ui.$email.value = commentData.email || ''
    this.ui.$link.value = commentData.link

    this.setContent(commentData.content)
    this.ui.$textarea.focus()

    this.refreshSendBtnText()
  }

  /** 取消编辑评论 */
  public cancelEditComment() {
    if (this.ui.$editCancelBtn) {
      this.ui.$editCancelBtn.remove()
      this.ui.$editCancelBtn = undefined
    }

    this.editComment = null
    this.travelBack()

    this.ui.$nick.value = this.user.data.nick
    this.ui.$email.value = this.user.data.email
    this.ui.$link.value = this.user.data.link

    this.setContent('')
    this.refreshSendBtnText()
    this.ui.$header.style.display = '' // TODO
  }

  public showNotify(msg: string, type: "i"|"s"|"w"|"e") {
    Ui.showNotify(this.ui.$notifyWrap, msg, type)
  }

  public showLoading() {
    Ui.showLoading(this.ui.$el)
  }

  public hideLoading() {
    Ui.hideLoading(this.ui.$el)
  }

  /** 点击评论提交按钮事件 */
  public async submit() {
    if (this.getFinalContent().trim() === '') {
      this.ui.$textarea.focus()
      return
    }

    this.ctx.trigger('editor-submit')

    this.showLoading()

    if (!this.isEditMode) {
      await this.submitAdd()
    } else {
      await this.submitEdit()
    }

    this.hideLoading()
  }

  /** 创建评论 */
  public async submitAdd() {
    try {
      const nComment = await this.ctx.getApi().comment.add({
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

      this.ctx.insertComment(nComment)

      // 复原编辑器
      this.setContent('')
      this.cancelReply()

      this.ctx.trigger('editor-submitted')
    } catch (err: any) {
      console.error(err)
      this.showNotify(`${this.$t('commentFail')}，${err.msg || String(err)}`, 'e')
    }
  }

  /** 修改评论 */
  public async submitEdit() {
    try {
      const saveData = {
        content: this.getFinalContent(),
        nick: this.ui.$nick.value,
        email: this.ui.$email.value,
        link: this.ui.$link.value,
      }
      const nComment = await this.ctx.getApi().comment.commentEdit({
        ...this.editComment, ...saveData
      })
      this.ctx.updateComment(nComment)

      // 复原编辑器
      this.setContent('')
      this.cancelEditComment()

      this.ctx.trigger('editor-submitted')
    } catch (err: any) {
      console.error(err)
      this.showNotify(`${this.$t('commentFail')}，${err.msg || String(err)}`, 'e')
    }
  }

  /** 关闭评论 */
  public close() {
    if (!this.ui.$textareaWrap.querySelector('.atk-comment-closed'))
      this.ui.$textareaWrap.prepend(Utils.createElement(`<div class="atk-comment-closed">${this.$t('onlyAdminCanReply')}</div>`))

    if (!this.user.data.isAdmin) {
      this.ui.$textarea.style.display = 'none'
      this.closePlugPanel()
      this.ui.$bottom.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.ui.$textarea.style.display = ''
      this.ui.$bottom.style.display = ''
    }
  }

  /** 打开评论 */
  public open() {
    this.ui.$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.ui.$textarea.style.display = ''
    this.ui.$bottom.style.display = ''
  }

  /** 移动评论框到置顶元素之后 */
  public travel($afterEl: HTMLElement) {
    if (this.isTraveling) return
    this.isTraveling = true
    this.ui.$el.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

    const $travelPlace = Utils.createElement('<div></div>')
    $afterEl.after($travelPlace)
    $travelPlace.replaceWith(this.ui.$el)

    this.ui.$el.classList.add('atk-fade-in') // 添加渐入动画
  }

  /** 评论框归位 */
  public travelBack() {
    if (!this.isTraveling) return
    this.isTraveling = false
    this.ctx.$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(this.ui.$el)

    // 取消回复
    if (this.replyComment !== null) this.cancelReply()
  }

  /** 插件初始化 */
  private initPlugs() {
    this.plugList = {}
    this.ui.$plugPanelWrap.innerHTML = ''
    this.ui.$plugPanelWrap.style.display = 'none'
    this.openedPlugName = null
    this.ui.$plugBtnWrap.innerHTML = ''

    const disabledPlugs: string[] = []
    if (!this.conf.imgUpload) disabledPlugs.push('upload')
    if (!this.conf.emoticons) disabledPlugs.push('emoticons')
    if (!this.conf.preview) disabledPlugs.push('preview')

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
        this.ui.$plugBtnWrap.appendChild($btn)
        $btn.onclick = $btn.onclick || (() => {
          // 其他按钮去除 Active
          this.ui.$plugBtnWrap.querySelectorAll('.active').forEach(item => item.classList.remove('active'))

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
          this.ui.$plugPanelWrap.appendChild($panel)
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

    this.ui.$plugPanelWrap.style.display = ''
    this.openedPlugName = plugName
  }

  /** 收起插件面板 */
  public closePlugPanel() {
    if (!this.openedPlugName) return

    const plug = this.plugList[this.openedPlugName]
    if (!plug) return

    if (plug.onPanelHide) plug.onPanelHide()

    this.ui.$plugPanelWrap.style.display = 'none'
    this.openedPlugName = null
  }
}

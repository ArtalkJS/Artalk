import type { CommentData } from '~/types/artalk-data'
import type Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import marked from '../lib/marked'
import User from '../lib/user'
import initEditorFuncs from './funcs'
import { render, EditorUI } from './ui'
import { PlugManager } from './plug-manager'
import { Mover } from './mover'
import { ReplyManager } from './reply'
import { EditModeManager } from './edit-mode'
import { SubmitManager } from './submit'

export default class Editor extends Component {
  /** 界面 */
  private ui: EditorUI
  public getUI() { return this.ui }

  /** 插件管理器 */
  private plugs?: PlugManager
  public getPlugs() { return this.plugs }
  public setPlugs(p: PlugManager) { this.plugs = p }

  /** 评论框移动 */
  private mover?: Mover
  public setMover(m: Mover) { this.mover = m }

  /** 回复评论 */
  private reply?: ReplyManager
  public setReplyManager(m: ReplyManager) { this.reply = m }
  public getReplyManager() { return this.reply }
  public get isReplyMode() { return !!this.reply?.comment }

  /** 编辑评论 */
  private editMode?: EditModeManager
  public setEditModeManager(m: EditModeManager) { this.editMode = m }
  public get isEditMode() { return !!this.editMode?.comment }

  /** 提交评论 */
  private submitManager?: SubmitManager
  public setSubmitManager(m: SubmitManager) { this.submitManager = m }
  public getSubmitManager() { return this.submitManager }

  /** 已加载功能的 unmount 函数 */
  private unmountFuncs: (() => void)[] = []

  public constructor(ctx: Context) {
    super(ctx)

    // init editor ui
    this.ui = render()
    this.$el = this.ui.$el

    // event listen
    this.ctx.on('conf-loaded', () => {
      // call unmount funcs
      // (this will only be called while conf reloaded, not be called at first time)
      this.unmountFuncs?.forEach((unmount) => !!unmount && unmount())

      // call init will return unmount funcs and save it
      this.unmountFuncs = initEditorFuncs(this) // init editor funcs
    })
  }

  public getHeaderInputEls() {
    return { nick: this.ui.$nick, email: this.ui.$email, link: this.ui.$link }
  }

  public saveToLocalStorage() {
    window.localStorage.setItem('ArtalkContent', this.getContentOriginal().trim())
  }

  public refreshSendBtnText() {
    if (this.isEditMode) this.ui.$submitBtn.innerText = this.$t('save')
    else this.ui.$submitBtn.innerText = this.ctx.conf.sendBtn || this.$t('send')
  }

  /** 最终用于 submit 的数据 */
  public getFinalContent() {
    let content = this.getContentOriginal()

    // plug hook: final content transformer
    if (this.plugs) content = this.plugs.getTransformedContent(content)

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

    // plug hook: content updated
    if (this.plugs) this.plugs.triggerContentUpdatedEvt(val)

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

  public adjustTextareaHeight() {
    const diff = this.ui.$textarea.offsetHeight - this.ui.$textarea.clientHeight
    this.ui.$textarea.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.ui.$textarea.style.height = `${this.ui.$textarea.scrollHeight + diff}px`
  }

  public focus() {
    this.ui.$textarea.focus()
  }

  public reset() {
    this.setContent('')
    this.cancelReply()
    this.cancelEditComment()
  }

  /** 设置回复评论 */
  public setReply(commentData: CommentData, $comment: HTMLElement, scroll = true) {
    this.reply?.setReply(commentData, $comment, scroll)
  }

  /** 取消回复评论 */
  public cancelReply() {
    this.reply?.cancelReply()
  }

  /** 设置编辑评论 */
  public setEditComment(commentData: CommentData, $comment: HTMLElement) {
    this.editMode?.setEdit(commentData, $comment)
  }

  /** 取消编辑评论 */
  public cancelEditComment() {
    this.editMode?.cancelEdit()
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
    if (!this.submitManager) throw Error('submitManger not initialized')
    await this.submitManager.do()
  }

  /** 关闭评论 */
  public close() {
    if (!this.ui.$textareaWrap.querySelector('.atk-comment-closed'))
      this.ui.$textareaWrap.prepend(Utils.createElement(`<div class="atk-comment-closed">${this.$t('onlyAdminCanReply')}</div>`))

    if (!User.data.isAdmin) {
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
    this.mover?.move($afterEl)
  }

  /** 评论框归位 */
  public travelBack() {
    this.mover?.back()
  }

  /** 展开插件面板 */
  public openPlugPanel(plugName: string) {
    this.plugs?.openPlugPanel(plugName)
  }

  /** 收起插件面板 */
  public closePlugPanel() {
    this.plugs?.closePlugPanel()
  }
}

import type { CommentData, ContextApi, EditorApi } from '@/types'
import Component from '../lib/component'
import * as Ui from '../lib/ui'
import marked from '../lib/marked'
import { render, EditorUI } from './ui'
import EditorStateManager from './state'

class Editor extends Component implements EditorApi {
  private ui: EditorUI
  private state: EditorStateManager

  getUI() {
    return this.ui
  }
  getPlugs() {
    return this.ctx.get('editorPlugs')
  }
  getState() {
    return this.state.get()
  }

  constructor(ctx: ContextApi) {
    super(ctx)

    // init editor ui
    this.ui = render()
    this.$el = this.ui.$el

    // init state manager
    this.state = new EditorStateManager(this)
  }

  getHeaderInputEls() {
    return { nick: this.ui.$nick, email: this.ui.$email, link: this.ui.$link }
  }

  getContentFinal() {
    let content = this.getContentRaw()

    // plug hook: final content transformer
    const plugs = this.getPlugs()
    if (plugs) content = plugs.getTransformedContent(content)

    return content
  }

  getContentRaw() {
    return this.ui.$textarea.value || ''
  }

  getContentMarked() {
    return marked(this.getContentFinal())
  }

  setContent(val: string) {
    this.ui.$textarea.value = val

    // plug hook: content updated
    this.getPlugs()?.getEvents().trigger('content-updated', val)
  }

  insertContent(val: string) {
    if ((document as any).selection) {
      this.ui.$textarea.focus()
      ;(document as any).selection.createRange().text = val
      this.ui.$textarea.focus()
    } else if (this.ui.$textarea.selectionStart || this.ui.$textarea.selectionStart === 0) {
      const sStart = this.ui.$textarea.selectionStart
      const sEnd = this.ui.$textarea.selectionEnd
      const sT = this.ui.$textarea.scrollTop
      this.setContent(
        this.ui.$textarea.value.substring(0, sStart) +
          val +
          this.ui.$textarea.value.substring(sEnd, this.ui.$textarea.value.length),
      )
      this.ui.$textarea.focus()
      this.ui.$textarea.selectionStart = sStart + val.length
      this.ui.$textarea.selectionEnd = sStart + val.length
      this.ui.$textarea.scrollTop = sT
    } else {
      this.ui.$textarea.focus()
      this.ui.$textarea.value += val
    }
  }

  focus() {
    this.ui.$textarea.focus()
  }

  reset() {
    this.setContent('')
    this.resetState()
  }

  resetState() {
    this.state.switch('normal')
  }

  setReply(comment: CommentData, $comment: HTMLElement) {
    this.state.switch('reply', { comment, $comment })
  }

  setEditComment(comment: CommentData, $comment: HTMLElement) {
    this.state.switch('edit', { comment, $comment })
  }

  showNotify(msg: string, type: any) {
    Ui.showNotify(this.ui.$notifyWrap, msg, type)
  }

  showLoading() {
    Ui.showLoading(this.ui.$el)
  }

  hideLoading() {
    Ui.hideLoading(this.ui.$el)
  }

  submit() {
    const next = () => this.ctx.trigger('editor-submit')
    if (this.ctx.conf.beforeSubmit) {
      this.ctx.conf.beforeSubmit(this, next)
    } else {
      next()
    }
  }
}

export default Editor

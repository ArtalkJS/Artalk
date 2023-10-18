import type { CommentData } from '~/types/artalk-data'
import type { EditorApi, EditorState } from '~/types/editor'
import type Context from '~/types/context'
import Component from '../lib/component'
import * as Ui from '../lib/ui'
import marked from '../lib/marked'
import { render, EditorUI } from './ui'
import EditorStateManager from './state'

class Editor extends Component implements EditorApi {
  private ui: EditorUI
  private state: EditorStateManager

  getUI() { return this.ui }
  getPlugs() { return this.ctx.get('editorPlugs') }

  constructor(ctx: Context) {
    super(ctx)

    // init editor ui
    this.ui = render()
    this.$el = this.ui.$el

    // init state manager
    this.state = new EditorStateManager(this)

    let confLoaded = false

    // event listen
    this.ctx.on('conf-loaded', () => {
      // trigger unmount event will call all plugs' unmount function
      // (this will only be called while conf reloaded, not be called at first time)
      confLoaded && this.getPlugs()?.getEvents().trigger('unmounted')

      // trigger event for plug initialization
      this.getPlugs()?.getEvents().trigger('mounted')

      confLoaded = true
    })
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
    return marked(this.ctx, this.getContentFinal())
  }

  setContent(val: string) {
    this.ui.$textarea.value = val

    // plug hook: content updated
    this.getPlugs()?.getEvents().trigger('content-updated', val)
  }

  insertContent(val: string) {
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
    this.ctx.trigger('editor-submit')
  }
}

export default Editor

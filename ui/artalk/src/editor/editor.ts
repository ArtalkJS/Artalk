import * as Ui from '../lib/ui'
import marked from '../lib/marked'
import { render, type EditorUI } from './ui'
import { EditorStateManager } from './state'
import type { ConfigManager, CommentData, Editor as IEditor, EventManager } from '@/types'
import type { PluginManager } from '@/plugins/editor-kit'

export interface EditorOptions {
  getEvents: () => EventManager
  getConf: () => ConfigManager
}

export class Editor implements IEditor {
  private opts: EditorOptions
  private $el: HTMLElement
  private ui: EditorUI
  private state: EditorStateManager
  private plugins?: PluginManager

  constructor(opts: EditorOptions) {
    this.opts = opts

    // init editor ui
    this.ui = render()
    this.$el = this.ui.$el

    // init state manager
    this.state = new EditorStateManager(this)
  }

  getOptions() {
    return this.opts
  }

  getEl() {
    return this.$el
  }

  getUI() {
    return this.ui
  }

  getPlugins() {
    return this.plugins
  }

  setPlugins(plugins: PluginManager) {
    this.plugins = plugins
  }

  getState() {
    return this.state.get()
  }

  getHeaderInputEls() {
    return { name: this.ui.$name, email: this.ui.$email, link: this.ui.$link }
  }

  getContentFinal() {
    let content = this.getContentRaw()

    // plugin hook: final content transformer
    const plugins = this.getPlugins()
    if (plugins) content = plugins.getTransformedContent(content)

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
    this.getPlugins()?.getEvents().trigger('content-updated', val)
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

  setReplyComment(comment: CommentData, $comment: HTMLElement) {
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
    const next = () => {
      this.getPlugins()?.getEvents().trigger('editor-submit')
      this.opts.getEvents().trigger('editor-submit')
    }
    const beforeSubmit = this.opts.getConf().get().beforeSubmit
    if (beforeSubmit) {
      beforeSubmit(this, next)
    } else {
      next()
    }
  }
}

import './preview-plug.scss'

import Editor from '../editor'
import EditorPlug from './editor-plug'

export default class PreviewPlug extends EditorPlug {
  public static Name = 'preview'
  declare protected $panel: HTMLElement

  private isBind = false

  public constructor(editor: Editor) {
    super(editor)

    this.registerPanel(`<div class="atk-editor-plug-preview"></div>`)

    let btnText = this.editor.$t('preview')
    if (this.ctx.getMarkedInstance()) btnText += ` <i title="Markdown is supported"><svg class="markdown" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M14.85 3H1.15C.52 3 0 3.52 0 4.15v7.69C0 12.48.52 13 1.15 13h13.69c.64 0 1.15-.52 1.15-1.15v-7.7C16 3.52 15.48 3 14.85 3zM9 11H7V8L5.5 9.92 4 8v3H2V5h2l1.5 2L7 5h2v6zm2.99.5L9.5 8H11V5h2v3h1.5l-2.51 3.5z"></path></svg></i>`
    this.registerBtn(btnText)

    this.registerContentUpdatedEvt((content) => {
      this.updateContent()
    })
  }

  public onPanelShow() {
    this.updateContent()

    if (!this.isBind) {
      const event = () => {
        this.updateContent()
      }
      this.editor.getUI().$textarea.addEventListener('input', event)
      this.editor.getUI().$textarea.addEventListener('change', event)
      this.isBind = true
    }
  }

  public updateContent() {
    if (this.$panel.style.display !== 'none') {
      this.$panel.innerHTML = this.editor.getContentMarked()
    }
  }
}

import Editor from '../editor'
import EditorPlug from '../editor-plug'

export default class TextareaPlug extends EditorPlug {
  constructor(editor: Editor) {
    super(editor)

    const onKeydown = (e: KeyboardEvent) => this.onKeydown(e)
    const onInput = () => this.onInput()

    this.kit.useMounted(() => {
      // 占位符
      editor.getUI().$textarea.placeholder = editor.ctx.conf.placeholder || editor.$t('placeholder')

      // bind the event
      editor.getUI().$textarea.addEventListener('keydown', onKeydown)
      editor.getUI().$textarea.addEventListener('input', onInput)
    })

    this.kit.useUnmounted(() => {
      // unmount the event
      editor.getUI().$textarea.removeEventListener('keydown', onKeydown)
      editor.getUI().$textarea.removeEventListener('input', onInput)
    })

    this.kit.useContentUpdated(() => {
      // delay 80ms to prevent invalid execution
      window.setTimeout(() => {
        this.adaptiveHeightByContent()
      }, 80)
    })
  }

  // 按下 Tab 输入内容，而不是失去焦距
  private onKeydown(e: KeyboardEvent) {
    const keyCode = e.keyCode || e.which

    if (keyCode === 9) {
      e.preventDefault()
      this.editor.insertContent('\t')
    }
  }

  private onInput() {
    this.editor.getPlugs()?.triggerContentUpdatedEvt(this.editor.getContentRaw())
  }

  // Resize the textarea height by content
  public adaptiveHeightByContent() {
    const diff = this.editor.getUI().$textarea.offsetHeight - this.editor.getUI().$textarea.clientHeight
    this.editor.getUI().$textarea.style.height = '0px' // it's a magic. 若不加此行，内容减少，高度回不去
    this.editor.getUI().$textarea.style.height = `${this.editor.getUI().$textarea.scrollHeight + diff}px`
  }
}

import Editor from '../editor'
import EditorPlug from '../editor-plug'
import EditPlug from './edit-plug'

export default class SubmitBtnPlug extends EditorPlug {
  constructor(editor: Editor) {
    super(editor)

    const onClick = () => {
      this.editor.submit()
    }

    this.kit.useMounted(() => {
      // apply the submit button text from user custom config
      this.editor.getUI().$submitBtn.innerText = this.editor.ctx.conf.sendBtn || this.editor.$t('send')

      // bind the event when click the submit button
      editor.getUI().$submitBtn.addEventListener('click', onClick)
    })

    this.kit.useUnmounted(() => {
      editor.getUI().$submitBtn.removeEventListener('click', onClick)
    })
  }
}

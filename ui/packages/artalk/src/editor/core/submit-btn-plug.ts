import $t from '@/i18n'
import EditorPlug from '../editor-plug'
import PlugKit from '../plug-kit'

export default class SubmitBtnPlug extends EditorPlug {
  constructor(kit: PlugKit) {
    super(kit)

    const onClick = () => {
      this.kit.useEditor().submit()
    }

    this.kit.useMounted(() => {
      // apply the submit button text from user custom config
      this.kit.useUI().$submitBtn.innerText = this.kit.useConf().sendBtn || $t('send')

      // bind the event when click the submit button
      this.kit.useUI().$submitBtn.addEventListener('click', onClick)
    })

    this.kit.useUnmounted(() => {
      this.kit.useUI().$submitBtn.removeEventListener('click', onClick)
    })
  }
}

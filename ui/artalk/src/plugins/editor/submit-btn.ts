import EditorPlugin from './_plug'
import type PlugKit from './_kit'
import $t from '@/i18n'

export default class SubmitBtn extends EditorPlugin {
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

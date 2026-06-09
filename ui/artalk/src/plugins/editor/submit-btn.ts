import EditorPlugin from './_plug'
import type PlugKit from './_kit'

const PAPER_PLANE_ICON =
  `<svg class="atk-send-icon" aria-hidden="true" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">` +
  `<path fill="currentColor" d="M565.747623 792.837176l260.819261 112.921839 126.910435-845.424882L66.087673 581.973678l232.843092 109.933785 562.612725-511.653099-451.697589 563.616588-5.996574 239.832274L565.747623 792.837176z"/>` +
  `</svg>`

export default class SubmitBtn extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onClick = () => {
      this.kit.useEditor().submit()
    }

    this.kit.useMounted(() => {
      // render paper-plane icon (custom sendBtn config still shows text)
      const sendBtn = this.kit.useConf().sendBtn
      if (sendBtn) this.kit.useUI().$submitBtn.innerText = sendBtn
      else this.kit.useUI().$submitBtn.innerHTML = PAPER_PLANE_ICON

      // bind the event when click the submit button
      this.kit.useUI().$submitBtn.addEventListener('click', onClick)
    })

    this.kit.useUnmounted(() => {
      this.kit.useUI().$submitBtn.removeEventListener('click', onClick)
    })
  }
}

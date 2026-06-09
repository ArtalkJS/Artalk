import EditorPlugin from './_plug'
import type PlugKit from './_kit'

/**
 * Editor focus state machine
 *
 * Toggles `.atk-editing` / `.atk-focus` on the input box when the textarea
 * gains or loses focus, driving the iDisqus-style expand/collapse animation.
 */
export default class Focus extends EditorPlugin {
  constructor(kit: PlugKit) {
    super(kit)

    const onFocus = () => this.onFocus()
    const onBlur = () => this.onBlur()

    this.kit.useMounted(() => {
      const $textarea = this.kit.useUI().$textarea
      $textarea.addEventListener('focus', onFocus)
      $textarea.addEventListener('blur', onBlur)
    })

    this.kit.useUnmounted(() => {
      const $textarea = this.kit.useUI().$textarea
      $textarea.removeEventListener('focus', onFocus)
      $textarea.removeEventListener('blur', onBlur)
    })
  }

  private onFocus() {
    const $box = this.kit.useUI().$inputBox
    $box.classList.add('atk-editing')
    $box.classList.add('atk-focus')
  }

  private onBlur() {
    // keep `.atk-editing` so the guest fields stay expanded once revealed,
    // only drop `.atk-focus` to remove the border highlight
    this.kit.useUI().$inputBox.classList.remove('atk-focus')
  }
}
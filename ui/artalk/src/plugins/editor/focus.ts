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
    // Prevent toolbar / floating panel clicks from stealing focus away from the
    // textarea. Without this the editor flickers out of `.atk-focus` whenever
    // the user clicks the markdown preview / emoji / upload / send buttons.
    // mousedown.preventDefault() suppresses the focus shift but still lets the
    // click event fire, so the buttons keep working. Real inputs inside a
    // panel (e.g. a search box) are skipped so they can still receive focus.
    const onToolMouseDown = (e: MouseEvent) => {
      const t = e.target as HTMLElement | null
      if (!t) return
      if (t.closest('input, textarea, [contenteditable="true"]')) return
      e.preventDefault()
    }

    this.kit.useMounted(() => {
      const ui = this.kit.useUI()
      ui.$textarea.addEventListener('focus', onFocus)
      ui.$textarea.addEventListener('blur', onBlur)
      ui.$bottom.addEventListener('mousedown', onToolMouseDown)
      ui.$plugPanelWrap.addEventListener('mousedown', onToolMouseDown)
    })

    this.kit.useUnmounted(() => {
      const ui = this.kit.useUI()
      ui.$textarea.removeEventListener('focus', onFocus)
      ui.$textarea.removeEventListener('blur', onBlur)
      ui.$bottom.removeEventListener('mousedown', onToolMouseDown)
      ui.$plugPanelWrap.removeEventListener('mousedown', onToolMouseDown)
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
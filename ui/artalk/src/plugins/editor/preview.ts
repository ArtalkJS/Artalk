import './preview.scss'

import EditorPlugin from './_plug'
import type PlugKit from './_kit'
import $t from '@/i18n'

export default class Preview extends EditorPlugin {
  private isPlugPanelShow = false

  constructor(kit: PlugKit) {
    super(kit)

    this.kit.useMounted(() => {
      this.usePanel(`<div class="atk-editor-plug-preview"></div>`)

      // initialize plug button
      this.useBtn(
        `<i aria-label="${$t('preview')}"><svg fill="currentColor" aria-hidden="true" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M14.85 3H1.15C.52 3 0 3.52 0 4.15v7.69C0 12.48.52 13 1.15 13h13.69c.64 0 1.15-.52 1.15-1.15v-7.7C16 3.52 15.48 3 14.85 3zM9 11H7V8L5.5 9.92 4 8v3H2V5h2l1.5 2L7 5h2v6zm2.99.5L9.5 8H11V5h2v3h1.5l-2.51 3.5z"></path></svg></i>`,
      )
    })
    this.kit.useUnmounted(() => {})

    // function to update content
    this.kit.useEvents().on('content-updated', (content) => {
      this.isPlugPanelShow && this.updateContent()
    })

    this.usePanelShow(() => {
      this.isPlugPanelShow = true
      this.updateContent()
    })
    this.usePanelHide(() => {
      this.isPlugPanelShow = false
    })
  }

  updateContent() {
    this.$panel!.innerHTML = this.kit.useEditor().getContentMarked()
  }
}

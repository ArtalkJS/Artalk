import * as Utils from '@/lib/utils'
import PlugKit from './plug-kit'

/**
 * Editor 插件
 */
interface EditorPlug {
  $btn?: HTMLElement
  $panel?: HTMLElement
  contentTransformer?(rawContent: string): string
}

class EditorPlug {
  constructor(
    protected kit: PlugKit
  ) {
  }

  /** Use plug btn will add a btn on the bottom of editor */
  useBtn(html: string = '<div></div>') {
    this.$btn = Utils.createElement(`<span class="atk-plug-btn">${html}</span>`)
    return this.$btn
  }

  /** Use plug panel will show the panel when btn is clicked */
  usePanel(html: string = '<div></div>') {
    this.$panel = Utils.createElement(html)
    return this.$panel
  }

  /** Use the content transformer to handle the content of the last submit by the editor */
  useContentTransformer(func: (raw: string) => string) {
    this.contentTransformer = func
  }

  /** Listen the event of panel show */
  usePanelShow(func: () => void) {
    this.kit.useEvents().on('panel-show', (aPlug) => {
      if (aPlug === this) func()
    })
  }

  /** Listen the event of panel hide */
  usePanelHide(func: () => void) {
    this.kit.useEvents().on('panel-hide', (aPlug) => {
      if (aPlug === this) func()
    })
  }
}

export default EditorPlug

import * as Utils from '@/lib/utils'
import EditorPlug from "./editor-plug"

/**
 * PlugKit provides a set of methods to help you develop editor plug
 */
export default class PlugKit {
  constructor(
    private plug: EditorPlug
  ) {
  }

  /** Use plug btn will add a btn on the bottom of editor */
  useBtn(html: string = '<div></div>') {
    this.plug.$btn = Utils.createElement(`<span class="atk-plug-btn">${html}</span>`)
    return this.plug.$btn
  }

  /** Use plug panel will show the panel when btn is clicked */
  usePanel(html: string = '<div></div>') {
    this.plug.$panel = Utils.createElement(html)
    return this.plug.$panel
  }

  /** Listen the event when plug is mounted */
  useMounted(func: () => void) {
    this.plug.onMounted = func
  }

  /** Listen the event when plug is unmounted */
  useUnmounted(func: () => void) {
    this.plug.onUnmounted = func
  }

  /** Listen the event of panel show */
  usePanelShow(func: () => void) {
    this.plug.onPanelShow = func
  }

  /** Listen the event of panel hide */
  usePanelHide(func: () => void) {
    this.plug.onPanelHide = func
  }

  /** Listen the event of header input is changed */
  useHeaderInput(func: (key: string, $input: HTMLElement) => void) {
    this.plug.onHeaderInput = func
  }

  /** Listen the event of editor content is updated */
  useContentUpdated(func: (content: string) => void) {
    this.plug.onContentUpdated = func
  }

  /** Use the content transformer to handle the content of the last submit by the editor */
  useContentTransformer(func: (raw: string) => string) {
    this.plug.contentTransformer = func
  }
}

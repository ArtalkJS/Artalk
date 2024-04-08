import type { PlugManager } from '../editor-kit'
import type EditorPlug from './_plug'

/**
 * PlugKit provides a set of methods to help you develop editor plug
 */
export default class PlugKit {
  constructor(private plugs: PlugManager) {}

  /** Use the editor */
  useEditor() {
    return this.plugs.editor
  }

  /**
   * Use the context of global
   *
   * @deprecated The calls to this function should be reduced as much as possible
   */
  useGlobalCtx() {
    return this.plugs.editor.ctx
  }

  /** Use the config of Artalk */
  useConf() {
    return this.plugs.editor.ctx.conf
  }

  /** Use the http api client */
  useApi() {
    return this.plugs.editor.ctx.getApi()
  }

  /** Use the user manager */
  useUser() {
    return this.plugs.editor.ctx.get('user')
  }

  /** Use the ui of editor */
  useUI() {
    return this.plugs.editor.getUI()
  }

  /** Use the events in editor scope */
  useEvents() {
    return this.plugs.getEvents()
  }

  /** Listen the event when plug is mounted */
  useMounted(func: () => void) {
    this.useEvents().on('mounted', func)
  }

  /** Listen the event when plug is unmounted */
  useUnmounted(func: () => void) {
    this.useEvents().on('unmounted', func)
  }

  /** Use the deps of other plug */
  useDeps<T extends typeof EditorPlug>(plug: T) {
    return this.plugs.get(plug)
  }
}

import EditorApi from '~/types/editor'
import * as Utils from '@/lib/utils'
import EditorPlug from "./editor-plug"
import PlugManager from './plug-manager'

/**
 * PlugKit provides a set of methods to help you develop editor plug
 */
export default class PlugKit {
  constructor(
    private plugs: PlugManager
  ) {
  }

  /** Use the editor */
  useEditor() {
    return this.plugs.editor
  }

  /** Use the context of global */
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

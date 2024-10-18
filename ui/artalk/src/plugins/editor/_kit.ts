import type { PluginManager } from '../editor-kit'
import type EditorPlugin from './_plug'

/**
 * PlugKit provides a set of methods to help you develop editor plug
 */
export default class PlugKit {
  constructor(private plugins: PluginManager) {}

  /** Use the editor */
  useEditor() {
    return this.plugins.getEditor()
  }

  /** Use the config of Artalk */
  useConf() {
    return this.plugins.getOptions().getConf().get()
  }

  /** Use the http api client */
  useApi() {
    return this.plugins.getOptions().getApi()
  }

  /** Use the data manager */
  useData() {
    return this.plugins.getOptions().getData()
  }

  /** Use the user manager */
  useUser() {
    return this.plugins.getOptions().getUser()
  }

  /** Use the checkers */
  useCheckers() {
    return this.plugins.getOptions().getCheckers()
  }

  /** Use the ui of editor */
  useUI() {
    return this.plugins.getEditor().getUI()
  }

  /** Use the events in editor scope */
  useEvents() {
    return this.plugins.getEvents()
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
  useDeps<T extends typeof EditorPlugin>(plug: T) {
    return this.plugins.get(plug)
  }

  /** Use the root element of artalk */
  useArtalkRootEl() {
    return this.plugins.getOptions().getArtalkRootEl()
  }
}

import type { CommentData, NotifyLevel, EventManager } from '.'
import type {
  EditorEventPayloadMap,
  PluginManager,
  PluginManagerOptions,
} from '@/plugins/editor-kit'
import type { EditorPlugin } from '@/plugins/editor/_plug'
import type { EditorOptions } from '@/editor/editor'
import type { EditorUI } from '@/editor/ui'

export type EditorState = 'reply' | 'edit' | 'normal'

export interface Editor {
  /**
   * Get editor options
   */
  getOptions(): EditorOptions

  /**
   * Get the editor element
   */
  getEl(): HTMLElement

  /**
   * Get the editor UI instance
   */
  getUI(): EditorUI

  /**
   * Get the header input elements
   */
  getHeaderInputEls(): Record<string, HTMLInputElement>

  /**
   * Set content
   */
  setContent(val: string): void

  /**
   * Insert content
   */
  insertContent(val: string): void

  /**
   * Get the final content
   *
   * This function returns the raw content or the content transformed through a plugin hook.
   */
  getContentFinal(): string

  /**
   * Get the raw content which is inputed by user
   */
  getContentRaw(): string

  /**
   * Get the HTML format content which is rendered by marked (a markdown parser)
   */
  getContentMarked(): string

  /**
   * Get editor current state
   */
  getState(): EditorState

  /**
   * Focus editor
   */
  focus(): void

  /**
   * Reset editor
   */
  reset(): void

  /**
   * Reset editor UI
   *
   * call it will move editor to the initial position
   */
  resetState(): void

  /**
   * Submit comment
   */
  submit(): void

  /**
   * Show notification message
   */
  showNotify(msg: string, type: NotifyLevel): void

  /**
   * Show loading on editor
   */
  showLoading(): void

  /**
   * Hide loading on editor
   */
  hideLoading(): void

  /**
   * Start replying a comment
   */
  setReplyComment(commentData: CommentData, $comment: HTMLElement, scroll?: boolean): void

  /**
   * Start editing a comment
   */
  setEditComment(commentData: CommentData, $comment: HTMLElement): void

  /**
   * Get plugin manager
   */
  getPlugins(): PluginManager | undefined

  /**
   * Set plugin manager
   */
  setPlugins(plugins: PluginManager): void
}

export interface EditorPluginManager {
  getPlugins: () => EditorPlugin[]
  getEvents: () => EventManager<EditorEventPayloadMap>
  getEditor: () => Editor
  getOptions: () => PluginManagerOptions
  get<T extends typeof EditorPlugin>(plug: T): InstanceType<T> | undefined
  openPluginPanel: (plug: EditorPlugin) => void
  closePluginPanel: () => void
  getTransformedContent: (rawContent: string) => string
}

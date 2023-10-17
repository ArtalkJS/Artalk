import type { CommentData } from './artalk-data'
import Component from '../src/lib/component'
import { EditorUI } from '../src/editor/ui'
import PlugManager from '../src/editor/plug-manager'

export default interface Editor extends Component {
  getUI(): EditorUI
  getPlugs(): PlugManager | undefined

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
  resetUI(): void

  /**
   * Submit comment
   */
  submit(): void

  /**
   * Close comment editor which prevent user from submitting (but admin excluded)
   */
  close(): void

  /**
   * Open comment editor which allow user to submit (only be called while editor is closed)
   */
  open(): void

  /**
   * Show notification message
   */
  showNotify(msg: string, type: "i" | "s" | "w" | "e"): void

  /**
   * Show loading on editor
   */
  showLoading(): void

  /**
   * Hide loading on editor
   */
  hideLoading(): void

  /**
   * Start replaying a comment
   */
  setReply(commentData: CommentData, $comment: HTMLElement, scroll?: boolean): void

  /**
   * Cancel replaying the comment
   */
  cancelReply(): void

  /**
   * Start editing a comment
   */
  setEditComment(commentData: CommentData, $comment: HTMLElement): void

  /**
   * Cancel editing the comment
   */
  cancelEditComment(): void
}

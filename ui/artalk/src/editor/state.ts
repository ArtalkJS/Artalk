import type { EditorState, CommentData } from '@/types'
import * as Ui from '@/lib/ui'
import type Editor from './editor'
import Mover from '../plugins/editor/mover'

export default class EditorStateManager {
  constructor(private editor: Editor) {}

  private stateCurt: EditorState = 'normal'
  private stateUnmountFn: (() => void) | null = null

  /** Get current state */
  get() {
    return this.stateCurt
  }

  /**
   * Switch editor state
   *
   * @param state The state to switch
   * @param payload The cause of state switch
   */
  switch(state: EditorState, payload?: { $comment: HTMLElement; comment: CommentData }) {
    // trigger unmount
    if (this.stateUnmountFn) {
      this.stateUnmountFn()
      this.stateUnmountFn = null

      // move editor back to the initial position
      this.editor.getPlugs()?.get(Mover)?.back()
    }

    // invoke effect function and save unmount function
    if (state !== 'normal' && payload) {
      // move editor position
      let moveAfterEl = payload.$comment
      if (!this.editor.conf.flatMode)
        moveAfterEl = moveAfterEl.querySelector<HTMLElement>('.atk-footer')!
      this.editor.getPlugs()?.get(Mover)?.move(moveAfterEl)

      const $relative =
        this.editor.ctx.conf.scrollRelativeTo && this.editor.ctx.conf.scrollRelativeTo()
      Ui.scrollIntoView(this.editor.getUI().$el, true, $relative)

      const plugin = this.editor
        .getPlugs()
        ?.getPlugs()
        .find((p) => p.editorStateEffectWhen === state)
      if (plugin && plugin.editorStateEffect) {
        this.stateUnmountFn = plugin.editorStateEffect(payload.comment)
      }
    }

    // change current state
    this.stateCurt = state
  }
}

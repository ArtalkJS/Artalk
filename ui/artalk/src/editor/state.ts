import Mover from '../plugins/editor/mover'
import type { EditorState, CommentData, Editor } from '@/types'
import * as Ui from '@/lib/ui'

export class EditorStateManager {
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
      this.editor.getPlugins()?.get(Mover)?.back()
    }

    // invoke effect function and save unmount function
    if (state !== 'normal' && payload) {
      // move editor position
      let moveAfterEl = payload.$comment
      if (!this.editor.getOptions().getConf().get().flatMode)
        moveAfterEl = moveAfterEl.querySelector<HTMLElement>('.atk-footer')!
      this.editor.getPlugins()?.get(Mover)?.move(moveAfterEl)

      const $relative = this.editor.getOptions().getConf().get().scrollRelativeTo?.()
      Ui.scrollIntoView(this.editor.getUI().$el, true, $relative)

      const plugin = this.editor
        .getPlugins()
        ?.getPlugins()
        .find((p) => p.editorStateEffectWhen === state)
      if (plugin && plugin.editorStateEffect) {
        this.stateUnmountFn = plugin.editorStateEffect(payload.comment)
      }
    }

    // change current state
    this.stateCurt = state
  }
}

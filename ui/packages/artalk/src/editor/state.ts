import type { EditorState } from '~/types/editor'
import type { CommentData } from '~/types/artalk-data'
import type Editor from './editor'
import MoverPlug from '../plugins/editor/mover'

export default class EditorStateManager {
  constructor(private editor: Editor) {}

  private stateCurt: EditorState = 'normal'
  private stateUnmountFn: (() => void)|null = null

  /**
   * Switch editor state
   *
   * @param state The state to switch
   * @param payload The cause of state switch
   */
  switch(state: EditorState, payload?: { $comment: HTMLElement, comment: CommentData }) {
    // trigger unmount
    if (this.stateUnmountFn) {
      this.stateUnmountFn()
      this.stateUnmountFn = null

      // move editor back to the initial position
      this.editor.getPlugs()?.get(MoverPlug)?.back()
    }

    // invoke effect function and save unmount function
    if (state !== 'normal' && payload) {
      // move editor position
      this.editor.getPlugs()?.get(MoverPlug)?.move(payload.$comment)

      const plugin = this.editor.getPlugs()?.getPlugs().find(p => p.editorStateEffectWhen === state)
      if (plugin && plugin.editorStateEffect) {
        this.stateUnmountFn = plugin.editorStateEffect(payload.comment)
      }
    }

    // change current state
    this.stateCurt = state
  }
}

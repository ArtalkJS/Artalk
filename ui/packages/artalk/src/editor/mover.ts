import Editor from './editor'
import * as Utils from '../lib/utils'

export interface Mover {
  editor: Editor
  isMoved: boolean
  move(afterEl: HTMLElement): void
  back(): void
}

export function createMover(editor: Editor): Mover {
  const m: Mover = {
    editor,
    isMoved: false,
    move: (el) => move(m, el),
    back: () => back(m)
  }

  return m
}

function move(m: Mover, afterEl: HTMLElement) {
  if (m.isMoved) return
  m.isMoved = true

  const editorEl = m.editor.getUI().$el

  editorEl.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

  const $travelPlace = Utils.createElement('<div></div>')
  afterEl.after($travelPlace)
  $travelPlace.replaceWith(editorEl)

  editorEl.classList.add('atk-fade-in') // 添加渐入动画
}

function back(m: Mover) {
  if (!m.isMoved) return
  m.isMoved = false
  m.editor.ctx.$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(m.editor.getUI().$el)

  m.editor.cancelReply()  // 取消回复
}

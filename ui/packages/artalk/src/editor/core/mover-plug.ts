import Editor from '../editor'
import * as Utils from '../../lib/utils'
import EditorPlug from '../editor-plug'

export default class MoverPlug extends EditorPlug {
  private isMoved = false

  constructor(editor: Editor) {
    super(editor)
  }

  move(afterEl: HTMLElement) {
    if (this.isMoved) return
    this.isMoved = true

    const editorEl = this.editor.getUI().$el

    editorEl.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

    const $travelPlace = Utils.createElement('<div></div>')
    afterEl.after($travelPlace)
    $travelPlace.replaceWith(editorEl)

    editorEl.classList.add('atk-fade-in') // 添加渐入动画
  }

  back() {
    if (!this.isMoved) return
    this.isMoved = false
    this.editor.ctx.$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(this.editor.getUI().$el)

    this.editor.cancelReply()  // 取消回复
  }
}

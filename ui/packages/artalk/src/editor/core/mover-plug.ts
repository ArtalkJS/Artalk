import * as Utils from '@/lib/utils'
import EditorPlug from '../editor-plug'
import PlugKit from '../plug-kit'

export default class MoverPlug extends EditorPlug {
  private isMoved = false

  constructor(kit: PlugKit) {
    super(kit)
  }

  move(afterEl: HTMLElement) {
    if (this.isMoved) return
    this.isMoved = true

    const editorEl = this.kit.useUI().$el

    editorEl.after(Utils.createElement('<div class="atk-editor-travel-placeholder"></div>'))

    const $travelPlace = Utils.createElement('<div></div>')
    afterEl.after($travelPlace)
    $travelPlace.replaceWith(editorEl)

    editorEl.classList.add('atk-fade-in') // 添加渐入动画
  }

  back() {
    if (!this.isMoved) return
    this.isMoved = false
    this.kit.useGlobalCtx().$root.querySelector('.atk-editor-travel-placeholder')?.replaceWith(this.kit.useUI().$el)

    this.kit.useEditor().cancelReply()  // 取消回复
  }
}

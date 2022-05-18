import * as Utils from '@/lib/utils'
import Editor from '../editor'

/**
 * Editor 插件
 *
 * @desc 使用 Interface x Abstract 合并声明
 * @see https://www.typescriptlang.org/docs/handbook/declaration-merging.html#merging-interfaces
 */
interface EditorPlug {
  onPanelShow?(): void
  onPanelHide?(): void
}

abstract class EditorPlug {
  protected editor: Editor
  protected get ctx() { return this.editor.ctx }
  protected $panel?: HTMLElement
  protected $btn?: HTMLElement
  public onHeaderInput?: (key: string, $input: HTMLElement) => void

  public constructor(editor: Editor) {
    this.editor = editor
  }

  public static Name: string

  protected registerPanel(html: string = '<div></div>') {
    this.$panel = Utils.createElement(html)
    return this.$panel
  }

  protected registerBtn(html: string) {
    this.$btn = Utils.createElement(`<span class="atk-plug-btn" data-plug-name="${this.constructor.name}">${html}</span>`)
    return this.$btn
  }

  protected registerHeaderInputEvt(action: (key: string, $input: HTMLElement) => void) {
    this.onHeaderInput = action
  }

  public getPanel() {
    return this.$panel
  }

  public getBtn() {
    return this.$btn
  }
}

export default EditorPlug

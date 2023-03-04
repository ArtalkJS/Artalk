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
  public contentTransformer?: (rawContent: string) => string
  public onContentUpdated?: (content: string) => void

  public constructor(editor: Editor) {
    this.editor = editor
  }

  /** Name of plug */
  public static Name: string

  /** Register plug panel will provide a plug el */
  protected registerPanel(html: string = '<div></div>') {
    this.$panel = Utils.createElement(html)
    return this.$panel
  }

  /** Register plug btn will add a btn on the bottom of editor */
  protected registerBtn(html: string) {
    this.$btn = Utils.createElement(`<span class="atk-plug-btn" data-plug-name="${this.constructor.name}">${html}</span>`)
    return this.$btn
  }

  /** Register the event of header input is changed */
  protected registerHeaderInputEvt(func: (key: string, $input: HTMLElement) => void) {
    this.onHeaderInput = func
  }

  /** Register the content transformer to handle the content of the last submit by the editor */
  protected registerContentTransformer(func: (raw: string) => string) {
    this.contentTransformer = func
  }

  /** Register the event of editor content is updated */
  protected registerContentUpdatedEvt(func: (content: string) => void) {
    this.onContentUpdated = func
  }

  /** Get panel element of plug */
  public getPanel() {
    return this.$panel
  }

  /** Get btn element of plug */
  public getBtn() {
    return this.$btn
  }
}

export default EditorPlug

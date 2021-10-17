import Context from "~/src/Context"
import Editor from "../Editor"

export default abstract class EditorPlug {
  protected editor: Editor
  protected ctx: Context
  public abstract el: HTMLElement

  constructor (editor: Editor) {
    this.editor = editor
    this.ctx = editor.ctx
  }

  public abstract initEl(): void
  public abstract getEl(): HTMLElement
  public abstract getName(): string
  public abstract getBtnHtml(): string
  public abstract onShow(): void
  public abstract onHide(): void
}

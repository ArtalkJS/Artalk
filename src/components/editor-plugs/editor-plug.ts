import Context from "~/src/context"
import Editor from "../editor"

export default abstract class EditorPlug {
  protected editor: Editor
  protected ctx: Context
  public abstract $el: HTMLElement

  constructor (editor: Editor) {
    this.editor = editor
    this.ctx = editor.ctx
  }

  public static Name: string
  public static BtnHTML: string

  public abstract getEl(): HTMLElement
  public abstract onShow(): void
  public abstract onHide(): void
}

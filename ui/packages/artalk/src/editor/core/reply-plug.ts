import type { CommentData } from '~/types/artalk-data'
import Editor from '../editor'
import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import EditorPlug from '../editor-plug'
import MoverPlug from './mover-plug'

export default class ReplyPlug extends EditorPlug {
  private comment?: CommentData

  constructor(editor: Editor) {
    super(editor)
  }

  getComment() {
    return this.comment
  }

  setReply(commentData: CommentData, $comment: HTMLElement, scroll = true) {
    this.editor.cancelEditComment()
    this.cancelReply()

    const ui = this.editor.getUI()
    if (!ui.$sendReply) {
      ui.$sendReply = Utils.createElement(
        `<div class="atk-send-reply">` +
          `${this.editor.$t('reply')} ` +
          `<span class="atk-text"></span><span class="atk-cancel">×</span>` +
        `</div>`
      )
      ui.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      ui.$sendReply.addEventListener('click', () => {
        this.editor.cancelReply()
      })
      ui.$textareaWrap.append(ui.$sendReply)
    }

    this.comment = commentData
    this.editor.getPlugs()?.get(MoverPlug)?.move($comment)

    if (scroll) Ui.scrollIntoView(ui.$el)

    ui.$textarea.focus()
  }

  cancelReply() {
    if (!this.comment) return

    const ui = this.editor.getUI()
    if (ui.$sendReply) {
      ui.$sendReply.remove()
      ui.$sendReply = undefined
    }
    this.comment = undefined

    this.editor.getPlugs()?.get(MoverPlug)?.back()
  }
}

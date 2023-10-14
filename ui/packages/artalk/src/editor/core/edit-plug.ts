import type { CommentData } from '~/types/artalk-data'
import * as Utils from '../../lib/utils'
import Editor from '../editor'
import User from '../../lib/user'
import EditorPlug from '../editor-plug'
import SubmitPlug from './submit-plug'

export default class EditPlug extends EditorPlug {
  private comment?: CommentData

  getComment() {
    return this.comment
  }

  getIsEditMode() {
    return !!this.comment
  }

  constructor(editor: Editor) {
    super(editor)

    this.kit.useMounted(() => {
      const submitPlug = this.editor.getPlugs()?.get(SubmitPlug)
      if (!submitPlug) throw Error("SubmitPlug not initialized")

      submitPlug.registerCustom({
        activeCond: () => !!this.comment, // active this custom submit when edit mode
        req: async () => {
          const saveData = {
            content: this.editor.getFinalContent(),
            nick: this.editor.getUI().$nick.value,
            email: this.editor.getUI().$email.value,
            link: this.editor.getUI().$link.value,
          }
          const nComment = await this.editor.ctx.getApi().comment.commentEdit({
            ...this.comment, ...saveData
          })
          return nComment
        },
        post: (nComment: CommentData) => {
          this.editor.ctx.updateComment(nComment)
        }
      })
    })
  }

  edit(comment: CommentData, $comment: HTMLElement) {
    this.cancelEdit()
    this.editor.cancelReply()

    const ui = this.editor.getUI()
    if (!ui.$editCancelBtn) {
      const $btn = Utils.createElement(
        `<div class="atk-send-reply">` +
          `${this.editor.$t('editCancel')} ` +
          `<span class="atk-cancel">×</span>` +
        `</div>`
      )
      $btn.onclick = () => {
        this.cancelEdit()
      }
      ui.$textareaWrap.append($btn)
      ui.$editCancelBtn = $btn
    }
    this.comment = comment

    ui.$header.style.display = 'none' // TODO support modify header information

    this.editor.move($comment)

    ui.$nick.value = comment.nick || ''
    ui.$email.value = comment.email || ''
    ui.$link.value = comment.link || ''

    this.editor.setContent(comment.content)
    ui.$textarea.focus()

    this.updateSubmitBtnText(this.editor.$t('save'))
  }

  cancelEdit() {
    if (!this.comment) return

    const ui = this.editor.getUI()

    if (ui.$editCancelBtn) {
      ui.$editCancelBtn.remove()
      ui.$editCancelBtn = undefined
    }

    this.comment = undefined
    this.editor.moveBack()

    const { nick, email, link } = User.data
    ui.$nick.value = nick
    ui.$email.value = email
    ui.$link.value = link

    this.editor.setContent('')
    this.restoreSubmitBtnText()

    ui.$header.style.display = '' // TODO support modify header information
  }

  // -------------------------------------------------------------------
  //  Submit Btn Text Modifier
  // -------------------------------------------------------------------

  private originalSubmitBtnText = 'Send'

  private updateSubmitBtnText(text: string) {
    this.originalSubmitBtnText = this.editor.getUI().$submitBtn.innerText
    this.editor.getUI().$submitBtn.innerText = text
  }

  private restoreSubmitBtnText() {
    this.editor.getUI().$submitBtn.innerText = this.originalSubmitBtnText
  }
}

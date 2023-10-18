import type { CommentData } from '~/types/artalk-data'
import $t from '@/i18n'
import * as Utils from '@/lib/utils'
import User from '@/lib/user'
import PlugKit from './_kit'
import EditorPlug from './_plug'
import SubmitPlug from './submit'
import MoverPlug from './mover'

export default class EditPlug extends EditorPlug {
  private comment?: CommentData

  getComment() {
    return this.comment
  }

  getIsEditMode() {
    return !!this.comment
  }

  constructor(kit: PlugKit) {
    super(kit)

    // add effect when state switch to `edit`
    this.useEditorStateEffect('edit', (comment) => {
      this.edit(comment)

      return () => {
        this.cancelEdit()
      }
    })

    // register submit preset
    this.kit.useMounted(() => {
      const submitPlug = this.kit.useDeps(SubmitPlug)
      if (!submitPlug) throw Error("SubmitPlug not initialized")

      submitPlug.registerCustom({
        activeCond: () => !!this.comment, // active this custom submit when edit mode
        req: async () => {
          const saveData = {
            content: this.kit.useEditor().getContentFinal(),
            nick: this.kit.useUI().$nick.value,
            email: this.kit.useUI().$email.value,
            link: this.kit.useUI().$link.value,
          }
          const nComment = await this.kit.useApi().comment.commentEdit({
            ...this.comment, ...saveData
          })
          return nComment
        },
        post: (nComment: CommentData) => {
          this.kit.useGlobalCtx().updateComment(nComment)
        }
      })
    })
  }

  private edit(comment: CommentData) {
    const ui = this.kit.useUI()
    if (!ui.$editCancelBtn) {
      const $btn = Utils.createElement(
        `<div class="atk-send-reply">` +
          `${$t('editCancel')} ` +
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

    ui.$nick.value = comment.nick || ''
    ui.$email.value = comment.email || ''
    ui.$link.value = comment.link || ''

    this.kit.useEditor().setContent(comment.content)
    ui.$textarea.focus()

    this.updateSubmitBtnText($t('save'))
  }

  private cancelEdit() {
    if (!this.comment) return

    const ui = this.kit.useUI()

    if (ui.$editCancelBtn) {
      ui.$editCancelBtn.remove()
      ui.$editCancelBtn = undefined
    }

    this.comment = undefined
    this.kit.useDeps(MoverPlug)?.back()

    const { nick, email, link } = User.data
    ui.$nick.value = nick
    ui.$email.value = email
    ui.$link.value = link

    this.kit.useEditor().setContent('')
    this.restoreSubmitBtnText()

    ui.$header.style.display = '' // TODO support modify header information
  }

  // -------------------------------------------------------------------
  //  Submit Btn Text Modifier
  // -------------------------------------------------------------------

  private originalSubmitBtnText = 'Send'

  private updateSubmitBtnText(text: string) {
    this.originalSubmitBtnText = this.kit.useUI().$submitBtn.innerText
    this.kit.useUI().$submitBtn.innerText = text
  }

  private restoreSubmitBtnText() {
    this.kit.useUI().$submitBtn.innerText = this.originalSubmitBtnText
  }
}

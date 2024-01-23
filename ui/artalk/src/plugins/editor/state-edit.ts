import type { CommentData } from '@/types'
import $t from '@/i18n'
import * as Utils from '@/lib/utils'
import type PlugKit from './_kit'
import EditorPlug from './_plug'
import Submit from './submit'

export default class StateEdit extends EditorPlug {
  private comment?: CommentData

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
      const submitPlug = this.kit.useDeps(Submit)
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
          const comment = this.comment!
          const nComment = await this.kit.useApi().comments.updateComment(comment.id, {
            ...comment, ...saveData
          })
          return nComment.data
        },
        post: (nComment: CommentData) => {
          this.kit.useGlobalCtx().getData().updateComment(nComment)
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
        this.kit.useEditor().resetState()
      }
      ui.$textareaWrap.append($btn)
      ui.$editCancelBtn = $btn
    }
    this.comment = comment

    ui.$header.style.display = 'none' // TODO: support modify header information

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

    const { nick, email, link } = this.kit.useUser().getData()
    ui.$nick.value = nick
    ui.$email.value = email
    ui.$link.value = link

    this.kit.useEditor().setContent('')
    this.restoreSubmitBtnText()

    ui.$header.style.display = '' // TODO: support modify header information
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

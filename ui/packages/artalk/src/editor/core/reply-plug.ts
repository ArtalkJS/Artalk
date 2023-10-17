import type { CommentData } from '~/types/artalk-data'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import $t from '@/i18n'
import EditorPlug from '../editor-plug'
import MoverPlug from './mover-plug'
import PlugKit from '../plug-kit'

export default class ReplyPlug extends EditorPlug {
  private comment?: CommentData

  constructor(kit: PlugKit) {
    super(kit)
  }

  getComment() {
    return this.comment
  }

  setReply(commentData: CommentData, $comment: HTMLElement, scroll = true) {
    this.kit.useEditor().cancelEditComment()
    this.cancelReply()

    const ui = this.kit.useUI()
    if (!ui.$sendReply) {
      ui.$sendReply = Utils.createElement(
        `<div class="atk-send-reply">` +
          `${$t('reply')} ` +
          `<span class="atk-text"></span><span class="atk-cancel">×</span>` +
        `</div>`
      )
      ui.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
      ui.$sendReply.addEventListener('click', () => {
        this.kit.useEditor().cancelReply()
      })
      ui.$textareaWrap.append(ui.$sendReply)
    }

    this.comment = commentData
    this.kit.useDeps(MoverPlug)?.move($comment)

    if (scroll) Ui.scrollIntoView(ui.$el)

    ui.$textarea.focus()
  }

  cancelReply() {
    if (!this.comment) return

    const ui = this.kit.useUI()
    if (ui.$sendReply) {
      ui.$sendReply.remove()
      ui.$sendReply = undefined
    }
    this.comment = undefined

    this.kit.useDeps(MoverPlug)?.back()
  }
}

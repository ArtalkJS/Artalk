import type { CommentData } from '~/types/artalk-data'
import Editor from './editor'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

export interface ReplyManager {
  editor: Editor
  comment?: CommentData
  setReply(comment: CommentData, $comment: HTMLElement, scroll?: boolean): void
  cancelReply(): void
}

export function createReplyManager(editor: Editor): ReplyManager {
  const m: ReplyManager = {
    editor,
    comment: undefined,
    setReply: (c, $, s) => setReply(m, c, $, s),
    cancelReply: () => cancelReply(m)
  }

  return m
}

function setReply(m: ReplyManager, commentData: CommentData, $comment: HTMLElement, scroll = true) {
  m.editor.cancelEditComment()
  cancelReply(m)

  const ui = m.editor.getUI()
  if (!ui.$sendReply) {
    ui.$sendReply = Utils.createElement(
      `<div class="atk-send-reply">` +
        `${m.editor.$t('reply')} ` +
        `<span class="atk-text"></span><span class="atk-cancel">×</span>` +
      `</div>`
    )
    ui.$sendReply.querySelector<HTMLElement>('.atk-text')!.innerText = `@${commentData.nick}`
    ui.$sendReply.addEventListener('click', () => {
      m.editor.cancelReply()
    })
    ui.$textareaWrap.append(ui.$sendReply)
  }

  m.comment = commentData
  m.editor.travel($comment)

  if (scroll) Ui.scrollIntoView(ui.$el)

  ui.$textarea.focus()
}

function cancelReply(m: ReplyManager) {
  if (!m.comment) return

  const ui = m.editor.getUI()
  if (ui.$sendReply) {
    ui.$sendReply.remove()
    ui.$sendReply = undefined
  }
  m.comment = undefined

  m.editor.travelBack()
}

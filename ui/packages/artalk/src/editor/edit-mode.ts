import type { CommentData } from '~/types/artalk-data'
import * as Utils from '../lib/utils'
import Editor from './editor'
import User from '../lib/user'

export interface EditModeManager {
  editor: Editor
  comment?: CommentData
  setEdit: (comment: CommentData, $comment: HTMLElement) => void
  cancelEdit: () => void
}

export function createEditModeManager(editor: Editor) {
  const m: EditModeManager = {
    editor,
    comment: undefined,
    setEdit: (c, e) => edit(m, c, e),
    cancelEdit: () => cancelEdit(m)
  }

  initEditModeSubmit(m)

  return m
}

function edit(m: EditModeManager, comment: CommentData, $comment: HTMLElement) {
  cancelEdit(m)
  m.editor.cancelReply()

  const ui = m.editor.getUI()
  if (!ui.$editCancelBtn) {
    const $btn = Utils.createElement(
      `<div class="atk-send-reply">` +
        `${m.editor.$t('editCancel')} ` +
        `<span class="atk-cancel">×</span>` +
      `</div>`
    )
    $btn.onclick = () => {
      cancelEdit(m)
    }
    ui.$textareaWrap.append($btn)
    ui.$editCancelBtn = $btn
  }
  m.comment = comment

  ui.$header.style.display = 'none' // TODO support modify header information

  m.editor.travel($comment)

  ui.$nick.value = comment.nick || ''
  ui.$email.value = comment.email || ''
  ui.$link.value = comment.link || ''

  m.editor.setContent(comment.content)
  ui.$textarea.focus()

  m.editor.refreshSendBtnText()
}

function cancelEdit(m: EditModeManager) {
  if (!m.comment) return

  const ui = m.editor.getUI()

  if (ui.$editCancelBtn) {
    ui.$editCancelBtn.remove()
    ui.$editCancelBtn = undefined
  }

  m.comment = undefined
  m.editor.travelBack()

  const { nick, email, link } = User.data
  ui.$nick.value = nick
  ui.$email.value = email
  ui.$link.value = link

  m.editor.setContent('')
  m.editor.refreshSendBtnText()

  ui.$header.style.display = '' // TODO support modify header information
}

function initEditModeSubmit(m: EditModeManager) {
  m.editor.getSubmitManager()!.registerCustom({
    activeCond: () => !!m.comment, // active this custom submit when edit mode
    req: async () => {
      const saveData = {
        content: m.editor.getFinalContent(),
        nick: m.editor.getUI().$nick.value,
        email: m.editor.getUI().$email.value,
        link: m.editor.getUI().$link.value,
      }
      const nComment = await m.editor.ctx.getApi().comment.commentEdit({
        ...m.comment, ...saveData
      })
      return nComment
    },
    post: (nComment: CommentData) => {
      m.editor.ctx.updateComment(nComment)
    }
  })
}

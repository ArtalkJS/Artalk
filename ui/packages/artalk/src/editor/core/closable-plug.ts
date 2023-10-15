import User from '@/lib/user'
import * as Utils from '@/lib/utils'
import Editor from '../editor'
import EditorPlug from '../editor-plug'

export default class ClosablePlug extends EditorPlug {
  constructor(editor: Editor) {
    super(editor)
  }

  close() {
    if (!this.editor.getUI().$textareaWrap.querySelector('.atk-comment-closed'))
      this.editor.getUI().$textareaWrap.prepend(Utils.createElement(`<div class="atk-comment-closed">${this.editor.$t('onlyAdminCanReply')}</div>`))

    if (!User.data.isAdmin) {
      this.editor.getUI().$textarea.style.display = 'none'
      this.editor.getPlugs()?.closePlugPanel()
      this.editor.getUI().$bottom.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.editor.getUI().$textarea.style.display = ''
      this.editor.getUI().$bottom.style.display = ''
    }
  }

  open() {
    this.editor.getUI().$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.editor.getUI().$textarea.style.display = ''
    this.editor.getUI().$bottom.style.display = ''
  }
}

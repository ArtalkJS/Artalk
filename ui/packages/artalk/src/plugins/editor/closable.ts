import User from '@/lib/user'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'
import EditorPlug from './_plug'
import PlugKit from './_kit'

export default class ClosablePlug extends EditorPlug {
  constructor(kit: PlugKit) {
    super(kit)

    this.kit.useEvents().on('editor-open', () => this.open())
    this.kit.useEvents().on('editor-close', () => this.close())
  }

  open() {
    this.kit.useUI().$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.kit.useUI().$textarea.style.display = ''
    this.kit.useUI().$bottom.style.display = ''
  }

  close() {
    if (!this.kit.useUI().$textareaWrap.querySelector('.atk-comment-closed'))
      this.kit.useUI().$textareaWrap.prepend(Utils.createElement(`<div class="atk-comment-closed">${$t('onlyAdminCanReply')}</div>`))

    if (!User.data.isAdmin) {
      this.kit.useUI().$textarea.style.display = 'none'
      this.kit.useEvents().trigger('panel-close')
      this.kit.useUI().$bottom.style.display = 'none'
    } else {
      // 管理员一直打开评论
      this.kit.useUI().$textarea.style.display = ''
      this.kit.useUI().$bottom.style.display = ''
    }
  }
}

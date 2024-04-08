import * as Utils from '@/lib/utils'
import $t from '@/i18n'
import EditorPlug from './_plug'
import type PlugKit from './_kit'

export default class Closable extends EditorPlug {
  constructor(kit: PlugKit) {
    super(kit)

    const onOpen = () => this.open()
    const onClose = () => this.close()

    this.kit.useMounted(() => {
      this.kit.useEvents().on('editor-open', onOpen)
      this.kit.useEvents().on('editor-close', onClose)
    })
    this.kit.useUnmounted(() => {
      this.kit.useEvents().off('editor-open', onOpen)
      this.kit.useEvents().off('editor-close', onClose)
    })
  }

  private open() {
    this.kit.useUI().$textareaWrap.querySelector('.atk-comment-closed')?.remove()
    this.kit.useUI().$textarea.style.display = ''
    this.kit.useUI().$bottom.style.display = ''
  }

  private close() {
    if (!this.kit.useUI().$textareaWrap.querySelector('.atk-comment-closed'))
      this.kit
        .useUI()
        .$textareaWrap.prepend(
          Utils.createElement(`<div class="atk-comment-closed">${$t('onlyAdminCanReply')}</div>`),
        )

    if (!this.kit.useUser().getData().isAdmin) {
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

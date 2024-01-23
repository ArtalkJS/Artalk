import type { CommentData } from '@/types'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import $t from '@/i18n'
import EditorPlug from './_plug'
import type PlugKit from './_kit'
import Submit from './submit'
import SubmitAddPreset from './submit-add'

export default class StateReply extends EditorPlug {
  private comment?: CommentData

  constructor(kit: PlugKit) {
    super(kit)

    // add effect when state switch to `reply`
    this.useEditorStateEffect('reply', (commentData) => {
      this.setReply(commentData)

      return () => {
        this.cancelReply()
      }
    })

    // register submit preset
    this.kit.useEvents().on('mounted', () => {
      const submitPlug = this.kit.useDeps(Submit)
      if (!submitPlug) throw Error("SubmitPlug not initialized")

      const defaultPreset = new SubmitAddPreset(this.kit)

      submitPlug.registerCustom({
        activeCond: () => !!this.comment, // active this custom submit when reply mode
        req: async () => {
          if (!this.comment) throw new Error('reply comment cannot be empty')

          const nComment = (await this.kit.useApi().comments.createComment({
            ...(await defaultPreset.getSubmitAddParams()),
            rid: this.comment.id,
            page_key: this.comment.page_key,
            page_title: undefined,
            site_name: this.comment.site_name
          })).data

          return nComment
        },
        post: (nComment: CommentData) => {
          // open another page when reply comment is not the same pageKey
          const conf = this.kit.useConf()
          if (nComment.page_key !== conf.pageKey) {
            window.open(`${nComment.page_url}#atk-comment-${nComment.id}`)
          }

          defaultPreset.postSubmitAdd(nComment)
        }
      })
    })
  }

  private setReply(commentData: CommentData) {
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
        this.kit.useEditor().resetState()
      })
      ui.$textareaWrap.append(ui.$sendReply)
    }

    this.comment = commentData

    ui.$textarea.focus()
  }

  private cancelReply() {
    if (!this.comment) return

    const ui = this.kit.useUI()
    if (ui.$sendReply) {
      ui.$sendReply.remove()
      ui.$sendReply = undefined
    }
    this.comment = undefined
  }
}

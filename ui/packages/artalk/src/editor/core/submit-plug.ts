import type { CommentData } from '~/types/artalk-data'
import Editor from '../editor'
import User from '../../lib/user'
import EditorPlug from '../editor-plug'
import ReplyPlug from './reply-plug'

interface CustomSubmit {
  activeCond: () => void
  pre?: () => void
  req?: () => Promise<CommentData>
  post?: (nComment: CommentData) => void
}

export default class SubmitPlug extends EditorPlug {
  customs: CustomSubmit[] = []

  constructor(editor: Editor) {
    super(editor)
  }

  async do() {
    if (this.editor.getFinalContent().trim() === '') {
      this.editor.focus()
      return
    }

    const custom = this.customs.find(o => o.activeCond())

    this.editor.ctx.trigger('editor-submit')
    this.editor.showLoading()

    try {
      // pre submit
      if (custom?.pre) custom.pre()

      let nComment: CommentData

      // submit request
      if (custom?.req) nComment = await custom.req()
      else nComment = await this.reqAdd()

      // post submit
      if (custom?.post) custom.post(nComment)
      else this.postSubmitAdd(nComment)
    } catch (err: any) {
      // submit error
      console.error(err)
      this.editor.showNotify(`${this.editor.$t('commentFail')}，${err.msg || String(err)}`, 'e')
      return
    } finally {
      this.editor.hideLoading()
    }

    this.editor.reset() // 复原编辑器
    this.editor.ctx.trigger('editor-submitted')
  }

  registerCustom(c: CustomSubmit) {
    this.customs.push(c)
  }

  // -------------------------------------------------------------------
  //  Submit CommentAdd
  // -------------------------------------------------------------------

  private async reqAdd() {
    const nComment = await this.editor.ctx.getApi().comment.add({
      ...this.getSubmitAddParams()
    })
    return nComment
  }

  private getSubmitAddParams() {
    const { nick, email, link } = User.data
    const conf = this.editor.ctx.conf
    const reply = this.editor.getPlugs()?.get(ReplyPlug)?.getComment()

    return {
      content: this.editor.getFinalContent(),
      nick, email, link,
      rid: (!reply) ? 0 : reply.id,
      page_key: (!reply) ? conf.pageKey : reply.page_key,
      page_title: (!reply) ? conf.pageTitle : undefined,
      site_name: (!reply) ? conf.site : reply.site_name
    }
  }

  private postSubmitAdd(commentNew: CommentData) {
    // 回复不同页面的评论，跳转到新页面
    const replyComment = this.editor.getPlugs()?.get(ReplyPlug)?.getComment()
    const conf = this.editor.ctx.conf
    if (!!replyComment && replyComment.page_key !== conf.pageKey) {
      window.open(`${replyComment.page_url}#atk-comment-${commentNew.id}`)
    }

    this.editor.ctx.insertComment(commentNew)
  }
}

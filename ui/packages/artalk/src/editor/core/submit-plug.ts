import type { CommentData } from '~/types/artalk-data'
import $t from '@/i18n'
import User from '@/lib/user'
import EditorPlug from '../editor-plug'
import ReplyPlug from './reply-plug'
import PlugKit from '../plug-kit'

interface CustomSubmit {
  activeCond: () => void
  pre?: () => void
  req?: () => Promise<CommentData>
  post?: (nComment: CommentData) => void
}

export default class SubmitPlug extends EditorPlug {
  customs: CustomSubmit[] = []

  constructor(kit: PlugKit) {
    super(kit)
  }

  async do() {
    if (this.kit.useEditor().getContentFinal().trim() === '') {
      this.kit.useEditor().focus()
      return
    }

    const custom = this.customs.find(o => o.activeCond())

    this.kit.useGlobalCtx().trigger('editor-submit')
    this.kit.useEditor().showLoading()

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
      this.kit.useEditor().showNotify(`${$t('commentFail')}，${err.msg || String(err)}`, 'e')
      return
    } finally {
      this.kit.useEditor().hideLoading()
    }

    this.kit.useEditor().reset() // 复原编辑器
    this.kit.useGlobalCtx().trigger('editor-submitted')
  }

  registerCustom(c: CustomSubmit) {
    this.customs.push(c)
  }

  // -------------------------------------------------------------------
  //  Submit CommentAdd
  // -------------------------------------------------------------------

  private async reqAdd() {
    const nComment = await this.kit.useApi().comment.add({
      ...this.getSubmitAddParams()
    })
    return nComment
  }

  private getSubmitAddParams() {
    const { nick, email, link } = User.data
    const conf = this.kit.useConf()
    const reply = this.kit.useDeps(ReplyPlug)?.getComment()

    return {
      content: this.kit.useEditor().getContentFinal(),
      nick, email, link,
      rid: (!reply) ? 0 : reply.id,
      page_key: (!reply) ? conf.pageKey : reply.page_key,
      page_title: (!reply) ? conf.pageTitle : undefined,
      site_name: (!reply) ? conf.site : reply.site_name
    }
  }

  private postSubmitAdd(commentNew: CommentData) {
    // 回复不同页面的评论，跳转到新页面
    const replyComment = this.kit.useDeps(ReplyPlug)?.getComment()
    const conf = this.kit.useConf()
    if (!!replyComment && replyComment.page_key !== conf.pageKey) {
      window.open(`${replyComment.page_url}#atk-comment-${commentNew.id}`)
    }

    this.kit.useGlobalCtx().insertComment(commentNew)
  }
}

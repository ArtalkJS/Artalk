import type { CommentData } from '@/types'
import * as Utils from '@/lib/utils'
import type PlugKit from './_kit'

export default class SubmitAddPreset {
  constructor(private kit: PlugKit) {}

  async reqAdd() {
    const nComment = (await this.kit.useApi().comments.createComment({
      ...(await this.getSubmitAddParams())
    })).data
    return nComment
  }

  async getSubmitAddParams() {
    const { nick, email, link } = this.kit.useUser().getData()
    const conf = this.kit.useConf()

    return {
      content: this.kit.useEditor().getContentFinal(),
      name: nick, email, link,
      rid: 0,
      page_key: conf.pageKey,
      page_title: conf.pageTitle,
      site_name: conf.site,
      ua: await Utils.getCorrectUserAgent(), // Get the corrected UA
    }
  }

  postSubmitAdd(commentNew: CommentData) {
    // insert the new comment to list
    this.kit.useGlobalCtx().getData().insertComment(commentNew)
  }
}

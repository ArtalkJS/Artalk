import type PlugKit from './_kit'
import type { CommentData } from '@/types'
import { getCorrectUserAgent } from '@/lib/utils'

export default class SubmitAddPreset {
  constructor(private kit: PlugKit) {}

  async reqAdd(replyTo?: CommentData) {
    const api = this.kit.useApi()
    const params = await this.getSubmitAddParams()
    return (
      await api.comments.createComment({
        ...params,
        ...(replyTo && {
          rid: replyTo.id,
          page_key: replyTo.page_key,
          page_title: undefined,
          site_name: replyTo.site_name,
        }),
      })
    ).data
  }

  async getSubmitAddParams() {
    const { name, email, link, token } = this.kit.useUser().getData()
    const conf = this.kit.useConf()
    const params = {
      content: this.kit.useEditor().getContentFinal(),
      name,
      email,
      link,
      rid: 0,
      page_key: conf.pageKey,
      page_title: conf.pageTitle,
      site_name: conf.site,
      ua: await getCorrectUserAgent(),
    }
    // Client hints can finish after the viewer or destination has changed.
    const currentUser = this.kit.useUser().getData()
    const currentConf = this.kit.useConf()
    if (
      currentUser.token !== token ||
      currentUser.name !== name ||
      currentUser.email !== email ||
      currentConf.server !== conf.server ||
      currentConf.site !== conf.site ||
      currentConf.pageKey !== conf.pageKey
    ) {
      throw new Error('Your account or page changed. Please submit the comment again.')
    }
    return params
  }

  postSubmitAdd(commentNew: CommentData) {
    this.kit.useData().insertComment(commentNew)
  }
}

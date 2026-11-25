import type { Context } from 'artalk'
import type { Anchor, AnnotationCommentData } from './types'

export interface AnnotationPage {
  comments: AnnotationCommentData[]
  rootsCount: number
}

export class AnnotationApi {
  constructor(private ctx: Context) {}

  private page() {
    const conf = this.ctx.getConf()
    return { page_key: conf.pageKey, site_name: conf.site }
  }

  private user() {
    const user = this.ctx.getUser().getData()
    return { name: user.name, email: user.email, link: user.link }
  }

  async listThreads(offset = 0, limit = 50): Promise<AnnotationPage> {
    const { data } = await this.ctx.getApi().comments.getComments({
      ...this.page(),
      has_anchor: 'true',
      flat_mode: false,
      offset,
      limit,
    })
    return {
      comments: data.comments as AnnotationCommentData[],
      rootsCount: data.roots_count,
    }
  }

  async createAnnotation(content: string, anchor: Anchor): Promise<AnnotationCommentData> {
    const { data } = await this.ctx.getApi().comments.createComment({
      ...this.page(),
      page_title: this.ctx.getConf().pageTitle,
      ...this.user(),
      content,
      anchor,
    })
    return data as AnnotationCommentData
  }

  async reply(rootId: number, content: string): Promise<AnnotationCommentData> {
    const { data } = await this.ctx.getApi().comments.createComment({
      ...this.page(),
      page_title: this.ctx.getConf().pageTitle,
      ...this.user(),
      content,
      rid: rootId,
    })
    return data as AnnotationCommentData
  }

  async updateAnchor(
    id: number,
    body: { resolved?: boolean; anchor?: Anchor },
  ): Promise<AnnotationCommentData> {
    const { data } = await this.ctx.getApi().comments.updateCommentAnchor(id, body)
    return data as AnnotationCommentData
  }

  async remove(id: number): Promise<void> {
    await this.ctx.getApi().comments.deleteComment(id)
  }
}

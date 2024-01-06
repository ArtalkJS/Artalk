import type { CommentData, ListData } from '@/types'
import * as Utils from '../lib/utils'
import ApiBase from './_base'

/**
 * 评论 API
 */
export default class CommentApi extends ApiBase {
  /** 评论 · 获取 */
  public get(offset: number, pageSize: number, flatMode?: boolean, paramsEditor?: (params: any) => void) {
    const params: any = {
      page_key: this.options.pageKey,
      site_name: this.options.siteName,
      limit: pageSize,
      offset,
    }

    if (flatMode) params.flat_mode = flatMode // 平铺模式

    this.withUserInfo(params)

    if (paramsEditor) paramsEditor(params)

    return this.fetch<ListData>('GET', '/comments', params)
  }

  /** 评论 · 创建 */
  public async add(comment: { nick: string, email: string, link: string, content: string, rid: number, page_key: string, page_title?: string, site_name?: string }) {
    const params: any = {
      name: comment.nick,
      email: comment.email,
      link: comment.link,
      content: comment.content,
      rid: comment.rid,
      page_key: comment.page_key,
      ua: await Utils.getCorrectUserAgent(), // 需要后端支持，获取修正后的 UA
    }

    if (comment.page_title) params.page_title = comment.page_title
    if (comment.site_name) params.site_name = comment.site_name

    const data = await this.fetch<CommentData>('POST', '/comments', params)

    return data
  }

  /** 评论 · 修改 */
  public async commentEdit(id: number, data: Partial<CommentData>) {
    const params: any = {
      ...data,
    }

    const d = await this.fetch<CommentData>('PUT', `/comments/${id}`, params)

    return d
  }

  /** 评论 · 删除 */
  public commentDel(id: number) {
    return this.fetch('DELETE', `/comments/${id}`)
  }

  /** 投票 */
  public async vote(targetID: number, type: 'comment_up'|'comment_down'|'page_up'|'page_down') {
    const data = await this.fetch<{
      up: number,
      down: number
    }>('POST', `/votes/${type}/${targetID}`, this.withUserInfo({}))

    return data
  }
}

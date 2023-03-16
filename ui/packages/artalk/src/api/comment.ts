import { CommentData, ListData } from '~/types/artalk-data'
import * as Utils from '../lib/utils'
import ApiBase from './api-base'
import User from '../lib/user'

/**
 * 评论 API
 */
export default class CommentApi extends ApiBase {
  /** 评论 · 获取 */
  public get(offset: number, pageSize: number, flatMode?: boolean, paramsEditor?: (params: any) => void) {
    const params: any = {
      page_key: this.ctx.conf.pageKey,
      site_name: this.ctx.conf.site || '',
      limit: pageSize,
      offset,
    }

    if (flatMode) params.flat_mode = flatMode // 平铺模式
    if (User.checkHasBasicUserInfo()) {
      params.name = User.data.nick
      params.email = User.data.email
    }

    if (paramsEditor) paramsEditor(params)

    return this.POST<ListData>('/get', params)
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

    const data = await this.POST<any>('/add', params)
    return (data.comment as CommentData)
  }

  /** 评论 · 修改 */
  public async commentEdit(data: Partial<CommentData>) {
    const params: any = {
      ...data,
    }

    const d = await this.POST<any>('/admin/comment-edit', params)
    return (d.comment as CommentData)
  }

  /** 评论 · 删除 */
  public commentDel(commentID: number, siteName?: string) {
    const params: any = {
      id: String(commentID),
      site_name: siteName || '',
    }

    return this.POST('/admin/comment-del', params)
  }

  /** 投票 */
  public async vote(targetID: number, type: string) {
    const params: any = {
      target_id: targetID,
      type,
    }

    if (User.checkHasBasicUserInfo()) {
      params.name = User.data.nick
      params.email = User.data.email
    }

    const data = await this.POST<any>('/vote', params)
    return (data as {up: number, down: number})
  }
}

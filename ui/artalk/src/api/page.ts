import type { PageData, CommentData } from '@/types'
import ApiBase from './_base'

/**
 * 页面 API
 */
export default class PageApi extends ApiBase {
  /** 页面 · 获取 */
  public async pageGet(siteName?: string, offset?: number, limit?: number) {
    const params: any = {
      site_name: siteName || '',
      offset: offset || 0,
      limit: limit || 15,
    }

    const d = await this.fetch<{
      pages: PageData[],
      count: number
    }>('GET', '/pages', params)

    return d
  }

  /** 页面 · 修改 */
  public async pageEdit(data: PageData) {
    const params: any = {
      key: data.key,
      title: data.title,
      admin_only: data.admin_only,
      site_name: data.site_name || this.options.siteName,
    }

    const d = await this.fetch<PageData>('PUT', `/pages/${params.id}`, params)

    return d
  }

  /** 页面 · 删除 */
  public pageDel(id: number) {
    return this.fetch('DELETE', `/pages/${id}`, id)
  }

  /** 页面 · 数据更新 */
  public async pageFetch(id?: number, siteName?: string, getStatus?: boolean) {
    const params: any = {}
    if (id) params.id = id
    if (siteName) params.site_name = siteName
    if (getStatus) params.get_status = getStatus

    const d = await this.fetch<PageData>('POST', `/pages/${id}/fetch`, params)
    return d
  }

  /** PV */
  public async pv() {
    const params: any = {
      page_key: this.options.pageKey,
      page_title: this.options.pageTitle,
      site_name: this.options.siteName
    }

    const p = await this.fetch<any>('POST', `/pages/pv`, params)
    return p.pv as number
  }

  /** 统计 */
  public async stat(
    type: 'latest_comments'|'latest_pages'|'pv_most_pages'|'comment_most_pages'|
          'page_pv'|'site_pv'|'page_comment'|'site_comment',
    pageKeys?: string|string[],
    limit?: number
  ) {
    const params: any = {}

    if (pageKeys) params.page_keys = Array.isArray(pageKeys) ? pageKeys.join(',') : pageKeys
    if (limit) params.limit = limit

    const data = await this.fetch<PageData[]|CommentData[]|object|number>('POST', `/stats/${type}`, params)
    return data
  }
}

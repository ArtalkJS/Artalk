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

    const d = await this.fetch<PageData>('PUT', `/pages/${data.id}`, params)

    return d
  }

  /** 页面 · 删除 */
  public pageDel(id: number) {
    return this.fetch('DELETE', `/pages/${id}`)
  }

  /** 页面 · 数据更新 */
  public async pageFetch(id: number) {
    const d = await this.fetch<PageData>('POST', `/pages/${id}/fetch`)
    return d
  }

  /** 页面 · 整站数据更新 */
  public async pagesAllFetch(siteName?: string) {
    const params: any = {}
    if (siteName) params.site_name = siteName
    const d = await this.fetch('POST', `/pages/fetch`, params)
    return d
  }

  /** 页面 · 整站数据更新 - 当前状态 */
  public async pagesAllFetchStatus() {
    const d = await this.fetch<{
      msg: string,
      is_progress: boolean,
      done: number,
      total: number
    }>('GET', `/pages/fetch/status`)
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

    params.site_name = this.options.siteName
    if (pageKeys) params.page_keys = Array.isArray(pageKeys) ? pageKeys.join(',') : pageKeys
    if (limit) params.limit = limit

    const data = await this.fetch<{
      data: PageData[]|CommentData[]|object|number
    }>('GET', `/stats/${type}`, params)
    return data
  }
}

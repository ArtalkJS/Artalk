import { PageData, CommentData } from '~/types/artalk-data'
import ApiBase from './api-base'

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

    const d = await this.POST<any>('/admin/page-get', params)
    return (d as { pages: PageData[], total: number })
  }

  /** 页面 · 修改 */
  public async pageEdit(data: PageData) {
    const params: any = {
      id: data.id,
      key: data.key,
      title: data.title,
      admin_only: data.admin_only,
      site_name: data.site_name || this.ctx.conf.site,
    }

    const d = await this.POST<any>('/admin/page-edit', params)
    return (d.page as PageData)
  }

  /** 页面 · 删除 */
  public pageDel(pageKey: string, siteName?: string) {
    const params: any = {
      key: String(pageKey),
      site_name: siteName || '',
    }

    return this.POST('/admin/page-del', params)
  }

  /** 页面 · 数据更新 */
  public async pageFetch(id?: number, siteName?: string, getStatus?: boolean) {
    const params: any = {}
    if (id) params.id = id
    if (siteName) params.site_name = siteName
    if (getStatus) params.get_status = getStatus

    const d = await this.POST<any>('/admin/page-fetch', params)
    return (d as any)
  }

  /** PV */
  public async pv() {
    const params: any = {
      page_key: this.ctx.conf.pageKey || '',
      page_title: this.ctx.conf.pageTitle || ''
    }

    const p = await this.POST<any>('/pv', params)
    return p.pv as number
  }

  /** 统计 */
  public async stat(
    type: 'latest_comments'|'latest_pages'|'pv_most_pages'|'comment_most_pages'|
          'page_pv'|'site_pv'|'page_comment'|'site_comment',
    pageKeys?: string|string[],
    limit?: number
  ) {
    const params: any = {
      type,
    }

    if (pageKeys) params.page_keys = Array.isArray(pageKeys) ? pageKeys.join(',') : pageKeys
    if (limit) params.limit = limit

    const data = await this.POST<PageData[]|CommentData[]|object|number>(`/stat`, params)
    return data
  }
}

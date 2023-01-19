import { SiteData } from '~/types/artalk-data'
import ApiBase from './api-base'

/**
 * 站点 API
 */
export default class SiteApi extends ApiBase {
  /** 站点 · 获取 */
  public async siteGet() {
    const params: any = {}

    const d = await this.POST<any>('/admin/site-get', params)
    return (d.sites as SiteData[])
  }

  /** 站点 · 创建 */
  public async siteAdd(name: string, urls: string) {
    const params: any = {
      name, urls,
      site_name: '' // 全局保留字段，当前站点名
    }

    const d = await this.POST<any>('/admin/site-add', params)
    return (d.site as SiteData)
  }

  /** 站点 · 修改 */
  public async siteEdit(data: SiteData) {
    const params: any = {
      id: data.id,
      name: data.name || '',
      urls: data.urls || '',
    }

    const d = await this.POST<any>('/admin/site-edit', params)
    return (d.site as SiteData)
  }

  /** 站点 · 删除 */
  public siteDel(id: number, delContent = false) {
    const params: any = { id, del_content: delContent }

    return this.POST('/admin/site-del', params)
  }

  /** 导出 */
  public async export() {
    const d = await this.Fetch(`/admin/export`, { method: 'POST' }, 0)
    return (d.data?.data || '' as string)
  }
}

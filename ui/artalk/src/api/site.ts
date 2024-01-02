import type { SiteData } from '@/types'
import ApiBase from './_base'

/**
 * 站点 API
 */
export default class SiteApi extends ApiBase {
  /** 站点 · 获取 */
  public async siteGet() {
    const params: any = {}

    const d = await this.fetch<any>('GET', '/sites', params)
    return (d.sites as SiteData[])
  }

  /** 站点 · 创建 */
  public async siteAdd(name: string, urls: string) {
    const params: any = { name, urls }

    const d = await this.fetch<any>('POST', '/sites', params)
    return (d.site as SiteData)
  }

  /** 站点 · 修改 */
  public async siteEdit(id: number, data: SiteData) {
    const params: any = {
      name: data.name || '',
      urls: data.urls || '',
    }

    const d = await this.fetch<any>('PUT', `/sites/${id}`, params)
    return (d.site as SiteData)
  }

  /** 站点 · 删除 */
  public siteDel(id: number) {
    return this.fetch('DELETE', `/sites/${id}`)
  }

  /** 导出 */
  public async export() {
    const d = await this.fetch('GET', `/artransfer/export`, undefined, { timeout: 0 })
    return (d.data?.data || '' as string)
  }
}

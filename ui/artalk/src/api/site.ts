import type { SiteData } from '@/types'
import ApiBase from './_base'

/**
 * 站点 API
 */
export default class SiteApi extends ApiBase {
  /** 站点 · 获取 */
  public async siteGet() {
    const params: any = {}

    const d = await this.fetch<{
      sites: SiteData[],
      count: number
    }>('GET', '/sites', params)
    return d
  }

  /** 站点 · 创建 */
  public async siteAdd(name: string, urls: string) {
    const params: any = { name, urls }

    const d = await this.fetch<SiteData>('POST', '/sites', params)
    return d
  }

  /** 站点 · 修改 */
  public async siteEdit(id: number, data: SiteData) {
    const params: any = {
      name: data.name || '',
      urls: data.urls || '',
    }

    const d = await this.fetch<SiteData>('PUT', `/sites/${id}`, params)
    return d
  }

  /** 站点 · 删除 */
  public siteDel(id: number) {
    return this.fetch('DELETE', `/sites/${id}`)
  }

  /** 导出 */
  public async export() {
    const d = await this.fetch<{ artrans: string }>('GET', `/transfer/export`, undefined, { timeout: 0 })
    return (d.artrans || '' as string)
  }
}

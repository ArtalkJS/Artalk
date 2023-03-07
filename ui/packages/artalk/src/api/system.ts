import ArtalkConfig from '~/types/artalk-config'
import ApiBase from './api-base'
import { handleBackendRefConf } from '../config'

/**
 * 系统 API
 */
export default class SystemApi extends ApiBase {
  /** 获取配置 */
  public async conf() {
    const data = await this.POST<any>(`/conf`)
    const conf = (data.frontend_conf || {}) as ArtalkConfig

    return handleBackendRefConf(conf)
  }

  /** 获取配置数据 */
  public async getSettings() {
    const data = await this.POST<{custom: string, template: string}>('/admin/setting-get')
    return data
  }

  /** 保存配置数据 */
  public async saveSettings(yamlStr: string) {
    const data = await this.POST<boolean>('/admin/setting-save', {
      data: yamlStr,
    })
    return data
  }

  /** 获取 API 版本信息 */
  public async version() {
    const resp = await fetch(`${this.api.baseURL}/version`, { method: 'POST' })
    const data = await resp.json()
    return data as { app: string, version: string, commit_hash: string, fe_min_version: string }
  }
}

import type { ArtalkConfig } from '@/types'
import ApiBase from './_base'

interface ApiVersionInfo {
  app: string
  version: string
  commit_hash: string
}

interface SystemConfResp {
  frontend_conf?: Partial<ArtalkConfig>
  version?: ApiVersionInfo
}

/**
 * 系统 API
 */
export default class SystemApi extends ApiBase {
  /** 获取配置 */
  public async conf() {
    const data = await this.fetch<SystemConfResp>('GET', `/conf`)

    return data
  }

  /** 获取配置数据 */
  public async getSettings() {
    const data = await this.fetch<{ yaml: string }>('GET', '/settings')
    return data
  }

  public async getSettingsTemplate(locale?: string) {
    let path = '/settings/template'
    if (locale) path += `/${locale}`
    const data = await this.fetch<{ yaml: string }>('GET', path)
    return data
  }

  /** 保存配置数据 */
  public async saveSettings(yaml: string) {
    const data = await this.fetch('POST', '/settings', {
      yaml
    })
    return data
  }

  /** 获取 API 版本信息 */
  public async version() {
    const resp = await fetch(`${this.options.baseURL}/version`, { method: 'GET' })
    const data = await resp.json()
    return data as ApiVersionInfo
  }

  /** 缓存清除 */
  public cacheFlushAll() {
    return this.fetch('POST', '/cache/flush')
  }

  /** 缓存预热 */
  public cacheWarmUp() {
    return this.fetch('POST', '/cache/warm_up')
  }
}

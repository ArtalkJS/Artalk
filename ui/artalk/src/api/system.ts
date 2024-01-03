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
  public async conf(): Promise<SystemConfResp> {
    const data = await this.fetch<any>('GET', `/conf`)

    return {
      frontend_conf: data.frontend_conf,
      version: data.version,
    }
  }

  /** 获取配置数据 */
  public async getSettings() {
    const data = await this.fetch<{custom: string, template: string}>('GET', '/settings')
    return data
  }

  /** 保存配置数据 */
  public async saveSettings(yamlStr: string) {
    const data = await this.fetch<boolean>('POST', '/settings', {
      data: yamlStr,
    })
    return data
  }

  /** 获取 API 版本信息 */
  public async version() {
    const resp = await fetch(`${this.options.baseURL}/version`, { method: 'POST' })
    const data = await resp.json()
    return data as { app: string, version: string, commit_hash: string, fe_min_version: string }
  }
}

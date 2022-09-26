import ArtalkConfig from '~/types/artalk-config'
import ApiBase from './api-base'

/**
 * 系统 API
 */
export default class SystemApi extends ApiBase {
  /** 获取配置 */
  public async conf() {
    const data = await this.POST<any>(`/conf`)
    const conf = (data.frontend_conf || {}) as ArtalkConfig

    // Patch: `emoticons` config string to json
    if (conf.emoticons && typeof conf.emoticons === "string") {
      conf.emoticons = conf.emoticons.trim()
      if (conf.emoticons.startsWith("[") || conf.emoticons.startsWith("{")) {
        conf.emoticons = JSON.parse(conf.emoticons) // pase json
      } else if (conf.emoticons === "false") {
        conf.emoticons = false
      }
    }

    return conf
  }

  /** 获取配置数据 */
  public async getSettings() {
    const data = await this.POST<string>('/admin/setting-get')
    return data
  }

  /** 保存配置数据 */
  public async saveSettings(yamlStr: string) {
    const data = await this.POST<boolean>('/admin/setting-save', {
      data: yamlStr,
    })
    return data
  }
}

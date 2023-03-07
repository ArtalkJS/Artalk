import type ArtalkConfig from '~/types/artalk-config'
import * as Utils from './lib/utils'
import Defaults from './defaults'

export default { handelBaseConf }

/** 基本配置初始化 */
export function handelBaseConf(customConf: Partial<ArtalkConfig>): ArtalkConfig {
  // 合并默认配置
  const conf: ArtalkConfig = Utils.mergeDeep(Defaults, customConf)

  // 绑定元素
  if (typeof conf.el === 'string' && !!conf.el) {
    try {
      const findEl = document.querySelector<HTMLElement>(conf.el)
      if (!findEl) throw Error(`Target element "${conf.el}" was not found.`)
      conf.el = findEl
    } catch (e) {
      console.error(e)
      throw new Error('Please check your Artalk `el` config.')
    }
  }

  // 服务器配置
  conf.server = conf.server.replace(/\/$/, '').replace(/\/api\/?$/, '')

  // 默认 pageKey
  if (!conf.pageKey) {
    // @link http://bl.ocks.org/abernier/3070589
    conf.pageKey = `${window.location.pathname}`
  }

  // 默认 pageTitle
  if (!conf.pageTitle) {
    conf.pageTitle = `${document.title}`
  }

  // 自适应语言
  if (conf.locale === 'auto') {
    conf.locale = navigator.language
  }

  return conf
}

export function handleBackendRefConf(conf: Partial<ArtalkConfig>) {
  const DisabledKeys: (keyof ArtalkConfig)[] = [
    'el', 'pageKey', 'pageTitle', 'server', 'site', 'darkMode'
  ]
  Object.keys(conf).forEach(k => {
    if (DisabledKeys.includes(k as any)) delete conf[k]
  })

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

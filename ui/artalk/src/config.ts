import type { ArtalkConfig, ContextApi } from '~/types'
import type { ApiOptions } from './api/_options'
import * as Utils from './lib/utils'
import Defaults from './defaults'

/**
 * Merges the user custom config with the default config
 *
 * @param customConf - The custom config object which is provided by the user
 * @returns The config for Artalk instance creation
 */
export function handelCustomConf(customConf: Partial<ArtalkConfig>): ArtalkConfig {
  // 合并默认配置
  const conf: ArtalkConfig = Utils.mergeDeep({}, customConf)

  // TODO the type of el options may HTMLElement, use it directly instead of from mergeDeep
  if (customConf.el) conf.el = customConf.el

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

  // flatMode
  if (conf.flatMode === true || Number(conf.nestMax) <= 1)
    conf.flatMode = true

  // 自动判断启用平铺模式
  if (conf.flatMode === 'auto')
    conf.flatMode = window.matchMedia("(max-width: 768px)").matches

  return conf
}

/**
 * Handle the config which is provided by the server
 *
 * @param conf - The Server response config for control the frontend of Artalk remotely
 * @returns The config for Artalk instance creation
 */
export function handleConfFormServer(conf: Partial<ArtalkConfig>) {
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

/**
 * Convert Artalk Config to ApiOptions for Api client
 *
 * @param conf - Artalk config
 * @param ctx - If `ctx` not provided, `checkAdmin` and `checkCaptcha` will be disabled
 * @returns ApiOptions for Api client instance creation
 */
export function convertApiOptions(conf: Partial<ArtalkConfig>, ctx?: ContextApi): ApiOptions {
  return {
    baseURL: `${conf.server}/api`,
    siteName: conf.site || '',
    pageKey: conf.pageKey || '',
    pageTitle: conf.pageTitle || '',
    timeout: conf.reqTimeout,
    apiToken: ctx?.get('user').getData().token,
    userInfo: ctx?.get('user').checkHasBasicUserInfo() ? {
      name: ctx?.get('user').getData().nick,
      email: ctx?.get('user').getData().email,
    } : undefined,

    onNeedCheckAdmin(payload) {
      ctx?.checkAdmin({
        onSuccess: () => { payload.recall() },
        onCancel: () => { payload.reject() },
      })
    },

    onNeedCheckCaptcha(payload) {
      ctx?.checkCaptcha({
        imgData: payload.data.imgData,
        iframe: payload.data.iframe,
        onSuccess: () => { payload.recall() },
        onCancel: () => { payload.reject() },
      })
    },
  }
}

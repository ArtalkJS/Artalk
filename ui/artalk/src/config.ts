import type { ArtalkConfig, ContextApi } from '@/types'
import type { ApiOptions } from './api/options'
import { mergeDeep } from './lib/merge-deep'
import { createApiHandlers } from './api'
import Defaults from './defaults'

/**
 * Handle the custom config which is provided by the user
 *
 * @param customConf - The custom config object which is provided by the user
 * @param full - If `full` is `true`, the return value will be the complete config for Artalk instance creation
 * @returns The config for Artalk instance creation
 */
export function handelCustomConf(customConf: Partial<ArtalkConfig>, full: true): ArtalkConfig
export function handelCustomConf(customConf: Partial<ArtalkConfig>, full?: false): Partial<ArtalkConfig>
export function handelCustomConf(customConf: Partial<ArtalkConfig>, full = false) {
  // 合并默认配置
  const conf: Partial<ArtalkConfig> = full ? mergeDeep(Defaults, customConf) : customConf

  // 绑定元素
  if (conf.el && typeof conf.el === 'string') {
    try {
      const findEl = document.querySelector<HTMLElement>(conf.el)
      if (!findEl) throw Error(`Target element "${conf.el}" was not found.`)
      conf.el = findEl
    } catch (e) {
      console.error(e)
      throw new Error('Please check your Artalk `el` config.')
    }
  }

  // 默认 pageKey
  if (conf.pageKey === '')
    conf.pageKey = `${window.location.pathname}` // @link http://bl.ocks.org/abernier/3070589

  // 默认 pageTitle
  if (conf.pageTitle === '')
    conf.pageTitle = `${document.title}`

  // 服务器配置
  if (conf.server)
    conf.server = conf.server.replace(/\/$/, '').replace(/\/api\/?$/, '')

  // 自适应语言
  if (conf.locale === 'auto')
    conf.locale = navigator.language

  // 自动判断启用平铺模式
  if (conf.flatMode === 'auto')
    conf.flatMode = window.matchMedia("(max-width: 768px)").matches

  // flatMode
  if (typeof conf.nestMax === 'number' && Number(conf.nestMax) <= 1)
    conf.flatMode = true

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
      conf.emoticons = JSON.parse(conf.emoticons) // parse json
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
    baseURL: `${conf.server}/api/v2`,
    siteName: conf.site || '',
    pageKey: conf.pageKey || '',
    pageTitle: conf.pageTitle || '',
    timeout: conf.reqTimeout,
    getApiToken: () => ctx?.get('user').getData().token,
    userInfo: ctx?.get('user').checkHasBasicUserInfo() ? {
      name: ctx?.get('user').getData().nick,
      email: ctx?.get('user').getData().email,
    } : undefined,
    handlers: ctx?.getApiHandlers(),
  }
}

export function createNewApiHandlers(ctx: ContextApi) {
  const h = createApiHandlers()
  h.add('need_captcha', (res) => ctx.checkCaptcha(res))
  h.add('need_login', () => ctx.checkAdmin({}))

  return h
}

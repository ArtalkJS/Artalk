import type { ApiOptions } from './api/options'
import { mergeDeep } from './lib/merge-deep'
import type { ApiHandlers } from './api'
import { Defaults } from './defaults'
import type { Config, ConfigPartial, UserManager } from '@/types'

/**
 * Handle the custom config which is provided by the user
 *
 * @param customConf - The custom config object which is provided by the user
 * @param full - If `full` is `true`, the return value will be the complete config for Artalk instance creation
 * @returns The config for Artalk instance creation
 */
export function handelCustomConf(customConf: ConfigPartial, full: true): Config
export function handelCustomConf(customConf: ConfigPartial, full?: false): ConfigPartial
export function handelCustomConf(customConf: ConfigPartial, full = false) {
  // Merge default config
  const conf: ConfigPartial = full ? mergeDeep(Defaults, customConf) : customConf

  // Default pageKey
  if (conf.pageKey === '') conf.pageKey = `${window.location.pathname}` // @see http://bl.ocks.org/abernier/3070589

  // Default pageTitle
  if (conf.pageTitle === '') conf.pageTitle = `${document.title}`

  // Server
  if (conf.server) conf.server = conf.server.replace(/\/$/, '').replace(/\/api\/?$/, '')

  // Language auto-detection
  if (conf.locale === 'auto') conf.locale = navigator.language

  // Flat mode auto-detection
  if (conf.flatMode === 'auto') conf.flatMode = window.matchMedia('(max-width: 768px)').matches

  // Change flatMode by nestMax
  if (typeof conf.nestMax === 'number' && Number(conf.nestMax) <= 1) conf.flatMode = true

  return conf
}

/**
 * Handle the config which is provided by the server
 *
 * @param conf - The Server response config for control the frontend of Artalk remotely
 * @returns The config for Artalk instance creation
 */
export function handleConfFormServer(conf: ConfigPartial): ConfigPartial {
  const ExcludedKeys: (keyof Config)[] = [
    'el',
    'pageKey',
    'pageTitle',
    'server',
    'site',
    'pvEl',
    'countEl',
    'statPageKeyAttr',
    'pageVote',
  ]
  Object.keys(conf).forEach((k) => {
    if (ExcludedKeys.includes(k as any)) delete conf[k]
    if (k === 'darkMode' && conf[k] !== 'auto') delete conf[k]
  })

  // Patch: `emoticons` config string to json
  if (conf.emoticons && typeof conf.emoticons === 'string') {
    conf.emoticons = conf.emoticons.trim()
    if (conf.emoticons.startsWith('[') || conf.emoticons.startsWith('{')) {
      conf.emoticons = JSON.parse(conf.emoticons) // parse json
    } else if (conf.emoticons === 'false') {
      conf.emoticons = false
    }
  }

  return conf
}

/**
 * Get the root element of Artalk
 *
 * @param conf - Artalk config
 * @returns The root element of Artalk
 */
export function getRootEl(conf: ConfigPartial): HTMLElement {
  let $root: HTMLElement
  if (typeof conf.el === 'string') {
    const el = document.querySelector<HTMLElement>(conf.el)
    if (!el) throw new Error(`Element "${conf.el}" not found.`)
    $root = el
  } else if (conf.el instanceof HTMLElement) {
    $root = conf.el
  } else {
    throw new Error('Please provide a valid `el` config for Artalk.')
  }
  return $root
}

/**
 * Convert Artalk Config to ApiOptions for Api client
 *
 * @param conf - Artalk config
 * @param ctx - If `ctx` not provided, `checkAdmin` and `checkCaptcha` will be disabled
 * @returns ApiOptions for Api client instance creation
 */
export function convertApiOptions(
  conf: ConfigPartial,
  user?: UserManager,
  handlers?: ApiHandlers,
): ApiOptions {
  return {
    baseURL: `${conf.server}/api/v2`,
    siteName: conf.site || '',
    pageKey: conf.pageKey || '',
    pageTitle: conf.pageTitle || '',
    timeout: conf.reqTimeout,
    getApiToken: () => user?.getData().token,
    userInfo: user?.checkHasBasicUserInfo()
      ? {
          name: user?.getData().name,
          email: user?.getData().email,
        }
      : undefined,
    handlers,
  }
}

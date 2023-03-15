import en from './en'
import zhCN from './zh-CN'
import { GLOBAL_LOCALES_KEY } from './external'

export type I18n = typeof en
export type I18nKeys = keyof I18n

// @note the key of language is followed by `ISO 639`
// https://en.wikipedia.org/wiki/ISO_639
// https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
export const internal = {
  'en': en,
  'en-US': en,
  'zh-CN': zhCN,
}

/**
 * find a locale object by language name
 */
export function findLocaleSet(lang: string): I18n {
  // normalize a key of language to `ISO 639`
  lang = lang.replace(
    /^([a-zA-Z]+)(-[a-zA-Z]+)?$/,
    (_, p1: string, p2: string) => (p1.toLowerCase() + (p2 || '').toUpperCase())
  )

  // internal finding
  if (internal[lang]) {
    return internal[lang]
  }

  // external finding
  if (window[GLOBAL_LOCALES_KEY] && window[GLOBAL_LOCALES_KEY][lang]) {
    return window[GLOBAL_LOCALES_KEY][lang]
  }

  // case when not found:
  // use `en` by default
  return internal.en
}

/**
 * System locale setting
 */
let LocaleConf: I18n|string = 'en'
let LocaleDict: I18n = findLocaleSet(LocaleConf) // en by default

/**
 * Set system locale
 */
export function setLocale(locale: I18n|string) {
  if (locale === LocaleConf) return
  LocaleConf = locale
  LocaleDict = (typeof locale === 'string') ? findLocaleSet(locale) : locale
}

/**
 * Get an i18n message by key
 */
export function t(key: I18nKeys, args: {[key: string]: string} = {}) {
  let str = LocaleDict?.[key] || key
  str = str.replace(/\{\s*(\w+?)\s*\}/g, (_, token) => args[token] || '')

  return str
}

export default t

import en from './en'
import zhCN from './zh-CN'

export type I18n = typeof en

// @note the key of language is followed by `ISO 639`
// https://en.wikipedia.org/wiki/ISO_639
// https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
export const internal = {
  'en': en,
  'en-US': en,
  'zh-CN': zhCN,
}

const GLOBAL_LOCALES_KEY = "ArtalkI18n";

export function defineLocaleExternal(lang: string, locale: I18n, aliases?: string[]) {
  if (!window[GLOBAL_LOCALES_KEY]) window[GLOBAL_LOCALES_KEY] = {}
  window[GLOBAL_LOCALES_KEY][lang] = locale
  if (aliases) aliases.forEach(l => { window[GLOBAL_LOCALES_KEY][l] = locale })
  return locale
}

/**
 * get a locale object by language name
 */
function getLocaleSet(lang: string): I18n {
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
 * get a i18n message by key
 */
function getI18n(locale: I18n|string, key: keyof I18n, args: {[key: string]: string} = {}) {
  if (typeof locale === 'string') {
    locale = getLocaleSet(locale)
  }

  let str = locale?.[key] || key
  str = str.replace(/\{\s*(\w+?)\s*\}/g, (_, token) => args[token] || '')

  return str
}

export { getLocaleSet }
export default getI18n

import en from './en'
import zhCN from './zh-CN'

export type I18n = typeof en

// @note the Key name is followed by `ISO 639`
// https://en.wikipedia.org/wiki/ISO_639
// https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
export const internal = {
  'en': en,
  'en-US': en,
  'zh-CN': zhCN,
}

export const external = {
  'jp-JP': {},
  'zh-TW': {},
  'fr-FR': {},
  'de-DE': {},
  'bn-IN': {},
}

import type { I18n } from '.'

export const GLOBAL_LOCALES_KEY = "ArtalkI18n"

export function defineLocaleExternal(lang: string, locale: I18n, aliases?: string[]) {
  if (!window[GLOBAL_LOCALES_KEY]) window[GLOBAL_LOCALES_KEY] = {}
  window[GLOBAL_LOCALES_KEY][lang] = locale
  if (aliases) aliases.forEach(l => { window[GLOBAL_LOCALES_KEY][l] = locale })
  return locale
}

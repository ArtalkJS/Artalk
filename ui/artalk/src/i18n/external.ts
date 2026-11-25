import type { I18n } from '.'

export const GLOBAL_LOCALES_KEY = 'ArtalkI18n'

export function defineLocaleExternal(lang: string, locale: Partial<I18n>, aliases?: string[]) {
  // Locale files may omit keys; `t()` falls back to English for missing entries.
  const merged = locale as I18n
  if (!window[GLOBAL_LOCALES_KEY]) window[GLOBAL_LOCALES_KEY] = {}
  window[GLOBAL_LOCALES_KEY][lang] = merged
  if (aliases)
    aliases.forEach((l) => {
      window[GLOBAL_LOCALES_KEY][l] = merged
    })
  return merged
}

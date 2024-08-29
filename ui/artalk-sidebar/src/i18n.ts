import { createI18n, type I18n, type Locale } from 'vue-i18n'
import { en } from './i18n-en'

export type MessageSchema = typeof en

export function setupI18n() {
  const i18n = createI18n({
    legacy: false, // use i18n in Composition API
    locale: 'en',
    fallbackLocale: 'en',
    messages: { en } as any,
  })

  const setLocale = async (value: string) => {
    await loadLocaleMessages(i18n, value)
    i18n.global.locale.value = value
  }

  return { i18n, setLocale }
}

export async function loadLocaleMessages(i18n: I18n, locale: Locale) {
  if (i18n.global.availableLocales.includes(locale)) return

  // Load locale messages with dynamic import
  // @see https://vitejs.dev/guide/features#dynamic-import
  const messages = await import(`./i18n/${locale}.ts`)
    .then((r: any) => r.default || r)
    .catch(() => {
      console.error(`Failed to load locale messages for "${locale}"`)
      return
    })

  // Set locale and locale message
  i18n.global.setLocaleMessage(locale, messages)
  return nextTick()
}

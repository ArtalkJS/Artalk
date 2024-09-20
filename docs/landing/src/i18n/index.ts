import i18n, { InitOptions } from 'i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import { initReactI18next } from 'react-i18next'
import { en } from './en'
import { zhCN } from './zh-CN'
import { zhTW } from './zh-TW'
import { fr } from './fr'
import { ja } from './ja'
import { ko } from './ko'
import { ru } from './ru'

export type MessageSchema = typeof en
export const DefaultNameSpace = 'translation'

const i18nConfig: InitOptions = {
  defaultNS: DefaultNameSpace,
  resources: {
    en: { translation: en },
    zh: { translation: zhCN },
    'zh-CN': { translation: zhCN },
    'zh-TW': { translation: zhTW },
    'zh-HK': { translation: zhTW },
    fr: { translation: fr },
    ja: { translation: ja },
    ko: { translation: ko },
    ru: { translation: ru },
  },
  fallbackLng: 'en',

  interpolation: {
    escapeValue: false, // react already safes from xss => https://www.i18next.com/translation-function/interpolation#unescape
  },
}

export function initI18n() {
  i18n.use(initReactI18next).use(LanguageDetector).init(i18nConfig)
}

export function initI18nSSR(locale: string) {
  return i18n.use(initReactI18next).init({
    ...i18nConfig,
    lng: locale,
    interpolation: {
      escapeValue: true,
    },
  })
}

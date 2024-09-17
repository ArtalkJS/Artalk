import 'i18next'
import type { MessageSchema } from './i18n'

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'translation'
    resources: {
      translation: MessageSchema
    }
  }
}

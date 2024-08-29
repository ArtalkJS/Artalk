import { DefineLocaleMessage } from 'vue-i18n'
import { MessageSchema } from './i18n'

declare module 'vue-i18n' {
  // define the locale messages schema
  export interface DefineLocaleMessage extends MessageSchema {}
}

import { ListData, CommentData } from '~/types/artalk-data'

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'list-load': ListData
  'list-error': string
  'list-loading': undefined
  'list-clear': undefined
  'list-refresh-ui': undefined
  'list-import': CommentData[]
  'list-insert': CommentData
  'sidebar-show': undefined
  'sidebar-hide': undefined
  'check-admin-show-el': undefined
  'editor-open-comment': undefined
  'editor-close-comment': undefined
  'editor-show-loading': undefined
  'editor-hide-loading': undefined
  'editor-notify': NotifyConf
  'editor-reply': CommentData
  'user-changed': undefined
  'checker-admin': CheckerConf
  'checker-captcha': CheckerCaptchaConf
}

export interface CheckerConf {
  onSuccess?: (inputVal: string, dialogEl?: HTMLElement) => void
  onMount?: (dialogEl: HTMLElement) => void
  onCancel?: () => void
}

export interface NotifyConf {
  msg: string
  type: any
}

export interface CheckerCaptchaConf extends CheckerConf {
  imgData?: string
}

// ============================================
export interface Listener<P> {
  (payload: P): void
}

export interface Event<K extends keyof EventPayloadMap = keyof EventPayloadMap> {
  name: keyof EventPayloadMap
  listener: Listener<EventPayloadMap[K]>
}

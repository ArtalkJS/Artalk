import { ListData, CommentData, NotifyData } from '~/types/artalk-data'

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'list-loading': undefined
  'list-clear': undefined
  'list-refresh-ui': undefined
  'list-import': CommentData[]
  'list-insert': CommentData
  'list-comment-del': CommentData
  'list-reload': undefined
  'list-update-data': (data: ListData | undefined) => void
  'sidebar-show'?: SidebarShowPayload
  'sidebar-hide': undefined
  'check-admin-show-el': undefined
  'editor-open-comment': undefined
  'editor-close-comment': undefined
  'editor-show-loading': undefined
  'editor-hide-loading': undefined
  'editor-notify': NotifyConf
  'editor-reply': CommentData
  'user-changed': undefined
  'unread-update': UnreadUpdatePayload
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

export interface SidebarShowPayload {
  viewName?: string
}

export interface UnreadUpdatePayload {
  notifies: NotifyData[]
}

// ============================================
export interface Listener<P> {
  (payload: P): void
}

export interface Event<K extends keyof EventPayloadMap = keyof EventPayloadMap> {
  name: keyof EventPayloadMap
  listener: Listener<EventPayloadMap[K]>
}

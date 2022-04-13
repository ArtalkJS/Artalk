import { ListData, CommentData, NotifyData } from './artalk-data'
import { LocalUser } from './artalk-config'

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'comments-load': undefined       // 评论加载时
  'comments-loaded': undefined     // 评论装载后
  'editor-submit': undefined       // 编辑器提交时
  'editor-submitted': undefined    // 编辑器提交后
  'user-changed': LocalUser        // 本地用户数据变更时
  'conf-updated': undefined        // Artalk 配置变更时

  // List 操作（外部：不建议 listen，仅 trigger）
  'list-reload': undefined         // 重新加载 List
  'list-import': CommentData[]     // 评论导入
  'list-insert': CommentData       // 评论添加
  'list-delete': CommentData       // 评论删除
  'list-update': ListUpdatePayload // 更新评论数据
  'list-refresh-ui': undefined     // 刷新 List UI

  // Sidebar
  'sidebar-show': undefined           // 侧边栏显示
  'sidebar-hide': undefined           // 侧边栏隐藏

  // Editor
  'editor-open': undefined            // 打开评论
  'editor-close': undefined           // 关闭评论
  'editor-show-loading': undefined    // 加载显示
  'editor-hide-loading': undefined    // 加载隐藏
  'editor-notify': NotifyConf         // 显示提示
  'editor-reply': EditorReplyPayload  // 设置回复
  'editor-reply-cancel': undefined    // 取消回复
  'editor-travel': HTMLElement
  'editor-travel-back': undefined

  // Notify
  'unread-update': UnreadUpdatePayload // 未读数据更新

  // Checker
  'checker-admin': CheckerPayload           // 检查管理员
  'checker-captcha': CheckerCaptchaPayload  // 检查验证码
  'check-admin-show-el': undefined          // 检查并更新仅管理员显示的元素
}

export interface CheckerPayload {
  onSuccess?: (inputVal: string, dialogEl?: HTMLElement) => void
  onMount?: (dialogEl: HTMLElement) => void
  onCancel?: () => void
}

export interface NotifyConf {
  msg: string
  type: any
}

export interface CheckerCaptchaPayload extends CheckerPayload {
  imgData?: string
}

export interface UnreadUpdatePayload {
  notifies: NotifyData[]
}

export type ListUpdatePayload = (data: ListData | undefined) => void

export interface EditorReplyPayload {
  data: CommentData
  $el: HTMLElement
  scroll?: boolean
}

// ============================================
export interface Handler<P> {
  (payload: P): void
}

export type EventScopeType = 'internal'|'external'

export interface Event<K extends keyof EventPayloadMap = keyof EventPayloadMap> {
  name: keyof EventPayloadMap
  handler: Handler<EventPayloadMap[K]>
  scope: EventScopeType
}

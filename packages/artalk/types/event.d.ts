import { CommentData } from './artalk-data'
import { LocalUser } from './artalk-config'

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'list-load': undefined           // 评论加载时
  'list-loaded': undefined         // 评论装载后
  'list-inserted': CommentData     // 评论插入后
  'editor-submit': undefined       // 编辑器提交时
  'editor-submitted': undefined    // 编辑器提交后
  'user-changed': LocalUser        // 本地用户数据变更时
  'conf-loaded': undefined         // Artalk 配置变更时
  'sidebar-show': undefined        // 侧边栏显示
  'sidebar-hide': undefined        // 侧边栏隐藏
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

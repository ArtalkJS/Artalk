import type {
  CommentData,
  ListData,
  ListFetchParams,
  NotifyData,
  PageData,
  LocalUser,
  Config,
} from '.'
import type { CommentNode } from '@/comment'

export interface ListErrorData {
  msg: string
  data?: any
}
export interface ListFetchedArgs {
  params: Partial<ListFetchParams>
  data?: ListData
  error?: ListErrorData
}

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  // Basic lifecycle
  created: undefined
  mounted: undefined
  updated: Config
  unmounted: undefined

  'list-fetch': Partial<ListFetchParams> // 评论列表请求时
  'list-fetched': ListFetchedArgs // 评论列表请求后
  'list-load': CommentData[] // 评论装载前 (list-load payload is partial comments)

  'list-loaded': CommentData[] // 评论装载后 (list-loaded payload is full comments)
  'list-failed': ListErrorData // 评论加载错误时

  'list-goto-first': undefined // 评论列表归位时
  'list-reach-bottom': undefined // 评论列表滚动到底部时

  'comment-inserted': CommentData // 评论插入后
  'comment-updated': CommentData // 评论更新后
  'comment-deleted': CommentData // 评论删除后
  'comment-rendered': CommentNode // 评论节点渲染后
  'notifies-updated': NotifyData[] // 消息列表变更时
  'list-goto': number // 评论跳转时
  'page-loaded': PageData // 页面数据更新后
  'editor-submit': undefined // 编辑器提交时
  'editor-submitted': undefined // 编辑器提交后
  'user-changed': LocalUser // 本地用户数据变更时
  'sidebar-show': undefined // 侧边栏显示
  'sidebar-hide': undefined // 侧边栏隐藏
}

export type EventHandler<T> = (payload: T) => void

export interface Event<T, K extends keyof T = keyof T> extends EventOptions {
  name: K
  handler: EventHandler<T[K]>
}

export interface EventOptions {
  once?: boolean
}

export interface EventManager<T = EventPayloadMap> {
  on<K extends keyof T>(name: K, handler: EventHandler<T[K]>, opts?: EventOptions): void
  off<K extends keyof T>(name: K, handler: EventHandler<T[K]>): void
  trigger<K extends keyof T>(name: K, payload?: T[K]): void
}

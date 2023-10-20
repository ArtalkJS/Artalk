import type { TPgMode } from './paginator'
import type { SortByType } from './list-nest'
import type Comment from '../comment/comment'

export interface ListOptions {
  /** Lite mode */
  liteMode?: boolean

  /** Flat mode */
  flatMode?: boolean

  /** Pagination mode */
  pageMode?: TPgMode

  /** Page size */
  pageSize?: number

  /** 监听指定元素上的滚动 */
  scrollListenerAt?: HTMLElement

  /** 翻页归位到指定元素 */
  repositionAt?: HTMLElement

  /** 启用列表未读高亮 */
  unreadHighlight?: boolean

  /** Sort condition in nest mode */
  nestSortBy?: SortByType

  /** Text to show when no comment */
  noCommentText?: string

  // 一些 Hook 函数
  // ----------------
  renderComment?: (comment: Comment) => void
  paramsEditor?: (params: any) => void
}

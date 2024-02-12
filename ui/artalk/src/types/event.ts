import type { CommentNode } from '@/comment'
import { CommentData, ListData, ListFetchParams, NotifyData, PageData } from './data'
import { ArtalkConfig, LocalUser } from './config'

export interface ListErrorData { msg: string, data?: any }
export interface ListFetchedArgs { params: Partial<ListFetchParams>, data?: ListData, error?: ListErrorData }

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  // Basic lifecycle
  'created': undefined
  'mounted': undefined
  'updated': ArtalkConfig
  'unmounted': undefined

  'list-fetch': Partial<ListFetchParams>    // 评论列表请求时
  'list-fetched': ListFetchedArgs           // 评论列表请求后
  'list-load': CommentData[]     // 评论装载前 (list-load payload is partial comments)

  'list-loaded': CommentData[]     // 评论装载后 (list-loaded payload is full comments)
  'list-failed': ListErrorData          // 评论加载错误时

  'list-goto-first': undefined    // 评论列表归位时
  'list-reach-bottom': undefined  // 评论列表滚动到底部时

  'comment-inserted': CommentData  // 评论插入后
  'comment-updated': CommentData   // 评论更新后
  'comment-deleted': CommentData   // 评论删除后
  'comment-rendered': CommentNode  // 评论节点渲染后
  'notifies-updated': NotifyData[] // 消息列表变更时
  'list-goto': number              // 评论跳转时
  'page-loaded': PageData          // 页面数据更新后
  'editor-submit': undefined       // 编辑器提交时
  'editor-submitted': undefined    // 编辑器提交后
  'user-changed': LocalUser        // 本地用户数据变更时
  'sidebar-show': undefined        // 侧边栏显示
  'sidebar-hide': undefined        // 侧边栏隐藏
}

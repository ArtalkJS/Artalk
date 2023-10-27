import type Comment from '~/src/comment'
import { CommentData, ListData, ListFetchParams, NotifyData, PageData } from './data'
import { ArtalkConfig, LocalUser } from './config'

interface ErrorData { msg: string, data?: any }
interface ListFetchedArgs { params: Partial<ListFetchParams>, data?: ListData, error?: ErrorData }

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'inited': undefined              // Artalk 初始化后
  'destroy': undefined             // Artalk 销毁时
  'list-fetch': Partial<ListFetchParams>    // 评论列表请求时
  'list-fetched': ListFetchedArgs           // 评论列表请求后
  'list-load': CommentData[]     // 评论装载前

  // TODO merge 'list-loaded' and 'list-error', for purpose of `once` 'list-loaded' and simplify
  // example need bind 'list-loaded' and 'list-error' to same handler
  // consider name it to 'list-fetched': { data: CommentData[], error?: { msg: string, data?: any } }
  // and remove `list-load`
  'list-loaded': CommentData[]     // 评论装载后
  'list-error': ErrorData          // 评论加载错误时

  'list-goto-first': undefined    // 评论列表归位时
  'list-reach-bottom': undefined  // 评论列表滚动到底部时

  'comment-inserted': CommentData  // 评论插入后
  'comment-updated': CommentData   // 评论更新后
  'comment-deleted': CommentData   // 评论删除后
  'comment-rendered': Comment      // 评论节点渲染后
  'unreads-updated': NotifyData[]  // 未读消息变更时
  'list-goto': number              // 评论跳转时
  'page-loaded': PageData          // 页面数据更新后
  'editor-submit': undefined       // 编辑器提交时
  'editor-submitted': undefined    // 编辑器提交后
  'user-changed': LocalUser        // 本地用户数据变更时
  'conf-loaded': ArtalkConfig      // Artalk 配置变更时
  'sidebar-show': undefined        // 侧边栏显示
  'sidebar-hide': undefined        // 侧边栏隐藏
}

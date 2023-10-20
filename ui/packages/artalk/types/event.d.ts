import { CommentData, NotifyData, PageData } from './artalk-data'
import { LocalUser } from './artalk-config'

/** EventName to EventPayload Type */
export interface EventPayloadMap {
  'inited': undefined              // Artalk 初始化后
  'destroy': undefined             // Artalk 销毁时
  'list-load': undefined           // 评论加载时
  'list-loaded': undefined         // 评论装载后
  'list-inserted': CommentData     // 评论插入后
  'list-deleted': CommentData      // 评论删除后
  'list-error': { msg: string, data?: any } // 评论加载错误时
  'list-goto': number              // 评论跳转时
  'page-loaded': PageData          // 页面数据更新后
  'editor-submit': undefined       // 编辑器提交时
  'editor-submitted': undefined    // 编辑器提交后
  'user-changed': LocalUser        // 本地用户数据变更时
  'conf-loaded': undefined         // Artalk 配置变更时
  'unread-updated': NotifyData[]   // 未读消息变更时
  'sidebar-show': undefined        // 侧边栏显示
  'sidebar-hide': undefined        // 侧边栏隐藏
}

export interface ArtalkConfig {
  /**
   * 装载元素
   */
  el: string

  /**
   * 评论框占位字符
   */
  placeholder?: string

  /**
   * 评论为空时显示字符
   */
  noComment?: string

  /**
   * 默认头像（URL or Gravatar Type）
   */
  defaultAvatar?: string

  /**
   * 页面唯一标识
   */
  pageKey: string

  /**
   * 服务器地址
   */
  serverUrl: string

  /**
   * 表情包
   */
  emoticons?: object|any

  /**
   * 头像
   */
  gravatar?: {
    /** CDN 地址 */
    cdn?: string
  }
}

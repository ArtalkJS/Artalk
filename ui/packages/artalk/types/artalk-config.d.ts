import { I18n } from '~/src/i18n'
import { CommentData } from './artalk-data'

export default interface ArtalkConfig {
  /** 装载元素 */
  el: string|HTMLElement

  /** 页面唯一标识（完整 URL） */
  pageKey: string

  /** 页面标题 */
  pageTitle: string

  /** 服务器地址 */
  server: string

  /** 站点名 */
  site: string

  /** 评论框占位字符 */
  placeholder: string

  /** 评论为空时显示字符 */
  noComment: string

  /** 发送按钮文字 */
  sendBtn: string

  /** 评论框穿梭（显示在待回复评论后面） */
  editorTravel: boolean

  /** 表情包 */
  emoticons: object|any[]|string|false

  /** Gravatar 头像 */
  gravatar: {
    /** API 地址 */
    mirror: string
    /** API 参数 */
    params: string
  }

  /** 头像链接生成器 */
  avatarURLBuilder?: (comment: CommentData) => string,

  /** 分页配置 */
  pagination: {
    /** 每次请求获取数量 */
    pageSize: number
    /** 阅读更多模式 */
    readMore: boolean
    /** 滚动到底部自动加载 */
    autoLoad: boolean
  }

  /** 内容限高 */
  heightLimit: {
    /** 评论内容限高 */
    content: number

    /** 子评论区域限高 */
    children: number

    /** 滚动限高 */
    scrollable: boolean
  }

  /** 评论投票按钮 */
  vote: boolean

  /** 评论投票反对按钮 */
  voteDown: boolean

  /** 评论预览功能 */
  preview: boolean

  /** 评论数绑定元素 Selector */
  countEl: string

  /** PV 数绑定元素 Selector */
  pvEl: string

  /** 夜间模式 */
  darkMode: boolean|'auto'

  /** 请求超时（单位：秒） */
  reqTimeout: number

  /** 平铺模式 */
  flatMode: boolean|'auto'

  /** 嵌套模式 · 最大层数 */
  nestMax: number

  /** 嵌套模式 · 排序方式 */
  nestSort: 'DATE_ASC'|'DATE_DESC'

  /** 显示 UA 徽标 */
  uaBadge: boolean

  /** 评论列表排序功能 (显示 Dropdown) */
  listSort: boolean

  /** 图片上传功能 */
  imgUpload: boolean

  /** 图片上传器 */
  imgUploader?: (file: File) => Promise<string>

  /** 版本检测 */
  versionCheck: boolean

  /** 引用后端配置 */
  useBackendConf: boolean

  /** 语言本地化 */
  locale: I18n|string
}

/**
 * 本地持久化用户数据
 * @note 始终保持一层结构，不支持多层结构
 */
export interface LocalUser {
  /** 昵称 */
  nick: string

  /** 邮箱 */
  email: string

  /** 链接 */
  link: string

  /** TOKEN */
  token: string

  /** 是否为管理员 */
  isAdmin: boolean
}

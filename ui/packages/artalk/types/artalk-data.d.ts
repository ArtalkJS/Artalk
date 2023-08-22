export interface CommentData {
  /** 评论 ID */
  id: number

  /** 评论正文 */
  content: string

  /** 用户昵称 */
  nick: string

  /** 用户邮箱 (该字段仅管理员可见) */
  email?: string

  /** 用户邮箱（已加密） */
  email_encrypted: string

  /** 用户链接 */
  link: string

  /** 回复目标评论 ID */
  rid: number

  /** User Agent */
  ua: string

  /** 评论日期 */
  date: string

  /** 是否折叠 */
  is_collapsed?: boolean

  /** 是否待审 */
  is_pending?: boolean

  /** 是否置顶 */
  is_pinned?: boolean

  /** 是否为管理员 */
  is_admin?: boolean

  /** 徽章文字 */
  badge_name?: string

  /** 徽章颜色 */
  badge_color?: string

  /** IP 归属地 */
  ip_region?: string

  /** 是否允许回复 */
  is_allow_reply?: boolean

  /** 评论页面 key */
  page_key: string

  /** 评论页面 url */
  page_url?: string

  /** 是否可见 */
  visible: boolean

  /** 站点名（用于隔离） */
  site_name: string

  /** 赞同数 */
  vote_up: number

  /** 反对数 */
  vote_down: number
}

export interface ListData {
  /** 评论数据 */
  comments: CommentData[]

  /** 根评论总数 */
  total_roots: number

  /** 评论总数（包括所有子评论） */
  total: number

  /** 页面信息 */
  page: PageData

  /** 站点信息 */
  site: SiteData

  /** 未读的 Notifies */
  unread: NotifyData[]

  /** 未读计数 */
  unread_count: number

  /** API 版本 */
  api_version: ApiVersionData

  /** 后端配置 */
  conf: any // TODO
}

export interface PageData {
  /** 页面 ID */
  id: number

  /** 页面唯一标识符 */
  key: string

  /** 页面标题 */
  title: string

  /** 页面 url */
  url: string

  /** 仅管理员可评 */
  admin_only: boolean

  /** 站点名（用于隔离） */
  site_name: string

  /** 赞同数 */
  vote_up: number

  /** 反对数 */
  vote_down: number
}

export interface SiteData {
  /** 站点 ID */
  id: number

  /** 站点名 */
  name: string

  /** 站点 URLs */
  urls: string[]

  /** 站点 URLs（原始字符串） */
  urls_raw: string

  /** 站点主 URL */
  first_url: string
}

export interface UserData {
  /** 用户 ID */
  id: number

  /** 用户名 */
  name: string

  /** 邮箱 */
  email: string

  /** 链接 */
  link: string

  /** 徽章名称 */
  badge_name: string

  /** 徽章颜色 */
  badge_color: string

  /** 是否属于管理员 */
  is_admin: boolean

  /** 所属站点名 */
  site_names: string[]

  /** 所属站点名 (字符串) */
  site_names_raw: string,

  /** 是否允许接收邮件通知 */
  receive_email: boolean
}

export interface UserDataForAdmin extends UserData {
  /** 最后一次操作 IP */
  last_ip: string

  /** 最后一次操作 UA */
  last_ua: string

  /** 是否存在于配置文件中 */
  is_in_conf: boolean

  /** 评论数 */
  comment_count: number

  /** 密码 */
  password: string
}

export interface NotifyData {
  /** 通知 ID */
  id: number

  /** 用户 ID */
  user_id: number

  /** 评论 ID */
  comment_id: number

  /** 是否已读 */
  is_read: boolean

  /** 是否已发送邮件通知 */
  is_emailed: boolean

  /** 标为已读地址 */
  read_link: string
}

export type EmoticonGrpData = {
  name: string,
  type: 'emoticon'|'image'|'emoji'
  items: {key: string, val: string}[]
}

export type EmoticonListData = EmoticonGrpData[]

export interface ApiVersionData {
  /** API 程序名 */
  app: string

  /** API 程序版本号 */
  version: string

  /** API 程序 CommitHash */
  commit_hash: string
}

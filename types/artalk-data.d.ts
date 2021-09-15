export interface CommentData {
  /** 评论 ID */
  id: number

  /** 评论正文 */
  content: string

  /** 用户昵称 */
  nick: string

  /** 用户邮箱 */
  email: string

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

  /** 是否为管理员 */
  is_admin?: boolean

  /** 徽章 */
  badge?: {
    name?: string
    color?: string
  }

  /** 是否允许回复 */
  is_allow_reply?: boolean
}

export interface ListData {
  /** 评论数据 */
  comments: CommentData[]

  /** 父级评论总数 */
  total_parents: number

  /** 评论总数（包括所有子评论） */
  total: number

  /** 页面信息 */
  page: {
    id: number
    page_key: string
    is_close_comment: boolean
  }
}

export interface UserData {
  /** 用户 ID */
  id: string

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

  /** 类型 */
  type: 'banned'|'admin'|''
}

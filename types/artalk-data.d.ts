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
}

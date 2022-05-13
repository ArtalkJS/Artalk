import ArtalkConfig, { LocalUser } from "~/types/artalk-config"
import Context from '~/types/context'

export default class User {
  public ctx: Context
  public data: LocalUser

  public constructor(ctx: Context) {
    this.ctx = ctx

    // 从 localStorage 导入
    const localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.data = {
      nick: localUser.nick || '',
      email: localUser.email || '',
      link: localUser.link || '',
      token: localUser.token || '',
      isAdmin: localUser.isAdmin || false
    }
  }

  /** 保存用户到 localStorage 中 */
  public update(obj: Partial<LocalUser> = {}) {
    Object.entries(obj).forEach(([key, value]) => {
      this.data[key] = value
    })

    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.data))
    this.ctx.trigger('user-changed', this.ctx.user.data)
  }

  /** 注销，清除用户登录状态 */
  public logout() {
    this.update({
      token: '',
      isAdmin: false
    })
  }

  /** 是否已填写基本用户信息 */
  public checkHasBasicUserInfo() {
    return !!this.data.nick && !!this.data.email
  }
}

import { ArtalkConfig } from "~/types/artalk-config"

export default class User {
  public data: {
    nick: string,
    email: string,
    link: string,
    token: string,
    isAdmin: boolean
  }

  constructor (conf: ArtalkConfig) {
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
  public save () {
    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.data))
  }

  /** 是否已填写基本用户信息 */
  public checkHasBasicUserInfo () {
    return !!this.data.nick && !!this.data.email
  }

  /** 根据请求数据判断 nick 是否为管理员 */
  public checkNickEmailIsAdmin (nick: string, email: string) {
    // TODO: checkNickEmailIsAdmin
    // if (!this.data || !this.data.admin_nicks || !this.data.admin_encrypted_emails) return false

    // return (this.data.admin_nicks.indexOf(nick) !== -1)
    //   && (this.data.admin_encrypted_emails.find(o => String(o).toLowerCase() === String(md5(email)).toLowerCase()))
  }
}

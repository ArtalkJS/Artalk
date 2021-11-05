import ArtalkConfig, { LocalUser } from "~/types/artalk-config"

export default class User {
  public data: LocalUser

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
}

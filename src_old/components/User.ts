import Artalk from '../Artalk';
import ArtalkContext from '../ArtalkContext';

export default class User extends ArtalkContext {
  public data: {
    nick: string|null,
    email: string|null,
    link: string|null,
    password: string|null,
    isAdmin: boolean
  }

  constructor (artalk: Artalk) {
    super(artalk)

    // 从 localStorage 导入
    const localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.data = {
      nick: localUser.nick || '',
      email: localUser.email || '',
      link: localUser.link || '',
      password: localUser.password || '',
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

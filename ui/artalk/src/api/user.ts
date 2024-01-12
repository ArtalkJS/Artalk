import type { UserData, UserDataForAdmin, UserInfoApiResponseData } from '@/types'
import ApiBase from './_base'

/**
 * 用户 API
 */
export default class UserApi extends ApiBase {
  /** 用户 · 登录 */
  public async login(name: string, email: string, password: string) {
    const params = { name, email, password }

    const data = await this.fetch<{
      token: string,
      user: UserData
    }>('POST', '/user/access_token', params)

    return data
  }

  /** 用户 · 获取  */
  public userGet(name: string, email: string) {
    const ctrl = new AbortController()
    const params: any = {
      name, email
    }

    const req = this.fetch<UserInfoApiResponseData>('GET', `/user/info`, params, {
      signal: ctrl.signal,
    })

    return {
      req,
      abort: () => { ctrl.abort() },
    }
  }

  /** 用户 · 登录状态 */
  public async loginStatus() {
    const data = await this.fetch<{
      is_login: boolean,
      is_admin: boolean
    }>('GET', '/user/status', this.withUserInfo({}))
    return data
  }

  /** 用户 · 列表 */
  public async userList(offset?: number, limit?: number, type?: 'all'|'admin'|'in_conf') {
    const params: any = {
      offset: offset || 0,
      limit: limit || 15,
    }

    let path = '/users'
    if (type) path += `/${type}`

    const d = await this.fetch<{
      users: UserDataForAdmin[],
      total: number
    }>('GET', path, params)

    return d
  }

  /** 用户 · 新增 */
  public async userAdd(user: Partial<UserDataForAdmin>, password?: string) {
    const params: any = {
      name: user.name || '',
      email: user.email || '',
      password: password || '',
      link: user.link || '',
      is_admin: user.is_admin || false,
      receive_email: user.receive_email || true,
      badge_name: user.badge_name || '',
      badge_color: user.badge_color || '',
    }

    const d = await this.fetch<UserDataForAdmin>('POST', '/users', params)
    return d
  }

  /** 用户 · 修改 */
  public async userEdit(id: number, user: Partial<UserDataForAdmin>, password?: string) {
    const params: any = {
      name: user.name || '',
      email: user.email || '',
      password: password || '',
      link: user.link || '',
      is_admin: user.is_admin || false,
      receive_email: user.receive_email || true,
      badge_name: user.badge_name || '',
      badge_color: user.badge_color || '',
    }

    const d = await this.fetch<UserDataForAdmin>('PUT', `/users/${id}`, params)

    return d
  }

  /** 用户 · 删除 */
  public userDel(id: number) {
    return this.fetch('DELETE', `/users/${id}`)
  }
}

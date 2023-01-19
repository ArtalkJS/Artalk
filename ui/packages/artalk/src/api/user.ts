import { UserData, NotifyData, UserDataForAdmin } from '~/types/artalk-data'
import ApiBase from './api-base'
import { ToFormData } from './request'

/**
 * 用户 API
 */
export default class UserApi extends ApiBase {
  /** 用户 · 登录 */
  public async login(name: string, email: string, password: string) {
    const params: any = {
      name, email, password
    }

    if (this.ctx.conf.site) params.site_name = this.ctx.conf.site

    const data = await this.POST<any>('/login', params)
    return data as { token: string, user: UserData }
  }

  /** 用户 · 获取  */
  public userGet(name: string, email: string) {
    const ctrl = new AbortController()
    const params: any = {
      name, email, site_name: (this.ctx.conf.site || '')
    }

    const req = this.Fetch(`/user-get`, {
      method: 'POST',
      body: ToFormData(params),
      signal: ctrl.signal,
    }).then((json) => ({
      user: json.data.user as UserData|null,
      is_login: json.data.is_login as boolean,
      unread: (json.data.unread || []) as NotifyData[],
      unread_count: json.data.unread_count || 0,
    }))

    return {
      req,
      abort: () => { ctrl.abort() },
    }
  }

  /** 用户 · 登录状态 */
  public async loginStatus() {
    const data = await this.POST<any>('/login-status', {
      name: this.ctx.user.data.nick,
      email: this.ctx.user.data.email
    })
    return (data || { is_login: false, is_admin: false }) as { is_login: boolean, is_admin: boolean }
  }

  /** 用户 · 注销 */
  public async logout() {
    return this.POST('/logout')
  }

  /** 已读标记 */
  public markRead(notifyKey: string, readAll = false) {
    const params: any = {
      site_name: this.ctx.conf.site || '',
      notify_key: notifyKey,
    }

    if (readAll) {
      delete params.notify_key
      params.read_all = true
      params.name = this.ctx.user.data.nick
      params.email = this.ctx.user.data.email
    }

    return this.POST(`/mark-read`, params)
  }

  /** 用户 · 列表 */
  public async userList(offset?: number, limit?: number, type?: 'all'|'admin'|'in_conf') {
    const params: any = {
      offset: offset || 0,
      limit: limit || 15,
    }

    if (type) params.type = type

    const d = await this.POST<any>('/admin/user-get', params)
    return (d as { users: UserDataForAdmin[], total: number })
  }

  /** 用户 · 新增 */
  public async userAdd(user: Partial<UserDataForAdmin>, password?: string) {
    const params: any = {
      name: user.name || '',
      email: user.email || '',
      password: password || '',
      link: user.link || '',
      is_admin: user.is_admin || false,
      site_names: user.site_names_raw || '',
      receive_email: user.receive_email || true,
      badge_name: user.badge_name || '',
      badge_color: user.badge_color || '',
    }

    const d = await this.POST<any>('/admin/user-add', params)
    return (d.user as UserDataForAdmin)
  }

  /** 用户 · 修改 */
  public async userEdit(user: Partial<UserDataForAdmin>, password?: string) {
    const params: any = {
      id: user.id,
      name: user.name || '',
      email: user.email || '',
      password: password || '',
      link: user.link || '',
      is_admin: user.is_admin || false,
      site_names: user.site_names_raw || '',
      receive_email: user.receive_email || true,
      badge_name: user.badge_name || '',
      badge_color: user.badge_color || '',
    }

    const d = await this.POST<any>('/admin/user-edit', params)
    return (d.user as UserDataForAdmin)
  }

  /** 用户 · 删除 */
  public userDel(userID: number) {
    return this.POST('/admin/user-del', {
      id: String(userID)
    })
  }
}

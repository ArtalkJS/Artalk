import { UserData, NotifyData } from '~/types/artalk-data'
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
    return (data.token as string)
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
}

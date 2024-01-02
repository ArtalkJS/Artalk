import { ApiOptions } from './_options'
import { Fetch } from './_request'

abstract class ApiBase {
  constructor(
    protected options: ApiOptions
  ) {}

  protected async fetch<T>(
    method: 'GET'|'POST'|'PUT'|'DELETE',
    path: string,
    payload?: any,
    opts?: RequestInit & { timeout?: number }
  ): Promise<any> {
    let url = this.options.baseURL + path
    let init: RequestInit = { method }
    if (method === 'POST' || method === 'PUT') init.body = JSON.stringify(payload)
    else if (payload) url = `${url}?${new URLSearchParams(payload)}`
    init = { ...init, ...opts }
    return Fetch<T>(this.options, url, init, opts?.timeout)
  }

  /**
   * Carry user info to request params
   *
   * @param params Request params
   * @returns Request params with user info
   */
  protected withUserInfo(params: any) {
    const user = this.options.userInfo
    if (user?.name) params.name = user.name
    if (user?.email) params.email = user.email
    return params
  }
}

export default ApiBase

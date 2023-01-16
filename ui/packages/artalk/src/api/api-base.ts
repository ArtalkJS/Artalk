import Context from '~/types/context'
import Api from '.'
import { POST, Fetch } from './request'

abstract class ApiBase {
  protected api: Api
  protected ctx: Context

  constructor(api: Api, ctx: Context) {
    this.api = api
    this.ctx = ctx
  }

  protected async POST<T>(path: string, data?: {[key: string]: any}) {
    return POST<T>(this.ctx, this.api.baseURL+path, data)
  }

  protected async Fetch(path: string, init: RequestInit, timeout?: number): Promise<any> {
    return Fetch(this.ctx, this.api.baseURL+path, init, timeout)
  }
}

interface ApiBase {

}

export default ApiBase

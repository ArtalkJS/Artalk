import type { ContextApi } from '~/types'
import Api from '.'
import { POST, Fetch } from './request'

abstract class ApiBase {
  protected api: Api
  protected ctx: ContextApi

  constructor(api: Api, ctx: ContextApi) {
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

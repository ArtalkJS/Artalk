import { ApiOptions } from './_options'
import { POST, Fetch } from './_request'

abstract class ApiBase {
  constructor(
    protected options: ApiOptions
  ) {}

  protected async POST<T>(path: string, data?: {[key: string]: any}) {
    return POST<T>(this.options, this.options.baseURL+path, data)
  }

  protected async Fetch(path: string, init: RequestInit, timeout?: number): Promise<any> {
    return Fetch(this.options, this.options.baseURL+path, init, timeout)
  }
}

export default ApiBase

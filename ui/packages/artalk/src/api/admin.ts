import ApiBase from './api-base'

/**
 * 管理员 API
 */
export default class AdminApi extends ApiBase {
  /** 缓存清除 */
  public cacheFlushAll() {
    const params: any = { flush_all: true }
    return this.POST('/admin/cache-flush', params)
  }

  /** 缓存预热 */
  public cacheWarmUp() {
    const params: any = {}
    return this.POST('/admin/cache-warm', params)
  }
}

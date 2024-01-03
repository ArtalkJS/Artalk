import ApiBase from './_base'

/**
 * 管理员 API
 */
export default class AdminApi extends ApiBase {
  /** 缓存清除 */
  public cacheFlushAll() {
    return this.fetch('POST', '/cache/flush')
  }

  /** 缓存预热 */
  public cacheWarmUp() {
    return this.fetch('POST', '/cache/warm_up')
  }
}

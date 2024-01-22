import { NotifyData } from '@/types'
import ApiBase from './_base'

/**
 * 通知 API
 */
export default class Notify extends ApiBase {
  /** 已读标记 */
  public markRead(commentID: number, notifyKey: string) {
    return this.fetch('POST', `/notifies/${commentID}/${notifyKey}/read`)
  }

  /** 全部标为已读 */
  public markAllRead() {
    return this.fetch('POST', `/notifies/read`, this.withUserInfo({}))
  }

  /** 获取未读通知 */
  public getNotifies() {
    return this.fetch<{
      notifies: NotifyData[],
      count: number,
    }>('GET', `/notifies`, this.withUserInfo({}))
  }
}

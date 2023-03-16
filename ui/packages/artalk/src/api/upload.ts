import { ToFormData } from './request'
import ApiBase from './api-base'
import User from '../lib/user'

/**
 * 上传 API
 */
export default class UploadApi extends ApiBase {
  /** 图片上传 */
  public async imgUpload(file: File) {
    const params: any = {
      name: User.data.nick,
      email: User.data.email,
      page_key: this.ctx.conf.pageKey,
    }

    const form = ToFormData(params)
    form.set('file', file)

    const init: RequestInit = {
      method: 'POST',
      body: form
    }

    const json = await this.Fetch('/img-upload', init)
    return ((json.data || {}) as any) as { img_file: string, img_url: string }
  }
}

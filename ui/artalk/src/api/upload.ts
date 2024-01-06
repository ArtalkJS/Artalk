import { ToFormData } from './_request'
import ApiBase from './_base'

/**
 * 上传 API
 */
export default class UploadApi extends ApiBase {
  /** 图片上传 */
  public async imgUpload(file: File) {
    const params: any = {
      page_key: this.options.pageKey,
    }

    this.withUserInfo(params)

    const form = ToFormData(params)
    form.set('file', file)

    const init: RequestInit = {
      method: 'POST',
      body: form
    }

    const json = await this.fetch<{
      img_file: string,
      img_url: string
    }>('POST', '/upload', init)
    return json
  }
}

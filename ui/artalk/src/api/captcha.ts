import ApiBase from './_base'

/**
 * 验证码 API
 */
export default class CaptchaApi extends ApiBase {
  /** 验证码 · 获取 */
  public async captchaGet() {
    const data = await this.fetch<{
      img_data: string
    }>('GET', '/captcha/get')
    return data
  }

  /** 验证码 · 检验 */
  public async captchaVerify(value: string) {
    const data = await this.fetch<{}>('POST', '/captcha/verify', { value })
    return data
  }

  /** 验证码 · 状态 */
  public async captchaStatus() {
    const data = await this.fetch<{
      is_pass: boolean
    }>('GET', '/captcha/status')
    return data
  }
}

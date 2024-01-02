import ApiBase from './_base'

/**
 * 验证码 API
 */
export default class CaptchaApi extends ApiBase {
  /** 验证码 · 获取 */
  public async captchaGet() {
    const data = await this.fetch<any>('GET', '/captcha/refresh')
    return (data.img_data || '') as string
  }

  /** 验证码 · 检验 */
  public async captchaVerify(value: string) {
    const data = await this.fetch<any>('POST', '/captcha/verify', { value })
    return (data.img_data || '') as string
  }

  /** 验证码 · 状态 */
  public async captchaStatus() {
    const data = await this.fetch<any>('GET', '/captcha/status')
    return (data || { is_pass: false }) as { is_pass: boolean }
  }
}

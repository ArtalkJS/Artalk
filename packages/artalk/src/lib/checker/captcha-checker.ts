import Api from '@/api'
import * as Utils from '../utils'
import { Checker } from '.'

const CaptchaChecker: Checker = {
  request(that, inputVal) {
    return new Api(that.ctx).captchaCheck(inputVal)
  },

  body(that) {
    const elem = Utils.createElement(
      `<span><img class="atk-captcha-img" src="${that.submitCaptchaImgData || ''}" alt="验证码">敲入验证码继续：</span>`
    );

    // 刷新验证码
    elem.querySelector<HTMLElement>('.atk-captcha-img')!.onclick = () => {
      const imgEl = elem.querySelector('.atk-captcha-img')
      new Api(that.ctx).captchaGet()
        .then((imgData) => {
          imgEl!.setAttribute('src', imgData)
        })
        .catch((err) => {
          console.error('验证码获取失败 ', err)
        })
    }
    return elem
  },

  onSuccess(that, data, inputVal, formEl) {
    that.submitCaptchaVal = inputVal
  },

  onError(that, err, inputVal, formEl) {
    formEl.querySelector<HTMLElement>('.atk-captcha-img')!.click() // 刷新验证码
  }
}

export default CaptchaChecker

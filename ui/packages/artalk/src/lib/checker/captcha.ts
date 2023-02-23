
import { Checker } from '.'
import * as render from './captcha-renders'

const CaptchaChecker: Checker = {
  request(checker, inputVal) {
    return checker.getApi().captcha.captchaCheck(inputVal)
  },

  body(checker) {
    if (checker.get('iframe')) return render.iframeBody(checker)
    return render.imgBody(checker)
  },

  onSuccess(checker, data, inputVal, formEl) {
    checker.set('val', inputVal)
  },

  onError(checker, err, inputVal, formEl) {
    formEl.querySelector<HTMLElement>('.atk-captcha-img')!.click() // 刷新验证码
    formEl.querySelector<HTMLInputElement>('input[type="text"]')!.value = '' // 清空输入框输入
  }
}

export default CaptchaChecker


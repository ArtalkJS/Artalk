import * as Utils from '../utils'
import * as Ui from '../ui'
import { Checker } from '.'

const CaptchaChecker: Checker = {
  request(that, ctx, inputVal) {
    return that.ctx.getApi().captcha.captchaCheck(inputVal)
  },

  body(that, ctx) {
    // Iframe 验证方式
    if (that.captchaConf.iframe) {
      const $iframeWrap = Utils.createElement(`<div class="atk-checker-iframe-wrap"></div>`)
      const $iframe = Utils.createElement<HTMLIFrameElement>(`<iframe class="atk-fade-in"></iframe>`)
      $iframe.style.display = 'none'
      Ui.showLoading($iframeWrap, { transparentBg: true })
      $iframe.src = `${that.ctx.conf.server}/api/captcha/get?t=${+new Date()}`
      $iframe.onload = () => {
        $iframe.style.display = ''
        Ui.hideLoading($iframeWrap)
      }
      $iframeWrap.append($iframe)

      const $closeBtn = Utils.createElement(`<div class="atk-close-btn"><i class="atk-icon atk-icon-close"></i></div>`)
      $iframeWrap.append($closeBtn)

      ctx.hideInteractInput()

      // 轮询状态
      let stop = false // 打断
      const sleep = (ms: number) => new Promise((resolve) => { window.setTimeout(() => { resolve(null) }, ms) })
      ;(async function queryStatus() {
        await sleep(1000)
        if (stop) return
        let isPass = false
        try {
          const resp = await that.ctx.getApi().captcha.captchaStatus()
          isPass = resp.is_pass
        } catch { isPass = false }
        if (isPass) {
          ctx.triggerSuccess()
        } else {
          queryStatus()
        }
      })()

      $closeBtn.onclick = () => {
        stop = true
        ctx.cancel()
      }

      return $iframeWrap
    }

    // 图片验证方式
    const elem = Utils.createElement(
      `<span><img class="atk-captcha-img" src="${that.captchaConf.imgData || ''}">${that.ctx.$t('captchaCheck')}</span>`
    );

    // 刷新验证码
    elem.querySelector<HTMLElement>('.atk-captcha-img')!.onclick = () => {
      const imgEl = elem.querySelector('.atk-captcha-img')
      that.ctx.getApi().captcha.captchaGet()
        .then((imgData) => {
          imgEl!.setAttribute('src', imgData)
        })
        .catch((err) => {
          console.error('Failed to get captcha image ', err)
        })
    }
    return elem
  },

  onSuccess(that, ctx, data, inputVal, formEl) {
    that.captchaConf.val = inputVal
  },

  onError(that, ctx, err, inputVal, formEl) {
    formEl.querySelector<HTMLElement>('.atk-captcha-img')!.click() // 刷新验证码
    formEl.querySelector<HTMLInputElement>('input[type="text"]')!.value = '' // 清空输入框输入
  }
}

export default CaptchaChecker

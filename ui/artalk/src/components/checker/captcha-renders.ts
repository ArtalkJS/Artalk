import type { CheckerCtx } from '.'
import $t from '@/i18n'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'

/** 图片验证码 */
export function imgBody(checker: CheckerCtx) {
  // 图片验证方式
  const elem = Utils.createElement(
    `<span><img class="atk-captcha-img" src="${checker.get('img_data') || ''}">${$t('captchaCheck')}</span>`,
  )

  // 刷新验证码
  elem.querySelector<HTMLElement>('.atk-captcha-img')!.onclick = () => {
    const imgEl = elem.querySelector('.atk-captcha-img')
    checker
      .getApi()
      .captcha.getCaptcha()
      .then((res) => {
        if (checker.isDisposed?.()) return
        imgEl!.setAttribute('src', res.data.img_data)
      })
      .catch((err) => {
        console.error('Failed to get captcha image ', err)
      })
  }
  return elem
}

/** iframe 形式的通用验证服务 */
export function iframeBody(checker: CheckerCtx) {
  const $iframeWrap = Utils.createElement(`<div class="atk-checker-iframe-wrap"></div>`)
  const $iframe = Utils.createElement<HTMLIFrameElement>(
    `<iframe class="atk-fade-in" referrerpolicy="strict-origin-when-cross-origin"></iframe>`,
  )
  $iframe.style.display = 'none'
  Ui.showLoading($iframeWrap, { transparentBg: true })
  $iframe.src = checker.getOpts().getCaptchaIframeURL()
  $iframe.onload = () => {
    $iframe.style.display = ''
    Ui.hideLoading($iframeWrap)
  }
  $iframeWrap.append($iframe)

  const $closeBtn = Utils.createElement(
    `<div class="atk-close-btn"><i class="atk-icon atk-icon-close"></i></div>`,
  )
  $iframeWrap.append($closeBtn)

  checker.hideInteractInput()

  let stopped = false
  let timer: number
  const stop = () => {
    stopped = true
    window.clearTimeout(timer)
  }
  checker.onDispose?.(stop)
  async function queryStatus() {
    if (stopped) return
    let isPass: boolean
    try {
      const resp = await checker.getApi().captcha.getCaptchaStatus()
      isPass = resp.data.is_pass
    } catch {
      isPass = false
    }
    if (stopped) return
    if (isPass) {
      stop()
      checker.triggerSuccess()
    } else {
      timer = window.setTimeout(queryStatus, 1000)
    }
  }
  timer = window.setTimeout(queryStatus, 1000)

  $closeBtn.onclick = () => {
    stop()
    checker.cancel()
  }

  return $iframeWrap
}

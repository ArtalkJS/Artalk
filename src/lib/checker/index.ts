import Context from '@/context'
import Layer from '@/components/layer'
import Dialog from '@/components/dialog'
import * as Utils from '../utils'
import * as Ui from '../ui'
import { CheckerPayload } from '~/types/event'
import CaptchaChecker from './captcha-checker'
import AdminChecker from './admin-checker'

/**
 * Checker 发射台
 */
export default class CheckerLauncher {
  public ctx: Context

  public launched: Checker[] = []
  public submitCaptchaVal?: string
  public submitCaptchaImgData?: string

  constructor(ctx: Context) {
    this.ctx = ctx

    this.initEventBind()
  }

  /** 初始化事件绑定 */
  public initEventBind() {
    this.ctx.on('checker-captcha', (conf) => {
      if (conf.imgData) {
        this.submitCaptchaImgData = conf.imgData
      }
      this.fire(CaptchaChecker, conf)
    })

    this.ctx.on('checker-admin', (conf) => {
      this.fire(AdminChecker, conf)
    })
  }

  public fire(checker: Checker, payload: CheckerPayload) {
    if (this.launched.includes(checker)) return // 阻止同时 fire 相同的 checker
    this.launched.push(checker)

    const layer = new Layer(this.ctx, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(this))

    const input = Utils.createElement<HTMLInputElement>(
      `<input id="check" type="${checker.inputType || 'text'}" autocomplete="off" required placeholder="">`
    )
    formEl.appendChild(input)
    setTimeout(() => input.focus(), 80) // 延迟 Focus

    input.onkeyup = (evt) => {
      if (evt.key === 'Enter' || evt.keyCode === 13) {
        // 按下回车键
        evt.preventDefault()
        layer
          .getEl()
          .querySelector<HTMLButtonElement>('button[data-action="confirm"]')!
          .click()
      }
    }

    let btnTextOrg: string | undefined
    const dialog = new Dialog(formEl)

    // 确认按钮
    dialog.setYes((btnEl) => {
      const inputVal = input.value.trim()
      if (!btnTextOrg) btnTextOrg = btnEl.innerText
      const btnTextSet = (btnText: string) => {
        btnEl.innerText = btnText
        btnEl.classList.add('error')
      }
      const btnTextRestore = () => {
        btnEl.innerText = btnTextOrg || ''
        btnEl.classList.remove('error')
      }

      btnEl.innerText = '加载中...'

      checker
        .request(this, inputVal)
        .then((data) => {
          // 请求成功
          this.done(checker, layer)
          if (checker.onSuccess) checker.onSuccess(this, data, inputVal, formEl)
          if (payload.onSuccess) payload.onSuccess(inputVal, dialog.$el)
        })
        .catch((err) => {
          // 请求失败
          btnTextSet(String(err.msg || String(err)))
          if (checker.onError) checker.onError(this, err, inputVal, formEl)

          const tf = setTimeout(() => btnTextRestore(), 3000)
          input.onfocus = () => {
            btnTextRestore()
            clearTimeout(tf)
          }
        })

      return false
    })

    // 取消按钮
    dialog.setNo(() => {
      this.done(checker, layer)
      if (payload.onCancel) payload.onCancel()
      return false
    })

    layer.getEl().append(dialog.$el)

    if (payload.onMount) payload.onMount(dialog.$el) // onMount
  }

  private done(checker: Checker, layer: Layer) {
    layer.disposeNow()
    this.launched = this.launched.filter(c => c !== checker)
  }
}

export interface Checker {
  el?: HTMLElement
  inputType?: 'password' | 'text'
  body: (launcher: CheckerLauncher) => HTMLElement
  request: (launcher: CheckerLauncher, inputVal: string) => Promise<string>
  onSuccess?: (launcher: CheckerLauncher, data: string, inputVal: string, formEl: HTMLElement) => void
  onError?: (launcher: CheckerLauncher, err: any, inputVal: string, formEl: HTMLElement) => void
}

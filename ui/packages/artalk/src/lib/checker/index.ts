import Context from '~/types/context'
import Dialog from '@/components/dialog'
import Layer from '@/layer'
import * as Utils from '../utils'
import * as Ui from '../ui'
import CaptchaChecker from './captcha-checker'
import AdminChecker from './admin-checker'

export interface CheckerCaptchaPayload extends CheckerPayload {
  imgData?: string
  iframe?: string
}

export interface CheckerPayload {
  onSuccess?: (inputVal: string, dialogEl?: HTMLElement) => void
  onMount?: (dialogEl: HTMLElement) => void
  onCancel?: () => void
}

/**
 * Checker 发射台
 */
export default class CheckerLauncher {
  public ctx: Context

  public launched: Checker[] = []
  public captchaConf: { val?: string, imgData?: string, iframe?: string } = {}

  constructor(ctx: Context) {
    this.ctx = ctx
  }

  public checkCaptcha(payload: CheckerCaptchaPayload) {
    this.captchaConf.imgData = payload.imgData
    this.captchaConf.iframe = payload.iframe

    this.fire(CaptchaChecker, payload)
  }

  public checkAdmin(payload: CheckerPayload) {
    this.fire(AdminChecker, payload)
  }

  public fire(checker: Checker, payload: CheckerPayload) {
    if (this.launched.includes(checker)) return // 阻止同时 fire 相同的 checker
    this.launched.push(checker)

    // 创建层
    const layer = new Layer(this.ctx, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    // Checker 的上下文
    let hideInteractInput = false
    const checkerCtx: CheckerCtx = {
      getLayer: () => layer,
      hideInteractInput: () => {
        hideInteractInput = true
      },
      triggerSuccess: () => {
        this.close(checker, layer)
        if (checker.onSuccess) checker.onSuccess(this, checkerCtx, "", "", formEl)
        if (payload.onSuccess) payload.onSuccess("", dialog.$el)
      },
      cancel: () => {
        this.close(checker, layer)
        if (payload.onCancel) payload.onCancel()
      }
    }

    // 创建表单
    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(this, checkerCtx))

    // 输入框
    const $input = Utils.createElement<HTMLInputElement>(
      `<input id="check" type="${checker.inputType || 'text'}" autocomplete="off" required placeholder="">`
    )
    formEl.appendChild($input)
    setTimeout(() => $input.focus(), 80) // 延迟 Focus

    // 绑定键盘事件
    $input.onkeyup = (evt) => {
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
    const dialog = new Dialog(this.ctx, formEl)

    // 确认按钮
    dialog.setYes((btnEl) => {
      const inputVal = $input.value.trim()

      // 按钮操作
      if (!btnTextOrg) btnTextOrg = btnEl.innerText
      const btnTextSet = (btnText: string) => {
        btnEl.innerText = btnText
        btnEl.classList.add('error')
      }
      const btnTextRestore = () => {
        btnEl.innerText = btnTextOrg || ''
        btnEl.classList.remove('error')
      }

      btnEl.innerText = `${this.ctx.$t('loading')}...`

      // 发送请求
      checker
        .request(this, checkerCtx, inputVal)
        .then((data) => {
          // 请求成功
          this.close(checker, layer)

          if (checker.onSuccess) checker.onSuccess(this, checkerCtx, data, inputVal, formEl)
          if (payload.onSuccess) payload.onSuccess(inputVal, dialog.$el)
        })
        .catch((err) => {
          // 请求失败
          btnTextSet(String(err.msg || String(err)))

          if (checker.onError) checker.onError(this, checkerCtx, err, inputVal, formEl)

          // 错误显示 3s 后恢复按钮
          const tf = setTimeout(() => btnTextRestore(), 3000)
          $input.onfocus = () => {
            btnTextRestore()
            clearTimeout(tf)
          }
        })

      return false
    })

    // 取消按钮
    dialog.setNo(() => {
      this.close(checker, layer)
      if (payload.onCancel) payload.onCancel()
      return false
    })

    if (hideInteractInput) {
      $input.style.display = 'none'
      dialog.$el.querySelector<HTMLElement>('.atk-layer-dialog-actions')!.style.display = 'none'
    }

    // 层装载 dialog 元素
    layer.getEl().append(dialog.$el)

    // onMount 回调
    if (payload.onMount) payload.onMount(dialog.$el)
  }

  // 关闭 checker 对话框
  private close(checker: Checker, layer: Layer) {
    layer.disposeNow()
    this.launched = this.launched.filter(c => c !== checker)
  }
}

export interface Checker {
  el?: HTMLElement
  inputType?: 'password' | 'text'
  body: (launcher: CheckerLauncher, ctx: CheckerCtx) => HTMLElement
  request: (launcher: CheckerLauncher, ctx: CheckerCtx, inputVal: string) => Promise<string>
  onSuccess?: (launcher: CheckerLauncher, ctx: CheckerCtx, data: string, inputVal: string, formEl: HTMLElement) => void
  onError?: (launcher: CheckerLauncher, ctx: CheckerCtx, err: any, inputVal: string, formEl: HTMLElement) => void
}

export interface CheckerCtx {
  getLayer(): Layer
  hideInteractInput(): void
  triggerSuccess(): void
  cancel(): void
}

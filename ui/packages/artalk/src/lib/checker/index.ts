import Context from '~/types/context'
import Dialog from '@/components/dialog'
import Layer from '@/layer'
import * as Utils from '../utils'
import CaptchaChecker from './captcha'
import AdminChecker from './admin'
import type Api from '../../api'
import type User from '../user'

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

  constructor(ctx: Context) {
    this.ctx = ctx
  }

  public checkCaptcha(payload: CheckerCaptchaPayload) {
    this.fire(CaptchaChecker, payload, (ctx) => {
      ctx.set('img_data', payload.imgData)
      ctx.set('iframe', payload.iframe)
    })
  }

  public checkAdmin(payload: CheckerPayload) {
    this.fire(AdminChecker, payload)
  }

  public fire(checker: Checker, payload: CheckerPayload, postFire?: (c: CheckerCtx) => void) {
    if (this.launched.includes(checker)) return // 阻止同时 fire 相同的 checker
    this.launched.push(checker)

    // 创建层
    const layer = new Layer(this.ctx, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    // 构建 Checker 的上下文
    const checkerStore: CheckerStore = {}
    let hideInteractInput = false
    const checkerCtx: CheckerCtx = {
      set: (key, val) => { checkerStore[key] = val },
      get: (key) => (checkerStore[key]),
      getCtx: () => (this.ctx),
      getApi: () => (this.ctx.getApi()),
      getLayer: () => layer,
      hideInteractInput: () => {
        hideInteractInput = true
      },
      triggerSuccess: () => {
        this.close(checker, layer)
        if (checker.onSuccess) checker.onSuccess(checkerCtx, "", "", formEl)
        if (payload.onSuccess) payload.onSuccess("", dialog.$el)
      },
      cancel: () => {
        this.close(checker, layer)
        if (payload.onCancel) payload.onCancel()
      }
    }

    if (postFire) postFire(checkerCtx)

    // 创建表单
    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(checkerCtx))

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
        .request(checkerCtx, inputVal)
        .then((data) => {
          // 请求成功
          this.close(checker, layer)

          if (checker.onSuccess) checker.onSuccess(checkerCtx, data, inputVal, formEl)
          if (payload.onSuccess) payload.onSuccess(inputVal, dialog.$el)
        })
        .catch((err) => {
          // 请求失败
          btnTextSet(String(err.msg || String(err)))

          if (checker.onError) checker.onError(checkerCtx, err, inputVal, formEl)

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
  body: (checker: CheckerCtx) => HTMLElement
  request: (checker: CheckerCtx, inputVal: string) => Promise<string>
  onSuccess?: (checker: CheckerCtx, respData: string, inputVal: string, formEl: HTMLElement) => void
  onError?: (checker: CheckerCtx, err: any, inputVal: string, formEl: HTMLElement) => void
}

interface CheckerStore {
  val?: string
  img_data?: string
  iframe?: string
}

export interface CheckerCtx {
  get<K extends keyof CheckerStore>(key: K): CheckerStore[K]
  set<K extends keyof CheckerStore>(key: K, val: CheckerStore[K]): void
  getCtx(): Context
  getApi(): Api
  getLayer(): Layer
  hideInteractInput(): void
  triggerSuccess(): void
  cancel(): void
}

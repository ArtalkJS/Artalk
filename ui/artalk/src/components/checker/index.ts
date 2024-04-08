import type { Api } from '@/api'
import Dialog from '@/components/dialog'
import $t from '@/i18n'
import type { ContextApi } from '@/types'
import type User from '@/lib/user'
import * as Utils from '@/lib/utils'
import CaptchaChecker from './captcha'
import AdminChecker from './admin'

export interface CheckerCaptchaPayload extends CheckerPayload {
  img_data?: string
  iframe?: string
}

export interface CheckerPayload {
  onSuccess?: () => void
  onMount?: (dialogEl: HTMLElement) => void
  onCancel?: () => void
}

export interface CheckerLauncherOptions {
  getCtx: () => ContextApi
  getApi: () => Api
  getCaptchaIframeURL: () => string
  onReload: () => void
}

function wrapPromise<P extends CheckerPayload = CheckerPayload>(fn: (p: P) => void) {
  return (payload: P) =>
    new Promise<void>((resolve, reject) => {
      const cancelFn = payload.onCancel
      payload.onCancel = () => {
        cancelFn && cancelFn()
        reject(new Error('user canceled the checker'))
      }
      const successFn = payload.onSuccess
      payload.onSuccess = () => {
        successFn && successFn()
        resolve()
      }
      fn(payload)
    })
}

/**
 * Checker 发射台
 */
export default class CheckerLauncher {
  constructor(private opts: CheckerLauncherOptions) {}

  public checkCaptcha: (payload: CheckerCaptchaPayload) => Promise<void> = wrapPromise((p) => {
    this.fire(CaptchaChecker, p, (ctx) => {
      ctx.set('img_data', p.img_data)
      ctx.set('iframe', p.iframe)
    })
  })

  public checkAdmin: (payload: CheckerPayload) => Promise<void> = wrapPromise((p) => {
    this.fire(AdminChecker, p)
  })

  public fire(checker: Checker, payload: CheckerPayload, postFire?: (c: CheckerCtx) => void) {
    // 显示层
    const layer = this.opts.getCtx().get('layerManager').create(`checker-${new Date().getTime()}`)
    layer.show()

    const close = () => {
      layer.destroy()
    }

    // 构建 Checker 的上下文
    const checkerStore: CheckerStore = {}
    let hideInteractInput = false
    const checkerCtx: CheckerCtx = {
      set: (key, val) => {
        checkerStore[key] = val
      },
      get: (key) => checkerStore[key],
      getOpts: () => this.opts,
      getUser: () => this.opts.getCtx().get('user'),
      getApi: () => this.opts.getApi(),
      hideInteractInput: () => {
        hideInteractInput = true
      },
      triggerSuccess: () => {
        close()
        if (checker.onSuccess) checker.onSuccess(checkerCtx, '', '', formEl)
        if (payload.onSuccess) payload.onSuccess()
      },
      cancel: () => {
        close()
        if (payload.onCancel) payload.onCancel()
      },
    }

    if (postFire) postFire(checkerCtx)

    // 创建表单
    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(checkerCtx))

    // 输入框
    const $input = Utils.createElement<HTMLInputElement>(
      `<input id="check" type="${checker.inputType || 'text'}" autocomplete="off" required placeholder="">`,
    )
    formEl.appendChild($input)
    setTimeout(() => $input.focus(), 80) // 延迟 Focus

    // 绑定键盘事件
    $input.onkeyup = (evt) => {
      if (evt.key === 'Enter' || evt.keyCode === 13) {
        // 按下回车键
        evt.preventDefault()
        layer.getEl().querySelector<HTMLButtonElement>('button[data-action="confirm"]')!.click()
      }
    }

    let btnTextOrg: string | undefined
    const dialog = new Dialog(formEl)

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

      btnEl.innerText = `${$t('loading')}...`

      // 发送请求
      checker
        .request(checkerCtx, inputVal)
        .then((data) => {
          // 请求成功
          close()

          if (checker.onSuccess) checker.onSuccess(checkerCtx, data, inputVal, formEl)
          if (payload.onSuccess) payload.onSuccess()
        })
        .catch((err) => {
          // 请求失败
          btnTextSet(String(err.message || String(err)))

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
      close()
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
}

export interface Checker<T = any> {
  el?: HTMLElement
  inputType?: 'password' | 'text'
  body: (checker: CheckerCtx) => HTMLElement
  request: (checker: CheckerCtx, inputVal: string) => Promise<T>
  onSuccess?: (checker: CheckerCtx, respData: T, inputVal: string, formEl: HTMLElement) => void
  onError?: (checker: CheckerCtx, errData: any, inputVal: string, formEl: HTMLElement) => void
}

interface CheckerStore {
  val?: string
  img_data?: string
  iframe?: string
}

export interface CheckerCtx {
  get<K extends keyof CheckerStore>(key: K): CheckerStore[K]
  set<K extends keyof CheckerStore>(key: K, val: CheckerStore[K]): void
  getOpts(): CheckerLauncherOptions
  getApi(): Api
  getUser(): User
  hideInteractInput(): void
  triggerSuccess(): void
  cancel(): void
}

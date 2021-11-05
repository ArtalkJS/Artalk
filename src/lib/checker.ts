import Context from "../context"
import * as Utils from "./utils"
import User from "./user"
import * as Ui from './ui'
import Layer from '../components/layer'
import Api from '../api'
import { CheckerCaptchaConf } from '~/types/event'

interface CheckerConf {
  el?: HTMLElement
  body: (that: Checker) => HTMLElement
  request: (that: Checker, inputVal: string) => Promise<string>
  onSuccess?: (that: Checker, data: string, inputVal: string, formEl: HTMLElement) => void
  onError?: (that: Checker, err: any, inputVal: string, formEl: HTMLElement) => void
}

const AdminChecker = {
  request: (that, inputVal) => {
    const data = {
      name: that.user.data.nick,
      email: that.user.data.email,
      password: inputVal
    }

    return new Api(that.ctx).login(data.name, data.email, data.password)
  },
  body: () => Utils.createElement('<span>敲入密码来验证管理员身份：</span>'),
  onSuccess: (that, userToken, inputVal, formEl) => {
    that.user.data.isAdmin = true
    that.user.data.token = userToken
    that.user.save()
    that.ctx.trigger('user-changed', that.ctx.user.data)
    that.ctx.trigger('list-reload')
  },
  onError: (that, err, inputVal, formEl) => {

  }

} as CheckerConf

const CaptchaChecker = {
  request: (that, inputVal) => (new Api(that.ctx).captchaCheck(inputVal)),
  body: (that) => {
    const elem = Utils.createElement(`<span><img class="atk-captcha-img" src="${that.submitCaptchaImgData || ''}" alt="验证码">敲入验证码继续：</span>`);

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
  onSuccess: (that, msg, data, inputVal, formEl) => {
    that.submitCaptchaVal = inputVal
  },
  onError: (that, err, inputVal, formEl) => {
    formEl.querySelector<HTMLElement>('.atk-captcha-img')!.click() // 刷新验证码
  }
} as CheckerConf

export default class Checker {
  public ctx: Context
  public user: User

  public submitCaptchaVal?: string
  public submitCaptchaImgData?: string

  constructor (ctx: Context) {
    this.ctx = ctx
    this.user = ctx.user

    this.ctx.on('checker-captcha', (conf) => {
      if (conf.imgData) {
        this.submitCaptchaImgData = conf.imgData
      }
      this.check('captcha', conf)
    })

    this.ctx.on('checker-admin', (conf) => {
      this.check('admin', conf)
    })
  }

  public check (name: 'admin'|'captcha', conf: CheckerCaptchaConf) {
    const checker = {
      admin: AdminChecker,
      captcha: CaptchaChecker
    }[name]

    const layer = new Layer(this.ctx, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(this))

    const input = Utils.createElement<HTMLInputElement>(`<input id="check" type="${(name === 'admin' ? 'password' : 'text')}" required placeholder="">`)
    formEl.appendChild(input)
    setTimeout(() => {
      input.focus() // 延迟以保证有效
    }, 80)

    input.onkeyup = (evt) => {
      if (evt.key === 'Enter' || evt.keyCode === 13) { // 按下回车键
        evt.preventDefault();
        (layer.getEl().querySelector<HTMLButtonElement>('button[data-action="confirm"]'))!.click()
      }
    }

    let btnRawText: string|undefined
    const dialogEl = Ui.buildDialog(formEl, (btnElem: HTMLElement) => {
      const inputVal = input.value.trim()
      if (!btnRawText) {
        btnRawText = btnElem.innerText
      }
      const btnTextSet = (btnText: string) => {
        btnElem.innerText = btnText
        btnElem.classList.add('error')
      }
      const btnTextRestore = () => {
        btnElem.innerText = btnRawText || ''
        btnElem.classList.remove('error')
      }

      btnElem.innerText = '加载中...'
      checker.request(this, inputVal)
        .then(data => {
          // 请求成功
          layer.disposeNow()
          if (checker.onSuccess) checker.onSuccess(this, data, inputVal, formEl)
          if (conf.onSuccess) conf.onSuccess(inputVal, dialogEl)
        })
        .catch(err => {
          btnTextSet(String(err.msg || String(err)))
          if (checker.onError)
            checker.onError(this, err, inputVal, formEl)

          const tf = setTimeout(() => {
            btnTextRestore()
          }, 3000)
          input.onfocus = () => {
            btnTextRestore()
            clearTimeout(tf)
          }
        })

      return false
    }, () => {
      layer.disposeNow()
      if (conf.onCancel) conf.onCancel()
      return false
    })

    layer.getEl().append(dialogEl)

    if (conf.onMount) conf.onMount(dialogEl) // onMount
  }
}

import Context from "../Context"
import * as Utils from "./utils"
import User from "./user"
import * as Ui from './ui'
import BuildLayer from "../components/Layer"
import Api from "./api"

interface CheckerConf {
  el?: HTMLElement
  body: (that: Checker) => HTMLElement
  api: (that: Checker) => Function
  reqData: (that: Checker, inputVal: string) => any
  onSuccess: (that: Checker, msg: string, data: any, inputVal: string) => void
}

const AdminChecker = {
  api: (that) => (new Api(that.ctx).login),
  body: () => Utils.createElement('<span>敲入密码来验证管理员身份：</span>'),
  reqData: (that, inputVal) => {
    const data = {
      user: that.user.data.nick,
      email: that.user.data.email,
      password: inputVal
    }

    return data
  },
  onSuccess: (that, msg, data, inputVal) => {
    that.user.data.isAdmin = true
    that.user.data.token = inputVal
    that.user.save()
  }
} as CheckerConf

const CaptchaChecker = {
  api: (that) => (new Api(that.ctx).captchaCheck),
  body: (that) => {
    const elem = Utils.createElement(`<span><img class="artalk-captcha-img" src="${that.submitCaptchaImgData || ''}" alt="验证码">敲入验证码继续：</span>`);

    // 刷新验证码
    elem.querySelector<HTMLElement>('.artalk-captcha-img')!.onclick = () => {
      const imgEl = elem.querySelector('.artalk-captcha-img')
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
  reqData: (inputVal) => {
    const data = {
      captcha: inputVal
    }

    return data
  },
  onSuccess: (that, msg, data, inputVal) => {
    that.submitCaptchaVal = inputVal
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
  }

  public check (name: 'admin'|'captcha', onSuccess: (inputVal: string, dialogEl: HTMLElement) => void, onMount?: (dialogEl: HTMLElement) => void) {
    const checker = {
      admin: AdminChecker,
      captcha: CaptchaChecker
    }[name]

    const formEl = Utils.createElement()
    formEl.appendChild(checker.body(this))

    const input = Utils.createElement<HTMLInputElement>(`<input id="check" type="${(name === 'admin' ? 'password' : 'text')}" required placeholder="">`)
    formEl.appendChild(input)
    setTimeout(() => {
      input.focus() // 延迟以保证有效
    }, 80)

    const layer = BuildLayer(this.ctx, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    input.onkeyup = (evt) => {
      if (evt.keyCode === 13) { // 按下回车键
        evt.preventDefault();
        (layer.getEl().querySelector<HTMLButtonElement>('button[data-action="confirm"]'))!.click()
      }
    }

    let btnRawText: string|undefined
    Ui.showDialog(layer.getEl(), formEl, (dialogElem, btnElem: HTMLElement) => {
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

      this.artalk.request(checker.reqAct, checker.reqData(this, inputVal), () => {
        btnElem.innerText = '加载中...'
      }, () => {

      }, (msg, data) => {
        // 请求成功
        checker.onSuccess(msg, data, inputVal)
        layer.disposeNow()
        onSuccess(inputVal, dialogElem)
      }, (msg, data) => {
        // 请求失败
        btnTextSet(msg)

        if (name === '管理员') {
          if ((typeof data === 'object') && data !== null && typeof data.need_captcha === 'boolean' && data.need_captcha === true) {
            // 验证码验证
            this.artalk.checker.submitCaptchaImgData = data.img_data
            layer.disposeNow()
            this.artalk.checker.action('验证码', () => { // 密码错误达到上限需输入验证码
              this.artalk.checker.action('管理员', (a, b) => { // 套娃 XD (虽然有点不优雅，但是懒得改了.... 555)
                this.submitCaptchaVal = null
                onSuccess(a, b)
              }, (el) => {
                // onMount
                el.querySelector<HTMLInputElement>('input')!.value = inputVal
                el.querySelector<HTMLButtonElement>('[data-action="confirm"]')!.click()
              })
            })
          }
        }

        if (name === '验证码') {
          checker.refresh(data.img_data)
        }
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
      return false
    }, (dialogEl) => { // onMount
      if (onMount) onMount(dialogEl)
    })
  }
}

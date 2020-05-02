import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'
import BuildLayer from './Layer'

interface CheckerItem {
  elem?: HTMLElement
  body: () => HTMLElement
  reqAct: string
  reqObj: (inputVal: string) => any
  onSuccess: (msg: string, data: any, inputVal: string) => void
  refresh?: Function
}

export default class Checker extends ArtalkContext {
  public submitCaptchaVal: string
  public submitCaptchaImgData: string

  private readonly LIST: {[key: string]: CheckerItem} = {
    '管理员': {
      body: () => Utils.createElement('<span>敲入密码来验证管理员身份：</span>'),
      reqAct: 'AdminCheck',
      reqObj: (inputVal) => {
        return {
          nick: this.artalk.user.data.nick,
          email: this.artalk.user.data.email,
          password: inputVal
        }
      },
      onSuccess: (msg, data, inputVal) => {
        this.artalk.user.data.isAdmin = true
        this.artalk.user.data.password = inputVal
        this.artalk.user.save()
      }
    },
    '验证码': {
      body: () => {
        const elem = Utils.createElement(`<span><img class="artalk-captcha-img" src="${this.submitCaptchaImgData || ''}" alt="验证码">敲入验证码继续：</span>`)
        this.LIST['验证码'].elem = elem;
        (elem.querySelector('.artalk-captcha-img') as HTMLElement).onclick = () => {
          this.LIST['验证码'].refresh()
        }
        return elem
      },
      reqAct: 'CaptchaCheck',
      reqObj: (inputVal) => {
        return {
          captcha: inputVal
        }
      },
      onSuccess: (msg, data, inputVal) => {
        this.submitCaptchaVal = inputVal
      },
      refresh: (imgData?: string) => {
        const { elem } = this.LIST['验证码']
        const imgEl = elem.querySelector('.artalk-captcha-img')
        if (!imgData) {
          this.artalk.request('CaptchaCheck', { refresh: true }, () => {}, () => {}, (msg, data) => {
            imgEl.setAttribute('src', data.img_data)
          }, () => {})
        } else {
          imgEl.setAttribute('src', imgData)
        }
      }
    }
  }

  public action (name: '管理员'|'验证码', action: () => void) {
    const checker = this.LIST[name]

    const formEl = Utils.createElement()
    formEl.appendChild(checker.body())

    const input = Utils.createElement(`<input id="check" type="${(name === '管理员' ? 'password' : 'text')}" required placeholder="">`) as HTMLInputElement
    formEl.appendChild(input)
    setTimeout(() => {
      input.focus() // 延迟以保证有效
    }, 80)

    const layer = BuildLayer(this.artalk, `checker-${new Date().getTime()}`)
    layer.setMaskClickHide(false)
    layer.show()

    input.onkeyup = (evt) => {
      if (evt.keyCode === 13) { // 按下回车键
        evt.preventDefault();
        (layer.getEl().querySelector('button[data-action="confirm"]') as HTMLButtonElement).click()
      }
    }

    this.artalk.ui.showDialog(layer.getEl(), formEl, (dialogElem, btnElem: HTMLElement) => {
      const inputVal = input.value.trim()
      const btnRawText = btnElem.innerText
      const btnTextSet = (btnText: string) => {
        btnElem.innerText = btnText
        btnElem.classList.add('error')
      }
      const btnTextRestore = () => {
        btnElem.innerText = btnRawText
        btnElem.classList.remove('error')
      }

      this.artalk.request(checker.reqAct, checker.reqObj(inputVal), () => {
        btnElem.innerText = '加载中...'
      }, () => {

      }, (msg, data) => {
        // 请求成功
        checker.onSuccess(msg, data, inputVal)
        layer.dispose()
        action()
      }, (msg, data) => {
        // 请求失败
        btnTextSet(msg)
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
      layer.dispose()
      return false
    })
  }
}

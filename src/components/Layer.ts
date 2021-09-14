import Context from '../Context'
import Component from '../lib/component'
import Constant from '../Constant'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

export function GetLayerWrap (ctx: Context): { wrapEl: HTMLElement, maskEl: HTMLElement } {
  let wrapEl = document.querySelector<HTMLElement>(`.artalk-layer-wrap#ctx-${ctx.cid}`)
  if (!wrapEl) {
    wrapEl = Utils.createElement(
      `<div class="artalk-layer-wrap" id="ctx-${ctx.cid}" style="display: none;"><div class="artalk-layer-mask"></div></div>`
    )
    document.body.appendChild(wrapEl)
  }

  const maskEl = wrapEl.querySelector<HTMLElement>('.artalk-layer-mask')!

  // dark mode
  if (wrapEl) {
    if (ctx.conf.darkMode) {
      wrapEl.classList.add(Constant.DARK_MODE_CLASSNAME)
    } else {
      wrapEl.classList.remove(Constant.DARK_MODE_CLASSNAME)
    }
  }

  return { wrapEl, maskEl }
}

export default function BuildLayer (ctx: Context, name: string, el?: HTMLElement) {
  return new Layer(ctx, name, el)
}

export class Layer extends Component {
  private name: string
  private wrapEl: HTMLElement
  private maskEl: HTMLElement

  private maskClickHideEnable: boolean = true

  constructor (ctx: Context, name: string, el?: HTMLElement) {
    super(ctx)

    this.name = name
    const { wrapEl, maskEl } = GetLayerWrap(ctx)
    this.wrapEl = wrapEl
    this.maskEl = maskEl

    this.el = this.wrapEl.querySelector(`[data-layer-name="${name}"].artalk-layer-item`)!
    if (this.el === null) {
      // 若传递 layer 元素为空
      if (!el) {
        this.el = Utils.createElement()
        this.el.classList.add('artalk-layer-item')
      } else {
        this.el = el
      }
    }
    this.el.setAttribute('data-layer-name', name)
    this.el.style.display = 'none'

    // 添加到 layers wrap 中
    this.wrapEl.append(this.el)
  }

  getName () {
    return this.name
  }

  getWrapEl () {
    return this.wrapEl
  }

  getEl () {
    return this.el
  }

  private static hideTimeoutList: number[] = []

  show () {
    Layer.hideTimeoutList.forEach(item => {
      clearTimeout(item)
    })
    Layer.hideTimeoutList = []

    this.wrapEl.style.display = 'block'
    this.maskEl.style.display = 'block'
    this.maskEl.classList.add('artalk-fade-in')
    this.el.style.display = ''

    this.maskEl.onclick = () => {
      if (this.maskClickHideEnable) this.hide()
    }

    document.body.style.overflow = 'hidden'
  }

  hide () {
    Layer.hideTimeoutList.push(window.setTimeout(() => {
      this.wrapEl.style.display = 'none'
      document.body.style.overflow = ''
    }, 450))

    this.wrapEl.classList.add('artalk-fade-out')
    Layer.hideTimeoutList.push(window.setTimeout(() => {
      this.wrapEl.style.display = 'none'
      this.wrapEl.classList.remove('artalk-fade-out')
    }, 200))

    this.el.style.display = 'none'
  }

  setMaskClickHide (enable: boolean) {
    this.maskClickHideEnable = enable
  }

  /** 销毁 - 无动画 */
  disposeNow () {
    document.body.style.overflow = ''
    this.el.remove()
    // this.el dispose
    this.checkCleanLayer()
  }

  /** 销毁 */
  dispose () {
    this.hide()
    this.el.remove()
    // this.el dispose
    this.checkCleanLayer()
  }

  checkCleanLayer () {
    if (this.getWrapEl().querySelectorAll('.artalk-layer-item').length === 0) {
      this.wrapEl.style.display = 'none'
    }
  }
}

import Artalk from '../Artalk'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

// TODO: for 一页多评论，该实例请用完及时销毁；待优化
export class Layer extends ArtalkContext {
  private el: HTMLElement
  private wrapEl: HTMLElement
  private maskEl: HTMLElement

  private maskClickHideEnable: boolean = true

  constructor (artalk: Artalk, private name: string, html?: HTMLElement) {
    super(artalk)
    this.initWrap()

    this.el = this.wrapEl.querySelector(`[data-layer-name="${name}"]`)
    if (this.el === null) {
      // 若 layer 元素未创建过
      if (!html) {
        this.el = Utils.createElement()
        this.el.classList.add('artalk-layer-item')
      } else {
        this.el = html
      }
    }
    this.el.setAttribute('data-layer-name', name)
    this.el.style.display = 'none'

    // 添加到 layers wrap 中
    this.wrapEl.append(this.el)
  }

  private initWrap () {
    this.wrapEl = document.querySelector(`.artalk-layer-wrap`)
    if (!this.wrapEl) {
      this.wrapEl = Utils.createElement(
        `<div class="artalk-layer-wrap" style="display: none;"><div class="artalk-layer-mask"></div></div>`
      )
      document.body.appendChild(this.wrapEl)
    }

    this.maskEl = this.wrapEl.querySelector('.artalk-layer-mask')
    this.artalk.ui.initDarkMode()
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
    Layer.hideTimeoutList.push(setTimeout(() => {
      this.wrapEl.style.display = 'none'
      document.body.style.overflow = ''
    }, 450))

    this.wrapEl.classList.add('artalk-fade-out')
    Layer.hideTimeoutList.push(setTimeout(() => {
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
    this.wrapEl.style.display = 'none'
    document.body.style.overflow = ''
    this.el.remove()
    this.el = null
  }

  /** 销毁 */
  dispose () {
    this.hide()
    this.el.remove()
    this.el = null
  }
}

export default function BuildLayer (artalk: Artalk, name: string, html?: HTMLElement) {
  return new Layer(artalk, name, html)
}

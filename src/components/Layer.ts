import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

export default class Layer extends ArtalkContext {
  private el: HTMLElement
  private wrapEl: HTMLElement
  private maskEl: HTMLElement

  private maskClickHideEnable: boolean = true

  constructor (private name: string, html?: HTMLElement) {
    super()
    this.initWrap()

    this.el = this.wrapEl.querySelector(
      `[data-layer-name="${name}"]`
    )
    if (this.el !== null) throw new Error(`layer "${name}" 已存在`)
    if (!html) {
      this.el = Utils.createElement()
      this.el.classList.add('artalk-layer-item')
    } else {
      this.el = html
    }
    this.el.setAttribute('data-layer-name', name)
    this.el.style.display = 'none'

    this.wrapEl.prepend(this.el)
  }

  private initWrap () {
    this.wrapEl = document.querySelector(`.artalk-layer-wrap`)
    if (!this.wrapEl) {
      this.wrapEl = Utils.createElement(
        '<div class="artalk-layer-wrap" style="display: none;"><div class="artalk-layer-mask"></div></div>'
      )
      document.body.appendChild(this.wrapEl)
    }

    this.maskEl = this.wrapEl.querySelector('.artalk-layer-mask')
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

  private static hideTimeoutList: number[] = []

  hide () {
    Layer.hideTimeoutList.push(setTimeout(() => {
      this.wrapEl.style.display = 'none'
    }, 450))

    this.wrapEl.classList.add('artalk-fade-out')
    Layer.hideTimeoutList.push(setTimeout(() => {
      this.wrapEl.style.display = 'none'
      this.wrapEl.classList.remove('artalk-fade-out')
    }, 200))

    this.el.style.display = 'none'
    document.body.style.overflow = ''
  }

  setMaskClickHide (enable: boolean) {
    this.maskClickHideEnable = enable
  }

  dispose () {
    this.hide()
    this.el.remove()
    this.el = null
  }
}

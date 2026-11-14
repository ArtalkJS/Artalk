import { getScrollbarHelper } from './scrollbar-helper'
import { Layer } from './layer'
import * as Utils from '@/lib/utils'

export class LayerWrap {
  private $wrap: HTMLElement
  private $mask: HTMLElement
  private items: Layer[] = []
  private scrollbar = getScrollbarHelper()
  private disposed = false

  constructor() {
    this.$wrap = Utils.createElement(
      `<div class="atk-layer-wrap" style="display: none;"><div class="atk-layer-mask"></div></div>`,
    )
    this.$mask = this.$wrap.querySelector<HTMLElement>('.atk-layer-mask')!
    this.$mask.onclick = () => {
      this.items.slice().forEach((layer) => {
        if (layer.getAllowMaskClose()) layer.hide()
      })
    }
  }

  createItem(name: string, el?: HTMLElement) {
    // init element
    el = el || this.createItemElement(name)
    el.setAttribute('data-layer-name', name)
    this.$wrap.appendChild(el)

    // create layer instance
    const layer = new Layer(el, {
      onHide: () => this.hideWrap(el!),
      onShow: () => this.showWrap(),
      onDestroy: () => {
        this.items = this.items.filter((item) => item !== layer)
      },
    })

    // add to items
    this.items.push(layer)
    if (this.disposed) layer.destroy()

    return layer
  }

  private createItemElement(name: string) {
    const el = document.createElement('div')
    el.classList.add('atk-layer-item')
    el.style.display = 'none'
    this.$wrap.appendChild(el)
    return el
  }

  destroy() {
    if (this.disposed) return
    this.disposed = true
    const errors: unknown[] = []
    try {
      this.items.slice().forEach((layer) => {
        try {
          layer.destroy()
        } catch (error) {
          errors.push(error)
        }
      })
    } finally {
      this.items = []
      this.scrollbar.unlock()
      this.$mask.onclick = null
      this.$wrap.remove()
    }
    if (errors.length) throw errors[0]
  }

  getWrap() {
    return this.$wrap
  }

  getMask() {
    return this.$mask
  }

  showWrap() {
    if (this.disposed) return
    this.$wrap.style.display = 'block'
    this.$mask.style.display = 'block'
    this.$mask.classList.add('atk-fade-in')
    this.scrollbar.lock()
  }

  hideWrap($el: HTMLElement) {
    // if wrap contains more than one item, do not hide entire wrap
    if (
      this.items
        .map((l) => l.getEl())
        .filter((e) => e !== $el && e.isConnected && e.style.display !== 'none').length > 0
    ) {
      return
    }

    this.$wrap.style.display = 'none'
    this.scrollbar.unlock()
  }
}

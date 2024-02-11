import * as Utils from '@/lib/utils'
import { getScrollbarHelper } from './scrollbar-helper'
import { Layer } from './layer'

export class LayerWrap {
  private $wrap: HTMLElement
  private $mask: HTMLElement
  private items: Layer[] = []

  constructor() {
    this.$wrap = Utils.createElement(
      `<div class="atk-layer-wrap" style="display: none;"><div class="atk-layer-mask"></div></div>`
    )
    this.$mask = this.$wrap.querySelector<HTMLElement>('.atk-layer-mask')!
  }

  createItem(name: string, el?: HTMLElement) {
    // init element
    el = el || this.createItemElement(name)
    el.setAttribute('data-layer-name', name)
    this.$wrap.appendChild(el)

    // create layer instance
    const layer = new Layer(el, {
      onHide: () => this.hideWrap(el!),
      onShow: () => this.showWrap()
    })

    // register mask click event
    this.getMask().addEventListener('click', () => {
      layer.getAllowMaskClose() && layer.hide()
    })

    // add to items
    this.items.push(layer)

    return layer
  }

  private createItemElement(name: string) {
    const el = document.createElement('div')
    el.classList.add('atk-layer-item')
    el.style.display = 'none'
    this.$wrap.appendChild(el)
    return el
  }

  getWrap() {
    return this.$wrap
  }

  getMask() {
    return this.$mask
  }

  showWrap() {
    this.$wrap.style.display = 'block'
    this.$mask.style.display = 'block'
    this.$mask.classList.add('atk-fade-in')
    getScrollbarHelper().lock()
  }

  hideWrap($el: HTMLElement) {
    // if wrap contains more than one item, do not hide entire wrap
    if (this.items.map(l => l.getEl()).filter(e => e !== $el && e.isConnected && e.style.display !== 'none').length > 0) {
      return
    }

    this.$wrap.style.display = 'none'
    getScrollbarHelper().unlock()
  }
}

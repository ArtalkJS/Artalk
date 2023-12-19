import * as Utils from '@/lib/utils'
import { getScrollbarHelper } from './scrollbar-helper'

export class LayerWrap {
  private $wrap: HTMLElement
  private $mask: HTMLElement
  private allowMaskClose: boolean = true
  private items: HTMLElement[] = []

  constructor() {
    this.$wrap = Utils.createElement(
      `<div class="atk-layer-wrap" style="display: none;"><div class="atk-layer-mask"></div></div>`,
    )
    this.$mask = this.$wrap.querySelector<HTMLElement>('.atk-layer-mask')!
  }

  createItem(name: string, el?: HTMLElement) {
    if (!el) {
      el = document.createElement('div')
      el.classList.add('atk-layer-item')
      el.setAttribute('data-layer-name', name)
      el.style.display = 'none'
    }
    this.$wrap.appendChild(el)
    this.items.push(el)
    return el
  }

  getWrap() {
    return this.$wrap
  }

  getMask() {
    return this.$mask
  }

  setMaskClose(enabled: boolean) {
    this.allowMaskClose = enabled
  }

  show() {
    this.$wrap.style.display = 'block'
    this.$mask.style.display = 'block'
    this.$mask.classList.add('atk-fade-in')
    this.$mask.onclick = () => {
      if (this.allowMaskClose) this.hide()
    }
    getScrollbarHelper().lock()
  }

  hide(callback?: () => void) {
    // if wrap contains more than one item, do not hide entire wrap
    if (
      this.items.filter((e) => e.isConnected && e.style.display !== 'none')
        .length > 1
    ) {
      callback && callback()
      return
    }

    const onAfterHide = () => {
      this.$wrap.style.display = 'none'
      this.$wrap.classList.remove('atk-fade-out')

      callback && callback()

      getScrollbarHelper().unlock()

      this.$wrap.onanimationend = null
    }

    // perform animation
    this.$wrap.classList.add('atk-fade-out')
    if (window.getComputedStyle(this.$wrap)['animation-name'] !== 'none') {
      this.$wrap.onanimationend = () => onAfterHide()
    } else {
      onAfterHide()
    }
  }
}

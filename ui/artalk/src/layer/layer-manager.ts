import type { ContextApi } from '@/types'
import { getScrollbarHelper } from './scrollbar-helper'
import { LayerWrap } from './wrap'

export class LayerManager {
  private wrap: LayerWrap

  constructor(ctx: ContextApi) {
    this.wrap = new LayerWrap()
    document.body.appendChild(this.wrap.getWrap())

    ctx.on('unmounted', () => {
      this.wrap.getWrap().remove()
    })

    // 记录页面原始 CSS 属性
    getScrollbarHelper().init()
  }

  getEl() {
    return this.wrap.getWrap()
  }

  create(name: string, el?: HTMLElement) {
    return this.wrap.createItem(name, el)
  }
}

import type { ContextApi } from '@/types'
import { getScrollbarHelper } from './scrollbar-helper'
import { LayerWrap } from './wrap'
import { Layer } from './layer'

export class LayerManager {
  private wrap: LayerWrap
  private ctx: ContextApi

  constructor(ctx: ContextApi) {
    this.ctx = ctx

    this.wrap = new LayerWrap()
    document.body.appendChild(this.wrap.getWrap())

    ctx.on('destroy', () => {
      this.wrap.getWrap().remove()
    })

    // 记录页面原始 CSS 属性
    getScrollbarHelper().init()
  }

  getEl() {
    return this.wrap.getWrap()
  }

  create(name: string, el?: HTMLElement | undefined) {
    const layer = new Layer(this.wrap, name, el)
    return layer
  }
}

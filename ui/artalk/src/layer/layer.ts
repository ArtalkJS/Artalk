import type { Layer as ILayer } from '@/types/layer'

export interface LayerOptions {
  onShow: () => void
  onHide: () => void
  /** @internal */
  onDestroy?: () => void
}

export class Layer implements ILayer {
  private allowMaskClose = true
  private destroyed = false
  private onAfterHide?: () => void

  constructor(
    private $el: HTMLElement,
    private opts: LayerOptions,
  ) {}

  setOnAfterHide(func: () => void) {
    this.onAfterHide = func
  }

  setAllowMaskClose(allow: boolean) {
    this.allowMaskClose = allow
  }

  getAllowMaskClose() {
    return this.allowMaskClose
  }

  getEl() {
    return this.$el
  }

  show() {
    if (this.destroyed) return
    this.opts.onShow()
    this.$el.style.display = ''
  }

  hide() {
    if (this.destroyed || this.$el.style.display === 'none') return
    this.opts.onHide()
    this.$el.style.display = 'none'
    this.onAfterHide && this.onAfterHide()
  }

  destroy() {
    if (this.destroyed) return
    this.destroyed = true
    this.opts.onHide()
    this.$el.remove()
    this.opts.onDestroy?.()
    this.onAfterHide && this.onAfterHide()
  }
}

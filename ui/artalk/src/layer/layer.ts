export interface LayerOptions {
  onShow: () => void
  onHide: () => void
}

export class Layer {
  private allowMaskClose = true
  private onAfterHide?: () => void

  constructor(
    private $el: HTMLElement,
    private opts: LayerOptions
  ) {
  }

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
    this.opts.onShow()
    this.$el.style.display = ''
  }

  async hide() {
    this.opts.onHide()
    this.$el.style.display = 'none'
    this.onAfterHide && this.onAfterHide()
  }

  async destroy() {
    this.opts.onHide()
    this.$el.remove()
    this.onAfterHide && this.onAfterHide()
  }
}

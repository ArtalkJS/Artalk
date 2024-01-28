export interface LayerOptions {
  showWrap: () => void
  hideWrap: () => void
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
    this.$el.style.display = ''
    this.opts.showWrap()
  }

  async hide() {
    this.opts.hideWrap()
    this.$el.style.display = 'none'
    this.onAfterHide && this.onAfterHide()
  }

  async destroy() {
    this.opts.hideWrap()
    this.$el.remove()
    this.onAfterHide && this.onAfterHide()
  }
}

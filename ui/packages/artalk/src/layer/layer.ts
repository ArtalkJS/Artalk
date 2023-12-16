import type { LayerWrap } from './wrap'

export class Layer {
  private $el: HTMLElement
  private wrap: LayerWrap
  private onAfterHide?: () => void

  constructor(wrap: LayerWrap, name: string, el?: HTMLElement) {
    this.wrap = wrap
    this.$el = this.wrap.createItem(name, el)
  }

  setOnAfterHide(func: () => void) {
    this.onAfterHide = func
  }

  getEl() {
    return this.$el
  }

  show() {
    this.$el.style.display = ''
    this.wrap.show()
  }

  hide() {
    this.wrap.hide(() => {
      this.$el.style.display = 'none'
      this.onAfterHide && this.onAfterHide()
    })
  }

  destroy() {
    this.wrap.hide(() => {
      this.$el.remove()
      this.onAfterHide && this.onAfterHide()
    })
  }
}

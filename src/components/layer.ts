import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

export default class Layer extends Component {
  private name: string
  private $wrap: HTMLElement
  private $mask: HTMLElement

  private maskClickHideEnable: boolean = true

  private bodyStyleOrgOverflow = ''
  private bodyStyleOrgPaddingRight = ''

  constructor (ctx: Context, name: string, el?: HTMLElement) {
    super(ctx)

    this.name = name
    const { $wrap, $mask } = GetLayerWrap(ctx)
    this.$wrap = $wrap
    this.$mask = $mask

    this.$el = this.$wrap.querySelector(`[data-layer-name="${name}"].atk-layer-item`)!
    if (this.$el === null) {
      // 若传递 layer 元素为空
      if (!el) {
        this.$el = Utils.createElement()
        this.$el.classList.add('atk-layer-item')
      } else {
        this.$el = el
      }
    }
    this.$el.setAttribute('data-layer-name', name)
    this.$el.style.display = 'none'

    // 添加到 layers wrap 中
    this.$wrap.append(this.$el)
  }

  getName () {
    return this.name
  }

  getWrapEl () {
    return this.$wrap
  }

  getEl () {
    return this.$el
  }

  private static hideTimeoutList: number[] = []

  show () {
    Layer.hideTimeoutList.forEach(item => {
      clearTimeout(item)
    })
    Layer.hideTimeoutList = []

    this.$wrap.style.display = 'block'
    this.$mask.style.display = 'block'
    this.$mask.classList.add('atk-fade-in')
    this.$el.style.display = ''

    this.$mask.onclick = () => {
      if (this.maskClickHideEnable) this.hide()
    }

    // body style 禁止滚动 + 防抖
    this.bodyStyleOrgOverflow = document.body.style.overflow
    this.bodyStyleOrgPaddingRight = document.body.style.paddingRight
    document.body.style.overflow = 'hidden'

    const bpr = parseInt(window.getComputedStyle(document.body, null).getPropertyValue('padding-right'), 10)
    document.body.style.paddingRight = `${Ui.getScrollBarWidth() + bpr || 0}px`
  }

  hide () {
    Layer.hideTimeoutList.push(window.setTimeout(() => {
      this.$wrap.style.display = 'none'
      // body style 禁止滚动解除
      document.body.style.overflow = this.bodyStyleOrgOverflow
      document.body.style.paddingRight = this.bodyStyleOrgPaddingRight
    }, 450))

    this.$wrap.classList.add('atk-fade-out')
    Layer.hideTimeoutList.push(window.setTimeout(() => {
      this.$wrap.style.display = 'none'
      this.$wrap.classList.remove('atk-fade-out')
    }, 200))

    this.$el.style.display = 'none'
  }

  setMaskClickHide (enable: boolean) {
    this.maskClickHideEnable = enable
  }

  /** 销毁 - 无动画 */
  disposeNow () {
    document.body.style.overflow = ''
    this.$el.remove()
    // this.$el dispose
    this.checkCleanLayer()
  }

  /** 销毁 */
  dispose () {
    this.hide()
    this.$el.remove()
    // this.$el dispose
    this.checkCleanLayer()
  }

  checkCleanLayer () {
    if (this.getWrapEl().querySelectorAll('.atk-layer-item').length === 0) {
      this.$wrap.style.display = 'none'
    }
  }
}

export function GetLayerWrap (ctx: Context): { $wrap: HTMLElement, $mask: HTMLElement } {
  let $wrap = document.querySelector<HTMLElement>(`.atk-layer-wrap#ctx-${ctx.cid}`)
  if (!$wrap) {
    $wrap = Utils.createElement(
      `<div class="atk-layer-wrap" id="ctx-${ctx.cid}" style="display: none;"><div class="atk-layer-mask"></div></div>`
    )
    document.body.appendChild($wrap)
  }

  const $mask = $wrap.querySelector<HTMLElement>('.atk-layer-mask')!

  return { $wrap, $mask }
}

import Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

export default class Layer extends Component {
  private name: string
  private $wrap: HTMLElement
  private $mask: HTMLElement

  private maskClickHideEnable: boolean = true

  public static BodyOrgOverflow: string
  public static BodyOrgPaddingRight: string

  public afterHide?: Function

  constructor (ctx: Context, name: string, el?: HTMLElement) {
    super(ctx)

    this.name = name
    const { $wrap, $mask } = getLayerWrap()
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

  getName() {
    return this.name
  }

  getWrapEl() {
    return this.$wrap
  }

  getEl() {
    return this.$el
  }

  show() {
    this.fireAllActionTimer()

    this.$wrap.style.display = 'block'
    this.$mask.style.display = 'block'
    this.$mask.classList.add('atk-fade-in')
    this.$el.style.display = ''

    this.$mask.onclick = () => {
      if (this.maskClickHideEnable) this.hide()
    }

    // body style 禁止滚动 + 防抖
    this.pageBodyScrollBarHide()
  }

  hide () {
    if (this.afterHide) this.afterHide()
    this.$wrap.classList.add('atk-fade-out')
    this.$el.style.display = 'none'

    // body style 禁止滚动解除
    this.pageBodyScrollBarShow()

    this.newActionTimer(() => {
      this.$wrap.style.display = 'none'
      this.checkCleanLayer()
    }, 450)
    this.newActionTimer(() => {
      this.$wrap.style.display = 'none'
      this.$wrap.classList.remove('atk-fade-out')
    }, 200)
  }

  setMaskClickHide (enable: boolean) {
    this.maskClickHideEnable = enable
  }

  // 页面滚动条隐藏
  pageBodyScrollBarHide() {
    document.body.style.overflow = 'hidden'

    const bpr = parseInt(window.getComputedStyle(document.body, null).getPropertyValue('padding-right'), 10)
    document.body.style.paddingRight = `${Ui.getScrollBarWidth() + bpr || 0}px`
  }

  // 页面滚动条显示
  pageBodyScrollBarShow() {
    document.body.style.overflow = Layer.BodyOrgOverflow
    document.body.style.paddingRight = Layer.BodyOrgPaddingRight
  }

  // Timers
  private static actionTimers: {act: Function, tid: number}[] = []

  private newActionTimer(func: Function, delay: number) {
    const act = () => {
      func() // 执行
      Layer.actionTimers = Layer.actionTimers.filter(o => o.act !== act) // 删除
    }

    const tid = window.setTimeout(() => act(), delay)

    Layer.actionTimers.push({ act, tid })
  }

  private fireAllActionTimer() {
    Layer.actionTimers.forEach(item => {
      clearTimeout(item.tid)
      item.act() // 立即执行
    })
  }

  /** 销毁 - 无动画 */
  disposeNow() {
    this.$el.remove()
    this.pageBodyScrollBarShow()
    // this.$el dispose
    this.checkCleanLayer()
  }

  /** 销毁 */
  dispose() {
    this.hide()
    this.$el.remove()
    // this.$el dispose
    this.checkCleanLayer()
  }

  checkCleanLayer() {
    if (this.getWrapEl().querySelectorAll('.atk-layer-item').length === 0) {
      this.$wrap.style.display = 'none'
    }
  }
}

export function getLayerWrap(): { $wrap: HTMLElement, $mask: HTMLElement } {
  let $wrap = document.querySelector<HTMLElement>(`.atk-layer-wrap`)
  if (!$wrap) {
    $wrap = Utils.createElement(
      `<div class="atk-layer-wrap" style="display: none;"><div class="atk-layer-mask"></div></div>`
    )
    document.body.appendChild($wrap)
  }

  const $mask = $wrap.querySelector<HTMLElement>('.atk-layer-mask')!

  return { $wrap, $mask }
}

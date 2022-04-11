import '../style/sidebar.less'

import Context from '@/context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import SidebarHTML from './html/sidebar.html?raw'
import * as Ui from '@/lib/ui'
import Layer from './layer'

export default class Sidebar extends Component {
  public layer?: Layer
  public $header: HTMLElement
  public $closeBtn: HTMLElement
  public $iframeWrap: HTMLIFrameElement

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(SidebarHTML)
    this.$header = this.$el.querySelector('.atk-sidebar-header')!
    this.$closeBtn = this.$header.querySelector('.atk-sidebar-close')!
    this.$iframeWrap = this.$el.querySelector('.atk-sidebar-iframe-wrap')!

    this.$closeBtn.onclick = () => {
      this.hide()
    }

    // event
    this.ctx.on('sidebar-show', () => (this.show()))
    this.ctx.on('sidebar-hide', () => (this.hide()))
    this.ctx.on('user-changed', () => { this.firstShow = true })

    console.log('hello Sidebar')
  }

  private firstShow = true

  /** 显示 */
  public async show() {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // 获取 Layer
    this.layer = new Layer(this.ctx, 'sidebar', this.$el)
    this.layer.afterHide = () => {
      // 防止评论框被吞
      if (this.ctx.conf.editorTravel === true) {
        this.ctx.trigger('editor-travel-back')
      }
    }
    this.layer.show()

    // viewWrap 滚动条归位
    // this.$viewWrap.scrollTo(0, 0)

    // 执行动画
    setTimeout(() => {
      this.$el.style.transform = 'translate(0, 0)'
    }, 20)

    // 第一次加载
    if (this.firstShow) {
      const $iframe = Utils.createElement<HTMLIFrameElement>('<iframe></iframe>')
      $iframe.src = "http://localhost:23366"
      this.$iframeWrap.append($iframe)
      this.firstShow = false
    }
  }

  /** 隐藏 */
  public hide() {
    // 执行动画
    this.$el.style.transform = ''

    this.layer?.hide()
  }
}

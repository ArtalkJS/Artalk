import Context from '~/types/context'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import SidebarHTML from './sidebar-layer.html?raw'
import User from '../lib/user'
import Layer from './layer'

export interface SidebarShowPayload {
  view?: 'comments'|'sites'|'pages'|'transfer'
}

export default class SidebarLayer extends Component {
  public layer?: Layer
  public $header: HTMLElement
  public $closeBtn: HTMLElement
  public $iframeWrap: HTMLElement
  public $iframe?: HTMLIFrameElement

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
    this.ctx.on('user-changed', () => { this.firstShow = true })
  }

  private firstShow = true

  /** 显示 */
  public async show(conf: SidebarShowPayload = {}) {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // 获取 Layer
    if (this.layer == null) {
      this.layer = new Layer(this.ctx, 'sidebar', this.$el)
      this.layer.afterHide = () => {
        // 防止评论框被吞
        if (this.ctx.conf.editorTravel === true) {
          this.ctx.editorTravelBack()
        }
      }
    }
    this.layer.show()

    // viewWrap 滚动条归位
    // this.$viewWrap.scrollTo(0, 0)

    // 管理员身份验证 (若身份失效，弹出验证窗口)
    ;(async () => {
      const resp = await this.ctx.getApi().user.loginStatus()
      if (resp.is_admin && !resp.is_login) {
        this.layer?.hide()
        this.firstShow = true

        this.ctx.checkAdmin({
          onSuccess: () => {
            setTimeout(() => {
              this.show(conf)
            }, 500)
          },
          onCancel: () => {}
        })
      }
    })()

    // 第一次加载
    if (this.firstShow) {
      this.$iframeWrap.innerHTML = ''
      this.$iframe = Utils.createElement<HTMLIFrameElement>('<iframe></iframe>')

      // 准备 Iframe 参数
      const baseURL = (import.meta.env.DEV)  ? 'http://localhost:23367/'
        : Utils.getURLBasedOnApi(this.ctx, '/sidebar/')

      const query: any = {
        pageKey: this.conf.pageKey,
        site: this.conf.site || '',
        user: JSON.stringify(User.data),
        time: +new Date()
      }

      if (conf.view) query.view = conf.view
      if (this.conf.darkMode) query.darkMode = '1'
      if (typeof this.conf.locale === 'string') query.locale = this.conf.locale

      const urlParams = new URLSearchParams(query);
      this.iframeLoad(`${baseURL}?${urlParams.toString()}`)

      this.$iframeWrap.append(this.$iframe)
      this.firstShow = false
    } else {
      // 夜间模式
      const isIframeSrcDarkMode = this.$iframe!.src.includes('darkMode=1')

      if (this.conf.darkMode && !isIframeSrcDarkMode)
        this.iframeLoad(`${this.$iframe!.src}&darkMode=1`)

      if (!this.conf.darkMode && isIframeSrcDarkMode)
        this.iframeLoad(this.$iframe!.src.replace('&darkMode=1', ''))
    }

    // 执行滑动显示动画
    setTimeout(() => {
      this.$el.style.transform = 'translate(0, 0)'
    }, 100)

    // 清空 unread
    setTimeout(() => {
      this.ctx.updateNotifies([])
    }, 0)

    this.ctx.trigger('sidebar-show')
  }

  /** 隐藏 */
  public hide() {
    // 执行动画
    this.$el.style.transform = ''

    this.layer?.hide()

    this.ctx.trigger('sidebar-hide')
  }

  private iframeLoad(src: string) {
    if (!this.$iframe) return

    this.$iframe.src = src

    // 加载动画
    Ui.showLoading(this.$iframeWrap)
    this.$iframe.onload = () => {
      Ui.hideLoading(this.$iframeWrap)
    }
  }
}

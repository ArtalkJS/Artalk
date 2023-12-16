import type { ContextApi, SidebarShowPayload } from '~/types'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import SidebarHTML from './sidebar-layer.html?raw'
import User from '../lib/user'
import type { Layer } from './layer'

export default class SidebarLayer extends Component {
  public layer?: Layer
  public $header: HTMLElement
  public $closeBtn: HTMLElement
  public $iframeWrap: HTMLElement
  public $iframe?: HTMLIFrameElement

  constructor(ctx: ContextApi) {
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
    this.initLayer()
    this.layer!.show()

    // 管理员身份验证 (若身份失效，弹出验证窗口)
    this.authCheck({
      onSuccess: () => this.show(conf) // retry show after auth check
    })

    // 第一次加载
    if (this.firstShow) {
      this.$iframeWrap.innerHTML = ''
      this.$iframe = this.createIframe(conf.view)
      this.$iframeWrap.append(this.$iframe)
      this.firstShow = false
    } else {
      const $iframe = this.$iframe!

      // 夜间模式
      const isIframeSrcDarkMode = $iframe.src.includes('darkMode=1')

      if (this.conf.darkMode && !isIframeSrcDarkMode)
        this.iframeLoad($iframe, `${this.$iframe!.src}&darkMode=1`)

      if (!this.conf.darkMode && isIframeSrcDarkMode)
        this.iframeLoad($iframe, this.$iframe!.src.replace('&darkMode=1', ''))
    }

    // 执行滑动显示动画
    setTimeout(() => {
      this.$el.style.transform = 'translate(0, 0)'
    }, 100)

    // 清空 unread
    setTimeout(() => {
      this.ctx.getData().updateUnreads([])
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

  // --------------------------------------------------

  private async authCheck(opts: { onSuccess: () => void }) {
    const resp = await this.ctx.getApi().user.loginStatus()
    if (resp.is_admin && !resp.is_login) {
      this.firstShow = true

      this.ctx.checkAdmin({
        onSuccess: () => {
          setTimeout(() => {
            opts.onSuccess()
          }, 500)
        },
        onCancel: () => {
          this.hide()
        }
      })
    }
  }

  private initLayer() {
    if (this.layer) return

    this.layer = this.ctx.get('layerManager').create('sidebar', this.$el)
    this.layer.setOnAfterHide(() => {
      // 防止评论框被吞
      this.ctx.editorResetState()
    })
  }

  private createIframe(view?: string) {
    const $iframe = Utils.createElement<HTMLIFrameElement>('<iframe></iframe>')

    // 准备 Iframe 参数
    const baseURL = (import.meta.env.DEV)  ? 'http://localhost:23367/'
      : Utils.getURLBasedOnApi({
        base: this.ctx.conf.server,
        path: '/sidebar/',
      })

    const query: any = {
      pageKey: this.conf.pageKey,
      site: this.conf.site || '',
      user: JSON.stringify(User.data),
      time: +new Date()
    }

    if (view) query.view = view
    if (this.conf.darkMode) query.darkMode = '1'
    if (typeof this.conf.locale === 'string') query.locale = this.conf.locale

    const urlParams = new URLSearchParams(query);
    this.iframeLoad($iframe, `${baseURL}?${urlParams.toString()}`)

    return $iframe
  }

  private iframeLoad($iframe: HTMLIFrameElement, src: string) {
    $iframe.src = src

    // 加载动画
    Ui.showLoading(this.$iframeWrap)
    $iframe.onload = () => {
      Ui.hideLoading(this.$iframeWrap)
    }
  }
}

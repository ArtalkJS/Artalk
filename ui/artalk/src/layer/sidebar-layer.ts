import SidebarHTML from './sidebar-layer.html?raw'
import type {
  SidebarShowPayload,
  ConfigManager,
  UserManager,
  LayerManager,
  Layer,
  SidebarLayer as ISidebarLayer,
  CheckerManager,
} from '@/types'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import type { Api } from '@/api'

export interface SidebarLayerOptions {
  onShow?: () => void
  onHide?: () => void

  getCheckers: () => CheckerManager
  getApi: () => Api
  getConf: () => ConfigManager
  getUser: () => UserManager
  getLayers: () => LayerManager
}

export class SidebarLayer implements ISidebarLayer {
  private opts: SidebarLayerOptions
  public $el: HTMLElement
  public layer?: Layer
  public $header: HTMLElement
  public $closeBtn: HTMLElement
  public $iframeWrap: HTMLElement
  public $iframe?: HTMLIFrameElement

  constructor(opts: SidebarLayerOptions) {
    this.opts = opts
    this.$el = Utils.createElement(SidebarHTML)
    this.$header = this.$el.querySelector('.atk-sidebar-header')!
    this.$closeBtn = this.$header.querySelector('.atk-sidebar-close')!
    this.$iframeWrap = this.$el.querySelector('.atk-sidebar-iframe-wrap')!

    this.$closeBtn.onclick = () => {
      this.hide()
    }
  }

  public async onUserChanged() {
    this.refreshWhenShow = true
  }

  /** Refresh iFrame when show */
  private refreshWhenShow = true

  /** Animation timer */
  private animTimer?: any = undefined

  /** 显示 */
  public async show(conf: SidebarShowPayload = {}) {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // init layer
    this.initLayer()
    this.layer!.show()

    // init iframe
    if (this.refreshWhenShow) {
      this.refreshWhenShow = false
      this.$iframeWrap.innerHTML = ''
      this.$iframe = this.createIframe(conf.view)
      this.$iframeWrap.append(this.$iframe)
    } else {
      // Sync Dark Mode (reload iframe if not match)
      const $iframe = this.$iframe!
      const iFrameSrc = $iframe.src
      if (this.getDarkMode() !== iFrameSrc.includes('&darkMode=1')) {
        this.iframeLoad(
          $iframe,
          iFrameSrc.replace(/&darkMode=\d/, `&darkMode=${Number(this.getDarkMode())}`),
        )
      }
    }

    // 管理员身份验证 (若身份失效，弹出验证窗口)
    this.authCheck({
      onSuccess: () => this.show(conf), // retry show after auth check
    })

    // 执行滑动显示动画
    this.animTimer = setTimeout(() => {
      this.animTimer = undefined
      this.$el.style.transform = 'translate(0, 0)'

      this.opts.onShow?.()
    }, 100)
  }

  /** 隐藏 */
  public hide() {
    this.layer?.hide()
  }

  // --------------------------------------------------

  private async authCheck(opts: { onSuccess: () => void }) {
    const data = (
      await this.opts.getApi().user.getUserStatus({
        ...this.opts.getApi().getUserFields(),
      })
    ).data
    if (data.is_admin && !data.is_login) {
      this.refreshWhenShow = true

      // show checker layer
      this.opts.getCheckers().checkAdmin({
        onSuccess: () => {
          setTimeout(() => {
            opts.onSuccess()
          }, 500)
        },
        onCancel: () => {
          this.hide()
        },
      })

      // hide sidebar layer
      this.hide()
    }
  }

  private initLayer() {
    if (this.layer) return

    this.layer = this.opts.getLayers().create('sidebar', this.$el)
    this.layer.setOnAfterHide(() => {
      // interrupt animation
      this.animTimer && clearTimeout(this.animTimer)

      // perform transform
      this.$el.style.transform = ''

      // trigger event
      this.opts.onHide?.()
    })
  }

  private createIframe(view?: string) {
    const $iframe = Utils.createElement<HTMLIFrameElement>(
      '<iframe referrerpolicy="strict-origin-when-cross-origin"></iframe>',
    )

    // 准备 Iframe 参数
    const baseURL = import.meta.env.DEV
      ? 'http://localhost:23367/'
      : Utils.getURLBasedOnApi({
          base: this.opts.getConf().get().server,
          path: '/sidebar/',
        })

    const query: any = {
      pageKey: this.opts.getConf().get().pageKey,
      site: this.opts.getConf().get().site || '',
      user: JSON.stringify(this.opts.getUser().getData()),
      time: +new Date(),
    }

    if (view) query.view = view
    query.darkMode = this.getDarkMode() ? '1' : '0'

    const urlParams = new URLSearchParams(query)
    this.iframeLoad($iframe, `${baseURL}?${urlParams.toString()}`)

    return $iframe
  }

  private getDarkMode() {
    return this.opts.getConf().get().darkMode === 'auto'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
      : this.opts.getConf().get().darkMode
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

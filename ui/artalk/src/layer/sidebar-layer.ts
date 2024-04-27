import type { ContextApi, SidebarShowPayload } from '@/types'
import Component from '@/lib/component'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import SidebarHTML from './sidebar-layer.html?raw'
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
    this.ctx.on('user-changed', () => {
      this.refreshOnShow = true
    })
  }

  /** Refresh iFrame on show */
  private refreshOnShow = true

  /** Animation timer */
  private animTimer?: any = undefined

  /** 显示 */
  public async show(conf: SidebarShowPayload = {}) {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // init layer
    this.initLayer()
    this.layer!.show()

    // init iframe
    if (this.refreshOnShow) {
      this.refreshOnShow = false
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
          this.getDarkMode()
            ? iFrameSrc.concat('&darkMode=1')
            : iFrameSrc.replace('&darkMode=1', ''),
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

      setTimeout(() => {
        this.ctx.getData().updateNotifies([])
      }, 0)

      this.ctx.trigger('sidebar-show')
    }, 100)
  }

  /** 隐藏 */
  public hide() {
    this.layer?.hide()
  }

  // --------------------------------------------------

  private async authCheck(opts: { onSuccess: () => void }) {
    const data = (
      await this.ctx.getApi().user.getUserStatus({
        ...this.ctx.getApi().getUserFields(),
      })
    ).data
    if (data.is_admin && !data.is_login) {
      this.refreshOnShow = true

      // show checker layer
      this.ctx.checkAdmin({
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

    this.layer = this.ctx.get('layerManager').create('sidebar', this.$el)
    this.layer.setOnAfterHide(() => {
      // 防止评论框被吞
      this.ctx.editorResetState()

      // interrupt animation
      this.animTimer && clearTimeout(this.animTimer)

      // perform transform
      this.$el.style.transform = ''

      // trigger event
      this.ctx.trigger('sidebar-hide')
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
          base: this.ctx.conf.server,
          path: '/sidebar/',
        })

    const query: any = {
      pageKey: this.conf.pageKey,
      site: this.conf.site || '',
      user: JSON.stringify(this.ctx.get('user').getData()),
      time: +new Date(),
    }

    if (view) query.view = view
    if (this.getDarkMode()) query.darkMode = '1'

    const urlParams = new URLSearchParams(query)
    this.iframeLoad($iframe, `${baseURL}?${urlParams.toString()}`)

    return $iframe
  }

  private getDarkMode() {
    return this.conf.darkMode === 'auto'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
      : this.conf.darkMode
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

import '../style/site-list.less'

import Context from 'artalk/src/context'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import { SiteData } from 'artalk/types/artalk-data'
import Api from 'artalk/src/api'

interface SiteListFloaterConf {
  onSwitchSite: (siteName: string) => boolean|void
  onClickSitesViewBtn: () => void
}

export default class SiteListFloater {
  ctx: Context
  conf: SiteListFloaterConf
  $el: HTMLElement
  sites: SiteData[] = []
  $sites: HTMLElement

  constructor(ctx: Context, conf: SiteListFloaterConf) {
    this.ctx = ctx
    this.conf = conf
    this.$el = Utils.createElement(
    `<div class="atk-site-list-floater" style="display: none;">
      <div class="atk-sites"></div>
    </div>`)
    this.$sites = this.$el.querySelector('.atk-sites')!
  }

  /** 装载 */
  public async load(selectedSite?: string) {
    this.$sites.innerHTML = ''

    const renderSiteItem = (siteName: string, siteLogo: string, siteTarget?: string, onclick?: Function) => {
      const $site = Utils.createElement(
        `<div class="atk-site-item">
          <div class="atk-site-logo"></div>
          <div class="atk-site-name"></div>
        </div>`)
        $site.onclick = !onclick ? () => this.switch(siteTarget || siteName) : () => onclick()
        $site.setAttribute('data-name', siteTarget || siteName)
        const $siteLogo = $site.querySelector<HTMLElement>('.atk-site-logo')!
        const $siteName = $site.querySelector<HTMLElement>('.atk-site-name')!
        $siteLogo.innerText = siteLogo
        $siteName.innerText = siteName
        if (selectedSite === (siteTarget || siteName)) $site.classList.add('atk-active')
        this.$sites.append($site)
    }

    renderSiteItem('所有站点', '_', '__ATK_SITE_ALL')

    const sites = await new Api(this.ctx).siteGet()
    sites.forEach((site) => {
      renderSiteItem(site.name, site.name.substring(0, 1))
    })

    renderSiteItem('站点管理', '+', '', () => { this.conf.onClickSitesViewBtn();this.hide() })
  }

  /** 切换站点 */
  private switch(siteName: string) {
    if (this.conf.onSwitchSite(siteName) === false) { return }

    // set active
    this.$sites.querySelectorAll('.atk-site-item').forEach(e => {
      if (e.getAttribute('data-name') !== siteName) {
        e.classList.remove('atk-active')
      } else {
        e.classList.add('atk-active')
      }
    })

    this.hide()
  }

  private outsideChecker?: (evt: MouseEvent) => void // outside checker

  /** 显示 */
  public show($trigger?: HTMLElement) {
    this.$el.style.display = ''

    // 点击外部隐藏
    if ($trigger) {
      this.outsideChecker = (evt: MouseEvent) => {
        const isClickInside = $trigger.contains(evt.target as any)||this.$el.contains(evt.target as any)
        if (!isClickInside) { this.hide() }
      }
      document.addEventListener('click', this.outsideChecker)
    }
  }

  /** 隐藏 */
  public hide() {
    this.$el.style.display = 'none'
    if (this.outsideChecker) document.removeEventListener('click', this.outsideChecker)
  }
}

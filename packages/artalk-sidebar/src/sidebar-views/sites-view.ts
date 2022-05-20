import SiteList from '../components/site-list'
import SidebarView from './sidebar-view'

export default class SitesView extends SidebarView {
  protected readonly viewName = 'sites'
  public readonly viewAdminOnly = true
  viewTitle() { return '站点' }

  protected tabs = {}
  protected activeTab = ''

  siteList!: SiteList

  mount() {
    if (!this.siteList) {
      this.siteList = new SiteList(this.ctx, this.sidebar)
      this.$el.append(this.siteList.$el)
    }

    this.reqSites()
  }

  switchTab(tab: string) {
    this.reqSites()
  }

  async reqSites() {
    const sites = await this.ctx.getApi().siteGet()
    this.siteList.loadSites(sites)
  }
}

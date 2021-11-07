import Api from '~/src/api'
import Context from '../../context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import SiteList from '../admin/site-list'
import Comment from '../comment'
import SidebarView from '../sidebar-view'

export default class SitesView extends SidebarView {
  static viewName = 'sites'
  static viewTitle = '站点'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  siteList!: SiteList

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  async mount(siteName: string) {
    // TODO 多次重复import issue
    if (!this.siteList) {
      this.siteList = new SiteList(this.ctx)
      this.$el.append(this.siteList.$el)
    }

    const sites = await new Api(this.ctx).siteGet()
    console.log(sites)
    this.siteList.importSites(sites)
  }

  switchTab(tab: string, siteName: string): boolean|void {}
}

import Api from '~/src/api'
import Context from '../../context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import PageList from '../admin/page-list'
import Comment from '../comment'
import SidebarView from '../sidebar-view'

export default class PagesView extends SidebarView {
  static viewName = 'pages'
  static viewTitle = '页面'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  pageList!: PageList

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  async mount() {
    if (!this.pageList) {
      this.pageList = new PageList(this.ctx)
      this.$el.append(this.pageList.$el)
    }

    // TODO for testing
    const pages = await new Api(this.ctx).pageGet('ArtalkDemo')
    console.log(pages)
    this.pageList.importPages(pages)
  }

  switch(tab: string): boolean|void {}
}

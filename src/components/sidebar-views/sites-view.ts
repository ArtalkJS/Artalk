import Context from '../../context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import Comment from '../comment'
import SidebarView from '../sidebar-view'

export default class SitesView extends SidebarView {
  static viewName = 'sites'
  static viewTitle = '站点'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  mount() {}

  switch(tab: string): boolean|void {}
}

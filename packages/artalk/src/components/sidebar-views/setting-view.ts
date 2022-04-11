import Api from '~/src/api'
import Context from '../../context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import SiteList from '../admin/site-list'
import Comment from '../comment'
import SidebarView from '../sidebar-view'

export default class SettingView extends SidebarView {
  static viewName = 'setting'
  static viewTitle = '配置'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  mount(siteName: string) {
  }

  switchTab(tab: string, siteName: string) {}
}

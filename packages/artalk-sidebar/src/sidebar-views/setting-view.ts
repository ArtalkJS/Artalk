import Api from 'artalk/src/api'
import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import Comment from 'artalk/src/comment'

import SiteList from '../admin/site-list'
import SidebarView from '../sidebar-view'

export default class SettingView extends SidebarView {
  static viewName = 'setting'
  static viewTitle = '配置'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  mount() {
  }

  switchTab(tab: string) {}
}

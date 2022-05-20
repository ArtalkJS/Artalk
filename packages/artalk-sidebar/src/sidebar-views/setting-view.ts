import SidebarView from './sidebar-view'

export default class SettingView extends SidebarView {
  protected readonly viewName = 'setting'
  public readonly viewAdminOnly = true
  viewTitle() { return '配置' }

  protected tabs: { [key: string]: string } = {}
  protected activeTab = ''

  mount() {
  }

  switchTab(tab: string) {}
}

import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import Comment from './comment'

export default class SidebarView extends Component {
  static viewName = ''
  static viewTitle = ''
  static viewAdminOnly = false

  viewTabs = {}
  viewActiveTab = ''

  protected $parent: HTMLElement

  constructor(ctx: Context, $parent: HTMLElement) {
    super(ctx)

    this.$parent = $parent
    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  mount(siteName: string) {}

  switchTab(tab: string, siteName: string): boolean|void {}
}

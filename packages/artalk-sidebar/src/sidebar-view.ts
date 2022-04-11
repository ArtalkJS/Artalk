import Context from 'artalk/src/context'
import Component from 'artalk/src/lib/component'
import * as Utils from 'artalk/src/lib/utils'
import Comment from 'artalk/src/components/comment'

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

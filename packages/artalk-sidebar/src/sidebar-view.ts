import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import Comment from 'artalk/src/comment'
import Component from './sidebar-component'
import { SidebarCtx } from './main'

export default class SidebarView extends Component {
  static viewName = ''
  static viewTitle = ''
  static viewAdminOnly = false

  viewTabs = {}
  viewActiveTab = ''

  protected $parent: HTMLElement

  constructor(ctx: Context, sidebar: SidebarCtx, $parent: HTMLElement) {
    super(ctx, sidebar)

    this.$parent = $parent
    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  mount() {}

  switchTab(tab: string): boolean|void {}
}

import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import Component from '../sidebar-component'
import { SidebarCtx } from '../main'

interface SidebarView {
  ctx: Context
  viewAdminOnly?: boolean
}

abstract class SidebarView extends Component {
  protected $parent: HTMLElement

  public constructor(ctx: Context, sidebar: SidebarCtx, $parent: HTMLElement) {
    super(ctx, sidebar)

    this.$parent = $parent
    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  protected abstract tabs: {[key: string]: string}
  protected abstract activeTab: string
  protected abstract readonly viewName: string
  public abstract viewTitle(): string

  public abstract mount(): void
  public abstract switchTab(tab: string): boolean|void

  public getName() {
    return this.viewName
  }

  public getTabs() {
    return this.tabs
  }

  public getActiveTab() {
    return this.activeTab
  }
}

export default SidebarView

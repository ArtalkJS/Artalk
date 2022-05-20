import ArtalkComponent from 'artalk/src/lib/component'
import Context from 'artalk/types/context'
import { SidebarCtx } from './main'

export default class Component extends ArtalkComponent {
  declare public ctx: Context
  sidebar: SidebarCtx

  constructor(ctx: Context, sidebar: SidebarCtx) {
    super(ctx)

    this.sidebar = sidebar
  }
}

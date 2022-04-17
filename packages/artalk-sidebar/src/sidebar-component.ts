import ArtalkComponent from 'artalk/src/lib/component'
import Context from 'artalk/src/context'
import { SidebarCtx } from './main'

export default class Component extends ArtalkComponent {
  sidebar: SidebarCtx

  constructor(ctx: Context, sidebar: SidebarCtx) {
    super(ctx)

    this.sidebar = sidebar
  }
}

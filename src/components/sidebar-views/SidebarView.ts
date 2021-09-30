import Context from '../../Context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'

export default class SidebarView extends Component {
  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`<div class="atk-sidebar-view-${this.name}"></div>`)
  }

  name = ''
  title = this.name
  actions: { [name: string]: string } = {}
  activeAction?: string

  adminOnly = false

  render(): HTMLElement {
    return this.el
  }

  switch(action: string): void {}
}

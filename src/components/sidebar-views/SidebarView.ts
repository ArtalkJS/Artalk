import Context from '../../Context'
import Component from '../../lib/component'
import * as Utils from '../../lib/utils'
import Comment from '../Comment'

export default class SidebarView extends Component {
  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  name = ''
  title = this.name
  actions: { [name: string]: string } = {}
  activeAction?: string

  adminOnly = false

  init (): void {}

  switch (action: string): void {}
}


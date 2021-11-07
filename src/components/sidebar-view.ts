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

  constructor(ctx: Context) {
    super(ctx)

    this.$el = Utils.createElement(`<div class="atk-sidebar-view"></div>`)
  }

  mount() {}

  switch(tab: string): boolean|void {}
}
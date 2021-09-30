import Context from '../../Context'
import SidebarView from './SidebarView'
import * as Utils from '../../lib/utils'
import Comment from '../Comment'
import ListLite from '../ListLite'

export default class AdminView extends SidebarView {
  name = 'admin'
  title = '控制台'
  actions = {
    comment: '评论',
    page: '页面',
    site: '网站',
    conf: '配置',
  }

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`<div class="atk-msg-center"></div>`)
  }

  mount(action: string) {

  }
}

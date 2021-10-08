import Context from '../../Context'
import SidebarView from './SidebarView'
import * as Utils from '../../lib/utils'
import Comment from '../Comment'
import ListLite from '../ListLite'
import { CreateCommentList } from './MessageView'

export default class AdminView extends SidebarView {
  name = 'admin'
  title = '控制台'
  actions = {
    comment: '评论',
    page: '页面',
    site: '网站',
    conf: '配置',
  }
  activeAction = 'comment'
  cList?: ListLite

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`<div class="atk-msg-center"></div>`)
  }

  render () {
    this.cList = CreateCommentList(this.ctx)

    this.switch('comment')

    return this.el
  }

  switch (action: string) {
    this.el.innerHTML = ''

    if (action === 'comment' && this.cList) {
      this.el.append(this.cList.el)

      const reqComments = (type = 'all') => {
        if (!this.cList) return
        this.cList.type = `admin_${type}` as any
        this.cList.isFirstLoad = true
        this.cList.reqComments()
      }

      this.showFilterBar({ all: '全部', pending: '待审' }, (name) => {
        reqComments(name)
      })
    }
  }

  showFilterBar (items: {[name: string]: string}, clickEvt: (item) => void) {
    const filterBarEl = Utils.createElement(`<div class="atk-filter-bar"></div>`)
    this.el.prepend(filterBarEl)

    Object.entries(items).forEach(([name, label]) => {
      const itemEl = Utils.createElement(`<span class="atk-filter-item"></span>`)
      itemEl.innerText = label
      itemEl.addEventListener('click', () => {
        clickEvt(name)
        filterBarEl.querySelectorAll('.atk-active')
          .forEach(item => item.classList.remove('atk-active')) // 删除其他 active
        itemEl.classList.add('atk-active')
      })
      filterBarEl.append(itemEl)
    })

    ;(filterBarEl.firstChild as HTMLElement).click() // 默认打开第一个项目
  }
}

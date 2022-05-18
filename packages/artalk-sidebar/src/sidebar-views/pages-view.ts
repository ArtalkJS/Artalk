import Api from 'artalk/src/api'
import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import Comment from 'artalk/src/comment'
import Pagination, { PaginationConf } from 'artalk/src/components/pagination'

import SidebarView from '../sidebar-view'
import PageList from '../admin/page-list'

const PAGE_SIZE = 20

export default class PagesView extends SidebarView {
  static viewName = 'pages'
  static viewTitle = '页面'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  isFirstLoad = true
  pageList!: PageList
  pagination!: Pagination

  mount() {
    if (!this.pageList) {
      this.pageList = new PageList(this.ctx, this.sidebar)
      this.$el.append(this.pageList.$el)
    }

    this.switchTab(this.viewActiveTab)
  }

  switchTab(tab: string): boolean|void {
    this.reqPages(0)
  }

  async reqPages(offset: number) {
    if (this.isFirstLoad) this.pageList.initPageList()
    ;(this.$el.parentNode as any)?.scrollTo(0, 0)

    if (this.isFirstLoad) Ui.showLoading(this.$el)
    else this.pagination.setLoading(true)

    const data = await new Api(this.ctx).pageGet(this.sidebar.curtSite, offset, PAGE_SIZE)
    this.pageList.$pageList!.innerHTML = ''
    this.pageList.importPages(data.pages || [])

    if (this.isFirstLoad) Ui.hideLoading(this.$el)
    else this.pagination.setLoading(false)

    if (this.isFirstLoad) {
      this.pagination = new Pagination(data.total, {
        pageSize: PAGE_SIZE,
        onChange: (o: number) => {
          this.reqPages(o)
        }
      })

      this.$el.append(this.pagination.$el)
      this.isFirstLoad = false
    } else {
      if (offset === 0) this.pagination.update(offset, data.total)
    }
  }
}

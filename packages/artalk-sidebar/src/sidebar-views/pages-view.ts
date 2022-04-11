import Api from 'artalk/src/api'
import Context from 'artalk/src/context'
import Component from 'artalk/src/lib/component'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'
import Comment from 'artalk/src/components/comment'
import Pagination, { PaginationConf } from 'artalk/src/components/pagination'

import PageList from '../admin/page-list'
import SidebarView from '../sidebar-view'

const PAGE_SIZE = 20

export default class PagesView extends SidebarView {
  static viewName = 'pages'
  static viewTitle = '页面'
  static viewAdminOnly = true

  viewTabs = {}
  viewActiveTab = ''

  pageList!: PageList
  pagination!: Pagination

  mount(siteName: string) {
    if (!this.pageList) {
      this.pageList = new PageList(this.ctx)
      this.$el.append(this.pageList.$el)
    }

    this.switchTab(this.viewActiveTab, siteName)
  }

  switchTab(tab: string, siteName: string): boolean|void {
    this.reqPages(siteName, 0)
  }

  async reqPages(siteName: string, offset: number) {
    this.pageList.clearAll()
    ;(this.$el.parentNode as any)?.scrollTo(0, 0)

    Ui.showLoading(this.$el)

    const data = await new Api(this.ctx).pageGet(siteName, offset, PAGE_SIZE)
    this.pageList.importPages(data.pages || [])

    Ui.hideLoading(this.$el)

    if (!this.pagination) {
      this.pagination = new Pagination(data.total, {
        pageSize: PAGE_SIZE,
        onChange: (o: number) => {
          this.reqPages(siteName, o)
        }
      })

      this.$el.append(this.pagination.$el)
    }
    if (this.pagination && offset === 0)
      this.pagination.update(offset, data.total)
  }
}

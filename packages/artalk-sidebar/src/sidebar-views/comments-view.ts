import ListLite from 'artalk/src/list/list-lite'
import SidebarView from './sidebar-view'

export default class MessageView extends SidebarView {
  protected readonly viewName = 'comments'
  viewTitle() { return '评论' }

  protected tabs = {}
  protected activeTab = ''

  private list!: ListLite

  mount() {
    // tabs
    if (this.ctx.user.data.isAdmin) {
      this.tabs = {
        admin_all: '全部',
        admin_pending: '待审',
        all: '个人',
      }
      this.activeTab = 'admin_all'
    } else {
      this.tabs = {
        mentions: '提及',
        all: '全部',
        mine: '我的',
        pending: '待审',
      }
      this.activeTab = 'mentions'
    }

    this.list = new ListLite(this.ctx)
    this.ctx.setList(this.list)
    this.list.flatMode = true
    this.list.unreadHighlight = true
    this.list.scrollListenerAt = this.$parent
    this.list.pageMode = 'pagination'
    this.list.noCommentText = '<div class="atk-sidebar-no-content">无内容</div>'
    this.list.renderComment = (comment) => {
      const pageURL = comment.getData().page_url
      comment.getRender().setOpenURL(`${pageURL}#atk-comment-${comment.getID()}`)
      comment.getConf().onReplyBtnClick = () => {
        this.ctx.replyComment(comment.getData(), comment.getEl(), true)
      }
    }
    this.list.paramsEditor = (params) => {
      params.site_name = this.sidebar.curtSite
    }
    this.ctx.on('list-inserted', (data) => {
      ;(this.$el.parentNode as any)?.scrollTo(0, 0)
    })

    this.$el.innerHTML = ''
    this.$el.append(this.list.$el)

    this.switchTab(this.activeTab)
  }

  switchTab(tab: string): boolean|void {
    //console.log(tab, siteName)
    this.activeTab = tab
    this.list.paramsEditor = (params) => {
      params.type = tab as any
      params.site_name = this.sidebar.curtSite
    }
    this.list.fetchComments(0)

    return true
  }
}


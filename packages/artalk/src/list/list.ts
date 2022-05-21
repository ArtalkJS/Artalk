import '../style/list.less'

import { ListData, NotifyData } from '~/types/artalk-data'
import Context from '~/types/context'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import ListHTML from './list.html?raw'
import ListLite from './list-lite'

export default class List extends ListLite {
  private $closeCommentBtn!: HTMLElement
  private $openSidebarBtn!: HTMLElement
  private $unreadBadge!: HTMLElement
  private $commentCount!: HTMLElement
  private $commentCountNum!: HTMLElement
  private $dropdownWrap?: HTMLElement

  constructor (ctx: Context) {
    const el = Utils.createElement(ListHTML)

    super(ctx)

    // 把 listLite $el 变为子元素
    el.querySelector('.atk-list-body')!.append(this.$el)
    this.$el = el

    // 平铺模式
    if (this.ctx.conf.flatMode === true) {
      this.flatMode = true // 遵循配置总是开启平铺模式
    } else if (this.conf.nestMax && this.conf.nestMax <= 1) {
      this.flatMode = true // 嵌套层数值无效时，强制开启平铺模式
    } else if (this.ctx.conf.flatMode === 'auto') {
      // 自动判断启用平铺模式
      if (window.matchMedia("(max-width: 768px)").matches)
        this.flatMode = true
    }

    // 分页模式
    this.pageMode = this.conf.pagination.readMore ? 'read-more' : 'pagination'
    this.pageSize = this.conf.pagination.pageSize || 20
    this.repositionAt = this.$el

    // 操作按钮
    this.initListActionBtn()

    this.$commentCount = this.$el.querySelector('.atk-comment-count')!
    this.$commentCount.innerHTML = this.$t('counter', {
      count: '<span class="atk-comment-count-num">0</span>',
    })
    this.$commentCountNum = this.$commentCount.querySelector('.atk-comment-count-num')!

    // 评论列表排序 Dropdown 下拉选择层
    if (this.ctx.conf.listSort) {
      this.initDropdown()
    }

    // copyright
    this.$el.querySelector<HTMLElement>('.atk-copyright')!.innerHTML = `Powered By <a href="https://artalk.js.org" target="_blank" title="Artalk v${ARTALK_VERSION}">Artalk</a>`

    // event listen
  }

  private initListActionBtn() {
    // 侧边栏呼出按钮
    this.$openSidebarBtn = this.$el.querySelector('[data-action="open-sidebar"]')!
    this.$closeCommentBtn = this.$el.querySelector('[data-action="admin-close-comment"]')!
    this.$unreadBadge = this.$el.querySelector('.atk-unread-badge')!

    this.$openSidebarBtn.addEventListener('click', () => {
      this.ctx.showSidebar()
    })

    // 关闭评论按钮
    this.$closeCommentBtn.addEventListener('click', () => {
      if (!this.data) return

      this.data.page.admin_only = !this.data.page.admin_only
      this.adminPageEditSave()
    })
  }

  /** 刷新界面 */
  public refreshUI() {
    super.refreshUI()

    this.$commentCountNum.innerText = String(Number(this.data?.total) || 0)

    // 已输入个人信息
    if (!!this.ctx.user.data.nick && !!this.ctx.user.data.email) {
      this.$openSidebarBtn.classList.remove('atk-hide')
    } else {
      this.$openSidebarBtn.classList.add('atk-hide')
    }

    // 仅管理员显示控制
    this.ctx.checkAdminShowEl()
    this.$openSidebarBtn.querySelector<HTMLElement>('.atk-text')!
      .innerText = (!this.ctx.user.data.isAdmin) ? this.$t('msgCenter') : this.$t('ctrlCenter')

    // 关闭评论
    if (!!this.data && !!this.data.page && this.data.page.admin_only === true) {
      this.ctx.editorClose()
      this.$closeCommentBtn.innerHTML = this.$t('openComment')
    } else {
      this.ctx.editorOpen()
      this.$closeCommentBtn.innerHTML = this.$t('closeComment')
    }
  }

  protected onLoad(data: ListData, offset: number) {
    super.onLoad(data, offset)

    // 检测锚点跳转
    if (!this.goToCommentFounded) this.checkGoToCommentByUrlHash()

    // 防止评论框被吞
    if (this.ctx.conf.editorTravel === true) {
      this.ctx.editorTravelBack()
    }
  }

  private goToCommentFounded = false
  public goToCommentDelay = true

  /** 跳到评论项位置 - 根据 `location.hash` */
  public checkGoToCommentByUrlHash() {
    let commentId = Number(Utils.getQueryParam('atk_comment')) // same as backend GetReplyLink()
    if (!commentId) {
      const match = window.location.hash.match(/#atk-comment-([0-9]+)/)
      if (!match || !match[1] || Number.isNaN(Number(match[1]))) return
      commentId = Number(match[1])
    }
    if (!commentId) return

    const comment = this.ctx.findComment(commentId)
    if (!comment) { // 若找不到评论
      // 自动翻页
      if (this.readMoreBtn) this.readMoreBtn.click()
      else if (this.pagination) this.pagination.next()
      return
    }

    // 已阅 API
    const notifyKey = Utils.getQueryParam('atk_notify_key')
    if (notifyKey) {
      this.ctx.getApi().markRead(notifyKey)
        .then(() => {
          this.unread = this.unread.filter(o => o.comment_id !== commentId)
          this.updateUnread(this.unread)
        })
    }

    // 若父评论存在 “子评论部分” 限高，取消限高
    comment.getParents().forEach((p) => {
      p.getRender().heightLimitRemoveForChildren()
    })

    const goTo = () => {
      Ui.scrollIntoView(comment.getEl(), false)

      comment.getEl().classList.remove('atk-flash-once')
      window.setTimeout(() => {
        comment.getEl().classList.add('atk-flash-once')
      }, 150)
    }

    if (!this.goToCommentDelay) goTo()
    else window.setTimeout(() => goTo(), 350)

    this.goToCommentFounded = true
  }

  /** 管理员设置页面信息 */
  public adminPageEditSave () {
    if (!this.data || !this.data.page) return

    this.ctx.editorShowLoading()
    this.ctx.getApi().pageEdit(this.data.page)
      .then((page) => {
        if (this.data)
          this.data.page = { ...page }
        this.refreshUI()
      })
      .catch(err => {
        this.ctx.editorShowNotify(`${this.$t('editFail')}: ${err.msg || String(err)}`, 'e')
      })
      .finally(() => {
        this.ctx.editorHideLoading()
      })
  }

  public showUnreadBadge(count: number) {
    if (count > 0) {
      this.$unreadBadge.innerText = `${Number(count || 0)}`
      this.$unreadBadge.style.display = 'block'
    } else {
      this.$unreadBadge.style.display = 'none'
    }
  }

  /** 初始化选择下拉层 */
  protected initDropdown() {
    this.$dropdownWrap = this.$commentCount
    this.$commentCount.classList.add('atk-dropdown-wrap')

    // 插入图标
    this.$dropdownWrap.append(Utils.createElement(`<span class="atk-arrow-down-icon"></span>`))

    const reloadUseParamsEditor = (func: (p: any) => void) => {
      this.paramsEditor = (p) => { func(p) }
      this.fetchComments(0)
    }

    // 下拉列表
    const dropdownList = [
      [this.$t('sortLatest'), () => { reloadUseParamsEditor(p => { p.sort_by = 'date_desc' }) }],
      [this.$t('sortBest'), () => { reloadUseParamsEditor(p => { p.sort_by = 'vote' }) }],
      [this.$t('sortOldest'), () => { reloadUseParamsEditor(p => { p.sort_by = 'date_asc' }) }],
      [this.$t('sortAuthor'), () => { reloadUseParamsEditor(p => { p.view_only_admin = true }) }],
    ]

    // 列表项点击事件
    let curtActive = 0 // 当前选中
    const onItemClick = (i: number, $item: HTMLElement, name: string, action: Function) => {
      action()

      // set active
      curtActive = i
      $dropdown.querySelectorAll('.active').forEach((e) => { e.classList.remove('active') })
      $item.classList.add('active')

      // 关闭层 (临时消失，取消 :hover)
      $dropdown.style.display = 'none'
      setTimeout(() => { $dropdown.style.display = '' }, 80)
    }

    // 生成列表元素
    const $dropdown = Utils.createElement(`<ul class="atk-dropdown atk-fade-in"></ul>`)
    dropdownList.forEach((item, i) => {
      const name = item[0] as string
      const action = item[1] as Function

      const $item = Utils.createElement(`<li class="atk-dropdown-item"><span></span></li>`)
      const $link = $item.querySelector<HTMLElement>('span')!
      $link.innerText = name
      $link.onclick = () => { onItemClick(i, $item, name, action) }
      $dropdown.append($item)

      if (i === curtActive) $item.classList.add('active') // 默认选中项
    })

    this.$dropdownWrap.append($dropdown)
  }

  public updateUnread(notifies: NotifyData[]): void {
    super.updateUnread(notifies)
    this.showUnreadBadge(notifies?.length || 0)
  }
}

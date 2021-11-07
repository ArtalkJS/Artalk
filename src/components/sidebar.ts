import '../style/sidebar.less'

import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment  from './comment'
import SidebarHTML from './html/sidebar.html?raw'
import Layer from './layer'

import SidebarView from './sidebar-view'
import CommentsView from './sidebar-views/comments-view'
import PagesView from './sidebar-views/pages-view'
import SitesView from './sidebar-views/sites-view'
import { SiteData } from '~/types/artalk-data'
import Pagination from './pagination'

export default class Sidebar extends Component {
  public layer?: Layer

  public $header: HTMLElement
  public $headerMenu: HTMLElement
  public $avatar: HTMLElement
  public $closeBtn: HTMLElement
  public $nav: HTMLElement
  public $curtViewBtn: HTMLElement
  public $curtViewBtnIcon: HTMLElement
  public $curtViewBtnText: HTMLElement
  public $navTabs: HTMLElement
  public $navViews: HTMLElement
  public $viewWrap: HTMLElement

  public siteList: SiteData[] = []

  public curtSite?: SiteData
  public curtView?: string

  public viewInstances: {[name: string]: SidebarView} = {}
  public registerViews: (typeof SidebarView)[] = [
    CommentsView, PagesView, SitesView
  ]

  private viewSwitcherShow = false

  constructor(ctx: Context) {
    super(ctx)

    // initial elements
    this.$el = Utils.createElement(SidebarHTML)
    this.$header = this.$el.querySelector('.atk-sidebar-header')!
    this.$headerMenu = this.$header.querySelector('.atk-menu')!
    this.$avatar = this.$header.querySelector('.atk-avatar')!
    this.$closeBtn = this.$header.querySelector('.atk-sidebar-close')!

    this.$nav = this.$el.querySelector('.atk-sidebar-nav')!
    this.$curtViewBtn = this.$nav.querySelector('.akt-curt-view-btn')!
    this.$curtViewBtnIcon = this.$curtViewBtn.querySelector('.atk-icon')!
    this.$curtViewBtnText = this.$curtViewBtn.querySelector('.atk-text')!
    this.$navTabs = this.$nav.querySelector('.atk-tabs')!
    this.$navViews = this.$nav.querySelector('.atk-views')!

    this.$viewWrap = this.$el.querySelector('.atk-sidebar-view-wrap')!

    // 关闭按钮
    this.$closeBtn.onclick = () => {
      this.hide()
    }

    // event
    this.ctx.on('sidebar-show', () => (this.show()))
    this.ctx.on('sidebar-hide', () => (this.hide()))

    // TODO for testing
    this.show()
  }

  private isFirstShow = true

  /** 显示 */
  public show() {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // 获取 Layer
    this.layer = new Layer(this.ctx, 'sidebar', this.$el)
    this.layer.show()

    // viewWrap 滚动条归位
    this.$viewWrap.scrollTo(0, 0)

    // 执行动画
    setTimeout(() => {
      this.$el.style.transform = 'translate(0, 0)'
    }, 20)

    // 第一次加载
    if (this.isFirstShow) {
      this.switchView('comments') // 打开默认 view
      this.initUI()
      this.isFirstShow = false
    }
  }

  private initUI() {
    this.initViewSwitcher()
  }

  /** 初始化页面切换器 */
  private initViewSwitcher() {
    this.$curtViewBtn.onclick = () => {
      this.toggleViewSwitcher()
    }

    this.$navViews.innerHTML = ''
    this.registerViews.forEach(view => {
      const $item = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navViews.append($item)
      $item.innerText = view.viewTitle
      if (view.viewName === this.curtView) {
        $item.classList.add('atk-active')
        this.$curtViewBtnText.innerText = view.viewTitle
      }
      $item.onclick = () => {
        // 切换 view
        this.switchView(view.viewName)

        this.$navViews.querySelectorAll('.atk-active').forEach((e) => e.classList.remove('atk-active'))
        $item.classList.add('atk-active')

        this.$curtViewBtnText.innerText = view.viewTitle

        this.toggleViewSwitcher()
      }
    })
  }

  private toggleViewSwitcher() {
    if (!this.viewSwitcherShow) {
      // 显示
      this.$navViews.style.display = ''
      this.$navTabs.style.display = 'none'
      this.$curtViewBtnIcon.classList.add('atk-arrow')
    } else {
      // 隐藏
      this.$navViews.style.display = 'none'
      this.$navTabs.style.display = ''
      this.$curtViewBtnIcon.classList.remove('atk-arrow')
    }

    this.viewSwitcherShow = !this.viewSwitcherShow
  }

  /** 隐藏 */
  public hide() {
    // 执行动画
    this.$el.style.transform = ''

    // 用完即销毁
    this.layer?.dispose()
  }

  /** 切换 View */
  public switchView(viewName: string) {
    let view = this.viewInstances[viewName]
    if (!view) {
      // 初始化 View
      const View = this.registerViews.find(o => o.viewName === viewName)!
      view = new View(this.ctx)
      this.viewInstances[viewName] = view
    }

    this.curtView = viewName

    // update tabs
    this.initNavTabs()

    // init view ui
    view.mount()
    this.$viewWrap.innerHTML = ''
    this.$viewWrap.append(view.$el)

    const p = new Pagination({
      total: 20,
      onChange: (offset) => {
        console.log(offset)
      }
    })
    this.$viewWrap.append(p.$el)
  }

  private initNavTabs() {
    if (!this.curtView) return

    const view = this.viewInstances[this.curtView]
    if (!view) return

    this.$navTabs.innerHTML = ''
    Object.entries<string>(view.viewTabs).forEach(([tabName, label]) => {
      const $tab = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navTabs.append($tab)
      $tab.innerText = label
      if (view.viewActiveTab === tabName) $tab.classList.add('atk-active')

      // 切换 tab
      $tab.onclick = () => {
        this.$navTabs.querySelectorAll('.atk-active').forEach(e => e.classList.remove('atk-active'))
        $tab.classList.add('atk-active')
        view.switch(tabName)
      }
    })
  }
}

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
import Api from '../api'

const DEFAULT_VIEW = 'comments'
const SITE_ALL_NAME = '__ATK_SITE_ALL'

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
  public $siteSwitcher: HTMLElement
  public $siteSwitcherSites: HTMLElement
  public $viewWrap: HTMLElement

  public curtSite: string = SITE_ALL_NAME
  public curtView: string = DEFAULT_VIEW
  public get curtViewInstance() {
    return this.curtView ? this.viewInstances[this.curtView] : undefined
  }
  public curtTab?: string

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
    this.$siteSwitcher = this.$el.querySelector('.atk-site-switcher')!
    this.$siteSwitcherSites = this.$siteSwitcher.querySelector('.atk-sites')!

    this.$viewWrap = this.$el.querySelector('.atk-sidebar-view-wrap')!

    // init UI
    this.initViewSwitcher()

    if (this.ctx.user.data.isAdmin) {
      this.$avatar.onclick = (evt) => this.showSiteSwitcher(evt.target as any)
    }

    this.$closeBtn.onclick = () => {
      this.hide()
    }

    // event
    this.ctx.on('sidebar-show', () => (this.show()))
    this.ctx.on('sidebar-hide', () => (this.hide()))

    // TODO for testing
    this.show()
  }

  /** 初始化 view 切换器 */
  private initViewSwitcher() {
    if (!this.ctx.user.data.isAdmin) {
      this.$curtViewBtn.remove()
      return
    }

    this.$curtViewBtn.onclick = () => {
      this.toggleViewSwitcher()
    }

    this.$navViews.innerHTML = ''
    this.registerViews.forEach(view => {
      const $item = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navViews.append($item)
      $item.setAttribute('data-name', view.viewName)
      $item.innerText = view.viewTitle
      if (view.viewName === this.curtView) {
        $item.classList.add('atk-active')
        this.$curtViewBtnText.innerText = view.viewTitle
      }
      $item.onclick = () => {
        // 切换 view
        this.switchView(view.viewName)

        this.toggleViewSwitcher()
      }
    })
  }

  /** 显示/隐藏 view 切换器 */
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

  private firstShow = true

  /** 显示 */
  public async show() {
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
    if (this.firstShow) {
      await this.loadSiteSwitcher()
      this.switchView(DEFAULT_VIEW) // 打开默认 view
      this.firstShow = false
    }
  }

  /** 隐藏 */
  public hide() {
    // 执行动画
    this.$el.style.transform = ''

    // 用完即销毁
    this.layer?.dispose()
  }

  /** 装载站点切换器 */
  private async loadSiteSwitcher() {
    this.$siteSwitcherSites.innerHTML = ''

    const renderSiteItem = (siteName: string, siteLogo: string, siteTarget?: string) => {
      const $site = Utils.createElement(
        `<div class="atk-site-item">
          <div class="atk-site-logo"></div>
          <div class="atk-site-name"></div>
        </div>`)
        $site.onclick = () => this.switchSite(siteTarget || siteName)
        $site.setAttribute('data-name', siteTarget || siteName)
        const $siteLogo = $site.querySelector<HTMLElement>('.atk-site-logo')!
        const $siteName = $site.querySelector<HTMLElement>('.atk-site-name')!
        $siteLogo.innerText = siteLogo
        $siteName.innerText = siteName
        if (this.curtSite === (siteTarget || siteName)) $site.classList.add('atk-active')
        this.$siteSwitcherSites.append($site)
    }

    renderSiteItem('所有站点', '', SITE_ALL_NAME)

    const sites = await new Api(this.ctx).siteGet()
    sites.forEach((site) => {
      renderSiteItem(site.name, site.name.substr(0, 1))
    })
  }

  private switchSite(siteName: string) {
    this.curtSite = siteName
    const curtView = this.curtViewInstance
    curtView?.switchTab(this.curtTab!, siteName)
    this.hideSiteSwitcher()

    // set active
    this.$siteSwitcher.querySelectorAll('.atk-site-item').forEach(e => {
      if (e.getAttribute('data-name') !== siteName) {
        e.classList.remove('atk-active')
      } else {
        e.classList.add('atk-active')
      }
    })
  }

  private siteSwitcherOC?: (evt: MouseEvent) => void // outside checker

  public showSiteSwitcher($trigger?: HTMLElement) {
    this.$siteSwitcher.style.display = ''

    // 点击外部隐藏
    if ($trigger) {
      this.siteSwitcherOC = (evt: MouseEvent) => {
        const isClickInside = $trigger.contains(evt.target as any)||this.$siteSwitcher.contains(evt.target as any)
        if (!isClickInside) { this.hideSiteSwitcher() }
      }
      document.addEventListener('click', this.siteSwitcherOC)
    }
  }

  public hideSiteSwitcher() {
    this.$siteSwitcher.style.display = 'none'
    if (this.siteSwitcherOC) document.removeEventListener('click', this.siteSwitcherOC)
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

    // update view indicator
    this.$curtViewBtnText.innerText = (view.constructor as typeof SidebarView).viewTitle
    this.$navViews.querySelectorAll('.atk-tab-item').forEach((e) => {
      if (e.getAttribute('data-name') === viewName) {
        e.classList.add('atk-active')
      } else {
        e.classList.remove('atk-active')
      }
    })

    // update tabs
    this.initNavTabs(view)

    // init view ui
    view.mount(this.curtSite)
    this.$viewWrap.innerHTML = ''
    this.$viewWrap.append(view.$el)
  }

  private initNavTabs(view: SidebarView) {
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
        view.switchTab(tabName, this.curtSite)
      }
    })
  }
}

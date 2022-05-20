import './style/sidebar.less'

import Context from 'artalk/types/context'
import * as Utils from 'artalk/src/lib/utils'
import * as Ui from 'artalk/src/lib/ui'

import Component from './sidebar-component'
import { SidebarCtx } from './main'
import SidebarHTML from './sidebar.html?raw'

import SidebarView from './sidebar-views/sidebar-view'
import CommentsView from './sidebar-views/comments-view'
import PagesView from './sidebar-views/pages-view'
import SitesView from './sidebar-views/sites-view'
import TransferView from './sidebar-views/transfer-view'
import SiteListFloater from './components/site-list-floater'
import SettingView from './sidebar-views/setting-view'

import MD5 from './lib/md5'

const REGISTER_VIEWS = [
  CommentsView, PagesView, SitesView, TransferView, // SettingView
]

export default class SidebarRoot extends Component {
  public static readonly DEFAULT_VIEW = 'comments'

  public $header: HTMLElement
  public $headerMenu: HTMLElement
  public $title: HTMLElement
  public $avatar: HTMLElement
  public $siteLogo?: HTMLElement
  public $nav: HTMLElement
  public $curtViewBtn: HTMLElement
  public $curtViewBtnIcon: HTMLElement
  public $curtViewBtnText: HTMLElement
  public $navTabs: HTMLElement
  public $navViews: HTMLElement
  public $viewWrap: HTMLElement
  public siteSwitcher?: SiteListFloater

  private get isAdmin() { return this.ctx.user.data.isAdmin }
  public curtView?: string
  public get curtViewInstance() {
    return this.curtView ? this.viewInstances[this.curtView] : undefined
  }
  public curtTab?: string

  public viewInstances: {[name: string]: SidebarView} = {}

  private viewSwitcherShow = false

  constructor(ctx: Context, sidebar: SidebarCtx) {
    super(ctx, sidebar)

    // initial elements
    this.$el = Utils.createElement(SidebarHTML)
    this.$header = this.$el.querySelector('.atk-sidebar-header')!
    this.$headerMenu = this.$header.querySelector('.atk-menu')!
    this.$title = this.$header.querySelector('.atk-sidebar-title')!
    this.$avatar = this.$header.querySelector('.atk-avatar')!

    this.$nav = this.$el.querySelector('.atk-sidebar-nav')!
    this.$curtViewBtn = this.$nav.querySelector('.akt-curt-view-btn')!
    this.$curtViewBtnIcon = this.$curtViewBtn.querySelector('.atk-icon')!
    this.$curtViewBtnText = this.$curtViewBtn.querySelector('.atk-text')!
    this.$navTabs = this.$nav.querySelector('.atk-tabs')!
    this.$navViews = this.$nav.querySelector('.atk-views')!

    this.$viewWrap = this.$el.querySelector('.atk-sidebar-view-wrap')!

    // init UI
    this.initViewInstances()
    this.initViewSwitcher()

    // 重载操作
    this.sidebar.reload = () => {
      this.init()
    }
  }

  /** 初始化 views 实例对象 */
  private initViewInstances() {
    REGISTER_VIEWS.forEach((View) => {
      const viewInstance = new View(this.ctx, this.sidebar, this.$viewWrap)
      this.viewInstances[viewInstance.getName()] = viewInstance
    })
  }

  /** 初始化 view 切换器 */
  private initViewSwitcher() {
    this.$curtViewBtn.onclick = () => {
      this.toggleViewSwitcher()
    }

    this.$navViews.innerHTML = ''
    Object.entries(this.viewInstances).forEach(([viewName, view]) => {
      const $item = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navViews.append($item)
      $item.setAttribute('data-name', viewName)
      $item.innerText = view.viewTitle()
      if (viewName === this.curtView) {
        $item.classList.add('atk-active')
        this.$curtViewBtnText.innerText = view.viewTitle()
      }
      $item.onclick = () => {
        // 切换 view
        this.switchView(viewName)

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
      // this.$curtViewBtnText.style.display = 'none'
    } else {
      // 隐藏
      this.$navViews.style.display = 'none'
      this.$navTabs.style.display = ''
      this.$curtViewBtnIcon.classList.remove('atk-arrow')
      // this.$curtViewBtnText.style.display = ''
    }

    this.viewSwitcherShow = !this.viewSwitcherShow
  }

  /** 初始化 */
  public async init() {
     // 用户权限检测
     if (this.isAdmin) {
      // 是管理员
      this.$title.innerText = '控制中心'
      this.$curtViewBtn.style.display = ''

      if (!this.siteSwitcher) {
        // 初始化站点切换器
        this.siteSwitcher = new SiteListFloater(this.ctx, {
          onSwitchSite: (siteName) => {
            this.switchSite(siteName)
          },
          onClickSitesViewBtn: () => {
            this.switchView('sites')
          },
        })
        this.$viewWrap.before(this.siteSwitcher.$el)
        this.$avatar.onclick = (evt) => {
          if (!this.isAdmin) return
          this.siteSwitcher?.show(evt.target as any)
        }
      }

      Ui.showLoading(this.$el)
      try {
        await this.siteSwitcher!.load(this.sidebar.curtSite)
      } catch (err: any) {
        const $err = Utils.createElement(
          `<span>加载失败：${err.msg || '网络错误'}<br/></span>`
        )
        const $retryBtn = Utils.createElement(
          '<span style="cursor:pointer;">点击重新获取</span>'
        )
        $err.appendChild($retryBtn)
        $retryBtn.onclick = () => {
          Ui.setError(this.$el, null)
          this.init()
        }
        Ui.setError(this.$el, $err)
        return
      } finally {
        Ui.hideLoading(this.$el)
      }

      // 站点图像
      this.$avatar.innerHTML = ''
      this.$siteLogo = Utils.createElement('<div class="atk-site-logo"></div>')
      this.$siteLogo.innerText = (this.sidebar.curtSite || '').substring(0, 1)
      this.$avatar.append(this.$siteLogo)
    } else {
      // 不是管理员
      this.$title.innerText = '通知中心'
      this.$curtViewBtn.style.display = 'none' // 隐藏 view 切换器
      this.sidebar.curtSite = this.conf.site // 第一次 show 使用当前站点数据

      // 头像
      const $avatarImg = document.createElement('img') as HTMLImageElement
      $avatarImg.src = Utils.getGravatarURL(
        this.ctx,
        MD5(this.ctx.user.data.email.toLowerCase())
      )
      this.$avatar.innerHTML = ''
      this.$avatar.append($avatarImg)
    }
  }

  /** 切换 View */
  public switchView(viewName: string) {
    const view = this.viewInstances[viewName]
    if (!view)
      throw new Error(`Sidebar View "${viewName}" not found`)

    // init view
    view.mount()

    this.curtView = viewName
    this.curtTab = view.getActiveTab()

    // update view indicator
    this.$curtViewBtnText.innerText = view.viewTitle()
    this.$navViews.querySelectorAll('.atk-tab-item').forEach((e) => {
      if (e.getAttribute('data-name') === viewName) {
        e.classList.add('atk-active')
      } else {
        e.classList.remove('atk-active')
      }
    })

    // update tabs
    this.loadViewTabs(view)

    // load element
    this.$viewWrap.innerHTML = ''
    this.$viewWrap.append(view.$el)
    this.$viewWrap.classList.forEach((c) => {
      if (c.startsWith('atk-view-name-'))
        this.$viewWrap.classList.remove(c)
    })
    this.$viewWrap.classList.add(`atk-view-name-${view.getName()}`)
  }

  private loadViewTabs(view: SidebarView) {
    this.$navTabs.innerHTML = ''
    Object.entries<string>(view.getTabs()).forEach(([tabName, label]) => {
      const $tab = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navTabs.append($tab)
      $tab.innerText = label
      if (view.getActiveTab() === tabName) $tab.classList.add('atk-active')

      // 切换 tab
      $tab.onclick = () => {
        if (view.switchTab(tabName) === false) { return }
        this.$navTabs.querySelectorAll('.atk-active').forEach(e => e.classList.remove('atk-active'))
        $tab.classList.add('atk-active')
        this.curtTab = tabName
      }
    })
  }

  /** 切换站点 */
  private switchSite(siteName: string) {
    this.sidebar.curtSite = siteName
    const curtView = this.curtViewInstance
    curtView?.switchTab(this.curtTab!)
    if (this.$siteLogo) this.$siteLogo.innerText = this.sidebar.curtSite.substring(0, 1)
  }
}

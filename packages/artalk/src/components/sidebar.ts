import '../style/sidebar.less'

import MD5 from '../lib/md5'
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
import TransferView from './sidebar-views/transfer-view'
import { SiteData } from '~/types/artalk-data'
import Api from '../api'
import SiteListFloater from './admin/site-list-floater'
import SettingView from './sidebar-views/setting-view'

const DEFAULT_VIEW = 'comments'
const REGISTER_VIEWS: (typeof SidebarView)[] = [
  CommentsView, PagesView, SitesView, TransferView, // SettingView
]

export default class Sidebar extends Component {
  public layer?: Layer

  public $header: HTMLElement
  public $headerMenu: HTMLElement
  public $title: HTMLElement
  public $avatar: HTMLElement
  public $siteLogo?: HTMLElement
  public $closeBtn: HTMLElement
  public $nav: HTMLElement
  public $curtViewBtn: HTMLElement
  public $curtViewBtnIcon: HTMLElement
  public $curtViewBtnText: HTMLElement
  public $navTabs: HTMLElement
  public $navViews: HTMLElement
  public $viewWrap: HTMLElement
  public siteSwitcher?: SiteListFloater

  private get isAdmin() { return this.ctx.user.data.isAdmin }
  public curtSite?: string
  public curtView: string = DEFAULT_VIEW
  public get curtViewInstance() {
    return this.curtView ? this.viewInstances[this.curtView] : undefined
  }
  public curtTab?: string

  public viewInstances: {[name: string]: SidebarView} = {}

  private viewSwitcherShow = false

  constructor(ctx: Context) {
    super(ctx)

    // initial elements
    this.$el = Utils.createElement(SidebarHTML)
    this.$header = this.$el.querySelector('.atk-sidebar-header')!
    this.$headerMenu = this.$header.querySelector('.atk-menu')!
    this.$title = this.$header.querySelector('.atk-sidebar-title')!
    this.$avatar = this.$header.querySelector('.atk-avatar')!
    this.$closeBtn = this.$header.querySelector('.atk-sidebar-close')!

    this.$nav = this.$el.querySelector('.atk-sidebar-nav')!
    this.$curtViewBtn = this.$nav.querySelector('.akt-curt-view-btn')!
    this.$curtViewBtnIcon = this.$curtViewBtn.querySelector('.atk-icon')!
    this.$curtViewBtnText = this.$curtViewBtn.querySelector('.atk-text')!
    this.$navTabs = this.$nav.querySelector('.atk-tabs')!
    this.$navViews = this.$nav.querySelector('.atk-views')!

    this.$viewWrap = this.$el.querySelector('.atk-sidebar-view-wrap')!

    // init UI
    this.initViewSwitcher()

    this.$closeBtn.onclick = () => {
      this.hide()
    }

    // event
    this.ctx.on('sidebar-show', () => (this.show()))
    this.ctx.on('sidebar-hide', () => (this.hide()))
    this.ctx.on('user-changed', () => { this.firstShow = true })
  }

  /** 初始化 view 切换器 */
  private initViewSwitcher() {
    this.$curtViewBtn.onclick = () => {
      this.toggleViewSwitcher()
    }

    this.$navViews.innerHTML = ''
    REGISTER_VIEWS.forEach(view => {
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

  private firstShow = true

  /** 显示 */
  public async show() {
    this.$el.style.transform = '' // 动画清除，防止二次打开失效

    // 获取 Layer
    this.layer = new Layer(this.ctx, 'sidebar', this.$el)
    this.layer.afterHide = () => {
      // 防止评论框被吞
      if (this.ctx.conf.editorTravel === true) {
        this.ctx.trigger('editor-travel-back')
      }
    }
    this.layer.show()

    // viewWrap 滚动条归位
    this.$viewWrap.scrollTo(0, 0)

    // 执行动画
    setTimeout(() => {
      this.$el.style.transform = 'translate(0, 0)'
    }, 20)

    // 第一次加载
    if (this.firstShow) {
      ///////////////////
      //// IMPORTANT ////
      ///////////////////

      // 用户权限检测
      if (this.isAdmin) {
        // 是管理员
        this.$title.innerText = '控制中心'
        this.$curtViewBtn.style.display = ''

        if (!this.siteSwitcher) {
          // 初始化站点切换器
          this.siteSwitcher = new SiteListFloater(this.ctx, {
            onSwitchSite: (siteName) => { this.switchSite(siteName) },
            onClickSitesViewBtn: () => { this.switchView('sites') }
          })
          this.$viewWrap.before(this.siteSwitcher.$el)
          this.$avatar.onclick = (evt) => {
            if (!this.isAdmin) return
            this.siteSwitcher?.show(evt.target as any)
          }

        }

        this.curtSite = this.conf.site

        Ui.showLoading(this.$el)
        try {
          await this.siteSwitcher!.load(this.curtSite)
        } catch (err: any) {
          const $err = Utils.createElement(`<span>加载失败：${err.msg || '网络错误'}<br/></span>`)
          const $retryBtn = Utils.createElement('<span style="cursor:pointer;">点击重新获取</span>')
          $err.appendChild($retryBtn)
          $retryBtn.onclick = () => { Ui.setError(this.$el, null); this.show() }
          Ui.setError(this.$el, $err)
          return
        } finally { Ui.hideLoading(this.$el) }

        // 站点图像
        this.$avatar.innerHTML = ''
        this.$siteLogo = Utils.createElement('<div class="atk-site-logo"></div>')
        this.$siteLogo.innerText = (this.curtSite || '').substr(0, 1)
        this.$avatar.append(this.$siteLogo)
      } else {
        // 不是管理员
        this.$title.innerText = '通知中心'
        this.$curtViewBtn.style.display = 'none' // 隐藏 view 切换器
        this.curtSite = this.conf.site // 第一次 show 使用当前站点数据

        // 头像
        const $avatarImg = document.createElement('img') as HTMLImageElement
        $avatarImg.src = Utils.getGravatarURL(this.ctx, MD5(this.ctx.user.data.email.toLowerCase()))
        this.$avatar.innerHTML = ''
        this.$avatar.append($avatarImg)
      }

      this.switchView(DEFAULT_VIEW) // 打开默认 view
      this.firstShow = false
    }
  }

  /** 隐藏 */
  public hide() {
    // 执行动画
    this.$el.style.transform = ''

    this.layer?.hide()
  }

  /** 切换 View */
  public switchView(viewName: string) {
    let view = this.viewInstances[viewName]
    if (!view) {
      // 初始化 View
      const View = REGISTER_VIEWS.find(o => o.viewName === viewName)!
      view = new View(this.ctx, this.$viewWrap)
      this.viewInstances[viewName] = view
    }

    // init view
    view.mount(this.curtSite!)

    this.curtView = viewName
    this.curtTab = view.viewActiveTab

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
    this.loadViewTabs(view)

    // load element
    this.$viewWrap.innerHTML = ''
    this.$viewWrap.append(view.$el)
    this.$viewWrap.classList.forEach((c) => {
      if (c.startsWith('atk-view-name-'))
        this.$viewWrap.classList.remove(c)
    })
    this.$viewWrap.classList.add(`atk-view-name-${(view.constructor as typeof SidebarView).viewName}`)
  }

  private loadViewTabs(view: SidebarView) {
    this.$navTabs.innerHTML = ''
    Object.entries<string>(view.viewTabs).forEach(([tabName, label]) => {
      const $tab = Utils.createElement(`<div class="atk-tab-item"></div>`)
      this.$navTabs.append($tab)
      $tab.innerText = label
      if (view.viewActiveTab === tabName) $tab.classList.add('atk-active')

      // 切换 tab
      $tab.onclick = () => {
        if (view.switchTab(tabName, this.curtSite!) === false) { return }
        this.$navTabs.querySelectorAll('.atk-active').forEach(e => e.classList.remove('atk-active'))
        $tab.classList.add('atk-active')
        this.curtTab = tabName
      }
    })
  }

  /** 切换站点 */
  private switchSite(siteName: string) {
    this.curtSite = siteName
    const curtView = this.curtViewInstance
    curtView?.switchTab(this.curtTab!, siteName)
    if (this.$siteLogo) this.$siteLogo.innerText = this.curtSite.substr(0, 1)
  }
}

import '../style/sidebar.less'

import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment  from './comment'
import SidebarHTML from './html/sidebar.html?raw'
import Layer from './layer'

import SidebarView from './sidebar-views/sidebar-view'
import MessageView from './sidebar-views/message-view'
import AdminView from './sidebar-views/admin-view'
import { SiteData } from '~/types/artalk-data'

export default class Sidebar extends Component {
  public el: HTMLElement
  public layer?: Layer

  public headerEl: HTMLElement
  public navEl: HTMLElement
  public viewWrapEl: HTMLElement

  public siteList: SiteData[] = []

  public site?: SiteData
  public view?: SidebarView
  private isFirstShow = true

  public viewInstances: {[name: string]: SidebarView} = {}
  public registerViews: (typeof SidebarView)[] = [
    MessageView, AdminView,
  ]

  constructor (ctx: Context) {
    super(ctx)

    // initial elements
    this.el = Utils.createElement(SidebarHTML)
    this.headerEl = this.el.querySelector('.atk-sidebar-header')!
    this.navEl = this.el.querySelector('.atk-sidebar-nav')!
    this.viewWrapEl = this.el.querySelector('.atk-sidebar-view-wrap')!

    this.el.querySelector('.atk-sidebar-close')!.addEventListener('click', () => {
      this.hide()
    })

    // event
    this.ctx.on('sidebar-show', (payload) => (this.show(payload?.viewName)))
    this.ctx.on('sidebar-hide', () => (this.hide()))

    // TODO
    this.show()

    // titles
    // this.registerViews.forEach(View => {
    //   const viewInstance = new View(this.ctx)
    //   this.viewInstances[viewInstance.name] = viewInstance
    //   viewInstance.el.classList.add(`atk-sidebar-view-${viewInstance.name}`)

    //   const titleEl = Utils.createElement(`
    //   <span class="atk-title-item" data-name="${viewInstance.name}">${viewInstance.title}</span>
    //   `)
    //   if (viewInstance.adminOnly) {
    //     titleEl.setAttribute('atk-only-admin-show', '')
    //     if (!this.ctx.user.data.isAdmin)
    //       titleEl.classList.add('atk-hide')
    //   }
    //   this.headerEl.append(titleEl)

    //   // title click
    //   titleEl.addEventListener('click', () => { this.switchView(viewInstance) })
    // })
  }

  show (viewName?: string) {
    this.el.style.transform = '' // 动画清除，防止二次打开失效

    this.layer = new Layer(this.ctx, 'sidebar', this.el)
    this.layer.show()
    this.viewWrapEl.scrollTo(0, 0)

    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)' // 执行动画
    }, 20)

    if (viewName) {
      this.switchViewByName(viewName) // 打开指定 view
    }

    if (this.isFirstShow) {
      if (!viewName) this.switchViewByName('message') // 打开默认 view
      this.isFirstShow = false
    }
  }

  hide () {
    this.el.style.transform = ''
    this.layer?.dispose() // 用完即销毁
  }

  switchViewByName (viewName: string) {
    if (!this.viewInstances[viewName]) {
      console.error(`未找到 view: ${viewName}`)
      return
    }

    this.switchView(this.viewInstances[viewName])
  }

  switchView (view: SidebarView) {
    this.view = view

    // focus title
    const titleEl = this.headerEl.querySelector<HTMLElement>(`[data-name=${view.name}]`)!
    this.headerEl.querySelector<HTMLElement>('.atk-active')!.innerText = view.title
    this.headerEl.querySelectorAll<HTMLElement>('.atk-title-item').forEach((item) => {
      if (!item.classList.contains('atk-active'))
        item.style.display = ''
    })
    titleEl.style.display = 'none'

    // init view ui
    view.init()
    this.viewWrapEl.innerHTML = ''
    this.viewWrapEl.append(view.el)

    // actions
    // this.actionsEl.innerHTML = ''
    // Object.entries(view.actions).forEach(([name, label]) => {
    //   const actionItemEl = Utils.createElement(`<span class="atk-action-item">${label}</span>`)
    //   this.actionsEl.append(actionItemEl)

    //   if (view.activeAction === name) actionItemEl.classList.add('atk-active')

    //   // action click
    //   actionItemEl.addEventListener('click', () => {
    //     view.switch(name)
    //     view.activeAction = name

    //     this.actionsEl.querySelectorAll('.atk-active').forEach((item) => {
    //       item.classList.remove('atk-active')
    //     })
    //     actionItemEl.classList.add('atk-active')
    //   })
    // })
  }
}

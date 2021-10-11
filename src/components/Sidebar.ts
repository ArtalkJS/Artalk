import '../style/sidebar.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment  from './Comment'
import SidebarHTML from './html/sidebar.html?raw'
import BuildLayer, { Layer } from './Layer'

import SidebarView from './sidebar-views/SidebarView'
import MessageView from './sidebar-views/MessageView'
import AdminView from './sidebar-views/AdminView'

export default class Sidebar extends Component {
  public el: HTMLElement
  public layer?: Layer
  public actionsEl: HTMLElement
  public contentEl: HTMLElement
  public titleWrapEl: HTMLElement

  public view?: SidebarView
  public action?: string
  private isFirstShow = true

  public viewInstances: {[name: string]: SidebarView} = {}
  public registerViews: (typeof SidebarView)[] = [
    MessageView, AdminView,
  ]

  constructor (ctx: Context) {
    super(ctx)

    // initial elements
    this.el = Utils.createElement(SidebarHTML)
    this.contentEl = this.el.querySelector('.atk-sidebar-content')!
    this.titleWrapEl = this.el.querySelector('.atk-sidebar-title-wrap')!
    this.actionsEl = this.el.querySelector('.atk-sidebar-actions')!

    this.el.querySelector('.atk-sidebar-close')!.addEventListener('click', () => {
      this.hide()
    })

    // event
    this.ctx.addEventListener('sidebar-show', (payload) => (this.show(payload?.viewName)))
    this.ctx.addEventListener('sidebar-hide', () => (this.hide()))

    // titles
    this.registerViews.forEach(View => {
      const viewInstance = new View(this.ctx)
      this.viewInstances[viewInstance.name] = viewInstance
      viewInstance.el.classList.add(`atk-sidebar-view-${viewInstance.name}`)

      const titleEl = Utils.createElement(`
      <span class="atk-title-item" data-name="${viewInstance.name}">${viewInstance.title}</span>
      `)
      if (viewInstance.adminOnly) {
        titleEl.setAttribute('atk-only-admin-show', '')
        if (!this.ctx.user.data.isAdmin)
          titleEl.classList.add('atk-hide')
      }
      this.titleWrapEl.append(titleEl)

      // title click
      titleEl.addEventListener('click', () => { this.switchView(viewInstance) })
    })
  }

  show (viewName?: string) {
    this.el.style.transform = '' // 动画清除，防止二次打开失效

    this.layer = BuildLayer(this.ctx, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)

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
    const titleEl = this.titleWrapEl.querySelector<HTMLElement>(`[data-name=${view.name}]`)!
    this.titleWrapEl.querySelector<HTMLElement>('.atk-active')!.innerText = view.title
    this.titleWrapEl.querySelectorAll<HTMLElement>('.atk-title-item').forEach((item) => {
      if (!item.classList.contains('atk-active'))
        item.style.display = ''
    })
    titleEl.style.display = 'none'

    // init view ui
    view.init()
    this.contentEl.innerHTML = ''
    this.contentEl.append(view.el)

    // actions
    this.actionsEl.innerHTML = ''
    Object.entries(view.actions).forEach(([name, label]) => {
      const actionItemEl = Utils.createElement(`<span class="atk-action-item">${label}</span>`)
      this.actionsEl.append(actionItemEl)

      if (view.activeAction === name) actionItemEl.classList.add('atk-active')

      // action click
      actionItemEl.addEventListener('click', () => {
        view.switch(name)
        view.activeAction = name

        this.actionsEl.querySelectorAll('.atk-active').forEach((item) => {
          item.classList.remove('atk-active')
        })
        actionItemEl.classList.add('atk-active')
      })
    })
  }
}

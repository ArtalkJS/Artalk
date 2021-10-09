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
  public adminMode: boolean = false

  public view?: SidebarView
  public action?: string

  public registerViews: (typeof SidebarView)[] = [
    MessageView, AdminView,
  ]

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(SidebarHTML)
    this.contentEl = this.el.querySelector('.atk-sidebar-content')!
    this.titleWrapEl = this.el.querySelector('.atk-sidebar-title-wrap')!
    this.actionsEl = this.el.querySelector('.atk-sidebar-actions')!

    this.el.querySelector('.atk-sidebar-close')!.addEventListener('click', () => {
      this.hide()
    })

    this.ctx.addEventListener('sidebar-show', () => (this.show()))
    this.ctx.addEventListener('sidebar-hide', () => (this.hide()))

    // titles
    this.registerViews.forEach(View => {
      const viewInstance = new View(this.ctx)
      viewInstance.el.classList.add(`atk-sidebar-view-${viewInstance.name}`)

      const titleEl = Utils.createElement(`
      <span class="atk-title-item" data-name="${viewInstance.name}">${viewInstance.title}</span>
      `)
      if (viewInstance.adminOnly) titleEl.setAttribute('atk-only-admin-show', '')
      this.titleWrapEl.append(titleEl)

      // title click
      titleEl.addEventListener('click', () => {
        this.titleWrapEl.querySelectorAll('.atk-active').forEach((item) => {
          item.classList.remove('atk-active')
        })
        titleEl.classList.add('atk-active')

        this.switchView(viewInstance)
      })
    })

    ;(this.titleWrapEl.firstChild as HTMLElement).click() // 打开第一个 view
  }

  show () {
    this.el.style.transform = '' // 动画清除，防止二次打开失效

    this.layer = BuildLayer(this.ctx, 'sidebar', this.el)
    this.layer.show()
    this.contentEl.scrollTo(0, 0)

    setTimeout(() => {
      this.el.style.transform = 'translate(0, 0)' // 执行动画
    }, 20)

    this.adminMode = this.ctx.user.data.isAdmin
  }

  hide () {
    this.el.style.transform = ''
    this.layer?.dispose() // 用完即销毁
  }

  switchView (view: SidebarView) {
    this.view = view

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

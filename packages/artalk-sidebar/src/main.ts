import Artalk from 'artalk'
import ArtalkConfig, { LocalUser } from 'artalk/types/artalk-config'
import 'artalk/dist/Artalk.css'
import SidebarRoot from './sidebar-root'

export default async function InitSidebar(conf: ArtalkConfig, user: LocalUser, view?: string) {
  // 初始化 Artalk 主程序
  conf.useBackendConf = true
  const artalk = await (new Artalk(conf))
  artalk.$root.style.display = 'none'

  // 初始化用户数据
  user = user || {}
  artalk.ctx.user.update({
    ...user
  })

  // 初始化 Sidebar
  const sidebarCtx = new SidebarCtx()
  sidebarCtx.curtSite = conf.site!

  const sidebar = new SidebarRoot(artalk.ctx, sidebarCtx)

  // 装载元素
  document.body.appendChild(sidebar.$el)

  // 暗黑模式
  if (conf.darkMode) sidebar.$el.classList.add('atk-dark-mode')

  // 初始化侧边栏
  await sidebar.init()

  // 打开 View
  let viewName = SidebarRoot.DEFAULT_VIEW
  let viewParams = {}
  if (view) {
    const splitted = view.split('|')
    if (splitted[0]) viewName = splitted[0]
    if (splitted[1]) viewParams = JSON.parse(splitted[1])
  }

  sidebar.switchView(viewName)

  if (viewName === 'sites') {
    // TODO 炮打翻山...
    artalk.ctx.trigger('sidebar-sites-create-form', viewParams)
  }

  console.log('hello artalk-sidebar')
}

export class SidebarCtx {
  reload!: () => void
  curtSite!: string
}

import Artalk from 'artalk'
import ArtalkConfig, { LocalUser } from 'artalk/types/artalk-config'
import 'artalk/dist/Artalk.css'
import SidebarRoot from './sidebar-root'

export default function InitSidebar(conf: ArtalkConfig, user: LocalUser) {
  const artalk = new Artalk(conf)

  artalk.$root.style.display = 'none'

  user = user || {}
  artalk.ctx.user.data = {
    ...artalk.ctx.user.data,
    ...user
  }
  artalk.ctx.user.save()

  const sidebarCtx = new SidebarCtx()
  const sidebar = new SidebarRoot(artalk.ctx, sidebarCtx)
  document.body.appendChild(sidebar.$el)

  if (conf.darkMode) sidebar.$el.classList.add('atk-dark-mode')

  sidebar.init()

  console.log('hello artalk-sidebar')
}

export class SidebarCtx {
  reload!: () => void
}

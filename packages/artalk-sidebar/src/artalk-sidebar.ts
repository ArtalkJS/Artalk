import Artalk from 'artalk'
import ArtalkConfig, { LocalUser } from 'artalk/types/artalk-config'
import 'artalk/dist/Artalk.css'
import Sidebar from './sidebar'

class ArtalkSidebar extends Artalk {
  constructor(customConf: ArtalkConfig, user: LocalUser) {
    super(customConf)

    this.$root.style.display = 'none'

    user = user || {}
    this.ctx.user.data = {
      ...this.ctx.user.data,
      ...user
    }
    this.ctx.user.save()

    const sidebar = new Sidebar(this.ctx)
    document.body.appendChild(sidebar.$el)

    if (customConf.darkMode) {
      sidebar.$el.classList.add('atk-dark-mode')
    }

    sidebar.show()
    console.log('hello artalk-sidebar')
  }
}

export default ArtalkSidebar

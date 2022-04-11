import Artalk from 'artalk'
import ArtalkConfig from 'artalk/types/artalk-config'
import 'artalk/dist/Artalk.css'
import Sidebar from './sidebar'

class ArtalkSidebar extends Artalk {
  constructor(customConf: ArtalkConfig) {
    super(customConf)

    this.$root.style.display = 'none'

    const sidebar = new Sidebar(this.ctx)
    document.body.appendChild(sidebar.$el)

    console.log('hello artalk-sidebar')
  }
}

export default ArtalkSidebar

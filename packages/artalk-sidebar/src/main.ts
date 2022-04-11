import Artalk from 'artalk'
import ArtalkConfig from 'artalk/types/artalk-config'

class ArtalkSidebar extends Artalk {
  constructor(customConf: ArtalkConfig) {
    super(customConf)

    console.log('hello artalk-sidebar')
  }
}

export default ArtalkSidebar

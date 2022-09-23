import type Artalk from 'artalk'
import type { LocalUser } from 'artalk/types/artalk-config'

export let artalk: Artalk|null = null

// 启动参数
export const bootParams = getBootParams()
function getBootParams() {
  const p = new URLSearchParams(document.location.search)

  return {
    pageKey:    p.get('pageKey') || '',
    site:       p.get('site') || '',
    user:       <LocalUser>JSON.parse(p.get('user') || '{}'),
    view:       p.get('view') || '',
    viewParams: <any>null,
    darkMode:   p.get('darkMode') === '1'
  }
}

export default {
  getArtalk: () => artalk,
  setArtalk: (artalkInstance: Artalk) => {
    artalk = artalkInstance
  },
  getBootParams: () => bootParams
}

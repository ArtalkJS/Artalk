import Artalk from 'artalk'
import type { LocalUser } from 'artalk/types/artalk-config'

export let artalk: Artalk|null = null

// 启动参数
export const bootParams = getBootParams()
export function getBootParams() {
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

export function createArtalkInstance() {
  const artalkEl = document.createElement('div')
  artalkEl.style.display = 'none'
  document.body.append(artalkEl)

  Artalk.DisabledComponents = ['list']
  return new Artalk({
    el: artalkEl,
    server: (import.meta.env.DEV) ? 'http://localhost:23366' : '/',
    pageKey: 'https://artalk.js.org/guide/intro.html',
    site: bootParams.site,
    darkMode: bootParams.darkMode,
    useBackendConf: true
  }) as unknown as Promise<Artalk>
}

export default {
  createArtalkInstance,
  getArtalk: () => artalk,
  setArtalk: (artalkInstance: Artalk) => {
    artalk = artalkInstance
  },
  getBootParams: () => bootParams
}

import Artalk from 'artalk'
import type { LocalUser } from 'artalk/types/artalk-config'
import { useUserStore } from './stores/user'

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
    server: (import.meta.env.DEV) ? 'http://localhost:23366' : '../',
    pageKey: bootParams.pageKey,
    site: bootParams.site,
    darkMode: bootParams.darkMode,
    useBackendConf: true
  })
}

export function importUserDataFromArtalkInstance() {
  if (!artalk) throw Error("the artalk instance is not exist")
  if (!artalk.ctx.user.data.email) throw Error("the user data in artalk instance is invalid")

  const userData = artalk.ctx.user.data
  useUserStore().$patch((state) => {
    state.site = '__ATK_SITE_ALL'
    state.name = userData.nick
    state.email = userData.email
    state.isAdmin = userData.isAdmin
    state.token = userData.token
  })
}

export default {
  createArtalkInstance,
  getArtalk: () => artalk!,
  setArtalk: (artalkInstance: Artalk) => {
    artalk = artalkInstance
  },
  getBootParams: () => bootParams,
  importUserDataFromArtalkInstance,
}

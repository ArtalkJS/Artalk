import Artalk from 'artalk'
import type { ArtalkType } from 'artalk'
import { useUserStore } from './stores/user'

export let artalk: Artalk|null = null

// 启动参数
export const bootParams = getBootParams()
export function getBootParams() {
  const p = new URLSearchParams(document.location.search)

  return {
    pageKey:    p.get('pageKey') || '',
    site:       p.get('site') || '',
    user:       <ArtalkType.LocalUser>JSON.parse(p.get('user') || '{}'),
    view:       p.get('view') || '',
    viewParams: <any>null,
    darkMode:   p.get('darkMode') === '1',
    locale:     p.get('locale') || '',
  }
}

export function createArtalkInstance() {
  const artalkEl = document.createElement('div')
  artalkEl.style.display = 'none'
  document.body.append(artalkEl)

  return Artalk.init({
    el: artalkEl,
    server: '../',
    pageKey: bootParams.pageKey,
    site: bootParams.site,
    darkMode: bootParams.darkMode,
    useBackendConf: true,
    remoteConfModifier: (conf) => {
      conf.noComment = `<div class="atk-sidebar-no-content">No Content</div>` // TODO i18n t('noComment')
      conf.flatMode = true
      conf.pagination = {
        pageSize: 20,
        readMore: false,
        autoLoad: false,
      }
      conf.listUnreadHighlight = true
    }
  })
}

export function importUserDataFromArtalkInstance() {
  if (!artalk) throw Error("the artalk instance is not exist")
  if (!artalk.ctx.get('user').getData().email) throw Error("the user data in artalk instance is invalid")

  const userData = artalk.ctx.get('user').getData()
  useUserStore().$patch((state) => {
    state.site = ''
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

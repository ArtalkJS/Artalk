import Artalk from 'artalk'
import type { ArtalkType } from 'artalk'

export let artalk: Artalk|null = null

export function setArtalk(artalkInstance: Artalk) {
  artalk = artalkInstance
}

export function getArtalk() {
  return artalk
}

// 启动参数
export const bootParams = getBootParams()

function getBootParams() {
  const p = new URLSearchParams(document.location.search)

  // call history api to clear search params
  // on purpose to prevent the params (e.g. user token)
  // from being leaked like from the referrer header or the browser history
  if (!!p.get('user') && window.history.replaceState) {
    window.history.replaceState({}, '', window.location.pathname)
  }

  return {
    pageKey:    p.get('pageKey') || '',
    site:       p.get('site') || '',
    user:       <ArtalkType.LocalUser>JSON.parse(p.get('user') || '{}'),
    view:       p.get('view') || '',
    viewParams: <any>null,
    darkMode:   p.get('darkMode') === '1',
  }
}

export function initArtalk() {
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
    pvAdd: false,
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

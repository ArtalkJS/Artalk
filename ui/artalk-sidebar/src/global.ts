import Artalk from 'artalk'
import type { LocalUser } from 'artalk'

export let artalk: Artalk | null = null

export function setArtalk(artalkInstance: Artalk) {
  artalk = artalkInstance
}

export function getArtalk() {
  return artalk
}

/**
 * Boot params from URL search params
 *
 * TODO: Refactor to a singleton store
 */
export const bootParams = getBootParams()

function getBootParams() {
  const p = new URLSearchParams(document.location.search)

  // call history api to clear search params
  // on purpose to prevent the params (e.g. user token)
  // from being leaked like from the referrer header or the browser history
  if (!!p.get('user') && window.history.replaceState) {
    window.history.replaceState({}, '', window.location.pathname)
  }

  const userFromURL = JSON.parse(p.get('user') || '{}')
  const user: LocalUser = {
    name: userFromURL.name || '',
    email: userFromURL.email || '',
    link: userFromURL.link || '',
    token: userFromURL.token || '',
    is_admin: userFromURL.is_admin || false,
  }

  let darkMode: boolean
  if (p.get('darkMode') != null) {
    darkMode = p.get('darkMode') == '1'
  } else {
    darkMode =
      localStorage.getItem('ATK_SIDEBAR_DARK_MODE') != null
        ? localStorage.getItem('ATK_SIDEBAR_DARK_MODE') == '1'
        : window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  return {
    user,
    pageKey: p.get('pageKey') || '',
    site: p.get('site') || '',
    view: p.get('view') || '',
    viewParams: <any>null,
    darkMode,
  }
}

export function isOpenFromSidebar() {
  return !!bootParams.user?.email
}

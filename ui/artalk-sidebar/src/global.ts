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
 * Boot params from URL search params and credential fragments
 *
 * TODO: Refactor to a singleton store
 */
export const bootParams = getBootParams()

function getBootParams() {
  const p = new URLSearchParams(document.location.search)
  const fragment = new URLSearchParams(document.location.hash.slice(1))
  const guestToken = fragment.get('guest_token') || fragment.get('atk_guest') || ''

  // Consume credentials before routing. Old clients may still put their token
  // in the user query; clearing it here cannot undo that initial HTTP request.
  if ((p.has('user') || fragment.has('token') || !!guestToken) && window.history.replaceState) {
    window.history.replaceState({}, '', window.location.pathname)
  }

  let userFromURL: Partial<LocalUser> = {}
  try {
    const parsed = JSON.parse(p.get('user') || '{}')
    if (parsed && typeof parsed === 'object') userFromURL = parsed
  } catch {
    // A malformed legacy profile must not prevent credential cleanup or boot.
  }
  const legacyToken = typeof userFromURL.token === 'string' ? userFromURL.token : ''
  const user: LocalUser = {
    name: typeof userFromURL.name === 'string' ? userFromURL.name : '',
    email: typeof userFromURL.email === 'string' ? userFromURL.email : '',
    link: typeof userFromURL.link === 'string' ? userFromURL.link : '',
    token: fragment.get('token') ?? legacyToken,
    is_admin: userFromURL.is_admin === true,
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
    guestToken,
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

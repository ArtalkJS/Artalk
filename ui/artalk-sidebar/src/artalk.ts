import Artalk, { type Context } from 'artalk'
import { Router } from 'vue-router'
import { bootParams } from './global'
import { useUserStore } from './stores/user'
import { getClientLocaleURLs } from './lib/forced-locale'

/**
 * @param locale - The locale forced by the parent page, its client locale set is loaded with the plugins
 */
export function setupArtalk(locale?: string) {
  // Create virtual element for Artalk
  const artalkEl = document.createElement('div')
  artalkEl.style.display = 'none'
  document.body.append(artalkEl)

  // Init Artalk
  return Artalk.init({
    el: artalkEl,
    server: '../',
    pageKey: bootParams.pageKey,
    site: bootParams.site,
    darkMode: bootParams.darkMode,
    pluginURLs: getClientLocaleURLs(locale || ''),
    pvAdd: false,
    noComment: `<div class="atk-sidebar-no-content"></div>`,
    flatMode: true,
    pagination: {
      pageSize: 20,
      readMore: false,
      autoLoad: false,
    },
    listUnreadHighlight: true,
    fetchCommentsOnInit: false,
  })
}

export async function syncArtalkUser(artalk: Context, router: Router) {
  const user = useUserStore()
  const logout = () => {
    user.logout()
    nextTick(() => {
      router.replace('/login')
    })
  }

  // Access from open sidebar or directly by url
  if (bootParams.user?.email) {
    // sync user from sidebar to artalk
    artalk.getUser().update(bootParams.user)
  } else {
    // Sync user from artalk to sidebar
    try {
      user.sync()
    } catch {
      logout()
      return
    }
  }

  // Async check user status (no await)
  checkUser(artalk, logout)
}

function checkUser(artalk: Context, logout: () => void) {
  // Get user info from artalk
  const { name, email } = artalk.getUser().getData()

  // Remove login failed dialog if sidebar
  artalk.getApiHandlers().remove('need_login')
  artalk.getApiHandlers().add('need_login', async () => {
    logout()
    throw new Error('Need login')
  })

  // Check user status
  artalk
    .getApi()
    .user.getUserStatus({ email, name })
    .then((res) => {
      if (res.data.is_admin && !res.data.is_login) {
        logout()
      } else {
        // Mark all notifications as read
        artalk.getApi().notifies.markAllNotifyRead({ email, name })
      }
    })
}

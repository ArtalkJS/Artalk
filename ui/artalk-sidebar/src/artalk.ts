import Artalk from 'artalk'
import { Router } from 'vue-router'
import { bootParams } from './global'
import { useUserStore } from './stores/user'

export function setupArtalk() {
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
    useBackendConf: true,
    pvAdd: false,
    remoteConfModifier: (conf) => {
      conf.noComment = `<div class="atk-sidebar-no-content"></div>`
      conf.flatMode = true
      conf.pagination = {
        pageSize: 20,
        readMore: false,
        autoLoad: false,
      }
      conf.listUnreadHighlight = true
    },
  })
}

export async function syncArtalkUser(artalk: Artalk, router: Router) {
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
    artalk.ctx.get('user').update(bootParams.user)
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

function checkUser(artalk: Artalk, logout: () => void) {
  // Get user info from artalk
  const { name, email } = artalk.ctx.get('user').getData()

  // Remove login failed dialog if sidebar
  artalk.ctx.getApiHandlers().remove('need_login')
  artalk.ctx.getApiHandlers().add('need_login', async () => {
    logout()
    throw new Error('Need login')
  })

  // Check user status
  artalk.ctx
    .getApi()
    .user.getUserStatus({ email, name })
    .then((res) => {
      if (res.data.is_admin && !res.data.is_login) {
        logout()
      } else {
        // Mark all notifications as read
        artalk.ctx.getApi().notifies.markAllNotifyRead({ email, name })
      }
    })
}

import { getGuestToken, setGuestToken } from './guest'
import { storageKey } from './storage'
import type { LocalUser, UserManager as IUserManager } from '@/types'

const LOCAL_USER_KEY = 'ArtalkUser'

interface UserOpts {
  onUserChanged?: (user: LocalUser) => void
  getServer?: () => string
}

export class UserManager implements IUserManager {
  private data!: LocalUser
  private key?: string
  private sessions = new Map<string, LocalUser>()
  private disposed = false

  constructor(private opts: UserOpts) {
    this.getData()
  }

  private load(key: string): LocalUser {
    // Import from localStorage
    let localUser: Record<string, any> = {}
    try {
      const stored = window.localStorage.getItem(key)
      localUser = JSON.parse(stored || window.localStorage.getItem(LOCAL_USER_KEY) || '{}') || {}
      // Legacy credentials have no backend binding. Retain profile fields, but
      // require a new login rather than send that token to an unknown server.
      if (!stored && key !== LOCAL_USER_KEY) {
        localUser.token = ''
        localUser.is_admin = localUser.isAdmin = false
      }
    } catch {
      // Private browsing or a third-party iframe may disable persistent storage.
    }

    // Initialize
    return {
      name: localUser.name || localUser.nick || '', // nick is deprecated (for historical compatibility)
      email: localUser.email || '',
      link: localUser.link || '',
      token: localUser.token || '',
      is_admin: localUser.is_admin || localUser.isAdmin || false,
    }
  }

  getGuestToken = getGuestToken
  setGuestToken = setGuestToken

  getData() {
    if (this.disposed) return this.data
    const key = this.opts.getServer
      ? storageKey(LOCAL_USER_KEY, this.opts.getServer())
      : LOCAL_USER_KEY
    if (key !== this.key) {
      const changed = this.key !== undefined
      this.key = key
      this.data = this.sessions.get(key) || this.load(key)
      this.sessions.set(key, this.data)
      if (changed) this.opts.onUserChanged?.(this.data)
    }
    return this.data
  }

  /** Update user data and save to localStorage */
  update(obj: Partial<LocalUser> = {}) {
    if (this.disposed) return
    this.getData()
    Object.entries(obj).forEach(([key, value]) => {
      this.data[key] = value
    })

    try {
      window.localStorage.setItem(this.key!, JSON.stringify(this.data))
    } catch {
      // Keep the in-memory user usable when persistence is unavailable.
    }
    this.opts.onUserChanged && this.opts.onUserChanged(this.data)
  }

  /**
   * Logout
   *
   * @description Logout will clear login status, but not clear user data (name, email, link)
   */
  logout() {
    this.update({
      token: '',
      is_admin: false,
    })
  }

  /** Check if user has filled basic data */
  checkHasBasicUserInfo() {
    const data = this.getData()
    return !!data.name && !!data.email
  }

  /** @internal Ignore late callbacks after the owning widget is destroyed. */
  dispose() {
    this.disposed = true
    this.sessions.clear()
  }
}

export default UserManager

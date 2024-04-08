import type { LocalUser } from '@/types'

const LOCAL_USER_KEY = 'ArtalkUser'

interface UserOpts {
  onUserChanged?: (user: LocalUser) => void
}

class User {
  private data: LocalUser

  constructor(private opts: UserOpts) {
    // Import from localStorage
    const localUser = JSON.parse(window.localStorage.getItem(LOCAL_USER_KEY) || '{}')

    // Initialize
    this.data = {
      nick: localUser.nick || '',
      email: localUser.email || '',
      link: localUser.link || '',
      token: localUser.token || '',
      isAdmin: localUser.isAdmin || false,
    }
  }

  getData() {
    return this.data
  }

  /** Update user data and save to localStorage */
  update(obj: Partial<LocalUser> = {}) {
    Object.entries(obj).forEach(([key, value]) => {
      this.data[key] = value
    })

    window.localStorage.setItem(LOCAL_USER_KEY, JSON.stringify(this.data))
    this.opts.onUserChanged && this.opts.onUserChanged(this.data)
  }

  /**
   * Logout
   *
   * @description Logout will clear login status, but not clear user data (nick, email, link)
   */
  logout() {
    this.update({
      token: '',
      isAdmin: false,
    })
  }

  /** Check if user has filled basic data */
  checkHasBasicUserInfo() {
    return !!this.data.nick && !!this.data.email
  }
}

export default User

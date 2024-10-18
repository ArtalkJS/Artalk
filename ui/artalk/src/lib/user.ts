import type { LocalUser, UserManager as IUserManager } from '@/types'

const LOCAL_USER_KEY = 'ArtalkUser'

interface UserOpts {
  onUserChanged?: (user: LocalUser) => void
}

export class UserManager implements IUserManager {
  private data: LocalUser

  constructor(private opts: UserOpts) {
    // Import from localStorage
    const localUser = JSON.parse(window.localStorage.getItem(LOCAL_USER_KEY) || '{}')

    // Initialize
    this.data = {
      name: localUser.name || localUser.nick || '', // nick is deprecated (for historical compatibility)
      email: localUser.email || '',
      link: localUser.link || '',
      token: localUser.token || '',
      is_admin: localUser.is_admin || localUser.isAdmin || false,
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
    return !!this.data.name && !!this.data.email
  }
}

export default UserManager

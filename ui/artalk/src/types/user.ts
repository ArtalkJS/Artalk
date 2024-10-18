/**
 * Local User Data (in localStorage)
 *
 * @note Keep flat for easy handling
 */
export interface LocalUser {
  /** Username (aka. Nickname) */
  name: string

  /** Email */
  email: string

  /** Link (aka. Website) */
  link: string

  /** Token (for authorization) */
  token: string

  /** Admin flag */
  is_admin: boolean
}

export interface UserManager {
  getData: () => LocalUser
  update: (dta: Partial<LocalUser>) => void
  logout: () => void
  checkHasBasicUserInfo: () => boolean
}

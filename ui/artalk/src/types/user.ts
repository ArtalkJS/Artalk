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
  getGuestToken?: (server: string, create?: boolean) => string | undefined
  setGuestToken?: (server: string, token: string) => void
  getData: () => LocalUser
  update: (dta: Partial<LocalUser>) => void
  logout: () => void
  checkHasBasicUserInfo: () => boolean
}

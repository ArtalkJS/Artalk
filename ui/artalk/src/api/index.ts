import { ApiOptions } from './options'
import { Api as ApiV2 } from './v2'
import { Fetch } from './fetch'

export class Api extends ApiV2<void> {
  private _opts: ApiOptions

  constructor(opts: ApiOptions) {
    super({
      baseUrl: opts.baseURL,
      customFetch: (input, init) => Fetch(opts, input, init)
    })

    this._opts = opts
  }

  /**
   * Get user info as params for request
   *
   * @returns Request params with user info
   */
  getUserFields() {
    const user = this._opts.userInfo
    if (!user?.name || !user?.email) return undefined
    return { name: user.name, email: user.email }
  }
}

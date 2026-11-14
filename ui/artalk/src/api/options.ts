import type { ApiHandlers } from './handler'

export interface ApiOptions {
  /** @internal Guard requests and retries after the widget's destination changes. */
  isActive?: () => boolean
  baseURL: string
  siteName: string
  pageKey: string
  pageTitle: string
  timeout?: number
  getGuestToken?: (create?: boolean) => string | undefined
  getApiToken?: () => string | undefined
  userInfo?: {
    name: string
    email: string
  }
  handlers?: ApiHandlers
}

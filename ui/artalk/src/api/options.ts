import type { ApiHandlers } from './handler'

export interface ApiOptions {
  baseURL: string
  siteName: string
  pageKey: string
  pageTitle: string
  timeout?: number
  getApiToken?: () => string | undefined
  userInfo?: {
    name: string
    email: string
  }
  handlers?: ApiHandlers
}

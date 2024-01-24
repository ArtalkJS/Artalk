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

  // -------------------------------------------------------------------
  //  Hooks
  // -------------------------------------------------------------------

  onNeedCheckCaptcha?: (payload: {
    data: {
      imgData?: string
      iframe?: string
    }
  }) => Promise<void>

  onNeedCheckAdmin?: (payload: {
  }) => Promise<void>
}

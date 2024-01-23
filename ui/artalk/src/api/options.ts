export interface ApiOptions {
  baseURL: string
  siteName: string
  pageKey: string
  pageTitle: string
  timeout?: number
  apiToken?: string
  userInfo?: {
    name: string
    email: string
  }

  // -------------------------------------------------------------------
  //  Hooks
  // -------------------------------------------------------------------

  onNeedCheckCaptcha?: (payload: {
    data: {
      imgData: string
      iframe: string
    }
  }) => Promise<void>

  onNeedCheckAdmin?: (payload: {
  }) => Promise<void>
}

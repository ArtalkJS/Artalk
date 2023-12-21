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
    recall: () => void
    reject: () => void
    data: {
      imgData: string
      iframe: string
    }
  }) => void

  onNeedCheckAdmin?: (payload: {
    recall: () => void
    reject: () => void
  }) => void
}

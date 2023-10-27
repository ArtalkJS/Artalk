export interface ApiOptions {
  baseURL: string
  siteName: string
  pageKey: string
  pageTitle: string
  timeout?: number

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

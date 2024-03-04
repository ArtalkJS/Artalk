export interface ApiHandlerPayload {
  'need_captcha': { img_data?: string; iframe?: string }
  'need_login': {}
  'need_auth_login': {}
}

type PayloadKey = keyof ApiHandlerPayload

export interface ApiHandler<T extends PayloadKey = PayloadKey> {
  action: T
  handler: (data: ApiHandlerPayload[T]) => Promise<void>
}

export interface ApiHandlers {
  add: <T extends PayloadKey>(action: T, handler: (data: ApiHandlerPayload[T]) => Promise<void>) => void
  get: () => ApiHandler[]
}

export function createApiHandlers(): ApiHandlers {
  const handlers: ApiHandler[] = []
  return {
    add: (action, handler) => { handlers.push({ action, handler }) },
    get: () => handlers
  }
}

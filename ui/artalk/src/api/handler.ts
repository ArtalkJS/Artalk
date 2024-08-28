export interface ApiHandlerPayload {
  need_captcha: { img_data?: string; iframe?: string }
  need_login: object
  need_auth_login: object
}

type PayloadKey = keyof ApiHandlerPayload

export interface ApiHandler<T extends PayloadKey = PayloadKey> {
  action: T
  handler: (data: ApiHandlerPayload[T]) => Promise<void>
}

export interface ApiHandlers {
  add: <T extends PayloadKey>(
    action: T,
    handler: (data: ApiHandlerPayload[T]) => Promise<void>,
  ) => void
  remove: (action: PayloadKey) => void
  get: () => ApiHandler[]
}

export function createApiHandlers(): ApiHandlers {
  const handlers: ApiHandler[] = []
  return {
    add: (action, handler) => {
      handlers.push({ action, handler })
    },
    remove: (action) => {
      const index = handlers.findIndex((h) => h.action === action)
      if (index !== -1) handlers.splice(index, 1)
    },
    get: () => handlers,
  }
}

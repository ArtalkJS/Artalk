import type { ArtalkPlugin } from '@/types'
import { Api, createApiHandlers } from '@/api'
import { convertApiOptions } from '@/config'

export const ApiService: ArtalkPlugin = (ctx) => {
  ctx.provide('apiHandlers', () => createApiHandlers(), [] as const)
  ctx.provide(
    'api',
    (user, config, apiHandlers) => {
      const conf = config.get()
      return new Api({
        ...convertApiOptions(conf, user, apiHandlers),
        isActive: () => !ctx.isDestroyed() && config.get().server === conf.server,
      })
    },
    ['user', 'config', 'apiHandlers'] as const,
    {
      lifecycle: 'transient',
    },
  )
}

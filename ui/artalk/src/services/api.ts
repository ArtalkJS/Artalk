import type { ArtalkPlugin } from '@/types'
import { Api, createApiHandlers } from '@/api'
import { convertApiOptions } from '@/config'

export const ApiService: ArtalkPlugin = (ctx) => {
  ctx.provide('apiHandlers', () => createApiHandlers(), [] as const)
  ctx.provide(
    'api',
    (user, config, apiHandlers) => new Api(convertApiOptions(config.get(), user, apiHandlers)),
    ['user', 'config', 'apiHandlers'] as const,
    {
      lifecycle: 'transient',
    },
  )
}

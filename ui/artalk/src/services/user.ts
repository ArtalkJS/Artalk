import type { ArtalkPlugin } from '@/types'
import { UserManager } from '@/lib/user'

export const UserService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'user',
    (config) => {
      const user = new UserManager({
        getServer: () => config.get().server,
        onUserChanged: (data) => {
          ctx.trigger('user-changed', data)
        },
      })
      config.watchConf(['server'], () => user.getData())
      ctx.on('unmounted', () => user.dispose())
      return user
    },
    ['config'] as const,
  )
}

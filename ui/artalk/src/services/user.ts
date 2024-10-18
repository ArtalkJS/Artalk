import type { ArtalkPlugin } from '@/types'
import { UserManager } from '@/lib/user'

export const UserService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'user',
    () => {
      return new UserManager({
        onUserChanged: (data) => {
          ctx.trigger('user-changed', data)
        },
      })
    },
    [] as const,
  )
}

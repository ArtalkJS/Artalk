import type { ArtalkPlugin } from '@/types'
import { CheckerLauncher } from '@/components/checker'

export const CheckersService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'checkers',
    (api, apiHandlers, layers, user, config) => {
      const checkers = new CheckerLauncher({
        getApi: () => api,
        getLayers: () => layers,
        getUser: () => user,
        onReload: () => ctx.reload(),

        // make sure suffix with a slash, because it will be used as a base url when call `fetch`
        getCaptchaIframeURL: () => `${config.get().server}/api/v2/captcha/?t=${+new Date()}`,
      })

      apiHandlers.add('need_captcha', (res) => checkers.checkCaptcha(res))
      apiHandlers.add('need_login', () => checkers.checkAdmin({}))

      return checkers
    },
    ['api', 'apiHandlers', 'layers', 'user', 'config'] as const,
  )
}

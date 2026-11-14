import type { ArtalkPlugin } from '@/types'
import { CheckerLauncher } from '@/components/checker'

export const CheckersService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'checkers',
    (apiHandlers, layers, user, config) => {
      const checkers = new CheckerLauncher({
        getApi: () => ctx.getApi(),
        getLayers: () => layers,
        getUser: () => user,
        onReload: () => ctx.reload(),

        // make sure suffix with a slash, because it will be used as a base url when call `fetch`
        getCaptchaIframeURL: () => `${config.get().server}/api/v2/captcha/?t=${+new Date()}`,
      })

      let server = config.get().server
      ctx.watchConf(['server'], (conf) => {
        if (server !== conf.server) checkers.cancelAll()
        server = conf.server
      })
      const identity = () => {
        const { name, email, token } = user.getData()
        return JSON.stringify([name, email, token])
      }
      let viewer = identity()
      ctx.on('user-changed', () => {
        const current = identity()
        if (viewer !== current) checkers.cancelAll()
        viewer = current
      })
      ctx.on('unmounted', () => checkers.dispose())

      apiHandlers.add('need_captcha', (res) => checkers.checkCaptcha(res))
      apiHandlers.add('need_login', () => checkers.checkAdmin({}))

      return checkers
    },
    ['apiHandlers', 'layers', 'user', 'config'] as const,
  )
}

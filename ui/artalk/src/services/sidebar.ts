import type { ArtalkPlugin } from '@/types'
import { SidebarLayer } from '@/layer/sidebar-layer'

export const SidebarService: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'sidebar',
    (events, data, editor, checkers, config, user, layers) => {
      const sidebarLayer = new SidebarLayer({
        onShow: () => {
          setTimeout(() => {
            data.updateNotifies([])
          }, 0)

          events.trigger('sidebar-show')
        },
        onHide: () => {
          // prevent comment box from being swallowed
          editor.resetState()

          events.trigger('sidebar-hide')
        },
        getCheckers: () => checkers,
        getApi: () => ctx.getApi(),
        getConf: () => config,
        getUser: () => user,
        getLayers: () => layers,
      })

      ctx.on('user-changed', () => {
        sidebarLayer?.onUserChanged()
      })

      return sidebarLayer
    },
    ['events', 'data', 'editor', 'checkers', 'config', 'user', 'layers'] as const,
  )
}

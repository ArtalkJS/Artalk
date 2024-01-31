import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

/** List scroll to the first comment */
export const GotoFirst: ArtalkPlugin = (ctx) => {
  const handler = () => {
    const list = ctx.get('list')

    const $relative = ctx.conf.scrollRelativeTo && ctx.conf.scrollRelativeTo()
    ;($relative || window).scroll({
      top: Utils.getOffset(list.$el, $relative).top,
      left: 0,
    })
  }

  ctx.on('mounted', () => {
    ctx.on('list-goto-first', handler)
  })

  ctx.on('unmounted', () => {
    ctx.off('list-goto-first', handler)
  })
}

import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

/** List scroll to the first comment */
export const GotoFirst: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  const conf = ctx.inject('config')

  const handler = () => {
    const $relative = conf.get().scrollRelativeTo?.()
    ;($relative || window).scroll({
      top: Utils.getOffset(list.getEl(), $relative).top,
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

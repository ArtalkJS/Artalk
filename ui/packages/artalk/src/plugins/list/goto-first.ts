import type ArtalkPlugin from '~/types/plugin'
import * as Utils from '@/lib/utils'

/** List scroll to the first comment */
export const GotoFirst: ArtalkPlugin = (ctx) => {
  const handler = () => {
    const list = ctx.get('list')
    if (!list) return

    // TODO support set custom value to replace (`window`, `list.$el`) with (`conf.xxxAt`, `conf.list.repositionAt`)
    ;(ctx.conf.listScrollListenerAt || window).scroll({
      top: Utils.getOffset(list.$el).top,
      left: 0,
    })
  }

  ctx.on('inited', () => {
    ctx.on('list-goto-first', handler)
  })

  ctx.on('destroy', () => {
    ctx.off('list-goto-first', handler)
  })
}

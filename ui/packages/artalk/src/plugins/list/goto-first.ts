import type { ArtalkPlugin } from '~/types'
import * as Utils from '@/lib/utils'

/** List scroll to the first comment */
export const GotoFirst: ArtalkPlugin = (ctx) => {
  const handler = () => {
    const list = ctx.get('list')

    // TODO support set custom value to replace (`window`, `list.$el`) with (`conf.xxxAt`, `conf.list.repositionAt`)
    const $relative = (ctx.conf.scrollRelativeTo && ctx.conf.scrollRelativeTo()) || window
    $relative.scroll({
      top: Utils.getOffset(list.$el, $relative).top,
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

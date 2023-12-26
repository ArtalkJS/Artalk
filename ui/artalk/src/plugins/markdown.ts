import type { ArtalkPlugin } from '@/types'
import * as marked from '@/lib/marked'

export const Markdown: ArtalkPlugin = (ctx) => {
  marked.initMarked()

  ctx.on('conf-loaded', (conf) => {
    if (conf.markedReplacers) marked.setReplacers(conf.markedReplacers)
  })
}

import type { ArtalkPlugin } from '@/types'
import * as marked from '@/lib/marked'

export const Markdown: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['markedReplacers', 'imgLazyLoad', 'markedOptions'], (conf) => {
    marked.initMarked({
      markedOptions: ctx.getConf().markedOptions,
      imgLazyLoad: ctx.getConf().imgLazyLoad,
    })

    if (conf.markedReplacers) marked.setReplacers(conf.markedReplacers)
  })
}

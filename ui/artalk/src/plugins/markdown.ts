import type { ArtalkPlugin } from '@/types'
import * as marked from '@/lib/marked'

export const Markdown: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['imgLazyLoad', 'markedOptions'], (conf) => {
    marked.initMarked({
      markedOptions: conf.markedOptions,
      imgLazyLoad: conf.imgLazyLoad,
    })
  })

  ctx.watchConf(['markedReplacers'], (conf) => {
    conf.markedReplacers && marked.setReplacers(conf.markedReplacers)
  })
}

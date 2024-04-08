import type { ArtalkPlugin } from '@/types'
import * as Ui from '@/lib/ui'

export const Loading: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', (p) => {
    const list = ctx.get('list')

    if (p.offset === 0)
      // only show loading when fetch first page
      Ui.setLoading(true, list.$el)
    // else if not first page, show loading in paginator (code not there)
  })

  ctx.on('list-fetched', () => {
    const list = ctx.get('list')
    Ui.setLoading(false, list.$el)
  })
}

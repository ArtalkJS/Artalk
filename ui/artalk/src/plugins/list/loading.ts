import type { ArtalkPlugin } from '@/types'
import * as Ui from '@/lib/ui'

export const Loading: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')

  ctx.on('list-fetch', (p) => {
    if (p.offset === 0)
      // only show loading when fetch first page
      Ui.setLoading(true, list.getEl())
    // else if not first page, show loading in paginator (code not there)
  })

  ctx.on('list-fetched', () => {
    Ui.setLoading(false, list.getEl())
  })
}

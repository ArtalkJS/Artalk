import { ArtalkPlugin } from '@/types'
import { DataManager } from '@/data'

export const DataService: ArtalkPlugin = (ctx) => {
  ctx.provide('data', (events) => new DataManager(events), ['events'] as const)

  ctx.on('mounted', () => {
    // Load comment list immediately after mounting
    if (ctx.getConf().fetchCommentsOnInit) {
      ctx.getData().fetchComments({ offset: 0 })
    }
  })
}

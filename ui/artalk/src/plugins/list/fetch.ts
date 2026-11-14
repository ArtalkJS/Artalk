import type { ArtalkPlugin } from '@/types'

// Compatibility for plugins that trigger list-fetch directly. Normal requests
// call the comment service through DataManager.
export const Fetch: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', (params) => ctx.inject('comments').fromEvent(params))
}

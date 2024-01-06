import type { ArtalkPlugin } from '@/types'

export const Notifies: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', () => {
    // Fetch all unread notifies
    ctx.getApi().notify.getUnread().then(res => {
      ctx.getData().updateUnreads(res.notifies)
    })
  })
}

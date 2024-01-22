import type { ArtalkPlugin } from '@/types'

export const Notifies: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', (params) => {
    if (params.offset !== 0) return

    // Fetch all unread notifies
    ctx.getApi().notify.getNotifies().then(res => {
      ctx.getData().updateNotifies(res.notifies)
    })
  })
}

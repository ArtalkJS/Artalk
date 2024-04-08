import type { ArtalkPlugin } from '@/types'

export const Notifies: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', (params) => {
    if (params.offset !== 0) return

    const user = ctx.getApi().getUserFields()
    if (!user) return

    // Fetch all unread notifies
    ctx
      .getApi()
      .notifies.getNotifies(user)
      .then((res) => {
        ctx.getData().updateNotifies(res.data.notifies)
      })
  })
}

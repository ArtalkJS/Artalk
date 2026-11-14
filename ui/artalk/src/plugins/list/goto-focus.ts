import type { ArtalkPlugin } from '@/types'

export const GotoFocus: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  let request: AbortController | undefined
  const cancel = () => {
    request?.abort()
    request = undefined
  }
  ctx.watchConf(['server', 'site', 'pageKey'], cancel)
  ctx.on('unmounted', cancel)

  ctx.on('list-goto', async (commentID) => {
    cancel()
    const controller = (request = new AbortController())
    const identity = () => {
      const { name, email, token } = ctx.getUser().getData()
      return JSON.stringify([name, email, token])
    }
    const viewer = identity()
    const active = () => !controller.signal.aborted && !ctx.isDestroyed() && identity() === viewer
    try {
      let comment = ctx.getCommentNodes().find((c) => c.getID() === commentID)
      if (!comment) {
        const { data } = await ctx
          .getApi()
          .comments.getComment(commentID, { signal: controller.signal })
        if (!active()) return
        list.getLayout({ forceFlatMode: true }).insert(data.comment, data.reply_comment)
        comment = ctx.getCommentNodes().find((c) => c.getID() === commentID)
      }
      if (active()) comment?.focus()
    } catch (error) {
      if (active()) console.error(error)
    } finally {
      if (request === controller) request = undefined
    }
  })
}

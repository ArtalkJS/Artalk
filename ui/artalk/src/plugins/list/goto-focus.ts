import type { ArtalkPlugin } from '@/types'

export const GotoFocus: ArtalkPlugin = (ctx) => {
  ctx.on('list-goto', async (commentID) => {
    // find the comment node
    let comment = ctx.getCommentNodes().find((c) => c.getID() === commentID)
    if (!comment) {
      // fetch and insert the comment from the server
      const data = (await ctx.getApi().comments.getComment(commentID)).data
      ctx
        .get('list')
        .getListLayout({ forceFlatMode: true })
        .insert(data.comment, data.reply_comment)
      comment = ctx.getCommentNodes().find((c) => c.getID() === commentID)
    }
    if (!comment) return
    comment.focus()
  })
}

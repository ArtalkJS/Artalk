import type { ArtalkPlugin } from '~/types'
import * as Utils from '@/lib/utils'

export const Unread: ArtalkPlugin = (ctx) => {
  ctx.on('comment-rendered', (comment) => {
    // comment unread highlight
    if (ctx.conf.listUnreadHighlight === true) {
      const unreads = ctx.getData().getUnreads()
      const notify = unreads.find((o) => o.comment_id === comment.getID())

      if (notify) {
        // if comment contains in unread list
        comment.getRender().setUnread(true)
        comment.getRender().setOpenAction(() => {
          window.open(notify.read_link)

          // remove notify which has been read
          ctx
            .getData()
            .updateUnreads(
              unreads.filter((o) => o.comment_id !== comment.getID()),
            )
        })
      } else {
        // comment not in unread list
        comment.getRender().setUnread(false)
      }
    }
  })

  ctx.on('list-goto', (commentID) => {
    const notifyKey = Utils.getQueryParam('atk_notify_key')
    if (notifyKey) {
      // mark as read
      ctx
        .getApi()
        .user.markRead(commentID, notifyKey)
        .then(() => {
          // remove from unread list
          ctx.getData().updateUnreads(
            ctx
              .getData()
              .getUnreads()
              .filter((o) => o.comment_id !== commentID),
          )
        })
    }
  })
}

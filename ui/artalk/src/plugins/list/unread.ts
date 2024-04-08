import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

export const Unread: ArtalkPlugin = (ctx) => {
  ctx.on('comment-rendered', (comment) => {
    // comment unread highlight
    if (ctx.conf.listUnreadHighlight === true) {
      const notifies = ctx.getData().getNotifies()
      const notify = notifies.find((o) => o.comment_id === comment.getID())

      if (notify) {
        // if comment contains in unread list
        comment.getRender().setUnread(true)
        comment.getRender().setOpenAction(() => {
          window.open(notify.read_link)

          // remove notify which has been read
          ctx.getData().updateNotifies(notifies.filter((o) => o.comment_id !== comment.getID()))
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
        .notifies.markNotifyRead(commentID, notifyKey)
        .then(() => {
          // remove from unread list
          ctx.getData().updateNotifies(
            ctx
              .getData()
              .getNotifies()
              .filter((o) => o.comment_id !== commentID),
          )
        })
    }
  })
}

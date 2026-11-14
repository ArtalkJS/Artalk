import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

export const Unread: ArtalkPlugin = (ctx) => {
  const conf = ctx.inject('config')

  ctx.on('comment-rendered', (comment) => {
    // comment unread highlight
    if (conf.get().listUnreadHighlight) {
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
      const scope = () => {
        const { server, site, pageKey } = ctx.getConf()
        const { name, email, token } = ctx.getUser().getData()
        return JSON.stringify([server, site, pageKey, name, email, token])
      }
      const snapshot = scope()
      const active = () => !ctx.isDestroyed() && scope() === snapshot
      // mark as read
      ctx
        .getApi()
        .notifies.markNotifyRead(commentID, notifyKey)
        .then(() => {
          if (!active()) return
          // remove from unread list
          ctx.getData().updateNotifies(
            ctx
              .getData()
              .getNotifies()
              .filter((o) => o.comment_id !== commentID),
          )
        })
        .catch((error) => {
          if (active()) console.error(error)
        })
    }
  })
}

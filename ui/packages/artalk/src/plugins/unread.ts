import ArtalkPlug from '~/types/plug'
import * as Utils from '@/lib/utils'

export const Unread: ArtalkPlug = (ctx) => {
  ctx.on('unread-updated', (unreadList) => {
    const list = ctx.get('list')
    if (!list) return

    // comment unread highlight
    if (list.unreadHighlight === true) {
      ctx.getCommentList().forEach((comment) => {
        const notify = unreadList.find(o => o.comment_id === comment.getID())

        if (notify) {
          // if comment contains in unread list
          comment.getRender().setUnread(true)
          comment.getRender().setOpenAction(() => {
            window.open(notify.read_link)

            // remove notify which has been read
            ctx.updateUnreadList(unreadList.filter(o => o.comment_id !== comment.getID()))
          })
        } else {
          // comment not in unread list
          comment.getRender().setUnread(false)
        }
      })
    }
  })

  ctx.on('list-goto', (commentID) => {
    const notifyKey = Utils.getQueryParam('atk_notify_key')
    if (notifyKey) {
      // mark as read
      ctx.getApi().user.markRead(commentID, notifyKey)
        .then(() => {
          // remove from unread list
          ctx.updateUnreadList(ctx.getUnreadList().filter(o => o.comment_id !== commentID))
        })
    }
  })
}

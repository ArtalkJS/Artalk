import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'

export const Goto: ArtalkPlugin = (ctx) => {
  let delayGoto = true

  const check = () => {
    const commentID = extractCommentID()
    if (!commentID) return

    // trigger event
    ctx.trigger('list-goto', commentID)

    // reset delayGoto
    delayGoto = true
  }

  const hashChangeHandler = () => {
    delayGoto = false
    check()
  }
  ctx.on('mounted', () => {
    window.addEventListener('hashchange', hashChangeHandler)
    ctx.on('list-loaded', check)
  })
  ctx.on('unmounted', () => {
    window.removeEventListener('hashchange', hashChangeHandler)
    ctx.off('list-loaded', check)
  })

  let foundID = 0
  ctx.on('list-goto', (commentID) => {
    if (foundID === commentID) return

    const comment = ctx.get('list').getCommentNodes().find(c => c.getID() === commentID)
    if (!comment) return

    foundID = commentID

    // 若父评论存在 “子评论部分” 限高，取消限高
    comment.getParents().forEach((p) => {
      p.getRender().heightLimitRemoveForChildren()
    })

    const goTo = () => {
      Ui.scrollIntoView(comment.getEl(), false)

      comment.getEl().classList.remove('atk-flash-once')
      window.setTimeout(() => {
        comment.getEl().classList.add('atk-flash-once')
      }, 150)
    }

    if (!delayGoto) goTo()
    else window.setTimeout(() => goTo(), 350)
  })
}

function extractCommentID(): number|null {
  // try get from query
  let commentId = Number(Utils.getQueryParam('atk_comment')) // same as backend GetReplyLink()

  // fail over to get from hash
  if (!commentId) {
    const match = window.location.hash.match(/#atk-comment-([0-9]+)/)
    if (!match || !match[1] || Number.isNaN(Number(match[1]))) return null
    commentId = Number(match[1])
  }

  return commentId || null
}

import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'

export const GotoDispatcher: ArtalkPlugin = (ctx) => {
  let lastID = 0
  const check = ({ locker }: { locker: boolean }) => {
    const commentID = extractCommentID()
    if (!commentID) return

    if (locker && lastID === commentID) return // if the commentID is the same as the last one, do nothing
    lastID = commentID // record the last commentID

    ctx.trigger('list-goto', commentID) // trigger event
  }

  const hashChangeHandler = () => check({ locker: false })
  const listLoadedHandler = () => check({ locker: true })
  ctx.on('mounted', () => {
    window.addEventListener('hashchange', hashChangeHandler)
    ctx.on('list-loaded', listLoadedHandler)
  })
  ctx.on('unmounted', () => {
    window.removeEventListener('hashchange', hashChangeHandler)
    ctx.off('list-loaded', listLoadedHandler)
  })
}

function extractCommentID(): number | null {
  // Try get from hash
  // Hash retrieval priority is higher than query,
  // Because click goto will change hash.
  const match = window.location.hash.match(/#atk-comment-([0-9]+)/)
  let commentId =
    match && match[1] && !Number.isNaN(parseFloat(match[1])) ? parseFloat(match[1]) : null

  // Fail over to get from query
  if (!commentId) {
    commentId = Number(Utils.getQueryParam('atk_comment')) // same as backend GetReplyLink()
  }

  return commentId || null
}

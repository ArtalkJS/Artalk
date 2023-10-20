import ArtalkPlugin from '~/types/plugin'
import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import Comment from '@/comment/comment'

export const ListGoto: ArtalkPlugin = (ctx) => {
  const check = (delayGoto = true) => {
    const list = ctx.get('list')
    if (!list) return

    const commentID = extractCommentID()
    if (!commentID) return

    // 自动翻页
    const comment = ctx.findComment(commentID)
    if (!comment) { // 若找不到评论
      // TODO 自动范围改为直接跳转到计算后的页面
      list.pgHolder?.next()
      return
    }

    // trigger event
    ctx.trigger('list-goto', commentID)

    // goto comment
    gotoComment(comment, delayGoto)
  }

  // bind events
  ctx.on('list-loaded', () => { check() })
  window.addEventListener('hashchange', () => { check(false) })
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

function gotoComment(comment: Comment, delayGoto: boolean = true) {
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
}

import type { CommentData } from '@/types'
import * as Ui from '@/lib/ui'
import type { LayoutStrategyCreator, LayoutOptions } from '.'

export const createFlatStrategy: LayoutStrategyCreator = (opts) => ({
  import: (comments) => {
    comments.forEach((comment: CommentData) => {
      const replyComment = (comment.rid === 0 ? undefined : comments.find(c => c.id === comment.rid))
      insertComment(opts, 'append', comment, replyComment)
    })
  },
  insert: (comment, replyComment) => {
    const node = insertComment(opts, 'prepend', comment, replyComment)
    node.scrollIntoView()
  }
})

/** 导入评论 · 平铺模式 */
function insertComment(opts: LayoutOptions, insertMode: 'append'|'prepend', comment: CommentData, replyComment?: CommentData | undefined) {
  if (comment.is_collapsed) comment.is_allow_reply = false
  const node = opts.createCommentNode(comment, replyComment)

  // 可见评论添加到界面
  // 注：不可见评论用于显示 “引用内容”
  if (comment.visible) {
    const $comment = node.getEl()
    const $listCommentsWrap = opts.$commentsWrap
    if (insertMode === 'append') $listCommentsWrap?.append($comment)
    if (insertMode === 'prepend') $listCommentsWrap?.prepend($comment)
    node.getRender().playFadeAnim()
  }

  // 平铺评论插入后 · 内容限高检测
  // 插入列表的评论 (visible=true) 和关联的评论 (visible=false) 都需要检测
  node.getRender().checkHeightLimit()

  return node
}

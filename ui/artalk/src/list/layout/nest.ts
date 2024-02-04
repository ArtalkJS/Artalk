import type { CommentNode } from '@/comment'
import * as Ui from '@/lib/ui'
import * as ListNest from '@/list/nest'
import type { LayoutStrategyCreator } from '.'

export const createNestStrategy: LayoutStrategyCreator = (opts) => ({
  import: (comments) => {
    // 遍历 root 评论
    const rootNodes = ListNest.makeNestCommentNodeList(comments, opts.nestSortBy, opts.nestMax)
    rootNodes.forEach((rootNode: ListNest.CommentNode) => {
      const rootC = opts.createCommentNode(rootNode.comment)

      // 显示并播放渐入动画
      opts.$commentsWrap?.appendChild(rootC.getEl())
      rootC.getRender().playFadeAnim()

      // 加载子评论
      const loadChildren = (parentC: CommentNode, parentNode: ListNest.CommentNode) => {
        parentNode.children.forEach((node: ListNest.CommentNode) => {
          const childD = node.comment
          const childC = opts.createCommentNode(childD, parentC.getData())

          // 插入到父评论中
          parentC.putChild(childC)

          // 递归加载子评论
          loadChildren(childC, node)
        })
      }
      loadChildren(rootC, rootNode)

      // 限高检测
      rootC.getRender().checkHeightLimit()
    })
  },
  insert: (comment, replyComment) => {
    // 嵌套模式
    const node = opts.createCommentNode(comment, replyComment)

    if (comment.rid === 0) {
      // root评论 新增
      opts.$commentsWrap?.prepend(node.getEl())
    } else {
      // 子评论 新增
      const parent = opts.findCommentNode(comment.rid)
      if (parent) {
        parent.putChild(node, (opts.nestSortBy === 'DATE_ASC' ? 'append' : 'prepend'))

        // 若父评论存在 “子评论部分” 限高，取消限高
        node.getParents().forEach((p) => {
          p.getRender().heightLimitRemoveForChildren()
        })
      }
    }

    node.getRender().checkHeightLimit()

    node.scrollIntoView() // 滚动到可以见
    node.getRender().playFadeAnim() // 播放评论渐出动画
  }
})

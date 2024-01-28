import type { CommentData } from '@/types'
import * as Ui from '@/lib/ui'
import type { CommentNode } from '@/comment'
import * as ListNest from './nest'

export interface LayoutOptions {
  $commentsWrap: HTMLElement
  nestSortBy: ListNest.SortByType
  nestMax: number
  flatMode: boolean

  createCommentNode(comment: CommentData, ctxComments: CommentData[]): CommentNode
  findCommentNode(id: number): CommentNode|undefined
  getCommentDataList(): CommentData[]
}

export default class ListLayout {
  constructor(private options: LayoutOptions) {}

  // TODO: refactor `if syntax` to strategy pattern
  import(comments: CommentData[]) {
    if (this.options.flatMode) {
      comments.forEach((commentData: CommentData) => {
        this.putCommentFlatMode(commentData, comments, 'append')
      })
    } else {
      this.importCommentsNestMode(comments)
    }
  }

  insert(comment: CommentData) {
    if (!this.options.flatMode) {
      this.insertCommentNest(comment)
    } else {
      this.insertCommentFlatMode(comment)
    }
  }

  // 导入评论 · 嵌套模式
  private importCommentsNestMode(srcData: CommentData[]) {
    // 遍历 root 评论
    const rootNodes = ListNest.makeNestCommentNodeList(srcData, this.options.nestSortBy, this.options.nestMax)
    rootNodes.forEach((rootNode: ListNest.CommentNode) => {
      const rootC = this.options.createCommentNode(rootNode.comment, srcData)

      // 显示并播放渐入动画
      this.options.$commentsWrap?.appendChild(rootC.getEl())
      rootC.getRender().playFadeAnim()

      // 加载子评论
      const loadChildren = (parentC: CommentNode, parentNode: ListNest.CommentNode) => {
        parentNode.children.forEach((node: ListNest.CommentNode) => {
          const childD = node.comment
          const childC = this.options.createCommentNode(childD, srcData)

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
  }

  /** 导入评论 · 平铺模式 */
  private putCommentFlatMode(cData: CommentData, ctxData: CommentData[], insertMode: 'append'|'prepend') {
    if (cData.is_collapsed) cData.is_allow_reply = false
    const comment = this.options.createCommentNode(cData, ctxData)

    // 可见评论添加到界面
    // 注：不可见评论用于显示 “引用内容”
    if (cData.visible) {
      if (insertMode === 'append') this.options.$commentsWrap?.append(comment.getEl())
      if (insertMode === 'prepend') this.options.$commentsWrap?.prepend(comment.getEl())
      comment.getRender().playFadeAnim()
    }

    // 平铺评论插入后 · 内容限高检测
    comment.getRender().checkHeightLimit()

    return comment
  }

  private insertCommentNest(commentData: CommentData) {
    // 嵌套模式
    const comment = this.options.createCommentNode(commentData, this.options.getCommentDataList())

    if (commentData.rid === 0) {
      // root评论 新增
      this.options.$commentsWrap?.prepend(comment.getEl())
    } else {
      // 子评论 新增
      const parent = this.options.findCommentNode(commentData.rid)
      if (parent) {
        parent.putChild(comment, (this.options.nestSortBy === 'DATE_ASC' ? 'append' : 'prepend'))

        // 若父评论存在 “子评论部分” 限高，取消限高
        comment.getParents().forEach((p) => {
          p.getRender().heightLimitRemoveForChildren()
        })
      }
    }

    comment.getRender().checkHeightLimit()

    Ui.scrollIntoView(comment.getEl()) // 滚动到可以见
    comment.getRender().playFadeAnim() // 播放评论渐出动画
  }

  private insertCommentFlatMode(commentData: CommentData) {
    // 平铺模式
    const comment = this.putCommentFlatMode(commentData, this.options.getCommentDataList(), 'prepend')
    Ui.scrollIntoView(comment.getEl()) // 滚动到可见
  }
}

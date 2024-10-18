import type { CommentNode } from '../comment'
import type { CommentData } from './data'

export interface List {
  getEl: () => HTMLElement
  getCommentsWrapEl: () => HTMLElement
  getLayout: (arg: { forceFlatMode?: boolean }) => ListLayout
  getCommentNodes: () => CommentNode[]
}

export interface ListLayout {
  import: (comments: CommentData[]) => void
  insert: (comment: CommentData, replyComment?: CommentData) => void
}

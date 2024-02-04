import type { CommentData } from '@/types'
import type { CommentNode } from '@/comment'
import * as ListNest from '../nest'
import { createNestStrategy } from './nest'
import { createFlatStrategy } from './flat'

export interface LayoutOptions {
  /** The comments wrap of list */
  $commentsWrap: HTMLElement
  nestSortBy: ListNest.SortByType
  nestMax: number
  flatMode: boolean

  createCommentNode(comment: CommentData, replyComment?: CommentData): CommentNode
  findCommentNode(id: number): CommentNode|undefined
}

export interface LayoutStrategy {
  import(comments: CommentData[]): void
  insert(comment: CommentData, replyComment?: CommentData): void
}

export type LayoutStrategyCreator = (opts: LayoutOptions) => LayoutStrategy

export class ListLayout {
  constructor(private options: LayoutOptions) {}

  private getStrategy() {
    return this.options.flatMode ? createFlatStrategy(this.options) : createNestStrategy(this.options)
  }

  import(comments: CommentData[]) {
    this.getStrategy().import(comments)
  }

  insert(comment: CommentData, replyComment?: CommentData) {
    this.getStrategy().insert(comment, replyComment)
  }
}

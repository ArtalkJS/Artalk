import ListHTML from './list.html?raw'
import { ListLayout } from './layout'
import { createCommentNode } from './comment'
import { initListPaginatorFunc } from './page'
import type { CommentData, EventManager, DataManager, ConfigManager, List as IList } from '@/types'
import * as Utils from '@/lib/utils'
import { CommentNode } from '@/comment'
import type { Api } from '@/api'

export interface ListOptions {
  getApi: () => Api
  getEvents: () => EventManager
  getConf: () => ConfigManager
  getData: () => DataManager

  replyComment: (c: CommentData, $el: HTMLElement) => void
  editComment: (c: CommentData, $el: HTMLElement) => void
  resetEditorState: () => void
  onListGotoFirst?: () => void
}

export class List implements IList {
  private opts: ListOptions
  private $el: HTMLElement
  getEl() {
    return this.$el
  }

  private $commentsWrap: HTMLElement
  getCommentsWrapEl() {
    return this.$commentsWrap
  }

  private commentNodes: CommentNode[] = []
  getCommentNodes() {
    return this.commentNodes
  }

  constructor(opts: ListOptions) {
    this.opts = opts

    // Init base element
    this.$el = Utils.createElement(ListHTML)
    this.$commentsWrap = this.$el.querySelector('.atk-list-comments-wrap')!

    // Init paginator
    initListPaginatorFunc({
      getList: () => this,
      ...opts,
    })

    // Bind events
    this.initCrudEvents()
  }

  getLayout({ forceFlatMode }: { forceFlatMode?: boolean } = {}) {
    return new ListLayout({
      $commentsWrap: this.$commentsWrap,
      nestSortBy: this.opts.getConf().get().nestSort,
      nestMax: this.opts.getConf().get().nestMax,
      flatMode:
        typeof forceFlatMode === 'boolean'
          ? forceFlatMode
          : (this.opts.getConf().get().flatMode as boolean),
      // flatMode must be boolean because it had been handled when Artalk.init
      createCommentNode: (d, r) => {
        const node = createCommentNode({ forceFlatMode, ...this.opts }, d, r)
        this.commentNodes.push(node) // store node instance
        return node
      },
      findCommentNode: (id) => this.commentNodes.find((c) => c.getID() === id),
    })
  }

  private initCrudEvents() {
    this.opts.getEvents().on('list-load', (comments) => {
      // 导入数据
      this.getLayout().import(comments)
    })

    this.opts.getEvents().on('list-loaded', (comments) => {
      if (comments.length === 0) {
        this.commentNodes = []
        this.$commentsWrap.innerHTML = ''
      }
    })

    // When comment insert
    this.opts.getEvents().on('comment-inserted', (comment) => {
      const replyComment = comment.rid
        ? this.commentNodes.find((c) => c.getID() === comment.rid)?.getData()
        : undefined
      this.getLayout().insert(comment, replyComment)
    })

    // When comment delete
    this.opts.getEvents().on('comment-deleted', (comment) => {
      const node = this.commentNodes.find((c) => c.getID() === comment.id)
      if (!node) {
        console.error(`comment node id=${comment.id} not found`)
        return
      }
      node.remove()
      this.commentNodes = this.commentNodes.filter((c) => c.getID() !== comment.id)
      // TODO: remove child nodes
    })

    // When comment update
    this.opts.getEvents().on('comment-updated', (comment) => {
      const node = this.commentNodes.find((c) => c.getID() === comment.id)
      node && node.setData(comment)
    })
  }
}

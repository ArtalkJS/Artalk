import type { NotifyData, PageData, CommentData, DataManagerApi, ListFetchParams, ListLastFetchData, EventPayloadMap } from '@/types'
import EventManager from './lib/event-manager'

export class DataManager implements DataManagerApi {
  private loading: boolean = false
  private listLastFetch?: ListLastFetchData
  private comments: CommentData[] = [] // Note: 无层级结构 + 无须排列
  private notifies: NotifyData[] = []
  private page?: PageData

  constructor(
    protected events: EventManager<EventPayloadMap>
  ) {}

  getLoading() {
    return this.loading
  }

  setLoading(val: boolean) {
    this.loading = val
  }

  getListLastFetch() {
    return this.listLastFetch
  }

  setListLastFetch(val: ListLastFetchData) {
    this.listLastFetch = val
  }

  // -------------------------------------------------------------------
  //  Comments
  // -------------------------------------------------------------------
  getComments() {
    return this.comments
  }

  fetchComments(params: Partial<ListFetchParams>) {
    this.events.trigger('list-fetch', params)
  }

  findComment(id: number) {
    return this.comments.find(c => c.id === id)
  }

  clearComments() {
    this.comments = []
    this.events.trigger('list-loaded', this.comments)
  }

  loadComments(partialComments: CommentData[]) {
    this.events.trigger('list-load', partialComments)

    this.comments.push(...partialComments)

    this.events.trigger('list-loaded', this.comments)
  }

  insertComment(comment: CommentData) {
    this.comments.push(comment)

    this.events.trigger('comment-inserted', comment)
    this.events.trigger('list-loaded', this.comments) // list-loaded should always keep the last
  }

  updateComment(comment: CommentData) {
    this.comments = this.comments.map(c => {
      if (c.id === comment.id) return comment
      return c
    })

    this.events.trigger('comment-updated', comment)
    this.events.trigger('list-loaded', this.comments)
  }

  deleteComment(id: number) {
    const comment = this.comments.find(c => c.id === id)
    if (!comment) throw new Error(`Comment ${id} not found`)
    this.comments = this.comments.filter(c => c.id !== id)

    this.events.trigger('comment-deleted', comment)
    this.events.trigger('list-loaded', this.comments)
  }

  // -------------------------------------------------------------------
  //  Notifies
  // -------------------------------------------------------------------
  getNotifies() {
    return this.notifies
  }

  updateNotifies(notifies: NotifyData[]) {
    this.notifies = notifies

    this.events.trigger('notifies-updated', this.notifies)
  }

  // -------------------------------------------------------------------
  // Page
  // -------------------------------------------------------------------
  getPage() {
    return this.page
  }

  updatePage(pageData: PageData) {
    this.page = pageData

    this.events.trigger('page-loaded', pageData)
  }
}

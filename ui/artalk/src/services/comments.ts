import type { ArtalkPlugin, Context, ListData, ListFetchParams } from '@/types'

/** Owns list loading, cancellation, and updates to the current view. */
export class CommentService {
  private request?: AbortController
  private disposed = false
  private notifying = false
  private identity = ''

  constructor(private ctx: Context) {}

  cancel() {
    const previous = this.request
    this.request?.abort()
    this.request = undefined
    this.ctx.getData().setLoading(false)
    if (previous) {
      this.ctx.trigger('list-fetched', {
        params: this.ctx.getData().getListLastFetch()?.params || {},
      })
    }
  }

  dispose() {
    this.disposed = true
    this.cancel()
  }

  fromEvent(input: Partial<ListFetchParams>) {
    if (!this.notifying) void this.load(input, false)
  }

  refreshIdentity() {
    if (!this.request || this.identity === this.getIdentity()) return
    // Resume an interrupted read with the same filters, without callbacks from
    // the previous viewer. Unrelated profile edits must not cancel the request.
    const { limit, flatMode, paramsModifier } = this.ctx.getData().getListLastFetch()!.params
    void this.load({ limit, flatMode, paramsModifier })
  }

  private getIdentity() {
    const { name, email, token } = this.ctx.getUser().getData()
    return JSON.stringify([name, email, token])
  }

  async load(input: Partial<ListFetchParams> = {}, notify = true): Promise<void> {
    if (this.disposed || this.ctx.isDestroyed()) return
    // A reload replaces an older request. Repeated load-more clicks are ignored.
    if (this.request && (input.offset ?? 0) !== 0) return
    this.cancel()
    // Cancellation callbacks may destroy the widget or start a newer request.
    if (this.disposed || this.ctx.isDestroyed() || this.request) return
    const conf = this.ctx.getConf()
    const params: ListFetchParams = {
      offset: 0,
      limit: conf.pagination.pageSize,
      flatMode: conf.flatMode as boolean,
      paramsModifier: conf.listFetchParamsModifier,
      ...input,
    }
    this.identity = this.getIdentity()
    const request = (this.request = new AbortController())
    const active = () =>
      this.request === request && !request.signal.aborted && !this.ctx.isDestroyed()
    const data = this.ctx.getData()
    data.setLoading(true)
    data.setListLastFetch({ params })

    try {
      // Preserve the existing notification event and its ordering for plugins.
      if (notify) {
        this.notifying = true
        try {
          this.ctx.trigger('list-fetch', params)
        } finally {
          this.notifying = false
        }
      }
      if (!active()) return
      const reqParams = {
        limit: params.limit,
        offset: params.offset,
        flat_mode: params.flatMode,
        page_key: conf.pageKey,
        site_name: conf.site,
      }
      params.paramsModifier?.(reqParams)
      if (!active()) return
      const api = this.ctx.getApi()
      const { data: result } = await api.comments.getComments(
        { ...reqParams, ...api.getUserFields() },
        { signal: request.signal },
      )
      if (!active()) return
      const scope = (reqParams as { scope?: string }).scope
      if ((!scope || scope === 'page') && !result.page) {
        throw new Error('Page data is missing from comment list response')
      }
      // User and site scopes have no page; keep the existing public ListData type.
      const listData = result as ListData
      data.setListLastFetch({ params, data: listData })
      data.loadComments(result.comments)
      if (!active()) return
      if (result.page) data.updatePage(result.page)
      if (!active()) return
      params.onSuccess?.(listData)
      if (!active()) return
      this.ctx.trigger('list-fetched', { params, data: listData })
    } catch (err: any) {
      if (!active()) return
      const error = { msg: err?.msg || String(err), data: err?.data }
      params.onError?.(error)
      if (!active()) return
      this.ctx.trigger('list-failed', error)
      if (!active()) return
      this.ctx.trigger('list-fetched', { params, error })
    } finally {
      if (active()) {
        this.request = undefined
        data.setLoading(false)
      }
    }
  }
}

export const CommentsService: ArtalkPlugin = (ctx) => {
  ctx.provide('comments', () => new CommentService(ctx), [] as const)
  // A changed destination or identity invalidates results already in flight.
  const destination = () => {
    const { server, site, pageKey } = ctx.getConf()
    return JSON.stringify([server, site, pageKey])
  }
  let previous = destination()
  ctx.watchConf(['server', 'site', 'pageKey'], () => {
    const current = destination()
    if (previous !== current) {
      ctx.inject('comments').cancel()
      const data = ctx.getData()
      if (data.reset) data.reset()
      else data.clearComments()
    }
    previous = current
  })
  ctx.on('user-changed', () => ctx.inject('comments').refreshIdentity())
  ctx.on('unmounted', () => ctx.inject('comments').dispose())
}

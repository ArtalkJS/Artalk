import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { comment, deferred, list, response, setup } from './lib/api'
import type { ListFetchedArgs } from '@/types'

beforeEach(() => localStorage.clear())
afterEach(() => vi.restoreAllMocks())

describe('comment request lifecycle', () => {
  it('does not restart loading if a cancellation callback destroys the widget', async () => {
    const { ctx, api, data, service } = setup()
    const pending = deferred<ReturnType<typeof list>>()
    vi.mocked(api.comments.getComments).mockImplementation(async () =>
      response(await pending.promise),
    )
    const first = service.load()
    ctx.on('list-fetched', () => ctx.destroy())
    await service.load()
    pending.resolve(list())
    await first
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    expect(data.getLoading()).toBe(false)
  })
  it('loads once through the public data manager and keeps event order', async () => {
    const { ctx, data, api } = setup()
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
    const events: string[] = []
    for (const name of [
      'list-fetch',
      'list-load',
      'list-loaded',
      'page-loaded',
      'list-fetched',
    ] as const) {
      ctx.on(name, () => events.push(name))
    }
    data.fetchComments({ offset: 0 })
    await vi.waitFor(() => expect(data.getLoading()).toBe(false))
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    expect(data.getComments()).toEqual([comment(1)])
    expect(events).toEqual([
      'list-fetch',
      'list-load',
      'list-loaded',
      'page-loaded',
      'list-fetched',
    ])
  })

  it('accepts the legacy list-fetch event without a second request or notification', async () => {
    const { ctx, data, api } = setup()
    const onStart = vi.fn()
    ctx.on('list-fetch', onStart)
    ctx.trigger('list-fetch', { offset: 0 })
    await vi.waitFor(() => expect(data.getLoading()).toBe(false))
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    expect(onStart).toHaveBeenCalledTimes(1)
  })

  it('ignores an old response even if the request completes after cancellation', async () => {
    const old = deferred<ReturnType<typeof list>>()
    const current = deferred<ReturnType<typeof list>>()
    const { api, data, service } = setup()
    const getComments = vi.mocked(api.comments.getComments)
    getComments
      .mockImplementationOnce(async () => response(await old.promise))
      .mockImplementationOnce(async () => response(await current.promise))
    const success = vi.fn()
    const oldLoad = service.load({ onSuccess: success })
    const signal = getComments.mock.calls[0][1]!.signal!
    const currentLoad = service.load()
    expect(signal.aborted).toBe(true)
    old.resolve(list([comment(1)]))
    await oldLoad
    expect(data.getLoading()).toBe(true)
    expect(data.getComments()).toEqual([])
    expect(success).not.toHaveBeenCalled()
    current.resolve(list([comment(2)]))
    await currentLoad
    expect(data.getLoading()).toBe(false)
    expect(data.getComments()).toEqual([comment(2)])
  })

  it('deduplicates load-more clicks while a request is running', async () => {
    const pending = deferred<ReturnType<typeof list>>()
    const { api, service } = setup()
    vi.mocked(api.comments.getComments).mockImplementation(async () =>
      response(await pending.promise),
    )
    const loading = service.load({ offset: 20 })
    await service.load({ offset: 20 })
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    pending.resolve(list())
    await loading
  })

  it.each(['destroy', 'pageKey', 'server', 'site'] as const)(
    'discards pending data after %s changes',
    async (change) => {
      const pending = deferred<ReturnType<typeof list>>()
      const { ctx, api, data, service } = setup()
      vi.mocked(api.comments.getComments).mockImplementation(async () =>
        response(await pending.promise),
      )
      ctx.updateConf({ fetchCommentsOnInit: false })
      ctx.trigger('mounted')
      const success = vi.fn()
      const loading = service.load({ onSuccess: success })
      if (change === 'destroy') ctx.destroy()
      else ctx.updateConf({ [change]: 'https://other.example/next' })
      expect(vi.mocked(api.comments.getComments).mock.calls[0][1]!.signal!.aborted).toBe(true)
      pending.resolve(list([comment(1)]))
      await loading
      expect(data.getLoading()).toBe(false)
      expect(data.getComments()).toEqual([])
      expect(success).not.toHaveBeenCalled()
    },
  )

  it('retries with the new identity and the same legacy query modifier', async () => {
    const old = deferred<ReturnType<typeof list>>()
    const current = deferred<ReturnType<typeof list>>()
    const { ctx, api, data, service } = setup()
    const getComments = vi.mocked(api.comments.getComments)
    getComments
      .mockImplementationOnce(async () => response(await old.promise))
      .mockImplementationOnce(async () => response(await current.promise))
    ctx.getUser().update({ name: 'Alice', email: 'alice@example.com' })
    const oldSuccess = vi.fn()
    const filter = (params: any) => {
      params.scope = 'user'
    }
    const loading = service.load({ paramsModifier: filter, onSuccess: oldSuccess })
    const signal = getComments.mock.calls[0][1]!.signal!
    ctx.getUser().update({ name: 'Bob', email: 'bob@example.com', token: 'new-session' })
    expect(signal.aborted).toBe(true)
    expect(getComments).toHaveBeenCalledTimes(2)
    expect(getComments).toHaveBeenLastCalledWith(
      expect.objectContaining({ scope: 'user', name: 'Bob', email: 'bob@example.com' }),
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    )
    old.resolve(list([comment(1)]))
    await loading
    expect(data.getLoading()).toBe(true)
    expect(data.getComments()).toEqual([])
    current.resolve(list([comment(2)]))
    await vi.waitFor(() => expect(data.getLoading()).toBe(false))
    expect(data.getComments()).toEqual([comment(2)])
    expect(oldSuccess).not.toHaveBeenCalled()
  })

  it('does not interrupt a list request when only the website field changes', async () => {
    const pending = deferred<ReturnType<typeof list>>()
    const { ctx, api, data, service } = setup()
    vi.mocked(api.comments.getComments).mockImplementation(async () =>
      response(await pending.promise),
    )
    const loading = service.load()
    const signal = vi.mocked(api.comments.getComments).mock.calls[0][1]!.signal!
    ctx.getUser().update({ link: 'https://example.com/about' })
    expect(signal.aborted).toBe(false)
    pending.resolve(list([comment(1)]))
    await loading
    expect(data.getComments()).toEqual([comment(1)])
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
  })

  it('reports failures, releases loading, and allows retry without an unhandled rejection', async () => {
    const { ctx, api, data, service } = setup()
    vi.mocked(api.comments.getComments).mockRejectedValueOnce(new Error('Unavailable'))
    const failure = vi.fn()
    const finished: ListFetchedArgs[] = []
    ctx.on('list-failed', failure)
    ctx.on('list-fetched', (args) => finished.push(args))
    await service.load()
    expect(data.getLoading()).toBe(false)
    expect(failure).toHaveBeenCalledWith({ msg: 'Error: Unavailable', data: undefined })
    expect(finished).toHaveLength(1)
    await service.load()
    expect(api.comments.getComments).toHaveBeenCalledTimes(2)
    expect(finished[1].data).toEqual(list())
  })

  it.each(['event', 'modifier'])(
    'does not send a request if the %s destroys the instance',
    async (stage) => {
      const { ctx, api, service } = setup()
      if (stage === 'event') ctx.on('list-fetch', () => ctx.destroy())
      await service.load({ paramsModifier: stage === 'modifier' ? () => ctx.destroy() : undefined })
      expect(api.comments.getComments).not.toHaveBeenCalled()
    },
  )

  it('does not apply an old failure after the error callback starts a new request', async () => {
    const pending = deferred<ReturnType<typeof list>>()
    const { ctx, api, data, service } = setup()
    vi.mocked(api.comments.getComments)
      .mockRejectedValueOnce(new Error('Retry'))
      .mockImplementationOnce(async () => response(await pending.promise))
    const failed = vi.fn()
    ctx.on('list-failed', failed)
    let retry: Promise<void> | undefined
    await service.load({
      onError: () => {
        retry = service.load()
      },
    })
    expect(failed).not.toHaveBeenCalled()
    expect(data.getLoading()).toBe(true)
    pending.resolve(list([comment(2)]))
    await retry
    expect(data.getComments()).toEqual([comment(2)])
    expect(data.getLoading()).toBe(false)
  })
})

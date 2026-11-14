import { afterEach, beforeEach, expect, it, vi } from 'vitest'
import { comment, deferred, list, response, setup } from './lib/api'
import { GotoFocus } from '@/plugins/list/goto-focus'
import { Unread } from '@/plugins/list/unread'
import { PageVoteWidget } from '@/plugins/page-vote'
import { ApiService } from '@/services/api'
import type { Context } from '@/types'

const contexts: Context[] = []
function initialize() {
  const state = setup()
  contexts.push(state.ctx)
  state.ctx.updateConf({ fetchCommentsOnInit: false })
  return state
}
beforeEach(() => localStorage.clear())
afterEach(() => {
  contexts.splice(0).forEach((ctx) => ctx.destroy())
  document.body.replaceChildren()
  window.history.replaceState(null, '', '/')
  vi.restoreAllMocks()
})

it.each(['server', 'site', 'pageKey', 'account', 'destroy'] as const)(
  'ignores comment lookup after %s changes',
  async (change) => {
    const { ctx, api } = initialize()
    const pending = deferred<Awaited<ReturnType<typeof api.comments.getComment>>>()
    vi.mocked(api.comments.getComment).mockReturnValue(pending.promise)
    const insert = vi.fn()
    ctx.provide(
      'list',
      () => ({ getCommentNodes: () => [], getLayout: () => ({ insert }) }) as any,
      [] as const,
    )
    GotoFocus(ctx)
    ctx.trigger('mounted')
    ctx.trigger('list-goto', 1)
    if (change === 'destroy') ctx.destroy()
    else if (change === 'account') ctx.getUser().update({ token: 'another-user' })
    else ctx.updateConf({ [change]: 'another-destination' })
    if (change !== 'account')
      expect(vi.mocked(api.comments.getComment).mock.calls[0][1]!.signal!.aborted).toBe(true)
    pending.resolve(response({ comment: comment(1), reply_comment: comment(0) }))
    await pending.promise
    await Promise.resolve()
    expect(insert).not.toHaveBeenCalled()
  },
)

it('does not remove notifications from a new account after an old read completes', async () => {
  const { ctx, api } = initialize()
  const pending = deferred<Awaited<ReturnType<typeof api.notifies.markNotifyRead>>>()
  vi.mocked(api.notifies.markNotifyRead).mockReturnValue(pending.promise)
  Unread(ctx)
  window.history.replaceState(null, '', '/?atk_notify_key=test')
  ctx.trigger('list-goto', 1)
  ctx.getUser().update({ token: 'another-user' })
  const notify = {
    id: 1,
    comment_id: 1,
    user_id: 2,
    is_read: false,
    is_emailed: false,
    read_link: '/',
  }
  ctx.getData().updateNotifies([notify])
  pending.resolve(response({}))
  await pending.promise
  await Promise.resolve()
  expect(ctx.getData().getNotifies()).toEqual([notify])
})

it('refetches votes for a new site with the same page ID and ignores old counts', async () => {
  const { ctx, api } = initialize()
  document.body.innerHTML = '<span class="artalk-page-vote-up-count"></span>'
  const old = deferred<Awaited<ReturnType<typeof api.votes.getVote>>>()
  vi.mocked(api.votes.getVote).mockReturnValueOnce(old.promise)
  PageVoteWidget(ctx)
  ctx.updateConf({ pageVote: true })
  ctx.trigger('mounted')
  ctx.trigger('list-fetched', { params: {}, data: list() })
  ctx.updateConf({ site: 'Another site' })
  vi.mocked(api.votes.getVote).mockResolvedValue(
    response({ up: 3, down: 0, is_up: false, is_down: false }),
  )
  ctx.trigger('list-fetched', { params: {}, data: list() })
  await Promise.resolve()
  expect(api.votes.getVote).toHaveBeenCalledTimes(2)
  old.resolve(response({ up: 99, down: 0, is_up: true, is_down: false }))
  await old.promise
  await Promise.resolve()
  expect(document.querySelector<HTMLElement>('.artalk-page-vote-up-count')!.innerText).toBe('3')
})

it('does not launch a challenge from an aborted request after its response arrived', async () => {
  const { ctx } = initialize()
  vi.mocked(ctx.getApi).mockRestore()
  ApiService(ctx)
  const body = deferred<any>()
  vi.spyOn(globalThis, 'fetch').mockResolvedValue({
    ok: false,
    status: 403,
    json: () => body.promise,
  } as Response)
  const challenge = vi.fn()
  // eslint-disable-next-line artalk/artalk-plugin -- Test composition root.
  ctx.inject('apiHandlers').add('need_login', challenge)
  const controller = new AbortController()
  const pending = ctx.getApi().user.getUserStatus({}, { signal: controller.signal })
  await Promise.resolve()
  controller.abort()
  body.resolve({ need_login: true })
  await expect(pending).rejects.toMatchObject({ name: 'AbortError' })
  expect(challenge).not.toHaveBeenCalled()
})

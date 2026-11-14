import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { comment, deferred, setup } from './lib/api'
import SubmitAddPreset from '@/plugins/editor/submit-add'
import type PlugKit from '@/plugins/editor/_kit'
import { ApiService } from '@/services/api'

beforeEach(() => localStorage.clear())
afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

function initialize(content = 'Reply') {
  const state = setup()
  vi.mocked(state.ctx.getApi).mockRestore()
  ApiService(state.ctx)
  state.ctx
    .getUser()
    .update({ name: 'Alice', email: 'alice@example.com', link: 'https://example.com', token: '' })
  const kit = {
    useApi: () => state.ctx.getApi(),
    useConf: () => state.ctx.getConf(),
    useUser: () => state.ctx.getUser(),
    useEditor: () => ({ getContentFinal: () => content }),
  } as unknown as PlugKit
  const fetch = vi.fn<typeof globalThis.fetch>(
    async () => new Response(JSON.stringify(comment(1)), { status: 200 }),
  )
  vi.stubGlobal('fetch', fetch)
  return {
    ...state,
    preset: new SubmitAddPreset(kit),
    fetch,
    handlers: state.ctx.inject('apiHandlers'),
  }
}

describe('comment submission compatibility', () => {
  it('keeps guest headers and the original reply destination', async () => {
    const { ctx, preset, fetch } = initialize()
    await preset.reqAdd({ ...comment(42), page_key: '/original', site_name: 'Another site' })
    const [url, init] = fetch.mock.calls[0]
    expect(url).toBe('https://example.com/api/v2/comments')
    const headers = new Headers(init?.headers)
    expect(headers.get('X-Artalk-Guest')).toMatch(/^[a-f0-9]{64}$/)
    expect(headers.has('Authorization')).toBe(false)
    const body = JSON.parse(init!.body as string)
    expect(body).toMatchObject({
      content: 'Reply',
      rid: 42,
      page_key: '/original',
      site_name: 'Another site',
      name: 'Alice',
      email: 'alice@example.com',
    })
    expect(body.page_title).toBeUndefined()
    expect(ctx.getUser().getData().token).toBe('')
  })

  it('keeps the existing submit parameter helper and new-comment fields', async () => {
    const { ctx, preset, fetch } = initialize('Hello')
    ctx.updateConf({ pageKey: '/post', pageTitle: 'Article', site: 'Blog' })
    expect(await preset.getSubmitAddParams()).toMatchObject({
      content: 'Hello',
      rid: 0,
      page_key: '/post',
      page_title: 'Article',
      site_name: 'Blog',
    })
    expect(fetch).not.toHaveBeenCalled()
    await preset.reqAdd()
    expect(JSON.parse(fetch.mock.calls[0][1]!.body as string)).toMatchObject({
      page_title: 'Article',
      rid: 0,
    })
  })

  it.each(['account', 'server', 'site', 'pageKey'] as const)(
    'does not send a prepared comment after %s changes',
    async (change) => {
      const { ctx, preset, fetch } = initialize('Private identity')
      const hints = deferred<{ platformVersion: string }>()
      vi.stubGlobal('navigator', {
        userAgent: 'Test',
        userAgentData: { platform: 'Windows', getHighEntropyValues: () => hints.promise },
      })
      const creating = preset.reqAdd()
      if (change === 'account') ctx.getUser().update({ token: 'another-account' })
      else ctx.updateConf({ [change]: 'https://other.example/next' })
      hints.resolve({ platformVersion: '13.0.0' })
      await expect(creating).rejects.toThrow('account or page changed')
      expect(fetch).not.toHaveBeenCalled()
    },
  )

  it.each(['need_login', 'need_captcha'] as const)(
    'preserves the existing %s retry',
    async (action) => {
      const { ctx, preset, fetch, handlers } = initialize('After verification')
      const challenge = vi.fn(async () => {
        if (action === 'need_login')
          ctx.getUser().update({ token: 'verified-token', is_admin: true })
      })
      handlers.add(action, challenge)
      fetch.mockResolvedValueOnce(new Response(JSON.stringify({ [action]: true }), { status: 403 }))
      expect((await preset.reqAdd()).id).toBe(1)
      expect(challenge).toHaveBeenCalledTimes(1)
      expect(fetch).toHaveBeenCalledTimes(2)
      const headers = new Headers(fetch.mock.calls[1][1]?.headers)
      if (action === 'need_login') {
        expect(headers.get('Authorization')).toBe('Bearer verified-token')
        expect(headers.has('X-Artalk-Guest')).toBe(false)
      } else expect(headers.has('X-Artalk-Guest')).toBe(true)
    },
  )
})

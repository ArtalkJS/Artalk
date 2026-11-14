import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { bindApi, comment, deferred, list, mockApi, response, version } from './lib/api'
import Artalk from '@/artalk'
import { GlobalPlugins, mount } from '@/mount'
import $t from '@/i18n'
import type { ArtalkPlugin, ConfigPartial } from '@/types'

const instances: Artalk[] = []
const innerText = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'innerText')

beforeEach(() => {
  localStorage.clear()
  // jsdom lacks innerText, which is used for button labels and HTML encoding.
  Object.defineProperty(HTMLElement.prototype, 'innerText', {
    configurable: true,
    get() {
      return this.textContent
    },
    set(value) {
      this.textContent = value
    },
  })
  vi.spyOn(window, 'scroll').mockImplementation(() => {})
  vi.stubGlobal(
    'fetch',
    vi.fn(() => {
      throw new Error('Unexpected HTTP request')
    }),
  )
})

afterEach(() => {
  instances.splice(0).forEach((instance) => instance.destroy())
  document.body.replaceChildren()
  if (innerText) Object.defineProperty(HTMLElement.prototype, 'innerText', innerText)
  else delete (HTMLElement.prototype as any).innerText
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

async function initialize(api = mockApi(), options: ConfigPartial = {}) {
  const el = document.createElement('div')
  document.body.append(el)
  const conf = {
    el,
    server: 'https://example.com',
    site: 'Test',
    pageKey: '/',
    locale: 'en',
    emoticons: false as const,
    ...options,
  }
  const artalk = new Artalk(conf)
  instances.push(artalk)
  bindApi(artalk.ctx, api)
  artalk.on('created', () =>
    artalk.ctx.getUser().update({ name: 'Visitor', email: 'visitor@example.com' }),
  )
  await mount(conf, artalk.ctx)
  await vi.waitFor(() => expect(artalk.ctx.getData().getLoading()).toBe(false))
  return artalk
}

describe('existing API and plugin compatibility', () => {
  it('uses the existing API for built-in features and public plugins', async () => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
    const plugin: ArtalkPlugin = (ctx) => {
      ctx.on('mounted', () => {
        void ctx.getApi().comments.getComment(42)
      })
    }
    Artalk.use(plugin)
    try {
      const artalk = await initialize(api, { listSort: true })
      expect(api.conf.conf).toHaveBeenCalledTimes(1)
      expect(api.comments.getComments).toHaveBeenCalledTimes(1)
      expect(api.comments.getComment).toHaveBeenCalledWith(42)
      expect(api.notifies.getNotifies).toHaveBeenCalledWith({
        name: 'Visitor',
        email: 'visitor@example.com',
      })
      expect(artalk.getEl().textContent).toContain('Comment 1')
      expect(artalk.ctx.inject('editor').getUI().$header.style.display).toBe('')
      expect(artalk.getEl().querySelector('.atk-dropdown')).not.toBeNull()
      expect(fetch).not.toHaveBeenCalled()
    } finally {
      GlobalPlugins.delete(plugin)
    }
  })

  it('ignores a moderation response after the comment is removed', async () => {
    const api = mockApi()
    const pending = deferred<Awaited<ReturnType<typeof api.comments.updateComment>>>()
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
    vi.mocked(api.comments.updateComment).mockReturnValue(pending.promise)
    const artalk = await initialize(api)
    artalk.ctx.getUser().update({ is_admin: true, token: 'admin-token' })
    const node = artalk.ctx.getCommentNodes()[0]
    const pin = [...node.getEl().querySelectorAll<HTMLElement>('.atk-common-action-btn')].find(
      (button) => button.innerText === $t('pin'),
    )!
    pin.click()
    expect(api.comments.updateComment).toHaveBeenCalledTimes(1)
    artalk.update({ pageKey: '/next' })
    expect(node.isRemoved()).toBe(true)
    const update = vi.spyOn(node, 'setData')
    pending.resolve(response({ ...comment(1), is_pinned: true }))
    await pending.promise
    await Promise.resolve()
    expect(update).not.toHaveBeenCalled()
    expect(artalk.ctx.getCommentNodes()).toHaveLength(0)
  })

  it('keeps separate API calls and comments for each widget', async () => {
    const firstApi = mockApi()
    const secondApi = mockApi()
    vi.mocked(firstApi.comments.getComments).mockResolvedValue(
      response(list([comment(1, 'First widget')])),
    )
    vi.mocked(secondApi.comments.getComments).mockResolvedValue(
      response(list([comment(2, 'Second widget')])),
    )
    const first = await initialize(firstApi)
    const second = await initialize(secondApi)
    expect(first.getEl().textContent).toContain('First widget')
    expect(first.getEl().textContent).not.toContain('Second widget')
    expect(second.getEl().textContent).toContain('Second widget')
    first.reload()
    await vi.waitFor(() => expect(firstApi.comments.getComments).toHaveBeenCalledTimes(2))
    expect(secondApi.comments.getComments).toHaveBeenCalledTimes(1)
  })

  it('loads initial notifications and discards them after the destination changes', async () => {
    const api = mockApi()
    const initial = deferred<Awaited<ReturnType<typeof api.notifies.getNotifies>>>()
    const stale = deferred<Awaited<ReturnType<typeof api.notifies.getNotifies>>>()
    vi.mocked(api.notifies.getNotifies)
      .mockReturnValueOnce(initial.promise)
      .mockReturnValueOnce(stale.promise)
    const artalk = await initialize(api)
    const update = vi.spyOn(artalk.ctx.getData(), 'updateNotifies')
    initial.resolve(response({ notifies: [], count: 0 }))
    await vi.waitFor(() => expect(update).toHaveBeenCalledTimes(1))
    update.mockClear()
    artalk.reload()
    artalk.update({ pageKey: '/next' })
    expect(update).toHaveBeenCalledWith([])
    update.mockClear()
    stale.resolve(response({ notifies: [], count: 0 }))
    await stale.promise
    await Promise.resolve()
    expect(update).not.toHaveBeenCalled()
  })

  it('clears old notifications without a request when profile fields are empty', async () => {
    const api = mockApi()
    const artalk = await initialize(api)
    artalk.ctx.getData().updateNotifies([
      {
        id: 1,
        comment_id: 1,
        user_id: 1,
        is_read: false,
        is_emailed: false,
        read_link: '/',
      },
    ])
    artalk.ctx.getUser().update({ name: '', email: '' })
    artalk.reload()
    expect(artalk.ctx.getData().getNotifies()).toEqual([])
    expect(api.notifies.getNotifies).toHaveBeenCalledTimes(1)
  })

  it.each(['account', 'destroy', 'website'] as const)(
    'checks the notification scope after a %s change',
    async (change) => {
      const api = mockApi()
      const pending = deferred<Awaited<ReturnType<typeof api.notifies.getNotifies>>>()
      vi.mocked(api.notifies.getNotifies).mockReturnValue(pending.promise)
      const artalk = await initialize(api)
      const update = vi.spyOn(artalk.ctx.getData(), 'updateNotifies')
      if (change === 'account') artalk.ctx.getUser().update({ token: 'another-viewer' })
      else if (change === 'destroy') artalk.destroy()
      else artalk.ctx.getUser().update({ link: 'https://example.com/about' })
      pending.resolve(response({ notifies: [], count: 0 }))
      await pending.promise
      await Promise.resolve()
      expect(update).toHaveBeenCalledTimes(change === 'website' ? 1 : 0)
    },
  )

  it('appends more comments and replaces them on reload', async () => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockImplementation(async ({ offset = 0 }) =>
      response({ ...list([comment(offset + 1)]), count: 2, roots_count: 2 }),
    )
    const artalk = await initialize(api, { pagination: { pageSize: 1, readMore: true } })
    await artalk.ctx.inject('comments').load({ offset: 1 })
    expect(artalk.ctx.getComments().map((item) => item.id)).toEqual([1, 2])
    expect(artalk.ctx.getCommentNodes()).toHaveLength(2)
    await artalk.ctx.inject('comments').load()
    expect(artalk.ctx.getComments().map((item) => item.id)).toEqual([1])
    expect(artalk.ctx.getCommentNodes()).toHaveLength(1)
  })

  it('keeps the legacy behavior of deprecated useBackendConf', async () => {
    const api = mockApi()
    vi.mocked(api.conf.conf).mockResolvedValue(
      response({ frontend_conf: { placeholder: 'Remote' }, version }),
    )
    const artalk = await initialize(api, {
      useBackendConf: false,
      placeholder: 'Local',
      preferRemoteConf: true,
    })
    expect(api.conf.conf).toHaveBeenCalledTimes(1)
    expect(artalk.getConf().placeholder).toBe('Remote')
  })

  it('retries initialization without duplicate lifecycle events or requests', async () => {
    const el = document.createElement('div')
    document.body.append(el)
    const conf = { el, server: 'https://example.com', emoticons: false as const }
    const artalk = new Artalk(conf)
    instances.push(artalk)
    const api = mockApi()
    bindApi(artalk.ctx, api)
    vi.mocked(api.conf.conf).mockRejectedValueOnce(new Error('Unavailable'))
    const created = vi.fn()
    const mounted = vi.fn()
    artalk.on('created', created)
    artalk.on('mounted', mounted)
    await expect(mount(conf, artalk.ctx)).rejects.toThrow('Unavailable')
    await mount(conf, artalk.ctx)
    await mount(conf, artalk.ctx)
    expect(created).toHaveBeenCalledTimes(1)
    expect(mounted).toHaveBeenCalledTimes(1)
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    expect(el.querySelectorAll('.atk-list')).toHaveLength(1)
  })

  it('does not mount or start requests after a created listener destroys the instance', async () => {
    const el = document.createElement('div')
    document.body.append(el)
    const conf = { el, emoticons: false as const }
    const artalk = new Artalk(conf)
    instances.push(artalk)
    const api = mockApi()
    bindApi(artalk.ctx, api)
    const mounted = vi.fn()
    artalk.on('created', () => artalk.destroy())
    artalk.on('mounted', mounted)
    await mount(conf, artalk.ctx)
    expect(mounted).not.toHaveBeenCalled()
    expect(api.comments.getComments).not.toHaveBeenCalled()
    expect(el.childElementCount).toBe(0)
  })
})

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

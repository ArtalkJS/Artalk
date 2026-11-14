import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { bindApi, mockApi, comment, deferred, list, response, version } from './lib/api'
import Artalk from '@/artalk'
import { GlobalPlugins, mount } from '@/mount'
import { defineLocaleExternal } from '@/i18n/external'
import $t, { setLocale } from '@/i18n'
import en from '@/i18n/en'
import zhCN from '@/i18n/zh-CN'
import Preview from '@/plugins/editor/preview'
import type { ArtalkPlugin, ConfigPartial } from '@/types'

const instances: Artalk[] = []
const innerText = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'innerText')

beforeEach(() => {
  localStorage.clear()
  setLocale('en')
  // jsdom lacks innerText. Match its plain-text behavior for these UI assertions.
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
  instances.splice(0).forEach((widget) => widget.destroy())
  setLocale('en')
  delete window['ArtalkI18n']?.test
  if (innerText) Object.defineProperty(HTMLElement.prototype, 'innerText', innerText)
  else delete (HTMLElement.prototype as any).innerText
  document.body.replaceChildren()
  vi.useRealTimers()
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

async function initialize(options: ConfigPartial = {}, source = mockApi()) {
  const el = document.createElement('div')
  document.body.append(el)
  const conf = {
    el,
    server: 'https://example.com',
    pageKey: '/',
    site: 'Test',
    emoticons: false as const,
    ...options,
  }
  const widget = new Artalk(conf)
  instances.push(widget)
  bindApi(widget.ctx, source)
  widget.on('created', () =>
    widget.ctx.getUser().update({ name: 'Visitor', email: 'visitor@example.com' }),
  )
  await mount(conf, widget.ctx)
  await vi.waitFor(() => expect(widget.ctx.getData().getLoading()).toBe(false))
  return widget
}

describe('shared language and instance Markdown', () => {
  it('applies the remote locale before initializing plugins', async () => {
    const translated: string[] = []
    const plugin: ArtalkPlugin = (ctx) => {
      translated.push($t('send'), ctx.$t('send'))
    }
    GlobalPlugins.add(plugin)
    try {
      const source = mockApi()
      vi.mocked(source.conf.conf).mockResolvedValue(
        response({ frontend_conf: { locale: 'zh-CN' }, version }),
      )
      await initialize({ locale: 'en', preferRemoteConf: true }, source)
      expect(translated).toEqual(['发送', '发送'])
    } finally {
      GlobalPlugins.delete(plugin)
    }
  })

  it('uses the local locale when the remote configuration fails to load', async () => {
    const source = mockApi()
    vi.mocked(source.conf.conf).mockRejectedValue(new Error('Offline'))
    await expect(initialize({ locale: 'zh-CN' }, source)).rejects.toThrow('Offline')
    expect(instances[0].getEl().textContent).toContain(zhCN.listLoadFailMsg)
    expect($t('send')).toBe('发送')
  })
})

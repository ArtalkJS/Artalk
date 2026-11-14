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

  it('scopes drafts to the backend, site, and page, and imports the legacy draft once', async () => {
    localStorage.setItem('ArtalkContent', 'Legacy draft')
    const first = await initialize()
    const second = await initialize({ pageKey: '/second' })
    const editor = first.ctx.inject('editor')
    expect(editor.getContentRaw()).toBe('Legacy draft')
    expect(second.ctx.inject('editor').getContentRaw()).toBe('')
    editor.setContent('First draft')
    first.update({ locale: 'zh-CN' })
    expect(editor.getContentRaw()).toBe('First draft')
    first.update({ pageKey: '/second' })
    expect(editor.getContentRaw()).toBe('')
    editor.setContent('Second draft')
    first.update({ pageKey: '/' })
    expect(editor.getContentRaw()).toBe('First draft')
    const reloaded = await initialize()
    expect(reloaded.ctx.inject('editor').getContentRaw()).toBe('First draft')
    expect((await initialize({ site: 'Another' })).ctx.inject('editor').getContentRaw()).toBe('')
    expect(
      (await initialize({ server: 'https://other.example' })).ctx.inject('editor').getContentRaw(),
    ).toBe('')
  })

  it('keeps editing usable when the browser blocks persistent storage', async () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    const widget = await initialize()
    const editor = widget.ctx.inject('editor')
    editor.setContent('In-memory draft')
    widget.update({ locale: 'zh-CN' })
    expect(editor.getContentRaw()).toBe('In-memory draft')
    editor.submit()
    await vi.waitFor(() => expect(widget.ctx.getComments()).toHaveLength(1))
  })

  it.each(['edit', 'reply'] as const)(
    'keeps the %s target and draft after a locale update',
    async (state) => {
      const source = mockApi()
      vi.mocked(source.comments.getComments).mockResolvedValue(response(list([comment(1)])))
      const widget = await initialize({}, source)
      const editor = widget.ctx.inject('editor')
      const node = widget.ctx.getCommentNodes()[0]
      if (state === 'edit') editor.setEditComment(node.getData(), node.getEl())
      else editor.setReplyComment(node.getData(), node.getEl())
      editor.setContent('Unsaved draft')
      widget.update({ locale: 'zh-CN' })
      expect(editor.getState()).toBe(state)
      expect(editor.getContentRaw()).toBe('Unsaved draft')
      expect(editor.getEl().isConnected).toBe(true)
      editor.submit()
      await vi.waitFor(() => expect(editor.getState()).toBe('normal'))
      if (state === 'edit') {
        expect(source.comments.updateComment).toHaveBeenCalledWith(
          1,
          expect.objectContaining({ content: 'Unsaved draft' }),
        )
        expect(source.comments.createComment).not.toHaveBeenCalled()
      } else {
        expect(source.comments.createComment).toHaveBeenCalledWith(
          expect.objectContaining({ content: 'Unsaved draft', rid: 1 }),
        )
      }
    },
  )

  it('does not rewrite a nickname while the visitor types a space', async () => {
    const widget = await initialize()
    const input = widget.ctx.inject('editor').getUI().$name
    input.value = 'First '
    input.dispatchEvent(new Event('input'))
    expect(input.value).toBe('First ')
    input.value = 'First Visitor'
    input.dispatchEvent(new Event('input'))
    expect(widget.ctx.getUser().getData().name).toBe('First Visitor')
  })

  it('does not duplicate listeners or submission hooks after config updates', async () => {
    const source = mockApi()
    const create = source.comments.createComment
    const widget = await initialize({}, source)
    const editor = widget.ctx.inject('editor')
    const submitted = vi.fn()
    widget.on('editor-submitted', submitted)
    for (const locale of ['zh-CN', 'en', 'zh-CN', 'en']) widget.update({ locale })
    editor.setContent('once')
    editor.getUI().$submitBtn.click()
    await vi.waitFor(() => expect(submitted).toHaveBeenCalledTimes(1))
    expect(create).toHaveBeenCalledTimes(1)
    widget.destroy()
    editor.getUI().$submitBtn.click()
    widget.update({ locale: 'zh-CN' })
    expect(create).toHaveBeenCalledTimes(1)
  })

  it('cancels emoticon preloads on destroy and leaves a sibling usable', async () => {
    vi.useFakeTimers()
    const first = await initialize({ emoticons: 'https://example.com/pack.json' })
    const second = await initialize({ locale: 'en' })
    first.destroy()
    await vi.advanceTimersByTimeAsync(1500)
    expect(fetch).not.toHaveBeenCalled()
    expect(second.ctx.$t('send')).toBe('Send')
  })

  it('ignores an upload result after its editor plugins are replaced', async () => {
    const upload = deferred<string>()
    const uploader = vi.fn(() => upload.promise)
    const widget = await initialize({ imgUpload: true, imgUploader: uploader })
    const editor = widget.ctx.inject('editor')
    const drop = new Event('drop', { cancelable: true })
    Object.defineProperty(drop, 'dataTransfer', {
      value: { files: [new File(['image'], 'test.png')] },
    })
    editor.getUI().$textarea.dispatchEvent(drop)
    expect(uploader).toHaveBeenCalledTimes(1)
    widget.update({ locale: 'zh-CN' })
    editor.setContent('new draft')
    upload.resolve('https://example.com/upload.png')
    await upload.promise
    await Promise.resolve()
    expect(editor.getContentRaw()).toBe('new draft')
  })
})

it.each(['server', 'site', 'pageKey'] as const)(
  'keeps an edit draft at its original %s when navigating',
  async (field) => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
    const artalk = await initialize({}, api)
    const previous = artalk.getConf()[field]
    const editor = artalk.ctx.inject('editor')
    const node = artalk.ctx.getCommentNodes()[0]
    artalk.ctx.editComment(node.getData(), node.getEl())
    editor.setContent('Unsaved edit')
    artalk.update({ [field]: 'https://other.example/next' })
    expect(editor.getState()).toBe('normal')
    expect(editor.getContentRaw()).toBe('')
    artalk.update({ [field]: previous })
    expect(editor.getState()).toBe('normal')
    expect(editor.getContentRaw()).toBe('Unsaved edit')
  },
)

it('keeps an edit draft when its widget is destroyed and mounted again', async () => {
  const api = mockApi()
  vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
  const artalk = await initialize({}, api)
  const node = artalk.ctx.getCommentNodes()[0]
  artalk.ctx.editComment(node.getData(), node.getEl())
  artalk.ctx.inject('editor').setContent('Unsaved edit')
  artalk.destroy()
  const replacement = await initialize()
  expect(replacement.ctx.inject('editor').getContentRaw()).toBe('Unsaved edit')
})

it('ignores a delayed beforeSubmit continuation after the page changes', async () => {
  let proceed!: () => void
  const api = mockApi()
  const artalk = await initialize(
    {
      beforeSubmit: (_editor, next) => {
        proceed = next
      },
    },
    api,
  )
  artalk.ctx.inject('editor').setContent('Old page')
  artalk.ctx.inject('editor').submit()
  artalk.update({ pageKey: '/new' })
  artalk.ctx.inject('editor').setContent('New page')
  proceed()
  await Promise.resolve()
  expect(api.comments.createComment).not.toHaveBeenCalled()
  expect(artalk.ctx.inject('editor').getContentRaw()).toBe('New page')
})

it.each(['link', 'token'] as const)(
  'keeps an upload when only the same viewer’s %s changes',
  async (field) => {
    const pending = deferred<string>()
    const upload = vi.fn(() => pending.promise)
    const artalk = await initialize({ imgUpload: true, imgUploader: upload })
    const editor = artalk.ctx.inject('editor')
    const drop = new Event('drop', { bubbles: true, cancelable: true })
    Object.defineProperty(drop, 'dataTransfer', {
      value: { files: [new File(['image'], 'image.png')] },
    })
    editor.getUI().$textarea.dispatchEvent(drop)
    await vi.waitFor(() => expect(upload).toHaveBeenCalledTimes(1))
    artalk.ctx.getUser().update({ [field]: 'updated' })
    pending.resolve('https://example.com/image.png')
    await vi.waitFor(() =>
      expect(editor.getContentRaw()).toContain('![](https://example.com/image.png)'),
    )
  },
)

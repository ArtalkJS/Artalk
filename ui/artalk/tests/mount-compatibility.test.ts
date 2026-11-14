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

  it('submits and replies with the original fields and event bridge', async () => {
    const api = mockApi()
    const artalk = await initialize(api)
    const submitted = vi.fn()
    artalk.on('editor-submitted', submitted)
    const editor = artalk.ctx.inject('editor')
    editor.setContent('Hello')
    editor.submit()
    await vi.waitFor(() => expect(artalk.ctx.getComments()).toHaveLength(1))
    expect(api.comments.createComment).toHaveBeenCalledWith(
      expect.objectContaining({
        content: 'Hello',
        rid: 0,
        name: 'Visitor',
        email: 'visitor@example.com',
        page_key: '/',
        site_name: 'Test',
      }),
    )
    expect(submitted).toHaveBeenCalledTimes(1)

    artalk.update({ preview: false })
    const first = artalk.ctx.getCommentNodes()[0]
    artalk.ctx.replyComment(first.getData(), first.getEl())
    editor.setContent('Reply')
    editor.submit()
    await vi.waitFor(() => expect(artalk.ctx.getComments()).toHaveLength(2))
    expect(api.comments.createComment).toHaveBeenLastCalledWith(
      expect.objectContaining({ content: 'Reply', rid: 1, page_key: '/', site_name: 'Test' }),
    )
    expect(submitted).toHaveBeenCalledTimes(2)
  })

  it('retains administrator controls and comment fields when editing', async () => {
    const api = mockApi()
    const original = {
      ...comment(1),
      nick: 'Author',
      email: 'author@example.com',
      link: 'https://author.example',
      is_pending: true,
      is_pinned: true,
    }
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([original])))
    const artalk = await initialize(api)
    const listEl = artalk.ctx.inject('list').getEl()
    // The sidebar moves this list outside the widget's original root.
    const sidebar = document.createElement('section')
    document.body.append(sidebar)
    sidebar.append(listEl)
    artalk.reload()
    await vi.waitFor(() => expect(api.comments.getComments).toHaveBeenCalledTimes(2))
    await vi.waitFor(() => expect(artalk.ctx.getData().getLoading()).toBe(false))
    const controls = () => [...listEl.querySelectorAll('.atk-comment [atk-only-admin-show]')]
    expect(controls()).toHaveLength(5)
    expect(controls().every((el) => el.classList.contains('atk-hide'))).toBe(true)
    artalk.ctx.getUser().update({ is_admin: true, token: 'admin-token' })
    expect(controls().every((el) => !el.classList.contains('atk-hide'))).toBe(true)
    artalk.ctx.editComment(original, artalk.ctx.getCommentNodes()[0].getEl())
    const editor = artalk.ctx.inject('editor')
    editor.setContent('Edited')
    // Existing plugins can edit author fields through the editor UI.
    editor.getUI().$name.value = 'Corrected author'
    editor.submit()
    await vi.waitFor(() => expect(artalk.ctx.getComments()[0].content).toBe('Edited'))
    expect(api.comments.updateComment).toHaveBeenCalledWith(1, {
      ...original,
      content: 'Edited',
      nick: 'Corrected author',
    })
    expect(editor.getUI().$header.style.display).toBe('')
    expect(editor.getUI().$name.value).toBe('Visitor')
    artalk.ctx.getUser().logout()
    expect(controls().every((el) => el.classList.contains('atk-hide'))).toBe(true)
  })

  it('preserves actions added by plugins when the viewer changes', async () => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockResolvedValue(response(list([comment(1)])))
    const artalk = await initialize(api)
    const actions = artalk.getEl().querySelector('.atk-actions')!
    const button = document.createElement('button')
    const clicked = vi.fn()
    button.textContent = 'Plugin action'
    button.addEventListener('click', clicked)
    actions.append(button)
    const count = actions.querySelectorAll('.atk-common-action-btn').length
    artalk.ctx.getUser().update({ token: 'signed-in', is_admin: true })
    expect(actions.contains(button)).toBe(true)
    expect(actions.querySelectorAll('.atk-common-action-btn')).toHaveLength(count)
    button.click()
    expect(clicked).toHaveBeenCalledTimes(1)
  })

  it('keeps moderation fields, votes, and the delete confirmation', async () => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockResolvedValue(
      response(list([{ ...comment(1), is_pending: true, is_pinned: true }])),
    )
    const artalk = await initialize(api)
    artalk.ctx.getUser().update({ is_admin: true, token: 'admin-token' })
    const node = artalk.ctx.getCommentNodes()[0]
    const action = (label: string) =>
      [...node.getEl().querySelectorAll<HTMLElement>('.atk-common-action-btn')].find(
        (button) => button.innerText === label,
      )!
    const original = { ...node.getData() }
    action($t('pending')).click()
    await vi.waitFor(() => expect(node.getData().is_pending).toBe(false))
    expect(api.comments.updateComment).toHaveBeenCalledWith(1, { ...original, is_pending: false })
    expect(node.getData().is_pinned).toBe(true)

    node.getRender().voteBtnUp!.$el.click()
    await vi.waitFor(() => expect(node.getData().vote_up).toBe(1))
    expect(api.votes.createVote).toHaveBeenCalledWith('comment', 1, 'up', {
      name: 'Visitor',
      email: 'visitor@example.com',
    })

    const remove = action($t('delete'))
    remove.click()
    expect(api.comments.deleteComment).not.toHaveBeenCalled()
    remove.click()
    await vi.waitFor(() => expect(artalk.ctx.getCommentNodes()).toHaveLength(0))
    expect(api.comments.deleteComment).toHaveBeenCalledExactlyOnceWith(1)
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

  it('uses a custom uploader and returns to the API when it is removed', async () => {
    const api = mockApi()
    const artalk = await initialize(api, { imgUpload: true })
    const editor = artalk.ctx.inject('editor')
    const uploader = vi.fn(async () => 'https://example.com/image.png')
    const file = new File(['image'], 'image.png', { type: 'image/png' })
    const drop = () => {
      const event = new Event('drop', { bubbles: true, cancelable: true })
      Object.defineProperty(event, 'dataTransfer', { value: { files: [file] } })
      editor.getUI().$textarea.dispatchEvent(event)
    }
    artalk.update({ imgUploader: uploader })
    drop()
    await vi.waitFor(() =>
      expect(editor.getContentRaw()).toContain('![](https://example.com/image.png)'),
    )
    expect(uploader).toHaveBeenCalledWith(file)
    expect(api.upload.upload).not.toHaveBeenCalled()
    artalk.update({ imgUploader: undefined })
    drop()
    await vi.waitFor(() =>
      expect(editor.getContentRaw()).toContain('![](https://example.com/upload.png)'),
    )
    expect(api.upload.upload).toHaveBeenCalledWith({ file })
    expect(uploader).toHaveBeenCalledTimes(1)
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

  it('keeps sorting through the original list parameter hook', async () => {
    const api = mockApi()
    const artalk = await initialize(api, { listSort: true })
    artalk.update({ locale: 'zh-CN' })
    expect(artalk.getEl().querySelectorAll('.atk-dropdown')).toHaveLength(1)
    artalk.getEl().querySelectorAll<HTMLElement>('.atk-dropdown-item span')[1].click()
    await vi.waitFor(() => expect(api.comments.getComments).toHaveBeenCalledTimes(2))
    expect(api.comments.getComments).toHaveBeenLastCalledWith(
      expect.objectContaining({ sort_by: 'vote' }),
      expect.anything(),
    )
    const params = {}
    artalk.getConf().listFetchParamsModifier!(params)
    expect(params).toEqual({ sort_by: 'vote' })
  })

  it('replaces numbered pages, bounds page input, and cancels pending input on destroy', async () => {
    const api = mockApi()
    vi.mocked(api.comments.getComments).mockImplementation(async ({ offset = 0 }) =>
      response({ ...list([comment(offset + 1)]), count: 3, roots_count: 3 }),
    )
    const artalk = await initialize(api, { pagination: { pageSize: 1, readMore: false } })
    const root = artalk.getEl()
    const input = root.querySelector<HTMLInputElement>('.atk-pagination input')!
    root.querySelector<HTMLElement>('.atk-btn-prev')!.click()
    expect(api.comments.getComments).toHaveBeenCalledTimes(1)
    root.querySelector<HTMLElement>('.atk-btn-next')!.click()
    await vi.waitFor(() => expect(artalk.ctx.getComments().map((item) => item.id)).toEqual([2]))
    expect(input.value).toBe('2')

    input.value = '99'
    input.dispatchEvent(new KeyboardEvent('keydown', { keyCode: 13 }))
    await vi.waitFor(() => expect(artalk.ctx.getComments().map((item) => item.id)).toEqual([3]))
    expect(input.value).toBe('3')
    root.querySelector<HTMLElement>('.atk-btn-next')!.click()
    expect(api.comments.getComments).toHaveBeenCalledTimes(3)

    vi.useFakeTimers()
    try {
      input.value = '1'
      input.dispatchEvent(new Event('input'))
      artalk.destroy()
      await vi.advanceTimersByTimeAsync(1000)
      expect(api.comments.getComments).toHaveBeenCalledTimes(3)
    } finally {
      vi.useRealTimers()
    }
  })

  it('inserts native and OwO emoticons without changing shared pack data', async () => {
    const packs = [
      { name: 'Emoji', type: 'emoji', items: [{ val: '🙂' }] },
      {
        Pictures: {
          type: 'image',
          container: [{ icon: '<img src="https://example.com/smile.png">', text: 'smile' }],
        },
      },
    ]
    const original = JSON.stringify(packs)
    const artalk = await initialize(mockApi(), { emoticons: packs })
    const root = artalk.getEl()
    root.querySelector<HTMLElement>('[aria-label="Emoji"]')!.click()
    await vi.waitFor(() =>
      expect(
        [...root.querySelectorAll('.atk-grp')].map((group) => group.getAttribute('data-grp-name')),
      ).toEqual(['Emoji', 'Pictures']),
    )
    root.querySelector<HTMLElement>('.atk-grp[data-index="0"] .atk-item')!.click()
    const editor = artalk.ctx.inject('editor')
    expect(editor.getContentRaw()).toBe('🙂')
    root.querySelector<HTMLElement>('.atk-grp-switcher [data-index="1"]')!.click()
    expect(root.querySelector<HTMLElement>('.atk-grp[data-index="0"]')!.style.display).toBe('none')
    root.querySelector<HTMLElement>('.atk-grp[data-index="1"] .atk-item')!.click()
    expect(editor.getContentRaw()).toBe('🙂:[smile]')
    root.querySelector<HTMLElement>('[aria-label="Preview"]')!.click()
    await vi.waitFor(() =>
      expect(root.querySelector('.atk-editor-plug-preview img')?.getAttribute('src')).toBe(
        'https://example.com/smile.png',
      ),
    )
    artalk.update({ locale: 'zh-CN' })
    expect(editor.getContentRaw()).toBe('🙂:[smile]')
    expect(JSON.stringify(packs)).toBe(original)
  })

  it('does not submit twice while the first write is pending', async () => {
    const api = mockApi()
    const pending = deferred<Awaited<ReturnType<typeof api.comments.createComment>>>()
    vi.mocked(api.comments.createComment).mockReturnValue(pending.promise)
    const artalk = await initialize(api)
    const editor = artalk.ctx.inject('editor')
    editor.setContent('Once')
    editor.submit()
    editor.submit()
    await vi.waitFor(() => expect(api.comments.createComment).toHaveBeenCalledTimes(1))
    pending.resolve(response(comment(1, 'Once')))
    await vi.waitFor(() => expect(artalk.ctx.getComments()).toHaveLength(1))
  })

  it.each(['page', 'account', 'destroy'] as const)(
    'does not apply a completed write after %s changes',
    async (change) => {
      const api = mockApi()
      const pending = deferred<Awaited<ReturnType<typeof api.comments.createComment>>>()
      vi.mocked(api.comments.createComment).mockReturnValue(pending.promise)
      const artalk = await initialize(api)
      const editor = artalk.ctx.inject('editor')
      const submitted = vi.fn()
      artalk.on('editor-submitted', submitted)
      editor.setContent('Original draft')
      editor.submit()
      await vi.waitFor(() => expect(api.comments.createComment).toHaveBeenCalledTimes(1))
      if (change === 'page') artalk.update({ pageKey: '/next' })
      else if (change === 'account')
        artalk.ctx.getUser().update({ name: 'Bob', email: 'bob@example.com', token: 'bob-token' })
      else artalk.destroy()
      pending.resolve(response(comment(1, 'Original draft')))
      await pending.promise
      // Allow both reqAdd and Submit to handle the completed response.
      await new Promise((resolve) => setTimeout(resolve, 0))
      expect(artalk.ctx.getComments()).toEqual([])
      if (change === 'page') {
        expect(editor.getContentRaw()).toBe('')
        artalk.update({ pageKey: '/' })
      }
      expect(editor.getContentRaw()).toBe('Original draft')
      expect(submitted).not.toHaveBeenCalled()
    },
  )

  it('completes a submission when login refreshes the same viewer token', async () => {
    const api = mockApi()
    const pending = deferred<Awaited<ReturnType<typeof api.comments.createComment>>>()
    vi.mocked(api.comments.createComment).mockReturnValue(pending.promise)
    const artalk = await initialize(api)
    const editor = artalk.ctx.inject('editor')
    editor.setContent('After login')
    editor.submit()
    await vi.waitFor(() => expect(api.comments.createComment).toHaveBeenCalledTimes(1))
    artalk.ctx.getUser().update({ token: 'authenticated-token', is_admin: true })
    pending.resolve(response(comment(1, 'After login')))
    await vi.waitFor(() => expect(artalk.ctx.getComments()).toHaveLength(1))
    expect(editor.getContentRaw()).toBe('')
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

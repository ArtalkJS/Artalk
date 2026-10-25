import { beforeEach, describe, expect, it, vi } from 'vitest'
import { Fetch } from '@/api/fetch'
import { convertApiOptions } from '@/config'

beforeEach(() => {
  localStorage.clear()
  vi.resetModules()
  vi.restoreAllMocks()
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('{}', { status: 200 })))
})

describe('anonymous comment credential', () => {
  it('is created only for a comment, survives reload and nickname changes without becoming a login', async () => {
    const { UserManager } = await import('@/lib/user')
    const user = new UserManager({})
    user.update({ name: 'Visitor', email: 'visitor@example.com' })
    const opts = convertApiOptions({ server: 'https://comments.example.com' }, user)
    await Fetch(opts, 'https://comments.example.com/api/v2/pages/pv', { method: 'POST' })
    expect(user.getGuestToken('https://comments.example.com')).toBeUndefined()
    await Fetch(opts, 'https://comments.example.com/api/v2/comments', { method: 'POST' })
    const secret = user.getGuestToken('https://comments.example.com')
    expect(secret).toMatch(/^[a-f0-9]{64}$/)
    expect(user.getData().token).toBe('')
    user.update({ name: 'Changed', email: 'changed@example.com' })
    user.logout()
    vi.resetModules() // simulate a new document, with only localStorage surviving
    const { UserManager: Reloaded } = await import('@/lib/user')
    const reloaded = new Reloaded({})
    expect(reloaded.getGuestToken('https://comments.example.com/')).toBe(secret)
    expect(reloaded.getGuestToken('https://other.example.com')).toBeUndefined()
    expect(reloaded.getGuestToken('https://comments.example.com/other')).toBeUndefined()
    expect(reloaded.getData().token).toBe('')
    await Fetch(
      convertApiOptions({ server: 'https://comments.example.com' }, reloaded),
      'https://comments.example.com/api/v2/comments/123',
    )
    const call = vi.mocked(fetch).mock.calls.at(-1)!
    expect(new Headers(call[1]?.headers).get('X-Artalk-Guest')).toBe(secret)
  })

  it('only sends the capability to ownership endpoints on the matching backend', async () => {
    const { UserManager } = await import('@/lib/user')
    const user = new UserManager({})
    const server = 'https://comments.example.com/artalk'
    const token = user.getGuestToken(server, true)
    const opts = convertApiOptions({ server }, user)
    for (const path of ['/comments', '/comments/1', '/notifies', '/notifies/read', '/user']) {
      await Fetch(opts, `${server}/api/v2${path}`)
      expect(
        new Headers(vi.mocked(fetch).mock.calls.at(-1)![1]?.headers).get('X-Artalk-Guest'),
      ).toBe(token)
    }
    for (const url of [
      `${server}/api/v2/pages/pv`,
      `${server}/api/v2/stat`,
      `${server}/api/v2/auth/email/login`,
      'https://evil.example.com/artalk/api/v2/comments',
      'https://comments.example.com/artalk/api/v2-other/comments',
    ]) {
      await Fetch(opts, url)
      expect(
        new Headers(vi.mocked(fetch).mock.calls.at(-1)![1]?.headers).has('X-Artalk-Guest'),
      ).toBe(false)
      expect(String(vi.mocked(fetch).mock.calls.at(-1)![0])).not.toContain(token)
    }
    user.update({ token: 'login-token' })
    await Fetch(opts, `${server}/api/v2/comments`, { method: 'POST' })
    const headers = new Headers(vi.mocked(fetch).mock.calls.at(-1)![1]?.headers)
    expect(headers.get('Authorization')).toBe('Bearer login-token')
    expect(headers.has('X-Artalk-Guest')).toBe(false)
  })

  it('keeps the current page usable when credential storage is blocked', async () => {
    const { getGuestToken, setGuestToken } = await import('@/lib/guest')
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    const token = getGuestToken('https://blocked.example.com', true)
    expect(token).toMatch(/^[a-f0-9]{64}$/)
    expect(getGuestToken('https://blocked.example.com')).toBe(token)
    setGuestToken('https://blocked.example.com', 'invalid')
    expect(getGuestToken('https://blocked.example.com')).toBe(token)
  })

  it('does not claim old comments from a legacy nickname/email record', async () => {
    localStorage.setItem('ArtalkUser', JSON.stringify({ nick: 'Legacy', email: 'old@example.com' }))
    const { UserManager } = await import('@/lib/user')
    const user = new UserManager({})
    expect(user.getData().name).toBe('Legacy')
    expect(user.getGuestToken('https://old.example.com')).toBeUndefined()
    expect(user.getData().token).toBe('')
  })

  it('keeps a handed-off credential when storage is readable but cannot be updated', async () => {
    const server = 'https://readonly.example.com'
    const oldToken = 'ab'.repeat(32)
    const incomingToken = 'cd'.repeat(32)
    localStorage.setItem(`ArtalkGuest:${server}`, oldToken)
    const { getGuestToken, setGuestToken } = await import('@/lib/guest')
    expect(getGuestToken(server)).toBe(oldToken)
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('quota exceeded')
    })
    setGuestToken(server, incomingToken)
    expect(getGuestToken(server)).toBe(incomingToken)
  })
})

describe('sidebar capability handoff', () => {
  it('keeps login credentials out of the iframe query while preserving the legacy profile', async () => {
    const { SidebarLayer } = await import('@/layer/sidebar-layer')
    const { UserManager } = await import('@/lib/user')
    const user = new UserManager({})
    user.update({
      name: 'Admin',
      email: 'admin@example.com',
      token: 'header.payload.signature',
      is_admin: true,
    })
    // Extra fields from integrations must not accidentally become URL parameters.
    Object.assign(user.getData(), { password: 'must-not-be-serialized' })
    const server = 'https://comments.example.com'
    const guest = user.getGuestToken(server, true)!
    const layer = new SidebarLayer({
      getConf: () => ({ get: () => ({ server, site: 'Test', pageKey: '/test', darkMode: false }) }),
      getUser: () => user,
    } as never)
    const frame = (layer as unknown as { createIframe(): HTMLIFrameElement }).createIframe()
    const url = new URL(frame.src)
    expect(JSON.parse(url.searchParams.get('user')!)).toEqual({
      name: 'Admin',
      email: 'admin@example.com',
      link: '',
      is_admin: true,
    })
    expect(url.search).not.toContain('header.payload.signature')
    expect(url.search).not.toContain(guest)
    expect(frame.src).not.toContain('must-not-be-serialized')
    const fragment = new URLSearchParams(url.hash.slice(1))
    expect(fragment.get('token')).toBe('header.payload.signature')
    expect(fragment.get('guest_token')).toBe(guest)
  })

  it('puts the guest secret only in the fragment, not the request query or user profile', async () => {
    const { SidebarLayer } = await import('@/layer/sidebar-layer')
    const { UserManager } = await import('@/lib/user')
    const user = new UserManager({})
    const server = 'https://comments.example.com'
    const secret = user.getGuestToken(server, true)!
    const layer = new SidebarLayer({
      getConf: () => ({ get: () => ({ server, site: 'Test', pageKey: '/test', darkMode: false }) }),
      getUser: () => user,
    } as never)
    const frame = (layer as unknown as { createIframe(): HTMLIFrameElement }).createIframe()
    const url = new URL(frame.src)
    expect(url.hash).toBe(`#guest_token=${secret}`)
    expect(url.search).not.toContain(secret)
    expect(JSON.stringify(user.getData())).not.toContain(secret)
  })
})

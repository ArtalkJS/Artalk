import { beforeEach, describe, expect, it, vi } from 'vitest'
import { UserManager } from '@/lib/user'
import Context from '@/context'
import { ConfigService } from '@/services/config'
import { EventsService } from '@/services/events'
import { UserService } from '@/services/user'
import { ApiService } from '@/services/api'

beforeEach(() => {
  vi.restoreAllMocks()
  localStorage.clear()
})

function user(server: string) {
  return new UserManager({ getServer: () => server })
}

describe('backend-scoped user state', () => {
  it('isolates live data and reloads credentials only for the same normalized backend', () => {
    const first = user('https://a.example/artalk/')
    const second = user('https://b.example/artalk')
    const sibling = user('https://a.example/artalk')
    first.update({
      name: 'Alice',
      email: 'alice@example.com',
      token: 'alice-token',
      is_admin: true,
    })
    expect(second.getData().token).toBe('')
    expect(sibling.getData().token).toBe('')
    expect(user('https://a.example:443/artalk').getData()).toEqual(first.getData())
    expect(user('https://a.example/other').getData().token).toBe('')
    expect(user('http://a.example/artalk').getData().token).toBe('')
    second.logout()
    expect(first.getData().token).toBe('alice-token')
    expect(user('https://a.example/artalk').getData().token).toBe('alice-token')
  })

  it('keeps legacy profile aliases, but does not guess a backend for a legacy token', () => {
    localStorage.setItem(
      'ArtalkUser',
      JSON.stringify({
        nick: 'Legacy',
        email: 'old@example.com',
        link: 'https://example.com',
        token: 'old-token',
        isAdmin: true,
      }),
    )
    expect(user('https://a.example').getData()).toEqual({
      name: 'Legacy',
      email: 'old@example.com',
      link: 'https://example.com',
      token: '',
      is_admin: false,
    })
    expect(user('https://b.example').getData().token).toBe('')
    // Standalone integrations without a backend retain the old constructor behavior.
    expect(new UserManager({}).getData().token).toBe('old-token')
  })

  it('supports backend changes and blocked storage without sharing live state', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('blocked')
    })
    let server = 'https://a.example'
    const changed = vi.fn()
    const instance = new UserManager({ getServer: () => server, onUserChanged: changed })
    instance.update({ name: 'Alice', token: 'alice-token' })
    server = 'https://b.example'
    expect(instance.getData().token).toBe('')
    instance.update({ name: 'Bob', token: 'bob-token' })
    server = 'https://a.example/'
    expect(instance.getData().token).toBe('alice-token')
    instance.dispose()
    instance.update({ token: 'late-token' })
    expect(instance.getData().token).toBe('alice-token')
    expect(changed).toHaveBeenCalledTimes(4)
  })

  it('does not send a new backend token through a previously obtained API client', async () => {
    const ctx = new Context(document.createElement('div'))
    EventsService(ctx)
    ConfigService(ctx)
    ctx.updateConf({ server: 'https://a.example' })
    UserService(ctx)
    ApiService(ctx)
    ctx.getUser().update({ token: 'alice-token' })
    const oldApi = ctx.getApi()
    ctx.updateConf({ server: 'https://b.example' })
    ctx.getUser().update({ token: 'bob-token' })
    const fetch = vi.spyOn(globalThis, 'fetch').mockImplementation(async () => new Response('{}'))
    await expect(oldApi.user.getUserStatus()).rejects.toMatchObject({ name: 'AbortError' })
    expect(fetch).not.toHaveBeenCalled()
    await ctx.getApi().user.getUserStatus()
    expect(String(fetch.mock.calls[0][0])).toContain('https://b.example/')
    expect(new Headers(fetch.mock.calls[0][1]?.headers).get('Authorization')).toBe(
      'Bearer bob-token',
    )
    const currentApi = ctx.getApi()
    ctx.destroy()
    await expect(currentApi.user.getUserStatus()).rejects.toMatchObject({ name: 'AbortError' })
    expect(fetch).toHaveBeenCalledTimes(1)
  })

  it.each(['server', 'account', 'destroy'] as const)(
    'does not launch a login challenge after the request becomes stale on %s changes',
    async (change) => {
      const ctx = new Context(document.createElement('div'))
      EventsService(ctx)
      ConfigService(ctx)
      ctx.updateConf({ server: 'https://a.example' })
      UserService(ctx)
      ApiService(ctx)
      ctx.getUser().update({ token: 'alice-token' })
      let resolve!: (response: Response) => void
      vi.spyOn(globalThis, 'fetch').mockImplementation(
        () =>
          new Promise((done) => {
            resolve = done
          }),
      )
      const challenge = vi.fn()
      // eslint-disable-next-line artalk/artalk-plugin -- Test composition root.
      ctx.inject('apiHandlers').add('need_login', challenge)
      const request = ctx.getApi().user.getUserStatus()
      if (change === 'server') ctx.updateConf({ server: 'https://b.example' })
      else if (change === 'account') ctx.getUser().update({ token: 'bob-token' })
      else ctx.destroy()
      resolve(new Response(JSON.stringify({ need_login: true }), { status: 403 }))
      await expect(request).rejects.toMatchObject({ name: 'AbortError' })
      expect(challenge).not.toHaveBeenCalled()
    },
  )
})

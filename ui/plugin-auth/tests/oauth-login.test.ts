import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { startOAuthLogin } from '../lib/oauth-login'
import type { AuthContext } from '../types'

describe('OAuth popup callback', () => {
  const origin = 'https://comments.example.com'
  const ctx = {} as AuthContext
  let popup: Window

  beforeEach(() => {
    vi.useFakeTimers()
    popup = { closed: false, close: vi.fn(), postMessage: vi.fn() } as unknown as Window
    vi.spyOn(window, 'open').mockReturnValue(popup)
  })

  afterEach(() => {
    vi.restoreAllMocks()
    vi.useRealTimers()
  })

  const send = (source: Window, messageOrigin: string, data: unknown) =>
    window.dispatchEvent(new MessageEvent('message', { source, origin: messageOrigin, data }))

  it.each(['wrong origin', 'wrong window'])('ignores a callback from the %s', async (kind) => {
    const login = startOAuthLogin(ctx, `${origin}/api/v2/auth/github`)
    const resolved = vi.fn()
    void login.then(resolved)
    send(
      kind === 'wrong window' ? window : popup,
      kind === 'wrong origin' ? 'https://attacker.example' : origin,
      {
        type: 'ATK_AUTH_CALLBACK',
        payload: 'attacker-token',
      },
    )
    await Promise.resolve()
    expect(resolved).not.toHaveBeenCalled()
    expect(popup.close).not.toHaveBeenCalled()
    send(popup, origin, { type: 'ATK_AUTH_CALLBACK', payload: 'real-token' })
    await expect(login).resolves.toEqual({ token: 'real-token' })
  })

  it('ignores malformed payloads and sends the handshake only to the backend origin', async () => {
    const login = startOAuthLogin(ctx, `${origin}/api/v2/auth/github`)
    const resolved = vi.fn()
    void login.then(resolved)
    for (const data of [
      null,
      {},
      { type: 'ATK_AUTH_CALLBACK', payload: {} },
      { type: 'ATK_AUTH_CALLBACK', payload: '' },
    ]) {
      send(popup, origin, data)
    }
    await Promise.resolve()
    expect(resolved).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1000)
    expect(popup.postMessage).toHaveBeenCalledWith({ type: 'ATK_LOGIN' }, origin)
    send(popup, origin, { type: 'ATK_AUTH_CALLBACK', payload: 'real-token' })
    await expect(login).resolves.toEqual({ token: 'real-token' })
    expect(vi.getTimerCount()).toBe(0)
  })

  it('rejects when the popup closes', async () => {
    const login = startOAuthLogin(ctx, `${origin}/api/v2/auth/github`)
    const rejected = expect(login).rejects.toThrow('Login canceled')
    Object.defineProperty(popup, 'closed', { value: true })
    vi.advanceTimersByTime(1000)
    await rejected
    expect(vi.getTimerCount()).toBe(0)
  })
  it('pins a separately configured callback origin and rejects the login entry origin', async () => {
    const callback = 'https://auth.example.com'
    const login = startOAuthLogin(ctx, `${origin}/api/v2/auth/github`, callback)
    const resolved = vi.fn()
    void login.then(resolved)
    send(popup, origin, { type: 'ATK_AUTH_CALLBACK', payload: 'wrong-origin-token' })
    await Promise.resolve()
    expect(resolved).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1000)
    expect(popup.postMessage).toHaveBeenCalledWith({ type: 'ATK_LOGIN' }, callback)
    send(popup, callback, { type: 'ATK_AUTH_CALLBACK', payload: 'real-token' })
    await expect(login).resolves.toEqual({ token: 'real-token' })
  })

  it('rejects a blocked popup without starting a timer', async () => {
    vi.mocked(window.open).mockReturnValue(null)
    await expect(startOAuthLogin(ctx, `${origin}/api/v2/auth/github`)).rejects.toThrow(
      'Login popup blocked',
    )
    expect(vi.getTimerCount()).toBe(0)
  })
})

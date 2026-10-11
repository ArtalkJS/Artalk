// @vitest-environment jsdom
import { afterEach, expect, it, vi } from 'vitest'
import type { Context } from 'artalk'
import { restoreGuestCredential } from './guest-handoff'

afterEach(() => vi.restoreAllMocks())

it.each(['direct', 'untrusted', 'missing', 'malformed', 'error', 'mismatched', 'trusted'])(
  'imports a guest credential only from a trusted iframe parent: %s',
  async (kind) => {
    vi.spyOn(window, 'parent', 'get').mockReturnValue(kind === 'direct' ? window : ({} as Window))
    vi.spyOn(document, 'referrer', 'get').mockReturnValue(
      kind === 'missing' ? '' : kind === 'malformed' ? 'about:blank' : 'https://blog.example/post',
    )
    const getDomain = vi.fn().mockResolvedValue({
      data: {
        is_trusted: kind !== 'untrusted',
        origin: kind === 'mismatched' ? 'https://other.example' : 'https://blog.example',
      },
    })
    if (kind === 'error') getDomain.mockRejectedValue(new Error('offline'))
    const setGuestToken = vi.fn()
    const ctx = {
      getApi: () => ({ conf: { getDomain } }),
      getUser: () => ({ setGuestToken }),
      getConf: () => ({ server: 'https://comments.example' }),
    } as unknown as Context
    await restoreGuestCredential(ctx, 'ab'.repeat(32))
    if (kind === 'trusted') {
      expect(getDomain).toHaveBeenCalledWith({ url: 'https://blog.example' })
      expect(setGuestToken).toHaveBeenCalledWith('https://comments.example', 'ab'.repeat(32))
    } else expect(setGuestToken).not.toHaveBeenCalled()
  },
)

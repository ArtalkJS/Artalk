import { expect, it, vi } from 'vitest'
import type { Api } from 'artalk'
import { fetchMethods } from '../lib/methods'

it.each([undefined, 'https://auth.example.com'])(
  'accepts provider responses with an optional callback origin: %s',
  async (origin) => {
    const api = {
      conf: {
        getSocialLoginProviders: vi.fn().mockResolvedValue({
          data: {
            callback_origin: origin,
            providers: [{ name: 'github', label: 'GitHub', icon: '', path: '/api/v2/auth/github' }],
          },
        }),
      },
    } as unknown as Api
    const methods = await fetchMethods(api)
    expect(methods[0].callbackOrigin).toBe(origin)
    expect(methods[0].link).toBe('/api/v2/auth/github')
  },
)

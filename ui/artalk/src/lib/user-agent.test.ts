import { afterEach, describe, expect, it, vi } from 'vitest'
import { getCorrectUserAgent } from './utils'
import SubmitAddPreset from '@/plugins/editor/submit-add'
import type PlugKit from '@/plugins/editor/_kit'

const windowsUA = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'
const macUA = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)'

function mockNavigator(userAgentData?: unknown, userAgent = windowsUA) {
  vi.stubGlobal('navigator', { userAgent, userAgentData })
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('User-Agent client hints', () => {
  it.each([undefined, null, {}, { getHighEntropyValues: true }])(
    'uses the original User-Agent when the client hints API is unavailable: %j',
    async (userAgentData) => {
      mockNavigator(userAgentData)
      await expect(getCorrectUserAgent()).resolves.toBe(windowsUA)
    },
  )

  it.each([
    undefined,
    null,
    {},
    { platformVersion: null },
    { platformVersion: 13 },
    { platformVersion: '' },
    { platformVersion: 'unknown' },
    { platformVersion: 'Infinity' },
    { platformVersion: '13.invalid' },
  ])('uses the original User-Agent for incomplete or invalid hints: %j', async (hints) => {
    mockNavigator({ platform: 'Windows', getHighEntropyValues: vi.fn().mockResolvedValue(hints) })
    await expect(getCorrectUserAgent()).resolves.toBe(windowsUA)
  })

  it('uses the original User-Agent if the browser rejects the hints request', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    mockNavigator({ getHighEntropyValues: vi.fn().mockRejectedValue(new Error('Denied')) })
    await expect(getCorrectUserAgent()).resolves.toBe(windowsUA)
  })

  it.each([
    ['Windows', '13.0.0', windowsUA, windowsUA.replace('10.0', '11.0')],
    ['Windows', '10.0.0', windowsUA, windowsUA],
    ['macOS', '12.3.0', macUA, macUA.replace('10_15_7', '12_3_0')],
    ['macOS', '10.15.7', macUA, macUA],
    ['Linux', '13.0.0', 'Linux browser', 'Linux browser'],
  ])('preserves platform correction for %s %s', async (platform, platformVersion, ua, expected) => {
    const getHighEntropyValues = vi.fn().mockResolvedValue({ platformVersion })
    mockNavigator({ platform, getHighEntropyValues }, ua)
    await expect(getCorrectUserAgent()).resolves.toBe(expected)
    expect(getHighEntropyValues).toHaveBeenCalledWith(['platformVersion'])
  })

  it('still submits a comment when a privacy extension omits platformVersion', async () => {
    mockNavigator({ platform: 'Windows', getHighEntropyValues: vi.fn().mockResolvedValue({}) })
    const comment = { id: 1, content: 'Hello' }
    const createComment = vi.fn().mockResolvedValue({ data: comment })
    const kit = {
      useUser: () => ({ getData: () => ({ name: 'Test', email: 'test@example.com', link: '' }) }),
      useConf: () => ({ pageKey: '/', pageTitle: 'Test page', site: 'Test site' }),
      useEditor: () => ({ getContentFinal: () => 'Hello' }),
      useApi: () => ({ comments: { createComment } }),
    } as unknown as PlugKit

    await expect(new SubmitAddPreset(kit).reqAdd()).resolves.toEqual(comment)
    expect(createComment).toHaveBeenCalledWith(
      expect.objectContaining({ content: 'Hello', ua: windowsUA }),
    )
  })
})

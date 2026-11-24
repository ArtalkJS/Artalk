// @vitest-environment jsdom
import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  applyForcedLocale,
  getClientLocaleURLs,
  hasClientLocale,
  matchLocale,
  resolveForcedLocale,
} from './forced-locale'

const available = ['en', 'fr', 'tr', 'zh-CN', 'zh-TW']
const KEY = 'atk_sidebar_forced_locale'

function createStorage(init: Record<string, string> = {}) {
  const data = new Map(Object.entries(init))
  return {
    data,
    getItem: (k: string) => data.get(k) ?? null,
    setItem: (k: string, v: string) => void data.set(k, v),
    removeItem: (k: string) => void data.delete(k),
  } as unknown as Storage & { data: Map<string, string> }
}

describe('matchLocale', () => {
  it.each([
    ['tr', 'tr'],
    ['tr-TR', 'tr'],
    ['EN-us', 'en'],
    ['zh-tw', 'zh-TW'],
    ['zh_CN', 'zh-CN'],
    ['zh', ''],
    ['de', ''],
    ['auto', ''],
    ['', ''],
    ['../en', ''],
  ])('%j -> %j', (tag, expected) => {
    expect(matchLocale(tag, available)).toBe(expected)
  })
})

describe('resolveForcedLocale', () => {
  it('uses and keeps a known locale from the URL param', () => {
    const storage = createStorage()
    expect(resolveForcedLocale('tr-TR', available, storage)).toBe('tr')
    expect(storage.data.get(KEY)).toBe('tr')
  })

  it.each(['', 'auto', 'unknown'])('forces nothing and clears the kept locale for %j', (param) => {
    const storage = createStorage({ [KEY]: 'tr' })
    expect(resolveForcedLocale(param, available, storage)).toBe('')
    expect(storage.data.has(KEY)).toBe(false)
  })

  it('restores the kept locale after a reload', () => {
    expect(resolveForcedLocale(null, available, createStorage({ [KEY]: 'fr' }))).toBe('fr')
    expect(resolveForcedLocale(null, available, createStorage({ [KEY]: 'xx' }))).toBe('')
    expect(resolveForcedLocale(null, available, createStorage())).toBe('')
  })

  it('works without an accessible storage', () => {
    const blocked = {
      getItem: () => {
        throw new Error('blocked')
      },
      setItem: () => {
        throw new Error('blocked')
      },
      removeItem: () => {
        throw new Error('blocked')
      },
    } as unknown as Storage

    expect(resolveForcedLocale('tr', available, blocked)).toBe('tr')
    expect(resolveForcedLocale(null, available, blocked)).toBe('')
    expect(resolveForcedLocale('tr', available, null)).toBe('tr')
  })
})

// `window.ArtalkI18n` is where the external client locale sets register themselves
const win = window as Window & { ArtalkI18n?: Record<string, unknown> }

describe('getClientLocaleURLs', () => {
  it('returns the URL the backend serves the external client locale set at', () => {
    // same format as the backend uses for its own locale (server/common/conf.go)
    expect(getClientLocaleURLs('tr')).toEqual(['dist/i18n/tr.js'])
    expect(getClientLocaleURLs('zh-TW')).toEqual(['dist/i18n/zh-TW.js'])
  })

  it.each(['', 'en', 'zh-CN'])('returns nothing for %j', (locale) => {
    expect(getClientLocaleURLs(locale)).toEqual([])
  })
})

describe('hasClientLocale', () => {
  afterEach(() => {
    delete win.ArtalkI18n
  })

  it.each(['en', 'zh-CN'])('is true for the bundled locale %j', (locale) => {
    expect(hasClientLocale(locale)).toBe(true)
  })

  it('is true only after the external locale set is registered', () => {
    expect(hasClientLocale('tr')).toBe(false)

    win.ArtalkI18n = { fr: {} }
    expect(hasClientLocale('tr')).toBe(false)

    win.ArtalkI18n.tr = {}
    expect(hasClientLocale('tr')).toBe(true)
  })
})

describe('applyForcedLocale', () => {
  afterEach(() => {
    delete win.ArtalkI18n
  })

  it('applies the forced locale when its client locale set is loaded', () => {
    win.ArtalkI18n = { tr: {} }
    const ctx = { updateConf: vi.fn() }

    expect(applyForcedLocale(ctx, 'tr')).toBe(true)
    expect(ctx.updateConf).toHaveBeenCalledWith({ locale: 'tr' })
  })

  it('applies a bundled forced locale without loading anything', () => {
    const ctx = { updateConf: vi.fn() }

    expect(applyForcedLocale(ctx, 'zh-CN')).toBe(true)
    expect(ctx.updateConf).toHaveBeenCalledWith({ locale: 'zh-CN' })
  })

  it('keeps the backend locale when the client locale set failed to load', () => {
    const ctx = { updateConf: vi.fn() }

    expect(applyForcedLocale(ctx, 'tr')).toBe(false)
    expect(ctx.updateConf).not.toHaveBeenCalled()
  })

  it('keeps the backend locale when no locale is forced', () => {
    const ctx = { updateConf: vi.fn() }

    expect(applyForcedLocale(ctx, '')).toBe(false)
    expect(ctx.updateConf).not.toHaveBeenCalled()
  })
})

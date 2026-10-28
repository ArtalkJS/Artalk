import { describe, expect, it } from 'vitest'
import { availableLocales } from './i18n'

describe('availableLocales', () => {
  it('lists the bundled and the on-demand locales', () => {
    expect(availableLocales[0]).toBe('en')
    expect(availableLocales).toEqual(expect.arrayContaining(['tr', 'zh-CN', 'zh-TW']))
    expect(availableLocales.every((l) => /^[a-z]+(-[A-Z]+)?$/.test(l))).toBe(true)
  })
})

import { describe, expect, expectTypeOf, it } from 'vitest'
import type { Config, ConfigPartial } from './config'

describe('ConfigPartial', () => {
  it('keeps callback, DOM, and array types intact', () => {
    expectTypeOf<NonNullable<ConfigPartial['beforeSubmit']>>().toEqualTypeOf<
      NonNullable<Config['beforeSubmit']>
    >()
    expectTypeOf<NonNullable<ConfigPartial['dateFormatter']>>().toEqualTypeOf<
      NonNullable<Config['dateFormatter']>
    >()
    expectTypeOf<NonNullable<ConfigPartial['imgUploader']>>().toEqualTypeOf<
      NonNullable<Config['imgUploader']>
    >()
    expectTypeOf<NonNullable<ConfigPartial['el']>>().toEqualTypeOf<NonNullable<Config['el']>>()
    expectTypeOf<NonNullable<ConfigPartial['pluginURLs']>>().toEqualTypeOf<
      NonNullable<Config['pluginURLs']>
    >()
  })

  it('still makes nested object fields optional', () => {
    const conf: ConfigPartial = {
      pagination: {
        pageSize: 20,
      },
    }

    expect(conf.pagination?.pageSize).toBe(20)
  })
})

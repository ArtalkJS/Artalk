import { describe, expectTypeOf, it } from 'vitest'
import type { Context, ContextApi, DataManager, DataManagerApi, Editor, EditorApi } from './index'

describe('deprecated public type aliases', () => {
  it('preserves the v2 public type names', () => {
    expectTypeOf<ContextApi>().toEqualTypeOf<Context>()
    expectTypeOf<DataManagerApi>().toEqualTypeOf<DataManager>()
    expectTypeOf<EditorApi>().toEqualTypeOf<Editor>()
  })
})

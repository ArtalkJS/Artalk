import { describe, expect, it } from 'vitest'
import { mergeDeep } from './merge-deep'

describe('Normal operations', () => {
  it('should merge objects (1 level)', () => {
    expect(mergeDeep({ a: 1 }, { b: 2 })).toEqual({ a: 1, b: 2 })
    expect(mergeDeep({ a: 1 }, { a: 2 })).toEqual({ a: 2 })
  })

  it('should merge objects (2 levels)', () => {
    expect(mergeDeep({ a: { b: 1 } }, { a: { c: 2 } })).toEqual({
      a: { b: 1, c: 2 },
    })
    expect(mergeDeep({ a: { b: 1 } }, { a: { b: 2 } })).toEqual({ a: { b: 2 } })
  })

  it('should merge objects (3 levels)', () => {
    expect(mergeDeep({ a: { b: { c: 1 } } }, { a: { b: { d: 2 } } })).toEqual({
      a: { b: { c: 1, d: 2 } },
    })
    expect(mergeDeep({ a: { b: { c: 1 } } }, { a: { b: { c: 2 } } })).toEqual({
      a: { b: { c: 2 } },
    })
  })
})

describe('Array merge', () => {
  it('should merge arrays (1 level)', () => {
    expect(mergeDeep({ a: [1] }, { a: [2] })).toEqual({ a: [1, 2] })
  })
  it('should merge arrays (2 levels)', () => {
    expect(mergeDeep({ a: { b: [1] } }, { a: { b: [2] } })).toEqual({
      a: { b: [1, 2] },
    })
  })
  it('should merge arrays (3 levels)', () => {
    expect(mergeDeep({ a: { b: { c: [1] } } }, { a: { b: { c: [2] } } })).toEqual({
      a: { b: { c: [1, 2] } },
    })
  })
})

describe('Prevent in-place modify: mergeDeep(a, b)', () => {
  it('should not modify a', () => {
    const a: any = { a: 1, arr: [1, 2, 3] }
    mergeDeep(a, { a: 2, arr: [4, 5, 6] })
    expect(a).toEqual({ a: 1, arr: [1, 2, 3] })
  })

  it('should not modify b', () => {
    const b: any = { a: 1, arr: [1, 2, 3] }
    mergeDeep({ a: 2, arr: [4, 5, 6] }, b)
    expect(b).toEqual({ a: 1, arr: [1, 2, 3] })
  })
})

describe('Merge special types', () => {
  const testItem = (name: string, val: any) => {
    it(name, () => {
      expect(mergeDeep({ a: val }, { b: val })).toEqual({ a: val, b: val })
    })
  }

  const dom = document.createElement('div')
  testItem('should can keep dom, not deep recursion', dom)

  const fn = () => {}
  testItem('should can keep function', fn)

  const date = new Date()
  testItem('should can keep date', date)

  const reg = /abc/
  testItem('should can keep regex', reg)
})

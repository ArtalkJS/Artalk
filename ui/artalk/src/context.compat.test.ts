import { describe, expect, it } from 'vitest'
import Context from './context'
import type { Services } from '@/types'

describe('deprecated Context service properties', () => {
  it('forwards old properties to the current service keys', () => {
    const ctx = new Context(document.createElement('div'))
    const services = {
      data: {},
      user: {},
      editor: {},
      list: {},
      layers: {},
      checkers: {},
      sidebar: {},
      editorPlugs: {},
    } as Pick<
      Services,
      'data' | 'user' | 'editor' | 'list' | 'layers' | 'checkers' | 'sidebar' | 'editorPlugs'
    >

    Object.entries(services).forEach(([key, service]) => {
      ctx.provide(key as keyof Services, () => service as never)
    })

    expect(ctx.data).toBe(services.data)
    expect(ctx.user).toBe(services.user)
    expect(ctx.editor).toBe(services.editor)
    expect(ctx.list).toBe(services.list)
    expect(ctx.layerManager).toBe(services.layers)
    expect(ctx.checkerLauncher).toBe(services.checkers)
    expect(ctx.sidebarLayer).toBe(services.sidebar)
    expect(ctx.editorPlugs).toBe(services.editorPlugs)
  })

  it('rejects direct config and root replacement', () => {
    const ctx = new Context(document.createElement('div'))

    expect(() => Reflect.set(ctx, 'conf', {})).toThrow(
      'Cannot replace ctx.conf directly; call ctx.updateConf() instead',
    )

    expect(() => Reflect.set(ctx, '$root', document.createElement('div'))).toThrow(
      'Cannot replace ctx.$root; create a new Artalk instance instead',
    )
  })
})

import { describe, expect, it } from 'vitest'
import { ArtalkAnnotationPlugin } from '../src/main'

function createCtx() {
  const handlers = new Map<string, Array<() => void>>()
  const ctx: any = {
    on: (name: string, handler: () => void) => {
      const list = handlers.get(name) || []
      list.push(handler)
      handlers.set(name, list)
    },
    watchConf: () => undefined,
    getConf: () => ({
      pageKey: '/doc',
      site: 'Test',
      locale: 'en',
      darkMode: false,
      annotation: { enabled: true, root: 'article', hideInList: true, allowGuest: true },
    }),
    getUser: () => ({
      getData: () => ({ name: 'alice', email: 'a@example.com', link: '' }),
    }),
    getApi: () => ({
      comments: {
        getComments: async () => ({ data: { comments: [], roots_count: 0 } }),
        createComment: async () => ({ data: {} }),
        updateCommentAnchor: async () => ({ data: {} }),
        deleteComment: async () => ({ data: {} }),
      },
    }),
    isDestroyed: () => false,
  }
  return { ctx, handlers }
}

describe('lifecycle', () => {
  it('mounts and destroys the shadow host', async () => {
    const { ctx, handlers } = createCtx()
    ArtalkAnnotationPlugin(ctx, { enabled: true })
    handlers.get('mounted')!.forEach((handler) => handler())
    await new Promise((resolve) => setTimeout(resolve, 0))

    expect(document.querySelector('.atk-annotation-host')).not.toBeNull()

    handlers.get('unmounted')!.forEach((handler) => handler())
    expect(document.querySelector('.atk-annotation-host')).toBeNull()
  })
})

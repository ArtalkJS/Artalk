import { describe, expect, it, vi } from 'vitest'
import { Fetch } from '@/plugins/list/fetch'

function setupFetch(scope?: 'page' | 'user' | 'site', page?: Record<string, unknown>) {
  const handlers = new Map<string, (params: Record<string, unknown>) => void>()
  const response = {
    comments: [],
    count: 0,
    roots_count: 0,
    ...(page ? { page } : {}),
  }
  const data = {
    getLoading: vi.fn(() => false),
    setLoading: vi.fn(),
    setListLastFetch: vi.fn(),
    loadComments: vi.fn(),
    updatePage: vi.fn(),
  }
  const trigger = vi.fn()
  const ctx = {
    inject: vi.fn(() => ({
      get: () => ({
        pagination: { pageSize: 20 },
        flatMode: true,
        pageKey: '/test',
        site: 'ArtalkDocs',
        listFetchParamsModifier: scope
          ? (params: { scope?: 'page' | 'user' | 'site' }) => {
              params.scope = scope
            }
          : undefined,
      }),
    })),
    on: vi.fn((name: string, handler: (params: Record<string, unknown>) => void) => {
      handlers.set(name, handler)
    }),
    getData: () => data,
    getApi: () => ({
      comments: {
        getComments: vi.fn().mockResolvedValue({ data: response }),
      },
      getUserFields: () => ({}),
    }),
    trigger,
  }

  Fetch(ctx as never)
  handlers.get('list-fetch')?.({})

  return { data, response, trigger }
}

describe('list fetch', () => {
  it.each(['site', 'user'] as const)(
    'accepts responses without page data for the %s scope',
    async (scope) => {
      const { data, response, trigger } = setupFetch(scope)

      await vi.waitFor(() => expect(data.setLoading).toHaveBeenLastCalledWith(false))

      expect(data.loadComments).toHaveBeenCalledWith(response.comments)
      expect(data.updatePage).not.toHaveBeenCalled()
      expect(trigger).toHaveBeenCalledWith(
        'list-fetched',
        expect.objectContaining({ data: response }),
      )
    },
  )

  it('updates page data for the default page scope', async () => {
    const page = { id: 1, key: '/test', site_name: 'ArtalkDocs' }
    const { data } = setupFetch(undefined, page)

    await vi.waitFor(() => expect(data.setLoading).toHaveBeenLastCalledWith(false))

    expect(data.updatePage).toHaveBeenCalledWith(page)
  })
})

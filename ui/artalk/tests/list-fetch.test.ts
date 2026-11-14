import { afterEach, describe, expect, it, vi } from 'vitest'
import { setup, list, page, response } from './lib/api'

afterEach(() => vi.restoreAllMocks())

function setupFetch(scope?: 'page' | 'user' | 'site', withPage = true) {
  const result = list()
  if (!withPage) Reflect.deleteProperty(result, 'page')
  const state = setup()
  vi.mocked(state.api.comments.getComments).mockResolvedValue(response(result))
  const query = {
    offset: 0,
    limit: 20,
    flatMode: true,
    paramsModifier: scope
      ? (params: any) => {
          params.scope = scope
        }
      : undefined,
  }
  return { ...state, query }
}

describe('v2 list compatibility', () => {
  it.each(['site', 'user'] as const)(
    'accepts responses without page data for the %s scope',
    async (scope) => {
      const { service, query, api, data } = setupFetch(scope, false)
      await service.load(query)
      expect(data.getListLastFetch()?.data?.page).toBeUndefined()
      expect(api.comments.getComments).toHaveBeenCalledWith(
        expect.objectContaining({ scope, flat_mode: true }),
        expect.anything(),
      )
    },
  )

  it('keeps page data for the default scope', async () => {
    const { service, query, data } = setupFetch()
    await service.load(query)
    expect(data.getPage()).toEqual(page)
  })

  it('reports malformed page responses through the existing error callback', async () => {
    const { service, query, data } = setupFetch(undefined, false)
    const onError = vi.fn()
    await service.load({ ...query, onError })
    expect(onError).toHaveBeenCalledWith(
      expect.objectContaining({ msg: expect.stringContaining('Page data is missing') }),
    )
    expect(data.getLoading()).toBe(false)
  })

  it('keeps custom sorting, filters, and the current viewer fields in v2 requests', async () => {
    const { ctx, service, api, query } = setupFetch()
    ctx.getUser().update({ name: 'Alice', email: 'alice@example.com' })
    await service.load({
      ...query,
      paramsModifier: (params) => {
        params.sort_by = 'vote'
        params.search = 'needle'
        params.name = 'Impersonated'
      },
    })
    expect(api.comments.getComments).toHaveBeenCalledWith(
      expect.objectContaining({
        sort_by: 'vote',
        search: 'needle',
        name: 'Alice',
        email: 'alice@example.com',
      }),
      expect.anything(),
    )
  })
})

// @vitest-environment jsdom
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, ref } from 'vue'
import Pages from '../pages/pages.vue'

const api = vi.hoisted(() => ({
  getPages: vi.fn().mockResolvedValue({ data: { count: 0, pages: [] } }),
  getPageFetchStatus: vi.fn().mockResolvedValue({ data: { is_progress: false } }),
}))
const nav = vi.hoisted(() => ({
  updateTabs: vi.fn(),
  enableSearch: vi.fn(),
  setPageLoading: vi.fn(),
  scrollPageToTop: vi.fn(),
  isPageLoading: false,
}))
vi.mock('../global', () => ({
  artalk: { ctx: { getApi: () => ({ pages: api }) } },
}))
vi.mock('../stores/nav', () => ({
  useNavStore: () => nav,
}))
vi.mock('../stores/user', () => ({
  useUserStore: () => ({ site: ref('Site A') }),
}))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (key: string) => key }) }))
vi.mock('../components/Pagination.vue', () => ({
  default: {
    setup: (_: unknown, { expose }: any) => expose({ reset: vi.fn() }),
    render: () => null,
  },
}))
vi.mock('../components/PageEditor.vue', () => ({ default: { render: () => null } }))

beforeEach(() => {
  api.getPages.mockResolvedValue({ data: { count: 0, pages: [] } })
  api.getPageFetchStatus.mockResolvedValue({ data: { is_progress: false } })
})
afterEach(() => {
  vi.clearAllMocks()
})

describe('page list comment filter', () => {
  it('requests only commented pages without an extra filter control', async () => {
    const container = document.createElement('div')
    document.body.append(container)
    const app = createApp(Pages)
    try {
      app.mount(container)
      await nextTick()
      expect(api.getPages).toHaveBeenLastCalledWith(
        { site_name: 'Site A', offset: 0, limit: 20, search: '', has_comments: true },
        { signal: expect.any(AbortSignal) },
      )
      expect(container.querySelector('input[type="checkbox"]')).toBeNull()
    } finally {
      app.unmount()
      container.remove()
    }
  })

  it('keeps the latest search result when an older request finishes last', async () => {
    let release!: (value: any) => void
    api.getPages.mockReturnValueOnce(new Promise((resolve) => (release = resolve)))
    const container = document.createElement('div')
    const app = createApp(Pages)
    try {
      app.mount(container)
      api.getPages.mockResolvedValueOnce({
        data: { count: 1, pages: [{ id: 2, title: 'Latest result', url: '/latest' }] },
      })
      nav.enableSearch.mock.calls[0][0]('latest')
      await vi.waitFor(() => expect(container.textContent).toContain('Latest result'))
      release({ data: { count: 1, pages: [{ id: 1, title: 'Stale result', url: '/stale' }] } })
      await new Promise((resolve) => setTimeout(resolve, 0))
      expect(container.textContent).toContain('Latest result')
      expect(container.textContent).not.toContain('Stale result')
    } finally {
      app.unmount()
    }
  })
})

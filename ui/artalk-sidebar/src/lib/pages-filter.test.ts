// @vitest-environment jsdom
import { afterEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, ref } from 'vue'
import Pages from '../pages/pages.vue'

const api = vi.hoisted(() => ({
  getPages: vi.fn().mockResolvedValue({ data: { count: 0, pages: [] } }),
  getPageFetchStatus: vi.fn().mockResolvedValue({ data: { is_progress: false } }),
}))
vi.mock('../global', () => ({
  artalk: { ctx: { getApi: () => ({ pages: api }) } },
}))
vi.mock('../stores/nav', () => ({
  useNavStore: () => ({
    updateTabs: vi.fn(),
    enableSearch: vi.fn(),
    setPageLoading: vi.fn(),
    scrollPageToTop: vi.fn(),
    isPageLoading: false,
  }),
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

afterEach(() => vi.clearAllMocks())

describe('page list comment filter', () => {
  it('requests only commented pages without an extra filter control', async () => {
    const container = document.createElement('div')
    document.body.append(container)
    const app = createApp(Pages)
    try {
      app.mount(container)
      await nextTick()
      expect(api.getPages).toHaveBeenLastCalledWith({
        site_name: 'Site A',
        offset: 0,
        limit: 20,
        search: '',
        has_comments: true,
      })
      expect(container.querySelector('input[type="checkbox"]')).toBeNull()
    } finally {
      app.unmount()
      container.remove()
    }
  })
})

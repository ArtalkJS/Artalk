// @vitest-environment jsdom
import { expect, it, vi } from 'vitest'
import { createApp, ref } from 'vue'
import Comments from '../pages/comments.vue'

const widget = vi.hoisted(() => ({
  getConf: () => ({ listFetchParamsModifier: undefined, scrollRelativeTo: undefined }),
  reload: vi.fn(),
  update: vi.fn(),
  ctx: {
    on: vi.fn(),
    off: vi.fn(),
    fetch: vi.fn(),
    updateConf: vi.fn(),
    editorResetState: vi.fn(),
    inject: vi.fn(),
  },
}))
vi.mock('../global', () => ({ artalk: widget }))
vi.mock('../stores/nav', () => ({
  useNavStore: () => ({ curtTab: ref('all'), updateTabs: vi.fn(), enableSearch: vi.fn() }),
}))
vi.mock('../stores/user', () => ({
  useUserStore: () => ({ site: ref('Test'), is_admin: true }),
}))

it('releases the comment listener, editor and callbacks when leaving the route', () => {
  const list = document.createElement('div')
  list.innerHTML = '<div class="atk-list-header"></div><div class="atk-list-footer"></div>'
  widget.ctx.inject.mockReturnValue({ getEl: () => list })
  for (let visit = 0; visit < 2; visit++) {
    const container = document.createElement('div')
    const app = createApp(Comments)
    app.mount(container)
    expect(container.contains(list)).toBe(true)
    const listener = widget.ctx.on.mock.lastCall![1]
    const setOpenURL = vi.fn()
    listener({
      getData: () => ({ page_url: 'https://example.com/post' }),
      getID: () => 42,
      getRender: () => ({ setOpenURL }),
    })
    expect(setOpenURL).toHaveBeenCalledWith('https://example.com/post#atk-comment-42')
    app.unmount()
    expect(widget.ctx.off).toHaveBeenLastCalledWith('comment-rendered', listener)
    expect(widget.ctx.editorResetState).toHaveBeenCalledTimes(visit + 1)
    expect(widget.update).toHaveBeenLastCalledWith({
      listFetchParamsModifier: undefined,
      scrollRelativeTo: undefined,
    })
  }
})

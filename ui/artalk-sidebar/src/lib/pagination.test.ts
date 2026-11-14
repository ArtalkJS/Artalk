// @vitest-environment jsdom
import { afterEach, expect, it, vi } from 'vitest'
import { createApp, nextTick } from 'vue'
import Pagination from '../components/Pagination.vue'

afterEach(() => vi.useRealTimers())

function mountPagination(total: number) {
  const change = vi.fn()
  const container = document.createElement('div')
  const app = createApp(Pagination, { total, pageSize: 20, onChange: change })
  const pagination = app.mount(container) as unknown as { reset: () => void }
  const input = container.querySelector('input')!
  return { change, app, input, pagination }
}

it('rejects fractional page numbers and keeps empty lists on page one', async () => {
  const { change, app, input } = mountPagination(0)
  try {
    for (const value of ['1.5', 'Infinity', '-1', '', 'text']) {
      input.value = value
      input.dispatchEvent(new Event('input'))
      input.dispatchEvent(new KeyboardEvent('keydown', { keyCode: 13 }))
      await nextTick()
      expect(input.value).toBe('1')
      expect(change).not.toHaveBeenCalled()
    }
    input.value = '99'
    input.dispatchEvent(new Event('input'))
    input.dispatchEvent(new KeyboardEvent('keydown', { keyCode: 13 }))
    await nextTick()
    expect(input.value).toBe('1')
    expect(change).toHaveBeenCalledExactlyOnceWith(0)
  } finally {
    app.unmount()
  }
})

it('cancels delayed page input on reset and unmount', async () => {
  vi.useFakeTimers()
  const { change, app, input, pagination } = mountPagination(60)
  input.value = '2'
  input.dispatchEvent(new Event('input'))
  pagination.reset()
  await vi.advanceTimersByTimeAsync(1000)
  expect(change).not.toHaveBeenCalled()
  expect(input.value).toBe('1')
  input.value = '3'
  input.dispatchEvent(new Event('input'))
  app.unmount()
  await vi.advanceTimersByTimeAsync(1000)
  expect(change).not.toHaveBeenCalled()
})

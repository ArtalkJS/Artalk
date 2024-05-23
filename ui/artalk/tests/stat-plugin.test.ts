import { describe, it, expect, vi } from 'vitest'

import { exportedForTesting } from '@/plugins/stat'
import type { CountOptions } from '@/plugins/stat'

const {
  incrementPvCount,
  loadCommentCount,
  loadPvCount,
  retrieveElements,
  getPageKeys,
  updateElementsText,
} = exportedForTesting

// Mocking API
const mockApi = {
  pages: {
    logPv: vi.fn().mockResolvedValue({ data: { pv: 100 } }),
  },
  stats: {
    getStats: vi.fn().mockResolvedValue({
      data: { data: { 'test-page-key': 5 } },
    }),
  },
}

// Mocking CountOptions
const mockOptions: CountOptions = {
  getApi: () => mockApi as any,
  siteName: 'test-site',
  pageKey: 'test-page-key',
  pageTitle: 'Test Page',
  countEl: '.count-element',
  pvEl: '.pv-element',
  pageKeyAttr: 'data-page-key',
  pvAdd: true,
}

describe('PvCountWidget', () => {
  it('should increment PV count and get cache data', async () => {
    const cacheData = await incrementPvCount(mockOptions)
    expect(mockApi.pages.logPv).toHaveBeenCalledWith({
      page_key: 'test-page-key',
      page_title: 'Test Page',
      site_name: 'test-site',
    })
    expect(cacheData).toEqual({ 'test-page-key': 100 })
  })

  it('should load comment count', async () => {
    document.body.innerHTML = '<div class="count-element" data-page-key="test-page-key"></div>'
    await loadCommentCount(mockOptions)
    const el = document.querySelector<HTMLElement>('.count-element')
    expect(el?.innerText).toBe('5')
  })

  it('should load PV count', async () => {
    document.body.innerHTML = '<div class="pv-element" data-page-key="test-page-key"></div>'
    await loadPvCount(mockOptions, { 'test-page-key': 100 })
    const el = document.querySelector<HTMLElement>('.pv-element')
    expect(el?.innerText).toBe('100')
  })

  it('should retrieve elements based on selectors', () => {
    document.body.innerHTML = `
      <div class="test-element"></div>
      <div class="test-element"></div>
    `
    const elements = retrieveElements(['.test-element'])
    expect(elements.size).toBe(2)
  })

  it('should get page keys to be queried', () => {
    document.body.innerHTML = `
      <div class="test-element" data-page-key="key1"></div>
      <div class="test-element" data-page-key="key2"></div>
    `
    const elements = retrieveElements(['.test-element'])
    const pageKeys = getPageKeys(elements, 'data-page-key', undefined, {})
    expect(pageKeys).toEqual(['key1', 'key2'])
  })

  it('should update elements text content with the count data', () => {
    document.body.innerHTML = `
      <div class="test-element" data-page-key="key1"></div>
      <div class="test-element" data-page-key="key2"></div>
      <div class="test-element"></div>
    `
    const elements = retrieveElements(['.test-element'])
    const data = { key1: 10, key2: 20 }
    updateElementsText(elements, data, 'defaultKey')
    const els = document.querySelectorAll<HTMLElement>('.test-element')
    expect(els[0].innerText).toBe('10')
    expect(els[1].innerText).toBe('20')
    expect(els[2].innerText).toBe('0')
  })
})

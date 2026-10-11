import { describe, expect, it } from 'vitest'
import { SidebarLayer, type SidebarLayerOptions } from './sidebar-layer'

function getIframeParams(locale: unknown) {
  const conf = {
    server: 'http://localhost:23366',
    pageKey: '/',
    site: 'Test',
    darkMode: false,
    locale,
  }
  const layer = new SidebarLayer({
    getConf: () => ({ get: () => conf }),
    getUser: () => ({ getData: () => ({}) }),
  } as unknown as SidebarLayerOptions)

  const $iframe: HTMLIFrameElement = (layer as any).createIframe()
  return new URL($iframe.src).searchParams
}

describe('sidebar iframe locale param', () => {
  it('passes the locale of the parent page', () => {
    expect(getIframeParams('tr').get('locale')).toBe('tr')
  })

  it.each([
    ['empty', ''],
    ['locale object', { lang: 'Custom' }],
  ])('omits the param when the locale is %s', (_, locale) => {
    expect(getIframeParams(locale).has('locale')).toBe(false)
  })
})

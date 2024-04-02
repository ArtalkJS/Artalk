import { describe, it, expect, vi, beforeAll } from 'vitest'
import Artalk from '@/artalk'

const InitConf = {
  pageTitle: 'Artalk DEMO',
  pageKey: '/unit_test_page.html?test=1',
  server: 'http://localhost:3000/api',
  site: 'Unit Test Page',
  darkMode: true,
}

const RemoteConf = {
  darkMode: false, // simulate response `false`, but the final should still be `true`, cannot override this
  gravatar: { mirror: 'https://test.avatar.com/unit_test', params: 'test=123' }
}

const ContainerID = 'artalk-container'

beforeAll(() => {
  // mock fetch
  global.fetch = vi.fn().mockImplementation((url: string, init: RequestInit) => Promise.resolve({
    ok: true,
    status: 200,
    json: () => {
      let resp: any = {}

      const map = {
        '/api/v2/conf': {
          frontend_conf: RemoteConf,
          version: {},
        },
        '/api/v2/stat': { data: { '/': 0 } },
        '/api/v2/pages/pv': { pv: 2 },
        '/api/v2/notifies': { notifies: [], count: 0 },
        '/api/v2/comments': {
          comments: [],
          count: 0,
          roots_count: 0,
          page: { id: 4, admin_only: false, key: '/', url: '/', title: 'Artalk DEMO', site_name: 'ArtalkDocs', vote_up: 0, vote_down: 0, pv: 1 },
        }
      }

      const path = new URL(url).pathname

      Object.entries(map).forEach(([k, v]) => {
        if (path.startsWith(k)) resp = v
      })

      return Promise.resolve(resp)
    }
  })) as any
})

describe('Artalk instance', () => {
  it('should be a class', () => {
    expect(Artalk).toBeInstanceOf(Function)
  })

  let artalk: Artalk

  it('create an instance (artalk.init)', () => {
    const el = document.createElement('div')
    el.setAttribute('id', ContainerID)
    document.body.appendChild(el)

    artalk = Artalk.init({
      ...InitConf,
      el,
    })

    expect(artalk).toBeInstanceOf(Artalk)
  })

  let confCopy: any

  it('should have correct config (artalk.getConf, artalk.getEl)', () => {
    const conf = artalk.getConf()
    expect(conf.pageKey).toBe(InitConf.pageKey)
    expect(conf.server).toBe(InitConf.server.replace('/api', ''))
    expect(conf.site).toBe(InitConf.site)
    expect(conf.darkMode).toBe(InitConf.darkMode)

    confCopy = JSON.parse(JSON.stringify(conf))
  })

  it('should can listen to events and the conf-remoter works (artalk.trigger, artalk.on, conf-remoter)', async () => {
    global.devLoadArtalk()

    const fn = vi.fn()

    await new Promise<void>(resolve => {
      artalk.on('mounted', (conf) => {
        resolve()
        fn()
      })
    })

    expect(fn).toBeCalledTimes(1)

    const conf = artalk.getConf()
    expect(conf.darkMode, "the darkMode is unmodifiable, should still false").toBe(true)
    expect(conf.gravatar, "the gravatar should be modified").toEqual(RemoteConf.gravatar)
  }, {
    timeout: 1000
  })

  it('should other config not changed after conf-remoter loaded (conf-remoter, artalk.update)', () => {
    const confNew = JSON.parse(JSON.stringify(artalk.getConf()))
    confCopy.gravatar = confNew.gravatar // exclude
    expect(confCopy).toMatchObject(confNew)
  })

  it('should can update config (artalk.update)', () => {
    const Placeholder = 'Test Placeholder'
    artalk.update({ placeholder: Placeholder })

    const conf = artalk.getConf()
    expect(conf.placeholder).toBe(Placeholder)
    expect(conf.gravatar, "the gravatar which not in update should keep the same").toEqual(RemoteConf.gravatar)
  })

  it('should can getEl after config updated (artalk.getEl)', () => {
    const target = document.getElementById(ContainerID)
    expect(target).not.toBe(null)

    const el = artalk.getEl()
    expect(el).toBe(target)

    const el2 = artalk.getConf().el
    expect(el2).toBe(target)
  })

  it('should can set dark mode (artalk.setDarkMode)', () => {
    const el = artalk.getEl()
    expect(artalk.getConf().darkMode).toBe(true)
    expect(el.classList.contains('atk-dark-mode')).toBe(true)
    artalk.setDarkMode(false)
    expect(artalk.getConf().darkMode).toBe(false)
    expect(el.classList.contains('atk-dark-mode')).toBe(false)
  })

  it('should can reload comments (artalk.reload)', async () => {
    artalk.reload()

    const fn = vi.fn()

    await new Promise(resolve => {
      artalk.on('list-loaded', () => {
        resolve(null)
        fn()
      })
    })

    expect(fn).toBeCalledTimes(1)
  })

  it('should can destroy (artalk.destroy)', () => {
    artalk.destroy()

    // detect if it is cleaned up
    expect(document.querySelector(`#${ContainerID}`)?.firstChild).toBe(null)
    expect(document.querySelector('.atk-layer-wrap')).toBe(null)
  })
})

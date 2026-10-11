import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import Artalk from '@/artalk'

const InitConf = {
  pageTitle: 'Artalk DEMO',
  pageKey: '/unit_test_page.html',
  server: 'http://localhost:3000/api',
  site: 'Unit Test Page',
}

const responses: Record<string, unknown> = {
  '/api/v2/conf': { frontend_conf: { placeholder: 'Remote placeholder' }, version: {} },
  '/api/v2/stat': { data: { '/': 0 } },
  '/api/v2/pages/pv': { pv: 1 },
  '/api/v2/notifies': { notifies: [], count: 0 },
  '/api/v2/comments': {
    comments: [],
    count: 0,
    roots_count: 0,
    page: { id: 1, key: '/', url: '/', title: 'Artalk DEMO', site_name: 'Unit Test Page' },
  },
}

const jsonResponse = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })

/** Hold back conf responses until `release` is called */
function holdConfResponses() {
  const pending: (() => void)[] = []
  let failed = false
  return {
    wait: () => new Promise<void>((resolve) => pending.push(resolve)),
    release: (opts: { fail?: boolean } = {}) => {
      failed = !!opts.fail
      pending.splice(0).forEach((resolve) => resolve())
    },
    get failed() {
      return failed
    },
  }
}

let confGate: ReturnType<typeof holdConfResponses> | null = null

beforeEach(() => {
  global.fetch = vi.fn().mockImplementation(async (url: string) => {
    const path = new URL(url).pathname
    if (path.startsWith('/api/v2/conf') && confGate) {
      const gate = confGate
      await gate.wait()
      if (gate.failed) return jsonResponse({ msg: 'conf request failed' }, 500)
    }

    const key = Object.keys(responses).find((k) => path.startsWith(k))
    return jsonResponse(key ? responses[key] : {})
  })
})

afterEach(() => {
  confGate = null
  document.body.innerHTML = ''
})

function createContainer() {
  const el = document.createElement('div')
  document.body.appendChild(el)
  return el
}

/** Create an instance and start mounting it, like `Artalk.init()` does outside of tests */
function initAndMount(el: HTMLElement) {
  const artalk = Artalk.init({ ...InitConf, el })
  const events: string[] = []
  artalk.on('created', () => events.push('created'))
  artalk.on('mounted', () => events.push('mounted'))
  const mounting: Promise<void> = global.devMountArtalk()
  return { artalk, events, mounting }
}

const countIn = (el: HTMLElement, selector: string) => el.querySelectorAll(selector).length

describe('destroy before mount has finished', () => {
  it('mounts normally when the instance is not destroyed', async () => {
    const el = createContainer()
    const { artalk, events, mounting } = initAndMount(el)

    await mounting

    expect(artalk.ctx.isDestroyed()).toBe(false)
    expect(artalk.getConf().placeholder).toBe('Remote placeholder')
    expect(events).toEqual(['created', 'mounted'])
    expect(countIn(el, '.atk-main-editor')).toBe(1)
    expect(countIn(el, '.atk-list')).toBe(1)
  })

  it('keeps a single editor and list when re-initialized on the same element', async () => {
    confGate = holdConfResponses()
    const el = createContainer()

    // e.g. React StrictMode: init -> destroy -> init on the same element
    const first = initAndMount(el)
    first.artalk.destroy()
    const second = initAndMount(el)

    // Resolve both conf requests; the destroyed instance resumes after its await
    confGate.release()
    await Promise.all([first.mounting, second.mounting])

    expect(countIn(el, '.atk-main-editor')).toBe(1)
    expect(countIn(el, '.atk-list')).toBe(1)
    expect(first.artalk.ctx.isDestroyed()).toBe(true)
    expect(first.events).toEqual([])
    expect(second.events).toEqual(['created', 'mounted'])
    expect(el.querySelector('.atk-main-editor')).toBe(second.artalk.ctx.inject('editor').getEl())
    expect(el.querySelector('.atk-list')).toBe(second.artalk.ctx.inject('list').getEl())
  })

  it('does not apply the remote config to a destroyed instance', async () => {
    confGate = holdConfResponses()
    const el = createContainer()

    const { artalk, mounting } = initAndMount(el)
    artalk.destroy()

    confGate.release()
    await mounting

    expect(artalk.getConf().placeholder).not.toBe('Remote placeholder')
    expect(el.childElementCount).toBe(0)
  })

  it('does not render the mount error when the conf request fails after destroy', async () => {
    confGate = holdConfResponses()
    const el = createContainer()

    const { artalk, events, mounting } = initAndMount(el)
    artalk.destroy()

    confGate.release({ fail: true })
    await expect(mounting).resolves.toBeUndefined()

    expect(events).toEqual([])
    expect(el.childElementCount).toBe(0)
  })

  it('does not clear a reused root element when destroyed twice', async () => {
    const el = createContainer()

    const first = initAndMount(el)
    await first.mounting
    first.artalk.destroy()

    const second = initAndMount(el)
    await second.mounting
    first.artalk.destroy()

    expect(countIn(el, '.atk-main-editor')).toBe(1)
    expect(countIn(el, '.atk-list')).toBe(1)
  })
})

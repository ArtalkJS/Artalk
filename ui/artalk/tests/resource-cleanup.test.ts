import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { deferred } from './lib/api'
import { CheckerLauncher } from '@/components/checker'
import { LayerManager } from '@/layer'
import { UserManager } from '@/lib/user'
import { EventManager } from '@/lib/event-manager'
import Context from '@/context'
import { EventsService } from '@/services/events'
import { ConfigService } from '@/services/config'

const managers: LayerManager[] = []
beforeEach(() => vi.useFakeTimers())
afterEach(() => {
  managers.splice(0).forEach((manager) => manager.destroy())
  document.body.replaceChildren()
  document.body.removeAttribute('style')
  vi.useRealTimers()
  vi.restoreAllMocks()
})

function layers() {
  const manager = new LayerManager()
  managers.push(manager)
  document.body.append(manager.getEl())
  return manager
}

function checker(getStatus: () => Promise<any>) {
  const manager = layers()
  const checkers = new CheckerLauncher({
    getApi: () => ({ captcha: { getCaptchaStatus: getStatus } }) as any,
    getCaptchaIframeURL: () => 'about:blank',
    getLayers: () => manager,
    getUser: () => new UserManager({}),
    onReload: vi.fn(),
  })
  return { manager, checkers }
}

describe('owned resources', () => {
  it('stops captcha polling before and after a request is in flight', async () => {
    const pending = deferred<{ data: { is_pass: boolean } }>()
    const status = vi.fn(() => pending.promise)
    const { checkers } = checker(status)
    const success = vi.fn()
    const result = checkers
      .checkCaptcha({ iframe: 'captcha', onSuccess: success })
      .catch((error) => error)
    await vi.advanceTimersByTimeAsync(1000)
    expect(status).toHaveBeenCalledTimes(1)
    checkers.dispose()
    expect(await result).toBeInstanceOf(Error)
    pending.resolve({ data: { is_pass: true } })
    await vi.advanceTimersByTimeAsync(5000)
    expect(success).not.toHaveBeenCalled()
    expect(status).toHaveBeenCalledTimes(1)

    const another = checker(status)
    const cancelled = another.checkers.checkCaptcha({ iframe: 'captcha' }).catch(() => {})
    another.checkers.dispose()
    await cancelled
    await vi.advanceTimersByTimeAsync(2000)
    expect(status).toHaveBeenCalledTimes(1)
  })

  it('cancels the checker when the layer mask closes it', async () => {
    const status = vi.fn(async () => ({ data: { is_pass: false } }))
    const { manager, checkers } = checker(status)
    const cancel = vi.fn()
    const result = checkers.checkCaptcha({ iframe: 'captcha', onCancel: cancel }).catch(() => {})
    manager.getEl().querySelector<HTMLElement>('.atk-layer-mask')!.click()
    await result
    await vi.advanceTimersByTimeAsync(5000)
    expect(cancel).toHaveBeenCalledTimes(1)
    expect(status).not.toHaveBeenCalled()
    manager.destroy()
    expect(cancel).toHaveBeenCalledTimes(1)
  })

  it('ignores a login result after the backend checker is cancelled', async () => {
    const pending = deferred<{ token: string }>()
    const { checkers, manager } = checker(vi.fn())
    const success = vi.fn()
    checkers.check(
      {
        body: () => document.createElement('div'),
        request: () => pending.promise,
        onSuccess: success,
      },
      {},
    )
    manager.getEl().querySelector<HTMLButtonElement>('[data-action="confirm"]')!.click()
    checkers.cancelAll()
    pending.resolve({ token: 'old-backend-token' })
    await pending.promise
    await Promise.resolve()
    expect(success).not.toHaveBeenCalled()
  })

  it('keeps page scrolling locked until the last instance closes its layers', () => {
    document.body.style.overflow = 'auto'
    document.body.style.paddingRight = '7px'
    const first = layers()
    first.create('first').show()
    const second = layers() // Created while the first already owns the lock.
    const layer = second.create('second')
    layer.show()
    layer.show()
    first.destroy()
    expect(document.body.style.overflow).toBe('hidden')
    second.destroy()
    expect(document.body.style.overflow).toBe('auto')
    expect(document.body.style.paddingRight).toBe('7px')
  })

  it('releases destroyed layers and does not retain their mask listeners', () => {
    const manager = layers()
    const closed = vi.fn()
    for (let i = 0; i < 20; i++) {
      const layer = manager.create(String(i))
      layer.setOnAfterHide(closed)
      layer.show()
      layer.destroy()
    }
    expect(closed).toHaveBeenCalledTimes(20)
    manager.getEl().querySelector<HTMLElement>('.atk-layer-mask')!.click()
    expect(closed).toHaveBeenCalledTimes(20)
    expect(manager.getEl().querySelectorAll('.atk-layer-item')).toHaveLength(0)
  })

  it('unlocks scrolling even if a layer callback throws, and cannot reopen a disposed manager', () => {
    const manager = layers()
    const layer = manager.create('failing-layer')
    layer.show()
    layer.setOnAfterHide(() => {
      throw new Error('Layer cleanup failed')
    })
    const another = manager.create('remaining-layer')
    const closed = vi.fn()
    another.setOnAfterHide(closed)
    another.show()
    expect(() => manager.destroy()).toThrow('Layer cleanup failed')
    expect(closed).toHaveBeenCalledTimes(1)
    another.show()
    expect(manager.getEl().isConnected).toBe(false)
    expect(document.body.style.overflow).toBe('')
    manager.create('late-layer').show()
    expect(document.body.style.overflow).toBe('')
  })

  it('cancels remaining checks even when a cancellation callback throws', () => {
    const { checkers, manager } = checker(vi.fn())
    const simple = { body: () => document.createElement('div'), request: async () => ({}) }
    checkers.check(simple, {
      onCancel: () => {
        throw new Error('Cancel failed')
      },
    })
    const cancelled = vi.fn()
    checkers.check(simple, { onCancel: cancelled })
    expect(() => checkers.dispose()).toThrow('Cancel failed')
    expect(cancelled).toHaveBeenCalledTimes(1)
    expect(manager.getEl().querySelectorAll('.atk-layer-item')).toHaveLength(0)
    expect(vi.getTimerCount()).toBe(0)
  })

  it('runs remaining cleanup after a plugin throws and makes the context inert', () => {
    const root = document.createElement('div')
    const ctx = new Context(root)
    EventsService(ctx)
    ConfigService(ctx)
    ctx.trigger('mounted')
    root.append(document.createElement('span'))
    const cleanup = vi.fn()
    const changed = vi.fn()
    ctx.on('updated', changed)
    ctx.on('unmounted', () => {
      throw new Error('Plugin cleanup failed')
    })
    ctx.on('unmounted', cleanup)
    expect(() => ctx.destroy()).toThrow('Plugin cleanup failed')
    expect(cleanup).toHaveBeenCalledTimes(1)
    expect(root.childElementCount).toBe(0)
    ctx.updateConf({ locale: 'zh-CN' })
    ctx.trigger('updated', ctx.getConf())
    // eslint-disable-next-line artalk/artalk-plugin -- Exercise direct access after teardown.
    ctx.inject('config').update({ locale: 'en' })
    expect(changed).not.toHaveBeenCalled()
    root.append(document.createElement('p'))
    ctx.destroy()
    expect(root.childElementCount).toBe(1)
  })

  it('stops the current event dispatch when an instance is disposed from a callback', () => {
    const events = new EventManager<{ mounted: undefined }>()
    const late = vi.fn()
    events.on('mounted', () => events.dispose())
    events.on('mounted', late)
    events.trigger('mounted')
    events.on('mounted', late)
    events.trigger('mounted')
    expect(late).not.toHaveBeenCalled()
  })
})

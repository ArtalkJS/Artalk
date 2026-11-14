import type { EventManager as IEventManager, Event, EventHandler, EventOptions } from '@/types'

export class EventManager<T> implements IEventManager<T> {
  private events: Event<T>[] = []
  private disposed = false

  /**
   * Add an event listener for a specific event name
   */
  public on<K extends keyof T>(name: K, handler: EventHandler<T[K]>, opts: EventOptions = {}) {
    if (this.disposed) return
    this.events.push({
      name,
      handler: handler as EventHandler<T[keyof T]>,
      ...opts,
    })
  }

  /**
   * Remove an event listener for a specific event name and handler
   */
  public off<K extends keyof T>(name: K, handler: EventHandler<T[K]>) {
    if (!handler) return // not allow remove all events with same name
    this.events = this.events.filter((evt) => !(evt.name === name && evt.handler === handler))
  }

  /**
   * Trigger an event with an optional payload
   */
  public trigger<K extends keyof T>(name: K, payload?: T[K]) {
    let failed = false
    let error: unknown
    this.events
      .slice(0) // make a copy, in case listeners are removed while iterating
      .filter((evt) => evt.name === name && typeof evt.handler === 'function')
      .forEach((evt) => {
        if (this.disposed) return
        if (evt.once) this.off(name, evt.handler)
        try {
          evt.handler(payload!)
        } catch (err) {
          if (name !== 'unmounted') throw err
          // One plugin must not prevent the remaining cleanup handlers.
          if (!failed) error = err
          failed = true
        }
      })
    if (failed) throw error
  }

  /** @internal */
  public dispose() {
    this.disposed = true
    this.events = []
  }
}

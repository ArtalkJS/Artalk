import type { EventManager as IEventManager, Event, EventHandler, EventOptions } from '@/types'

export class EventManager<T> implements IEventManager<T> {
  private events: Event<T>[] = []

  /**
   * Add an event listener for a specific event name
   */
  public on<K extends keyof T>(name: K, handler: EventHandler<T[K]>, opts: EventOptions = {}) {
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
    this.events
      .slice(0) // make a copy, in case listeners are removed while iterating
      .filter((evt) => evt.name === name && typeof evt.handler === 'function')
      .forEach((evt) => {
        if (evt.once) this.off(name, evt.handler)
        evt.handler(payload!)
      })
  }
}

export type EventHandler<T> = (payload: T) => void
export interface Event<PayloadMap, K extends keyof PayloadMap = keyof PayloadMap>
  extends EventOptions {
  name: K
  handler: EventHandler<PayloadMap[K]>
}
export interface EventOptions {
  once?: boolean
}

export interface EventManagerFuncs<PayloadMap> {
  on<K extends keyof PayloadMap>(
    name: K,
    handler: EventHandler<PayloadMap[K]>,
    opts?: EventOptions,
  ): void
  off<K extends keyof PayloadMap>(name: K, handler: EventHandler<PayloadMap[K]>): void
  trigger<K extends keyof PayloadMap>(name: K, payload?: PayloadMap[K]): void
}

export default class EventManager<PayloadMap> implements EventManagerFuncs<PayloadMap> {
  private events: Event<PayloadMap>[] = []

  /**
   * Add an event listener for a specific event name
   */
  public on<K extends keyof PayloadMap>(
    name: K,
    handler: EventHandler<PayloadMap[K]>,
    opts: EventOptions = {},
  ) {
    this.events.push({
      name,
      handler: handler as EventHandler<PayloadMap[keyof PayloadMap]>,
      ...opts,
    })
  }

  /**
   * Remove an event listener for a specific event name and handler
   */
  public off<K extends keyof PayloadMap>(name: K, handler: EventHandler<PayloadMap[K]>) {
    if (!handler) return // not allow remove all events with same name
    this.events = this.events.filter((evt) => !(evt.name === name && evt.handler === handler))
  }

  /**
   * Trigger an event with an optional payload
   */
  public trigger<K extends keyof PayloadMap>(name: K, payload?: PayloadMap[K]) {
    this.events
      .slice(0) // make a copy, in case listeners are removed while iterating
      .filter((evt) => evt.name === name && typeof evt.handler === 'function')
      .forEach((evt) => {
        if (evt.once) this.off(name, evt.handler)
        evt.handler(payload!)
      })
  }
}

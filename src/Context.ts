import { ArtalkConfig } from '~/types/artalk-config'
import User from './lib/user'
import { EventPayloadMap, Event, EventScopeType, Handler } from '~/types/event'

export default class Context {
  public cid: number // Context 唯一标识
  public rootEl: HTMLElement
  public conf: ArtalkConfig
  public user: User

  private eventList: Event[] = []

  public constructor (rootEl: HTMLElement, conf: ArtalkConfig) {
    this.cid = +new Date()
    this.rootEl = rootEl
    this.conf = conf
    this.user = new User(this.conf)
  }

  public on<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>, scope: EventScopeType = 'internal'): void {
    this.eventList.push({ name, handler: handler as any, scope })
  }

  public off<K extends keyof EventPayloadMap>(name: K, handler?: Handler<EventPayloadMap[K]>, scope: EventScopeType = 'internal'): void {
    this.eventList = this.eventList.filter((evt) => {
      if (handler) return !(evt.name === name && evt.handler === handler && evt.scope === scope)
      return !(evt.name === name && evt.scope === scope) // 删除全部相同 name event
    })
  }

  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K], scope?: EventScopeType): void {
    this.eventList
      .filter((evt) => evt.name === name && (scope ? (evt.scope === scope) : true))
      .map((evt) => evt.handler)
      .forEach((handler) => handler(payload as any))
  }
}

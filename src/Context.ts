import { ArtalkConfig } from '~/types/artalk-config'
import User from './lib/user'
import { EventPayloadMap, Listener, Event } from '~/types/event'

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

  public on<K extends keyof EventPayloadMap>(name: K, listener: Listener<EventPayloadMap[K]>): void {
    this.eventList.push({ name, listener: listener as any })
  }

  public off<K extends keyof EventPayloadMap>(name: K, listener?: Listener<EventPayloadMap[K]>): void {
    this.eventList = this.eventList.filter((evt) => {
      if (listener) return !(evt.name === name && evt.listener === listener)
      return !(evt.name === name) // 删除全部相同 name event
    })
  }

  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]): void {
    this.eventList
      .filter((evt) => evt.name === name)
      .map((evt) => evt.listener)
      .forEach((listener) => listener(payload as any))
  }
}

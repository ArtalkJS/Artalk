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

  public addEventListener<K extends keyof EventPayloadMap>(name: K, listener: Listener<EventPayloadMap[K]>): void {
    this.eventList.push({ name, listener: listener as any })
  }

  public removeEventListener<K extends keyof EventPayloadMap>(name: K): void {
    this.eventList = this.eventList.filter((event) => event.name !== name)
  }

  public dispatchEvent<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]): void {
    this.eventList
      .filter((event) => event.name === name)
      .map((event) => event.listener)
      .forEach((listener) => listener(payload as any))
  }
}

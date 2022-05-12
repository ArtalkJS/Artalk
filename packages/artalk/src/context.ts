import { marked as libMarked } from 'marked'
import ArtalkConfig from '~/types/artalk-config'
import { EventPayloadMap, Event, EventScopeType, Handler } from '~/types/event'
import { internal as internalLocales, I18n } from './i18n'
import User from './lib/user'

/**
 * Artalk Context
 */
export default class Context {
  public cid: number // Context 唯一标识
  public $root: HTMLElement
  public conf: ArtalkConfig
  public user: User

  private eventList: Event[] = []

  public constructor (rootEl: HTMLElement, conf: ArtalkConfig) {
    this.cid = +new Date()
    this.$root = rootEl
    this.conf = conf
    this.user = new User(this.conf)

    this.$root.setAttribute('atk-run-id', this.cid.toString())
  }

  public on<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>, scope: EventScopeType = 'internal') {
    this.eventList.push({ name, handler: handler as any, scope })
  }

  public off<K extends keyof EventPayloadMap>(name: K, handler?: Handler<EventPayloadMap[K]>, scope: EventScopeType = 'internal') {
    this.eventList = this.eventList.filter((evt) => {
      if (handler) return !(evt.name === name && evt.handler === handler && evt.scope === scope)
      return !(evt.name === name && evt.scope === scope) // 删除全部相同 name event
    })
  }

  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K], scope?: EventScopeType) {
    this.eventList
      .filter((evt) => evt.name === name && (scope ? (evt.scope === scope) : true))
      .map((evt) => evt.handler)
      .forEach((handler) => handler(payload as any))
  }

  public markedInstance!: typeof libMarked
  public markedReplacers: ((raw: string) => string)[] = []

  public $t(key: keyof I18n, args: {[key: string]: string} = {}): string {
    let locales = this.conf.i18n
    if (typeof locales === 'string') {
      locales = internalLocales[locales]
    }

    let str = locales?.[key] || key
    str = str.replace(/\{\s*(\w+?)\s*\}/g, (_, token) => args[token] || '')

    return str
  }
}

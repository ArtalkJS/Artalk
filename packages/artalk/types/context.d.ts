import { marked as libMarked } from 'marked'
import ArtalkConfig from './artalk-config'
import { EventPayloadMap, Event, EventScopeType, Handler } from './event'
import { internal as internalLocales, I18n } from '../src/i18n'
import User from '../src/lib/user'

/**
 * Context 接口
 * @desc 面向接口的编程
 */
export default interface ContextApi {
  cid: number
  $root: HTMLElement
  conf: ArtalkConfig
  user: User
  on<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>, scope?: EventScopeType): void
  off<K extends keyof EventPayloadMap>(name: K, handler?: Handler<EventPayloadMap[K]>, scope?: EventScopeType): void
  trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K], scope?: EventScopeType): void
  $t(key: keyof I18n, args?: {[key: string]: string}): string
  markedInstance: typeof libMarked
  markedReplacers: ((raw: string) => string)[]
}

import { marked as libMarked } from 'marked'
import ArtalkConfig from '~/types/artalk-config'
import { Event } from '~/types/event'
import { internal as internalLocales, I18n } from './i18n'
import User from './lib/user'
import ContextApi from '../types/context'

/**
 * Artalk Context
 */
export default class Context implements ContextApi {
  public cid: number // Context 唯一标识
  public conf: ArtalkConfig
  public user: User
  public $root: HTMLElement

  private eventList: Event[] = []

  public constructor($root: HTMLElement, conf: ArtalkConfig) {
    this.cid = +new Date()
    this.conf = conf
    this.user = new User(this)

    this.$root = $root
    this.$root.setAttribute('atk-run-id', this.cid.toString())
  }

  public on(name: any, handler: any, scope: any = 'internal') {
    this.eventList.push({ name, handler, scope })
  }

  public off(name: any, handler: any, scope: any = 'internal') {
    this.eventList = this.eventList.filter((evt) => {
      if (handler) return !(evt.name === name && evt.handler === handler && evt.scope === scope)
      return !(evt.name === name && evt.scope === scope) // 删除全部相同 name event
    })
  }

  public trigger(name: any, payload?: any, scope?: any) {
    this.eventList
      .filter((evt) => evt.name === name && (scope ? (evt.scope === scope) : true))
      .map((evt) => evt.handler)
      .forEach((handler) => handler(payload))
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

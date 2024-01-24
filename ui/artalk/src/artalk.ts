import './style/main.scss'

import type { ArtalkConfig, EventPayloadMap, ArtalkPlugin, ContextApi } from '@/types'
import type { EventHandler } from './lib/event-manager'
import Context from './context'
import { handelCustomConf, convertApiOptions } from './config'
import Services from './service'
import { DefaultPlugins } from './plugins'
import * as Stat from './plugins/stat'
import { Api } from './api'
import type { TInjectedServices } from './service'

/** Global Plugins for all instances */
const GlobalPlugins: ArtalkPlugin[] = [ ...DefaultPlugins ]

/**
 * Artalk
 *
 * @see https://artalk.js.org
 */
export default class Artalk {
  public ctx!: ContextApi

  /** Plugins */
  protected plugins: ArtalkPlugin[] = [ ...GlobalPlugins ]

  constructor(conf: Partial<ArtalkConfig>) {
    // Init Config
    const handledConf = handelCustomConf(conf, true)

    // Init Context
    this.ctx = new Context(handledConf)

    // Init Services
    Object.entries(Services).forEach(([name, initService]) => {
      const obj = initService(this.ctx)
      if (obj) this.ctx.inject(name as keyof TInjectedServices, obj) // auto inject deps to ctx
    })

    // Init Plugins
    this.plugins.forEach(plugin => {
      if (typeof plugin === 'function') plugin(this.ctx)
    })

    // Trigger inited event
    this.ctx.trigger('inited')
  }

  /** Get the config of Artalk */
  public getConf() {
    return this.ctx.getConf()
  }

  /** Get the root element of Artalk */
  public getEl() {
    return this.ctx.$root
  }

  /** Update config of Artalk */
  public update(conf: Partial<ArtalkConfig>) {
    this.ctx.updateConf(conf)
    return this
  }

  /** Reload comment list of Artalk */
  public reload() {
    this.ctx.reload()
  }

  /** Destroy instance of Artalk */
  public destroy() {
    this.ctx.trigger('destroy')
    this.ctx.$root.remove()
  }

  /** Add an event listener */
  public on<K extends keyof EventPayloadMap>(name: K, handler: EventHandler<EventPayloadMap[K]>) {
    this.ctx.on(name, handler)
  }

  /** Remove an event listener */
  public off<K extends keyof EventPayloadMap>(name: K, handler: EventHandler<EventPayloadMap[K]>) {
    this.ctx.off(name, handler)
  }

  /** Trigger an event */
  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]) {
    this.ctx.trigger(name, payload)
  }

  /** Set dark mode */
  public setDarkMode(darkMode: boolean) {
    this.ctx.setDarkMode(darkMode)
  }

  // ===========================
  //       Static Members
  // ===========================

  /** Init Artalk */
  public static init(conf: Partial<ArtalkConfig>): Artalk {
    return new Artalk(conf)
  }

  /** Use plugin, the plugin will be used when Artalk.init */
  public static use(plugin: ArtalkPlugin) {
    if (GlobalPlugins.includes(plugin)) return
    GlobalPlugins.push(plugin)
  }

  /** Load count widget */
  public static loadCountWidget(c: Partial<ArtalkConfig>) {
    const conf = handelCustomConf(c, true)

    Stat.initCountWidget({
      getApi: () => new Api(convertApiOptions(conf)),
      siteName: conf.site,
      pageKey: conf.pageKey,
      countEl: conf.countEl,
      pvEl: conf.pvEl,
      pvAdd: false
    })
  }

  // ===========================
  //         Deprecated
  // ===========================

  /** @deprecated Please use `getEl()` instead */
  public get $root() { return this.ctx.$root }

  /** @description Please use `getConf()` instead */
  public get conf() { return this.ctx.getConf() }
}

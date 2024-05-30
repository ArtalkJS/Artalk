import './style/main.scss'

import type { ArtalkConfig, EventPayloadMap, ArtalkPlugin, ContextApi } from '@/types'
import type { EventHandler } from './lib/event-manager'
import Context from './context'
import { handelCustomConf, convertApiOptions } from './config'
import Services from './service'
import * as Stat from './plugins/stat'
import { Api } from './api'
import type { TInjectedServices } from './service'
import { GlobalPlugins, load } from './load'

/**
 * Artalk
 *
 * @see https://artalk.js.org
 */
export default class Artalk {
  public ctx!: ContextApi

  constructor(conf: Partial<ArtalkConfig>) {
    // Init Config
    const handledConf = handelCustomConf(conf, true)

    // Init Context
    this.ctx = new Context(handledConf)

    // Init Services
    Object.entries(Services).forEach(([name, initService]) => {
      const obj = initService(this.ctx)
      obj && this.ctx.inject(name as keyof TInjectedServices, obj) // auto inject deps to ctx
    })

    if (import.meta.env.DEV && import.meta.env.VITEST) {
      global.devLoadArtalk = () => load(this.ctx)
    } else {
      load(this.ctx)
    }
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
    this.ctx.trigger('unmounted')
    while (this.ctx.$root.firstChild) {
      this.ctx.$root.removeChild(this.ctx.$root.firstChild)
    }
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
    GlobalPlugins.add(plugin)
  }

  /** Load count widget */
  public static loadCountWidget(c: Partial<ArtalkConfig>) {
    const conf = handelCustomConf(c, true)

    Stat.initCountWidget({
      getApi: () => new Api(convertApiOptions(conf)),
      siteName: conf.site,
      countEl: conf.countEl,
      pvEl: conf.pvEl,
      pageKeyAttr: conf.statPageKeyAttr,
      pvAdd: false,
    })
  }

  // ===========================
  //         Deprecated
  // ===========================

  /** @deprecated Please use `getEl()` instead */
  public get $root() {
    return this.ctx.$root
  }

  /** @description Please use `getConf()` instead */
  public get conf() {
    return this.ctx.getConf()
  }
}

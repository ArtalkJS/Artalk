import './style/main.scss'

import Context from './context'
import { handelCustomConf, convertApiOptions, getRootEl } from './config'
import * as Stat from './plugins/stat'
import { Api } from './api'
import { GlobalPlugins, PluginOptions, mount } from './mount'
import { ConfigService } from './services/config'
import { EventsService } from './services/events'
import type {
  ConfigPartial,
  EventPayloadMap,
  ArtalkPlugin,
  Context as IContext,
  EventHandler,
} from '@/types'

/**
 * Artalk
 *
 * @see https://artalk.js.org
 */
export default class Artalk {
  public ctx: IContext

  constructor(conf: ConfigPartial) {
    // Init Root Element
    const $root = getRootEl(conf)
    $root.classList.add('artalk')
    $root.innerHTML = ''
    conf.darkMode == true && $root.classList.add('atk-dark-mode')

    // Init Context
    const ctx = (this.ctx = new Context($root))

    // Init required services
    ;(() => {
      // Init event manager
      EventsService(ctx)

      // Init config service
      ConfigService(ctx)
    })()

    // Apply local conf first
    ctx.updateConf(conf)

    // Trigger created event
    ctx.trigger('created')

    // Load plugins and remote config, then mount Artalk
    const mountArtalk = async () => {
      await mount(conf, ctx)

      // Trigger mounted event
      ctx.trigger('mounted')
    }

    if (import.meta.env.DEV && import.meta.env.VITEST) {
      global.devMountArtalk = mountArtalk
    } else {
      mountArtalk()
    }
  }

  /** Get the config of Artalk */
  public getConf() {
    return this.ctx.getConf()
  }

  /** Get the root element of Artalk */
  public getEl() {
    return this.ctx.getEl()
  }

  /** Update config of Artalk */
  public update(conf: ConfigPartial) {
    this.ctx.updateConf(conf)
  }

  /** Reload comment list of Artalk */
  public reload() {
    this.ctx.reload()
  }

  /** Destroy instance of Artalk */
  public destroy() {
    this.ctx.destroy()
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
  public static init(conf: ConfigPartial): Artalk {
    return new Artalk(conf)
  }

  /** Use plugin, the plugin will be used when Artalk.init */
  public static use<T = any>(plugin: ArtalkPlugin<T>, options?: T) {
    GlobalPlugins.add(plugin)
    PluginOptions.set(plugin, options)
  }

  /** Load count widget */
  public static loadCountWidget(c: ConfigPartial) {
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
    console.warn('`$root` is deprecated, please use `getEl()` instead')
    return this.getEl()
  }

  /** @description Please use `getConf()` instead */
  public get conf() {
    console.warn('`conf` is deprecated, please use `getConf()` instead')
    return this.getConf()
  }
}

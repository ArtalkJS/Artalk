import './style/main.scss'

import type ArtalkConfig from '~/types/artalk-config'
import type { EventPayloadMap, Handler } from '~/types/event'
import type ArtalkPlug from '~/types/plug'
import type Context from '~/types/context'
import ConcreteContext from './context'
import defaults from './defaults'
import { handelBaseConf } from './config'
import Services from './service'
import * as Stat from './lib/stat'
import ListLite from './list/list-lite'

/**
 * Artalk
 *
 * @see https://artalk.js.org
 */
export default class Artalk {
  private static instance?: Artalk

  public static ListLite = ListLite
  public static readonly defaults: ArtalkConfig = defaults

  public conf!: ArtalkConfig
  public ctx!: Context
  public $root!: HTMLElement

  /** Plugins */
  protected static plugins: ArtalkPlug[] = [ Stat.PvCountWidget ]

  /** 禁用的组件 */
  public static DisabledComponents: string[] = []

  /** 构造函数 */
  constructor(conf: Partial<ArtalkConfig>) {
    if (Artalk.instance) return Artalk.instance.update(conf)

    // 初始化基本配置
    this.conf = handelBaseConf(conf)
    if (this.conf.el instanceof HTMLElement) this.$root = this.conf.el

    // 初始化 Context
    this.ctx = new ConcreteContext(this.conf, this.$root)

    // 内建服务初始化
    Object.entries(Services).forEach(([name, initService]) => {
      if (Artalk.DisabledComponents.includes(name)) return
      const obj = initService(this.ctx)
      if (obj) this.ctx.inject(name as any, obj) // auto inject deps to ctx
    })

    // 插件初始化 (global scope)
    Artalk.plugins.forEach(plugin => {
      if (typeof plugin === 'function')
        plugin(this.ctx)
    })
  }

  /** Init Artalk (单例模式) */
  public static init(conf: Partial<ArtalkConfig>): Artalk {
    if (this.instance) return this.instance.update(conf)
    this.instance = new Artalk(conf)
    return this.instance
  }

  /** 设置暗黑模式 */
  public setDarkMode(darkMode: boolean) {
    this.ctx.setDarkMode(darkMode)
  }

  /** Use Plugin (plugin will be called in instance `use` func) */
  public use(plugin: ArtalkPlug) {
    Artalk.plugins.push(plugin)
    if (typeof plugin === 'function') plugin(this.ctx)
  }

  /** Use Plugin (static method) */
  public static use(plugin: ArtalkPlug) {
    this.plugins.push(plugin)
  }

  /** Update config of Artalk */
  public update(conf: Partial<ArtalkConfig>) {
    if (!Artalk.instance) throw Error('cannot call `update` function before call `load`')
    Artalk.instance.ctx.updateConf(conf)
    return Artalk.instance
  }

  /** Reload comment list of Artalk */
  public reload() {
    this.ctx.listReload()
  }

  /** Destroy instance of Artalk */
  public destroy() {
    if (!Artalk.instance) throw Error('cannot call `destroy` function before call `load`')
    //...
    delete Artalk.instance
  }

  /** 监听事件 */
  public on<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>) {
    this.ctx.on(name, handler, 'external')
  }

  /** 解除监听事件 */
  public off<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>) {
    this.ctx.off(name, handler, 'external')
  }

  /** 触发事件 */
  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]) {
    this.ctx.trigger(name, payload, 'external')
  }

  /** 装载数量统计元素 */
  public static loadCountWidget(conf: Partial<ArtalkConfig>) {
    const ctx = new ConcreteContext(handelBaseConf(conf))
    Stat.initCountWidget({ ctx, pvAdd: false })
  }

  /** @deprecated Please use `loadCountWidget` instead */
  public static LoadCountWidget(conf: Partial<ArtalkConfig>) {
    console.warn('The method `LoadCountWidget` is deprecated, please use `loadCountWidget` instead.')
    this.loadCountWidget(conf)
  }
}

import './style/main.less'

import ArtalkConfig from '~/types/artalk-config'
import { EventPayloadMap, Handler } from '~/types/event'
import ArtalkPlug from '~/types/plug'
import Context from '~/types/context'
import ConcreteContext from './context'
import defaults from './defaults'
import { internal as internalLocales } from './i18n'

import CheckerLauncher from './lib/checker'
import Editor from './editor'
import List from './list'
import SidebarLayer from './layer/sidebar-layer'

import Layer from './layer'
import * as Utils from './lib/utils'
import * as Ui from './lib/ui'
import * as Stat from './lib/stat'

/**
 * Artalk
 *
 * @website https://artalk.js.org
 */
export default class Artalk {
  public static readonly defaults: ArtalkConfig = defaults

  public conf!: ArtalkConfig
  public ctx!: Context
  public $root!: HTMLElement

  public checkerLauncher!: CheckerLauncher
  public editor!: Editor
  public list!: List
  public sidebarLayer!: SidebarLayer

  /** Plugins (in global scope)  */
  protected static Plugins: ArtalkPlug[] = [ Stat.PvCountWidget ]

  /** Plugins (in a instance scope) */
  protected instancePlugins: ArtalkPlug[] = []

  constructor(customConf: Partial<ArtalkConfig>) {
    /* 初始化基本配置 */
    this.conf = Artalk.HandelBaseConf(customConf)
    if (this.conf.el instanceof HTMLElement) this.$root = this.conf.el

    /* 初始化 Context */
    this.ctx = new ConcreteContext(this.conf, this.$root)

    // @ts-ignore 远程加载引用后端的配置
    if (this.conf.useBackendConf) return this.loadConfRemoteAndInitComponents()

    /* 初始化组件 */
    this.initComponents()
  }

  /** 组件初始化 */
  private initComponents() {
    this.initLocale()
    this.initLayer()
    this.initDarkMode()
    Utils.initMarked(this.ctx)

    // CheckerLauncher
    this.checkerLauncher = new CheckerLauncher(this.ctx)
    this.ctx.setCheckerLauncher(this.checkerLauncher)

    // 编辑器
    this.editor = new Editor(this.ctx)
    this.ctx.setEditor(this.editor)
    this.$root.appendChild(this.editor.$el)

    // 评论列表
    this.list = new List(this.ctx)
    this.ctx.setList(this.list)
    this.$root.appendChild(this.list.$el)

    // 侧边栏
    this.sidebarLayer = new SidebarLayer(this.ctx)
    this.ctx.setSidebarLayer(this.sidebarLayer)
    this.$root.appendChild(this.sidebarLayer.$el)

    // 评论获取
    this.list.fetchComments(0)

    // 事件绑定初始化
    this.initEventBind()

    // 插件初始化 (global scope)
    Artalk.Plugins.forEach(plugin => {
      if (typeof plugin === 'function')
        plugin(this.ctx)
    })
  }

  /** 基本配置初始化 */
  private static HandelBaseConf(customConf: Partial<ArtalkConfig>): ArtalkConfig {
    // 合并默认配置
    const conf: ArtalkConfig = Utils.mergeDeep(Artalk.defaults, customConf)

    // 绑定元素
    if (typeof conf.el === 'string' && !!conf.el) {
      try {
        const findEl = document.querySelector<HTMLElement>(conf.el)
        if (!findEl) throw Error(`Target element "${conf.el}" was not found.`)
        conf.el = findEl
      } catch (e) {
        console.error(e)
        throw new Error('Please check your Artalk `el` config.')
      }
    }

    // 服务器配置
    conf.server = conf.server.replace(/\/$/, '').replace(/\/api\/?$/, '')

    // 默认 pageKey
    if (!conf.pageKey) {
      // @link http://bl.ocks.org/abernier/3070589
      conf.pageKey = `${window.location.pathname}`
    }

    // 默认 pageTitle
    if (!conf.pageTitle) {
      conf.pageTitle = `${document.title}`
    }

    return conf
  }

  /** 初始化组件根据远程加载的配置 */
  private async loadConfRemoteAndInitComponents() {
    await this.loadConfRemote()
    this.initComponents()
    return this
  }

  /** 远程加载配置 */
  private async loadConfRemote() {
    Ui.showLoading(this.$root)
    let backendConf = {}
    try {
      backendConf = await this.ctx.getApi().conf()
    } catch (err) { console.error("Load config from remote err", err) }
    this.ctx.conf = Utils.mergeDeep(this.ctx.conf, backendConf)
    Ui.hideLoading(this.$root)
  }

  /** 事件绑定初始化 */
  private initEventBind() {
    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.list.goToCommentDelay = false
      this.list.checkGoToCommentByUrlHash()
    })

    // 本地用户数据变更
    this.ctx.on('user-changed', () => {
      this.ctx.checkAdminShowEl()
      this.ctx.listRefreshUI()
    })
  }

  /** 语言初始化 */
  private initLocale() {
    if (typeof this.conf.locale === 'string') {
      if (this.conf.locale === 'auto') { // 自动切换
        this.conf.locale = ((navigator.languages) ? navigator.languages[0] : navigator.language)
      }

      this.conf.locale = this.conf.locale.replace(
        /^([a-zA-Z]+)(-[a-zA-Z]+)?$/,
        (_, p1: string, p2: string) => (p1.toLowerCase() + (p2 || '').toUpperCase())
      )

      if (!internalLocales[this.conf.locale]) { // 语言不存在时，使用 en
        console.log(`Locale "${this.conf.locale}" not found.`)
        this.conf.locale = 'en'
      }
    }
  }

  /** Layer 初始化 */
  private initLayer() {
    // 记录页面原始 Styles
    Layer.BodyOrgOverflow = document.body.style.overflow
    Layer.BodyOrgPaddingRight = document.body.style.paddingRight
  }

  /** 暗黑模式初始化 */
  private initDarkMode() {
    if (this.conf.darkMode === 'auto') {
      // 自动切换暗黑模式，事件监听
      const darkModeMedia = window.matchMedia('(prefers-color-scheme: dark)')
      darkModeMedia.addEventListener('change', (e) => { this.setDarkMode(e.matches) })
      this.setDarkMode(darkModeMedia.matches)
    } else {
      this.setDarkMode(this.conf.darkMode || false)
    }
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

  /** 重新加载 */
  public reload() {
    this.ctx.listReload()
  }

  /** 设置暗黑模式 */
  public setDarkMode(darkMode: boolean) {
    this.ctx.setDarkMode(darkMode)
  }

  /** Use Plugin (specific instance) */
  public use(plugin: ArtalkPlug) {
    this.instancePlugins.push(plugin)
    if (typeof plugin === 'function') plugin(this.ctx)
  }

  /** Use Plugin (static method for global scope) */
  public static use(plugin: ArtalkPlug) {
    this.Plugins.push(plugin)
  }

  /** @deprecated Please replace it with lowercase function name `use(...)`. */
  public static Use(plugin: ArtalkPlug) {
    this.use(plugin)
    console.warn('`Use(...)` is deprecated, replace it with lowercase `use(...)`.')
  }

  /** 装载数量统计元素 */
  public static LoadCountWidget(customConf: Partial<ArtalkConfig>) {
    const conf = this.HandelBaseConf(customConf)
    const ctx = new ConcreteContext(conf)
    Stat.initCountWidget({ ctx, pvAdd: false })
  }
}

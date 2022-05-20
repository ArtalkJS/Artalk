import './style/main.less'

import ArtalkConfig from '~/types/artalk-config'
import { EventPayloadMap, Handler } from '~/types/event'
import Context from '~/types/context'
import ConcreteContext from './context'
import defaults from './defaults'

import CheckerLauncher from './lib/checker'
import Editor from './editor'
import List from './list'
import SidebarLayer from './layer/sidebar-layer'

import Layer, { GetLayerWrap } from './layer'
import Api from './api'
import * as Utils from './lib/utils'
import * as Ui from './lib/ui'

/**
 * Artalk
 *
 * @website https://artalk.js.org
 */
export default class Artalk {
  public static readonly defaults: ArtalkConfig = defaults

  public ctx!: Context
  public conf!: ArtalkConfig
  public $root!: HTMLElement

  public api!: Api
  public checkerLauncher!: CheckerLauncher
  public editor!: Editor
  public list!: List
  public sidebarLayer!: SidebarLayer

  constructor(customConf: Partial<ArtalkConfig>) {
    // 配置
    this.conf = Utils.mergeDeep(Artalk.defaults, customConf)
    this.conf.server = this.conf.server.replace(/\/$/, '').replace(/\/api\/?$/, '')

    // 默认 pageKey
    if (!this.conf.pageKey) {
      // @link http://bl.ocks.org/abernier/3070589
      this.conf.pageKey = `${window.location.pathname}`
    }

    // 默认 pageTitle
    if (!this.conf.pageTitle) {
      this.conf.pageTitle = `${document.title}`
    }

    // 列表显示模式
    if (this.conf.nestMax && this.conf.nestMax <= 1) {
      this.conf.flatMode = true
    }

    // 装载元素
    if (!!this.conf.el && this.conf.el instanceof HTMLElement) {
      this.$root = this.conf.el
    } else {
      try {
        const $root = document.querySelector<HTMLElement>(this.conf.el)
        if (!$root) throw Error(`Sorry, target element "${this.conf.el}" was not found.`)
        this.$root = $root
      } catch (e) {
        console.error(e)
        throw new Error('Please check your Artalk `el` config.')
      }
    }

    // Context 初始化
    this.ctx = new ConcreteContext(this.$root, this.conf)

    // API 初始化
    this.api = new Api(this.ctx)
    this.ctx.setApi(this.api)

    // 界面初始化
    this.$root.classList.add('artalk')
    this.$root.innerHTML = ''

    // 初始化组件
    if (this.conf.useBackendConf) {
      // 复用后端的配置
      // @ts-ignore
      return (async () => {
        await this.loadConfRemote()
        this.initCore()
        return this // Promise<Artalk>
      })()
    }

    this.initCore()
  }

  private initCore() {
    /* 组件初始化 */
    this.initLayer()
    this.initDarkMode()

    // CheckerLauncher
    this.checkerLauncher = new CheckerLauncher(this.ctx)
    this.ctx.setCheckerLauncher(this.checkerLauncher)

    // 初始化 marked
    Utils.initMarked(this.ctx)

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

    // 其他
    this.initPV()

    // 插件初始化
    Artalk.Plugins.forEach(plugin => {
      if (typeof plugin === 'function') {
        plugin(this.ctx)
      }
    })
  }

  /** 获取远程配置 */
  private async loadConfRemote() {
    Ui.showLoading(this.$root)
    let backendConf = {}
    try {
      backendConf = await this.ctx.getApi().conf()
    } catch (err) { console.error("Load config from remote err", err) }
    this.ctx.conf = Utils.mergeDeep(this.ctx.conf, backendConf)
    Ui.hideLoading(this.$root)
  }

  /** 事件绑定 · 初始化 */
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

  /** 重新加载 */
  public reload() {
    this.list.fetchComments(0)
  }

  /** Layer 初始化 */
  initLayer() {
    // 记录页面原始 Styles
    Layer.BodyOrgOverflow = document.body.style.overflow
    Layer.BodyOrgPaddingRight = document.body.style.paddingRight
  }

  /** 暗黑模式 · 初始化 */
  public initDarkMode() {
    if (this.conf.darkMode === 'auto') {
      // 自动切换暗黑模式，事件监听
      const darkModeMedia = window.matchMedia('(prefers-color-scheme: dark)')
      darkModeMedia.addEventListener('change', (e) => { this.setDarkMode(e.matches) })
      this.setDarkMode(darkModeMedia.matches)
    } else {
      this.setDarkMode(this.conf.darkMode || false)
    }
  }

  /** 暗黑模式 · 设定 */
  public setDarkMode(darkMode: boolean) {
    const darkModeClassName = 'atk-dark-mode'

    this.ctx.conf.darkMode = darkMode
    this.ctx.trigger('conf-updated')

    if (this.conf.darkMode) {
      this.$root.classList.add(darkModeClassName)
    } else {
      this.$root.classList.remove(darkModeClassName)
    }

    // for Layer
    const { $wrap: $layerWrap } = GetLayerWrap(this.ctx)
    if ($layerWrap) {
      if (this.conf.darkMode) {
        $layerWrap.classList.add(darkModeClassName)
      } else {
        $layerWrap.classList.remove(darkModeClassName)
      }
    }
  }

  /** PV */
  public async initPV() {
    const curtPagePvNum = await this.ctx.getApi().pv()
    const curtPageKey = this.ctx.conf.pageKey

    // 界面更新数据
    const pvEl = this.conf.pvEl
    if (!pvEl || !document.querySelector(pvEl)) return

    // 当前页面 PV 数
    let pagePvs: {[key: string]: number} = {}
    pagePvs[curtPageKey] = curtPagePvNum

    // 查询其他页面的 PV 数
    const queryPageKeys = Array.from(document.querySelectorAll(pvEl))
      .filter(e => {
        const pageKeyAttr = e.getAttribute('data-page-key');
        return !!pageKeyAttr && pageKeyAttr !== curtPageKey
      })
      .map(e => e.getAttribute('data-page-key')!)

    if (queryPageKeys.length > 0) {
      const pvs: any = await this.ctx.getApi().stat('page_pv', queryPageKeys)
      pagePvs = { ...pagePvs, ...pvs }
    }

    this.updatePvEls(pagePvs)
  }

  private updatePvEls(pvs: {[key: string]: number}) {
    document.querySelectorAll(this.conf.pvEl).forEach((el) => {
      const pageKey = el.getAttribute('data-page-key') || this.ctx.conf.pageKey
      el.innerHTML = `${Number(pvs[pageKey] || 0)}`
    })
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

  /** Plugins */
  protected static Plugins: ((ctx: Context) => void)[] = []

  /** Enable Plugin */
  public static Use(plugin: ((ctx: Context) => void)) {
    this.Plugins.push(plugin)
  }
}

import './style/main.less'

import ArtalkConfig from '~/types/artalk-config'
import { EventPayloadMap, Handler } from '~/types/event'
import Context from './context'
import defaults from './defaults'

import CheckerLauncher from './lib/checker'
import Editor from './components/editor'
import List from './components/list'
import SidebarLayer from './components/sidebar-layer'

import Layer, { GetLayerWrap } from './components/layer'
import Api from './api'
import * as Utils from './lib/utils'
import * as Ui from './lib/ui'

/**
 * Artalk
 * @link https://artalk.js.org
 */
export default class Artalk {
  public static readonly defaults: ArtalkConfig = defaults

  public ctx!: Context
  public conf!: ArtalkConfig
  public $root!: HTMLElement

  public checkerLauncher!: CheckerLauncher
  public editor!: Editor
  public list!: List
  public sidebarLayer!: SidebarLayer

  constructor (customConf: ArtalkConfig) {
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
    this.ctx = new Context(this.$root, this.conf)

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
    // 组件初始化
    this.initLayer()
    this.initDarkMode()
    this.checkerLauncher = new CheckerLauncher(this.ctx)
    Utils.initMarked(this.ctx) // 初始化 marked

    // 编辑器
    this.editor = new Editor(this.ctx)
    this.$root.appendChild(this.editor.$el)

    // 评论列表
    this.list = new List(this.ctx)
    this.$root.appendChild(this.list.$el)

    // 侧边栏
    this.sidebarLayer = new SidebarLayer(this.ctx)
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
      backendConf = await (new Api(this.ctx)).conf()
    } catch (err) { console.error("配置远程获取失败", err) }
    this.ctx.conf = Utils.mergeDeep(this.ctx.conf, backendConf)
    Ui.hideLoading(this.$root)
  }

  /** 事件绑定 · 初始化 */
  private initEventBind() {
    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.list.checkGoToCommentByUrlHash()
    })

    // 仅管理员显示控制
    this.ctx.on('check-admin-show-el', () => {
      const items: HTMLElement[] = []

      this.$root.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

      // for layer
      const { $wrap: $layerWrap } = GetLayerWrap(this.ctx)
      if ($layerWrap) $layerWrap.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

      // for sidebar
      // TODO: 这个其实应该写在 packages/artalk-sidebar 里面的
      const $sidebarEl = document.querySelector<HTMLElement>('.atk-sidebar')
      if ($sidebarEl) $sidebarEl.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

      items.forEach(($item: HTMLElement) => {
        if (this.ctx.user.data.isAdmin)
          $item.classList.remove('atk-hide')
        else
          $item.classList.add('atk-hide')
      })
    })

    // 本地用户数据变更
    this.ctx.on('user-changed', () => {
      this.ctx.trigger('check-admin-show-el')
      this.ctx.trigger('list-refresh-ui')
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
    const darkModeClassName = 'atk-dark-mode'

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

  /** 暗黑模式 · 设定 */
  public setDarkMode(darkMode: boolean) {
    this.ctx.conf.darkMode = darkMode
    this.ctx.trigger('conf-updated')
    this.initDarkMode()
  }

  /** PV */
  public async initPV() {
    const pvNum = await new Api(this.ctx).pv()
    if (Number.isNaN(Number(pvNum)))
      return

    if (!this.conf.pvEl || !document.querySelector(this.conf.pvEl))
      return

    const $pv = document.querySelector<HTMLElement>(this.conf.pvEl)!
    // $pv.innerText = '-' // 默认占位符

    $pv.innerText = String(pvNum)
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

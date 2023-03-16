import Context from '~/types/context'
import CheckerLauncher from './lib/checker'
import Api from './api'
import Editor from './editor'
import List from './list'
import Layer from './layer'
import SidebarLayer from './layer/sidebar-layer'
import { initMarked } from './lib/marked'
import User from './lib/user'
import ListLite from './list/list-lite'
import * as DarkMode from './lib/dark-mode'
import * as I18n from './i18n'

/**
 * Services
 *
 * @description
 * 当函数有返回值时，该值将自动添加到 Context 类作为一个 private 成员的值，
 * 这个成员的名称与函数名对应。可在 Context 类中访问该对象（同事类）。
 */
const services = {
  // I18n
  i18n(ctx: Context) {
    I18n.setLocale(ctx.conf.locale)

    ctx.on('conf-loaded', () => {
      I18n.setLocale(ctx.conf.locale)
    })
  },

  // Markdown 组件
  markdown() {
    initMarked()
  },

  // User Store
  user(ctx: Context) {
    User.setContext(ctx)
    return User
  },

  // HTTP API client
  api(ctx: Context) {
    const api = new Api(ctx)
    return api
  },

  // CheckerLauncher
  checkerLauncher(ctx: Context) {
    const checkerLauncher = new CheckerLauncher(ctx)
    return checkerLauncher
  },

  // 编辑器
  editor(ctx: Context) {
    const editor = new Editor(ctx)
    ctx.$root.appendChild(editor.$el)
    return editor
  },

  // 评论列表
  list(ctx: Context): ListLite|undefined {
    // 评论列表
    const list = new List(ctx)
    ctx.$root.appendChild(list.$el)

    // 评论获取
    list.fetchComments(0)

    return list
  },

  // 弹出层
  layer(ctx: Context) {
    // 记录页面原始 CSS 属性
    Layer.BodyOrgOverflow = document.body.style.overflow
    Layer.BodyOrgPaddingRight = document.body.style.paddingRight
  },

  // 侧边栏 Layer
  sidebarLayer(ctx: Context) {
    const sidebarLayer = new SidebarLayer(ctx)
    return sidebarLayer
  },

  // 默认事件绑定
  eventsDefault(ctx: Context) {
    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      ctx.listHashGotoCheck()
    })

    // 本地用户数据变更
    ctx.on('user-changed', () => {
      ctx.checkAdminShowEl()
      ctx.listRefreshUI()
    })
  },

  // 夜间模式
  darkMode(ctx: Context) {
    DarkMode.syncDarkModeConf(ctx)

    ctx.on('conf-loaded', () => {
      DarkMode.syncDarkModeConf(ctx)
    })
  }
}

export default services

// type tricks for dependency injection
type TServiceImps = typeof services
type TObjectWithFuncs = {[k: string]: (...args: any) => any}
type TKeysOnlyReturn<T extends TObjectWithFuncs, V> = {[K in keyof T]: ReturnType<T[K]> extends V ? K : never}[keyof T]
type TOmitConditions = TKeysOnlyReturn<TServiceImps, void>
type TServiceInjectors = Omit<TServiceImps, TOmitConditions>
export type TInjectedServices = {[K in keyof TServiceInjectors]: ReturnType<TServiceInjectors[K]>}

import type { ContextApi } from '~/types'
import CheckerLauncher from './lib/checker'
import Api from './api'
import Editor from './editor/editor'
import Layer from './layer'
import SidebarLayer from './layer/sidebar-layer'
import User from './lib/user'
import List from './list/list'
import * as DarkMode from './lib/dark-mode'
import * as I18n from './i18n'
import { PlugManager } from './plugins/editor-kit'

/**
 * Services
 *
 * @description Call these services by `ctx.get('serviceName')` or `ctx.serviceName`
 */
const services = {
  // I18n
  i18n(ctx: ContextApi) {
    I18n.setLocale(ctx.conf.locale)

    ctx.on('conf-loaded', () => {
      I18n.setLocale(ctx.conf.locale)
    })
  },

  // User Store
  user(ctx: ContextApi) {
    User.setContext(ctx)
    return User
  },

  // HTTP API client
  api(ctx: ContextApi) {
    const api = new Api(ctx)
    return api
  },

  // CheckerLauncher
  checkerLauncher(ctx: ContextApi) {
    const checkerLauncher = new CheckerLauncher(ctx)
    return checkerLauncher
  },

  // 编辑器
  editor(ctx: ContextApi) {
    const editor = new Editor(ctx)
    ctx.$root.appendChild(editor.$el)
    return editor
  },

  // 评论列表
  list(ctx: ContextApi): List|undefined {
    // 评论列表
    const list = new List(ctx)
    ctx.$root.appendChild(list.$el)

    return list
  },

  // 弹出层
  layer(ctx: ContextApi) {
    // 记录页面原始 CSS 属性
    Layer.BodyOrgOverflow = document.body.style.overflow
    Layer.BodyOrgPaddingRight = document.body.style.paddingRight
  },

  // 侧边栏 Layer
  sidebarLayer(ctx: ContextApi) {
    const sidebarLayer = new SidebarLayer(ctx)
    return sidebarLayer
  },

  // 默认事件绑定
  eventsDefault(ctx: ContextApi) {
    // 本地用户数据变更
    ctx.on('user-changed', () => {
      ctx.checkAdminShowEl()
    })
  },

  // 夜间模式
  darkMode(ctx: ContextApi) {
    DarkMode.syncDarkModeConf(ctx)

    ctx.on('conf-loaded', () => {
      DarkMode.syncDarkModeConf(ctx)
    })
  },

  // Extra Service
  // -----------------------------------------
  // Only for type check
  // Not inject to ctx immediately,
  // but can be injected by other occasions

  editorPlugs(): PlugManager|undefined {
    return undefined
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

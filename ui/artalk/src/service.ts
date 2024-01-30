import type { ContextApi } from '@/types'
import CheckerLauncher from './components/checker'
import Editor from './editor/editor'
import SidebarLayer from './layer/sidebar-layer'

import List from './list/list'

import * as I18n from './i18n'
import { PlugManager } from './plugins/editor-kit'
import { LayerManager } from './layer/layer-manager'
import User from './lib/user'

/**
 * Services
 *
 * @description Call these services by `ctx.get('serviceName')` or `ctx.serviceName`
 */
const services = {
  // I18n
  i18n(ctx: ContextApi) {
    I18n.setLocale(ctx.conf.locale)

    ctx.watchConf(['locale'], (conf) => {
      I18n.setLocale(conf.locale)
    })
  },

  // User Store
  user(ctx: ContextApi) {
    const user = new User({
      onUserChanged: (data) => {
        ctx.trigger('user-changed', data)
      }
    })
    return user
  },

  // 弹出层
  layerManager(ctx: ContextApi) {
    return new LayerManager(ctx)
  },

  // CheckerLauncher
  checkerLauncher(ctx: ContextApi) {
    const checkerLauncher = new CheckerLauncher({
      getCtx: () => ctx,
      getApi: () => ctx.getApi(),
      onReload: () => ctx.reload(),

      // make sure suffix with a slash, because it will be used as a base url when call `fetch`
      getCaptchaIframeURL: () => `${ctx.conf.server}/api/v2/captcha/?t=${+new Date()}`
    })
    return checkerLauncher
  },

  // 编辑器
  editor(ctx: ContextApi) {
    const editor = new Editor(ctx)
    ctx.$root.appendChild(editor.$el)
    return editor
  },

  // 评论列表
  list(ctx: ContextApi): List {
    const list = new List(ctx)
    ctx.$root.appendChild(list.$el)
    return list
  },

  // 侧边栏 Layer
  sidebarLayer(ctx: ContextApi) {
    const sidebarLayer = new SidebarLayer(ctx)
    return sidebarLayer
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

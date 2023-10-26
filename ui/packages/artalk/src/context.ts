import type ArtalkConfig from '~/types/artalk-config'
import type { CommentData, ListFetchParams } from '~/types/artalk-data'
import type { EventPayloadMap } from '~/types/event'
import type ContextApi from '~/types/context'
import type { TInjectedServices } from './service'

import * as Utils from './lib/utils'
import * as DarkMode from './lib/dark-mode'
import * as marked from './lib/marked'
import { CheckerCaptchaPayload, CheckerPayload } from './lib/checker'

import { DataManager } from './data'
import * as I18n from './i18n'
import { getLayerWrap } from './layer'
import { SidebarShowPayload } from './layer/sidebar-layer'
import EventManager from './lib/event-manager'
import { handelBaseConf } from './config'

// Auto dependency injection
interface Context extends TInjectedServices { }

/**
 * Artalk Context
 */
class Context implements ContextApi {
  /* 运行参数 */
  public conf: ArtalkConfig
  public data: DataManager
  public $root: HTMLElement
  public markedReplacers: ((raw: string) => string)[] = []

  /* Event Manager */
  private events = new EventManager<EventPayloadMap>()

  public constructor(conf: ArtalkConfig, $root?: HTMLElement) {
    this.conf = conf

    this.$root = $root || document.createElement('div')
    this.$root.classList.add('artalk')
    this.$root.innerHTML = ''

    this.data = new DataManager(this.events)
  }

  public inject(depName: string, obj: any) {
    this[depName] = obj
  }

  public get(depName: string) {
    return this[depName]
  }

  public getApi() {
    return this.api
  }

  public getData() {
    return this.data
  }

  public replyComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setReply(commentData, $comment)
  }

  public editComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setEditComment(commentData, $comment)
  }

  public fetch(params: Partial<ListFetchParams>): void {
    this.data.fetchComments(params)
  }

  public reload(): void {
    this.data.fetchComments({
      offset: 0,
    })
  }

  public listGotoFirst(): void {
    this.trigger('list-goto-first')
  }

  /* 编辑器 */
  public editorShowLoading(): void {
    this.editor.showLoading()
  }

  public editorHideLoading(): void {
    this.editor.hideLoading()
  }

  public editorShowNotify(msg, type): void {
    this.editor.showNotify(msg, type)
  }

  public editorResetState(): void {
    this.editor.resetState()
  }

  /* 侧边栏 */
  public showSidebar(payload?: SidebarShowPayload): void {
    this.sidebarLayer.show(payload)
  }

  public hideSidebar(): void {
    this.sidebarLayer.hide()
  }

  /* 权限检测 */
  public checkAdmin(payload: CheckerPayload): void {
    this.checkerLauncher.checkAdmin(payload)
  }

  public checkCaptcha(payload: CheckerCaptchaPayload): void {
    this.checkerLauncher.checkCaptcha(payload)
  }

  public checkAdminShowEl() {
    const items: HTMLElement[] = []

    this.$root.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

    // for layer
    const { $wrap: $layerWrap } = getLayerWrap()
    if ($layerWrap) $layerWrap.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

    // for sidebar
    // TODO: 这个其实应该写在 packages/artalk-sidebar 里面的
    const $sidebarEl = document.querySelector<HTMLElement>('.atk-sidebar')
    if ($sidebarEl) $sidebarEl.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

    items.forEach(($item: HTMLElement) => {
      if (this.user.data.isAdmin) $item.classList.remove('atk-hide')
      else $item.classList.add('atk-hide')
    })
  }

  /* Events */
  public on(name: any, handler: any) {
    this.events.on(name, handler)
  }

  public off(name: any, handler: any) {
    this.events.off(name, handler)
  }

  public trigger(name: any, payload?: any) {
    this.events.trigger(name, payload)
  }

  /* i18n */
  public $t(key: I18n.I18nKeys, args: {[key: string]: string} = {}): string {
    return I18n.t(key, args)
  }

  public setDarkMode(darkMode: boolean): void {
    DarkMode.setDarkMode(this, darkMode)
  }

  public updateConf(nConf: Partial<ArtalkConfig>): void {
    this.conf = Utils.mergeDeep(this.conf, handelBaseConf(nConf))
    this.trigger('conf-loaded', this.conf)
  }

  public getMarkedInstance() {
    return marked.getInstance()
  }
}

export default Context

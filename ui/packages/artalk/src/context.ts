import type ArtalkConfig from '~/types/artalk-config'
import type { CommentData, NotifyData, PageData } from '~/types/artalk-data'
import type { EventPayloadMap } from '~/types/event'
import type ContextApi from '~/types/context'
import type { TInjectedServices } from './service'

import * as Utils from './lib/utils'
import * as DarkMode from './lib/dark-mode'
import * as marked from './lib/marked'
import { CheckerCaptchaPayload, CheckerPayload } from './lib/checker'

import * as I18n from './i18n'
import { getLayerWrap } from './layer'
import { SidebarShowPayload } from './layer/sidebar-layer'
import Comment from './comment'
import EventManager from './lib/event-manager'

// Auto dependency injection
interface Context extends TInjectedServices { }

/**
 * Artalk Context
 */
class Context implements ContextApi {
  /* 运行参数 */
  public conf: ArtalkConfig
  public $root: HTMLElement
  public markedReplacers: ((raw: string) => string)[] = []

  private commentList: Comment[] = [] // Note: 无层级结构 + 无须排列
  private page?: PageData
  private unreadList: NotifyData[] = []

  /* Event Manager */
  private events = new EventManager<EventPayloadMap>()

  public constructor(conf: ArtalkConfig, $root?: HTMLElement) {
    this.conf = conf

    this.$root = $root || document.createElement('div')
    this.$root.classList.add('artalk')
    this.$root.innerHTML = ''
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

  /* 评论操作 */
  public getCommentList() {
    return this.commentList
  }

  public clearCommentList() {
    this.commentList = []
  }

  public getCommentDataList() {
    return this.commentList.map(c => c.getData())
  }

  public findComment(id: number): Comment|undefined {
    return this.commentList.find(c => c.getData().id === id)
  }

  public deleteComment(id: number) {
    this.list?.deleteComment(id)
  }

  public clearAllComments() {
    this.list?.clearAllComments()
  }

  public insertComment(commentData: CommentData) {
    this.list?.insertComment(commentData)
  }

  public updateComment(commentData: CommentData): void {
    this.list?.updateComment(commentData)
  }

  public replyComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setReply(commentData, $comment)
  }

  public editComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setEditComment(commentData, $comment)
  }

  /** 未读通知 */
  public getUnreadList() {
    return this.unreadList
  }

  public updateUnreadList(notifies: NotifyData[]): void {
    this.unreadList = notifies
    this.trigger('unread-updated', notifies)
  }

  /** 页面数据 */
  getPage(): PageData|undefined {
    return this.page
  }

  updatePage(pageData: PageData): void {
    this.page = pageData
    this.trigger('page-loaded', pageData)
  }

  /* 评论列表 */
  public listReload(): void {
    this.list?.reload()
  }

  public reload(): void {
    this.listReload()
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
    this.conf = Utils.mergeDeep(this.conf, nConf)
    this.trigger('conf-loaded')
  }

  public getMarkedInstance() {
    return marked.getInstance()
  }
}

export default Context

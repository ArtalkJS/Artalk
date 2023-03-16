import type ArtalkConfig from '~/types/artalk-config'
import type { CommentData, NotifyData } from '~/types/artalk-data'
import type { Event } from '~/types/event'
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
import Api from './api'
import List from './list'

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

  /* 订阅者模式 */
  private eventList: Event[] = []

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

  public getCommentDataList() {
    return this.commentList.map(c => c.getData())
  }

  public findComment(id: number): Comment|undefined {
    return this.commentList.find(c => c.getData().id === id)
  }

  public deleteComment(_comment: number|Comment) {
    let comment: Comment
    if (typeof _comment === 'number') {
      const findComment = this.findComment(_comment)
      if (!findComment) throw Error(`Comment ${_comment} cannot be found`)
      comment = findComment
    } else comment = _comment

    comment.getEl().remove()
    this.commentList.splice(this.commentList.indexOf(comment), 1)

    if (this.list) {
      const listData = this.list.getData()
      if (listData) listData.total -= 1 // 评论数减 1

      this.list.refreshUI()
    }
  }

  public clearAllComments() {
    if (this.list) {
      this.list.getCommentsWrapEl().innerHTML = ''
      this.list.clearData()
    }

    this.commentList = []
  }

  public insertComment(commentData: CommentData) {
    this.list?.insertComment(commentData)
  }

  public updateComment(commentData: CommentData): void {
    this.list?.updateComment(commentData)
  }

  public replyComment(commentData: CommentData, $comment: HTMLElement, scroll?: boolean): void {
    this.editor.setReply(commentData, $comment, scroll)
  }

  public cancelReplyComment(): void {
    this.editor.cancelReply()
  }

  public editComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setEditComment(commentData, $comment)
  }

  public cancelEditComment(): void {
    this.editor.cancelEditComment()
  }

  public updateNotifies(notifies: NotifyData[]): void {
    this.list?.updateUnread(notifies)
  }

  /* 评论列表 */
  public listReload(): void {
    this.list?.reload()
  }

  public reload(): void {
    this.listReload()
  }

  public listRefreshUI(): void {
    this.list?.refreshUI()
  }

  public listHashGotoCheck(): void {
    if (!this.list || !(this.list instanceof List)) return
    const list = this.list as List

    list.goToCommentDelay = false
    list.checkGoToCommentByUrlHash()
  }

  /* 编辑器 */
  public editorOpen(): void {
    this.editor.open()
  }

  public editorClose(): void {
    this.editor.close()
  }

  public editorShowLoading(): void {
    this.editor.showLoading()
  }

  public editorHideLoading(): void {
    this.editor.hideLoading()
  }

  public editorShowNotify(msg, type): void {
    this.editor.showNotify(msg, type)
  }

  public editorTravel($el: HTMLElement): void {
    this.editor.travel($el)
  }

  public editorTravelBack(): void {
    this.editor.travelBack()
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

  /* 订阅模式 */
  public on(name: any, handler: any, scope: any = 'internal') {
    this.eventList.push({ name, handler, scope })
  }

  public off(name: any, handler: any, scope: any = 'internal') {
    this.eventList = this.eventList.filter((evt) => {
      if (handler) return !(evt.name === name && evt.handler === handler && evt.scope === scope)
      return !(evt.name === name && evt.scope === scope) // 删除全部相同 name event
    })
  }

  public trigger(name: any, payload?: any, scope?: any) {
    this.eventList
      .filter((evt) => evt.name === name && (scope ? (evt.scope === scope) : true))
      .map((evt) => evt.handler)
      .forEach((handler) => handler(payload))
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

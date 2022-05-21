import { marked as libMarked } from 'marked'
import ArtalkConfig from '~/types/artalk-config'
import { CommentData, NotifyData } from '~/types/artalk-data'
import { Event } from '~/types/event'
import { internal as internalLocales, I18n } from './i18n'
import User from './lib/user'
import ContextApi from '../types/context'
import Editor from './editor'
import Comment from './comment'
import ListLite from './list/list-lite'
import SidebarLayer, { SidebarShowPayload } from './layer/sidebar-layer'
import CheckerLauncher, { CheckerCaptchaPayload, CheckerPayload } from './lib/checker'
import Layer, { GetLayerWrap } from './layer'
import Api from './api'

/**
 * Artalk Context
 */
export default class Context implements ContextApi {
  /* 持有同事类 */
  private api!: Api
  private editor!: Editor
  private list!: ListLite
  private sidebarLayer!: SidebarLayer
  private checkerLauncher!: CheckerLauncher

  /* 运行参数 */
  public cid: number // Context 唯一标识
  public conf: ArtalkConfig
  public user: User
  public $root: HTMLElement
  public markedInstance?: typeof libMarked
  public markedReplacers: ((raw: string) => string)[] = []

  private commentList: Comment[] = [] // Note: 无层级结构 + 无须排列

  /* 订阅者模式 */
  private eventList: Event[] = []

  public constructor(conf: ArtalkConfig, $root?: HTMLElement) {
    this.cid = +new Date()
    this.conf = conf
    this.user = new User(this)

    this.$root = $root || document.createElement('div')
    this.$root.setAttribute('atk-run-id', this.cid.toString())
    this.$root.classList.add('artalk')
    this.$root.innerHTML = ''

    this.api = new Api(this)
  }

  /* 设置持有的同事类 */
  public setApi(api: Api): void {
    this.api = api
  }

  public setEditor(editor: Editor): void {
    this.editor = editor
  }

  public setList(list: ListLite): void {
    this.list = list
  }

  public setSidebarLayer(sidebarLayer: SidebarLayer): void {
    this.sidebarLayer = sidebarLayer
  }

  public setCheckerLauncher(checkerLauncher: CheckerLauncher): void {
    this.checkerLauncher = checkerLauncher
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

    const listData = this.list.getData()
    if (listData) listData.total -= 1 // 评论数减 1

    this.list.refreshUI()
  }

  public clearAllComments() {
    this.list.getCommentsWrapEl().innerHTML = ''
    this.list.clearData()
    this.commentList = []
  }

  public insertComment(commentData: CommentData) {
    this.list.insertComment(commentData)
  }

  public replyComment(commentData: CommentData, $comment: HTMLElement, scroll?: boolean): void {
    this.editor.setReply(commentData, $comment, scroll)
  }

  public cancelReplyComment(): void {
    this.editor.cancelReply()
  }

  public updateNotifies(notifies: NotifyData[]): void {
    this.list.updateUnread(notifies)
  }

  /* 评论列表 */
  public listReload(): void {
    this.list.reload()
  }

  public listRefreshUI(): void {
    this.list.refreshUI()
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
    const { $wrap: $layerWrap } = GetLayerWrap(this)
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
  public $t(key: keyof I18n, args: {[key: string]: string} = {}): string {
    let locales = this.conf.locale
    if (typeof locales === 'string') {
      locales = internalLocales[locales]
    }

    let str = locales?.[key] || key
    str = str.replace(/\{\s*(\w+?)\s*\}/g, (_, token) => args[token] || '')

    return str
  }
}

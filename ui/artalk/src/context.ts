import type { ArtalkConfig, CommentData, ListFetchParams, ContextApi, EventPayloadMap, SidebarShowPayload } from '@/types'
import type { TInjectedServices } from './service'
import { Api } from './api'
import type CommentNode from './comment'

import * as marked from './lib/marked'
import { mergeDeep } from './lib/merge-deep'
import { CheckerCaptchaPayload, CheckerPayload } from './components/checker'

import { DataManager } from './data'
import * as I18n from './i18n'

import EventManager from './lib/event-manager'
import { convertApiOptions, handelCustomConf } from './config'

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

  /* Event Manager */
  private events = new EventManager<EventPayloadMap>()

  public constructor(conf: ArtalkConfig) {
    this.conf = conf

    this.$root = conf.el as HTMLElement
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
    return new Api(convertApiOptions(this.conf, this))
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
    this.data.fetchComments({ offset: 0 })
  }

  public listGotoFirst(): void {
    this.events.trigger('list-goto-first')
  }

  public getCommentNodes(): CommentNode[] {
    return this.list.getCommentNodes()
  }

  public getComments(): CommentData[] {
    return this.data.getComments()
  }

  public getCommentList = this.getCommentNodes
  public getCommentDataList = this.getComments

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
  public checkAdmin(payload: CheckerPayload): Promise<void> {
    return this.checkerLauncher.checkAdmin(payload)
  }

  public checkCaptcha(payload: CheckerCaptchaPayload): Promise<void> {
    return this.checkerLauncher.checkCaptcha(payload)
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

  public setDarkMode(darkMode: boolean|'auto'): void {
    // prevent trigger 'conf-loaded' to improve performance
    // this.updateConf({ ...this.conf, darkMode })
    this.conf.darkMode = darkMode
    this.events.trigger('dark-mode-changed', darkMode)
  }

  public updateConf(nConf: Partial<ArtalkConfig>): void {
    this.conf = mergeDeep(this.conf, handelCustomConf(nConf, false))
    this.events.trigger('conf-loaded', this.conf)
  }

  public getConf(): ArtalkConfig {
    return this.conf
  }

  public getEl(): HTMLElement {
    return this.$root
  }

  public getMarked() {
    return marked.getInstance()
  }
}

export default Context

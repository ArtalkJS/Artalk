import type {
  ArtalkConfig,
  CommentData,
  ListFetchParams,
  ContextApi,
  EventPayloadMap,
  SidebarShowPayload,
} from '@/types'
import type { TInjectedServices } from './service'
import { Api, ApiHandlers } from './api'

import * as marked from './lib/marked'
import { mergeDeep } from './lib/merge-deep'
import { CheckerCaptchaPayload, CheckerPayload } from './components/checker'

import { DataManager } from './data'
import * as I18n from './i18n'

import EventManager from './lib/event-manager'
import { convertApiOptions, createNewApiHandlers, handelCustomConf } from './config'
import { watchConf } from './lib/watch-conf'

// Auto dependency injection
interface Context extends TInjectedServices {}

/**
 * Artalk Context
 */
class Context implements ContextApi {
  /* 运行参数 */
  conf: ArtalkConfig
  data: DataManager
  $root: HTMLElement

  /* Event Manager */
  private events = new EventManager<EventPayloadMap>()
  private mounted = false

  constructor(conf: ArtalkConfig) {
    this.conf = conf

    this.$root = conf.el as HTMLElement
    this.$root.classList.add('artalk')
    this.$root.innerHTML = ''

    this.data = new DataManager(this.events)

    this.on('mounted', () => {
      this.mounted = true
    })
  }

  inject(depName: string, obj: any) {
    this[depName] = obj
  }

  get(depName: string) {
    return this[depName]
  }

  getApi() {
    return new Api(convertApiOptions(this.conf, this))
  }

  private apiHandlers = <ApiHandlers | null>null
  getApiHandlers() {
    if (!this.apiHandlers) this.apiHandlers = createNewApiHandlers(this)
    return this.apiHandlers
  }

  getData() {
    return this.data
  }

  replyComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setReply(commentData, $comment)
  }

  editComment(commentData: CommentData, $comment: HTMLElement): void {
    this.editor.setEditComment(commentData, $comment)
  }

  fetch(params: Partial<ListFetchParams>): void {
    this.data.fetchComments(params)
  }

  reload(): void {
    this.data.fetchComments({ offset: 0 })
  }

  /* List */
  listGotoFirst(): void {
    this.events.trigger('list-goto-first')
  }

  getCommentNodes() {
    return this.list.getCommentNodes()
  }

  getComments() {
    return this.data.getComments()
  }

  getCommentList = this.getCommentNodes
  getCommentDataList = this.getComments

  /* Editor */
  editorShowLoading(): void {
    this.editor.showLoading()
  }

  editorHideLoading(): void {
    this.editor.hideLoading()
  }

  editorShowNotify(msg, type): void {
    this.editor.showNotify(msg, type)
  }

  editorResetState(): void {
    this.editor.resetState()
  }

  /* Sidebar */
  showSidebar(payload?: SidebarShowPayload): void {
    this.sidebarLayer.show(payload)
  }

  hideSidebar(): void {
    this.sidebarLayer.hide()
  }

  /* Checker */
  checkAdmin(payload: CheckerPayload): Promise<void> {
    return this.checkerLauncher.checkAdmin(payload)
  }

  checkCaptcha(payload: CheckerCaptchaPayload): Promise<void> {
    return this.checkerLauncher.checkCaptcha(payload)
  }

  /* Events */
  on(name: any, handler: any) {
    this.events.on(name, handler)
  }

  off(name: any, handler: any) {
    this.events.off(name, handler)
  }

  trigger(name: any, payload?: any) {
    this.events.trigger(name, payload)
  }

  /* i18n */
  $t(key: I18n.I18nKeys, args: { [key: string]: string } = {}): string {
    return I18n.t(key, args)
  }

  setDarkMode(darkMode: boolean | 'auto'): void {
    this.updateConf({ darkMode })
  }

  updateConf(nConf: Partial<ArtalkConfig>): void {
    this.conf = mergeDeep(this.conf, handelCustomConf(nConf, false))
    this.mounted && this.events.trigger('updated', this.conf)
  }

  getConf(): ArtalkConfig {
    return this.conf
  }

  getEl(): HTMLElement {
    return this.$root
  }

  getMarked() {
    return marked.getInstance()
  }

  watchConf<T extends (keyof ArtalkConfig)[]>(
    keys: T,
    effect: (conf: Pick<ArtalkConfig, T[number]>) => void,
  ): void {
    watchConf(this, keys, effect)
  }
}

export default Context

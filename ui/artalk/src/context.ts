import type {
  Config,
  ConfigPartial,
  CommentData,
  ListFetchParams,
  Context as IContext,
  SidebarShowPayload,
  Services,
} from './types'
import * as I18n from './i18n'
import * as marked from './lib/marked'
import { createInjectionContainer } from './lib/injection'
import type { CheckerCaptchaPayload, CheckerPayload } from './components/checker'

/**
 * Artalk Context
 */
class Context implements IContext {
  private _deps = createInjectionContainer<Services>()

  constructor(private _$root: HTMLElement) {}

  getEl(): HTMLElement {
    return this._$root
  }

  destroy(): void {
    this.trigger('unmounted')
    while (this._$root.firstChild) {
      this._$root.removeChild(this._$root.firstChild)
    }
  }

  // -------------------------------------------------------------------
  //  Dependency Injection
  // -------------------------------------------------------------------
  provide: IContext['provide'] = (key, impl, deps, opts) => {
    this._deps.provide(key, impl, deps, opts)
  }

  inject: IContext['inject'] = (key) => {
    return this._deps.inject(key)
  }
  get = this.inject

  // -------------------------------------------------------------------
  //  Event Manager
  // -------------------------------------------------------------------
  on: IContext['on'] = (name, handler) => {
    this.inject('events').on(name, handler)
  }

  off: IContext['off'] = (name, handler) => {
    this.inject('events').off(name, handler)
  }

  trigger: IContext['trigger'] = (name, payload) => {
    this.inject('events').trigger(name, payload)
  }

  // -------------------------------------------------------------------
  //  Configurations
  // -------------------------------------------------------------------
  getConf(): Config {
    return this.inject('config').get()
  }

  updateConf(conf: ConfigPartial): void {
    this.inject('config').update(conf)
  }

  watchConf<T extends (keyof Config)[]>(
    keys: T,
    effect: (conf: Pick<Config, T[number]>) => void,
  ): void {
    this.inject('config').watchConf(keys, effect)
  }

  getMarked() {
    return marked.getInstance()
  }

  setDarkMode(darkMode: boolean | 'auto'): void {
    this.updateConf({ darkMode })
  }

  get conf() {
    return this.getConf()
  }
  set conf(val) {
    console.error('Cannot set config directly, please call updateConf()')
  }
  get $root() {
    return this.getEl()
  }
  set $root(val) {
    console.error('set $root is prohibited')
  }

  // -------------------------------------------------------------------
  //  I18n: Internationalization
  // -------------------------------------------------------------------
  $t(key: I18n.I18nKeys, args: { [key: string]: string } = {}): string {
    return I18n.t(key, args)
  }

  // -------------------------------------------------------------------
  //  HTTP API Client
  // -------------------------------------------------------------------
  getApi() {
    return this.inject('api')
  }

  getApiHandlers() {
    return this.inject('apiHandlers')
  }

  // -------------------------------------------------------------------
  //  User Manager
  // -------------------------------------------------------------------
  getUser() {
    return this.inject('user')
  }

  // -------------------------------------------------------------------
  //  Data Manager
  // -------------------------------------------------------------------
  getData() {
    return this.inject('data')
  }

  fetch(params: Partial<ListFetchParams>): void {
    this.getData().fetchComments(params)
  }

  reload(): void {
    this.getData().fetchComments({ offset: 0 })
  }

  // -------------------------------------------------------------------
  //  List
  // -------------------------------------------------------------------
  listGotoFirst(): void {
    this.trigger('list-goto-first')
  }

  getCommentList = this.getCommentNodes
  getCommentNodes() {
    return this.inject('list').getCommentNodes()
  }

  getCommentDataList = this.getComments
  getComments() {
    return this.getData().getComments()
  }

  // -------------------------------------------------------------------
  //  Editor
  // -------------------------------------------------------------------
  replyComment(commentData: CommentData, $comment: HTMLElement): void {
    this.inject('editor').setReplyComment(commentData, $comment)
  }

  editComment(commentData: CommentData, $comment: HTMLElement): void {
    this.inject('editor').setEditComment(commentData, $comment)
  }

  editorShowLoading(): void {
    this.inject('editor').showLoading()
  }

  editorHideLoading(): void {
    this.inject('editor').hideLoading()
  }

  editorShowNotify(msg, type): void {
    this.inject('editor').showNotify(msg, type)
  }

  editorResetState(): void {
    this.inject('editor').resetState()
  }

  // -------------------------------------------------------------------
  //  Sidebar
  // -------------------------------------------------------------------
  showSidebar(payload?: SidebarShowPayload): void {
    this.inject('sidebar').show(payload)
  }

  hideSidebar(): void {
    this.inject('sidebar').hide()
  }

  // -------------------------------------------------------------------
  //  Checker
  // -------------------------------------------------------------------
  checkAdmin(payload: CheckerPayload): Promise<void> {
    return this.inject('checkers').checkAdmin(payload)
  }

  checkCaptcha(payload: CheckerCaptchaPayload): Promise<void> {
    return this.inject('checkers').checkCaptcha(payload)
  }
}

export default Context

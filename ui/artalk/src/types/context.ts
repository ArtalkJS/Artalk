import type { Marked } from 'marked'
import type {
  SidebarShowPayload,
  EventPayloadMap,
  ConfigPartial,
  Config,
  CommentData,
  DataManager,
  ListFetchParams,
  NotifyLevel,
  Services,
  EventManager,
  UserManager,
} from '.'
import type { CheckerCaptchaPayload, CheckerPayload } from '@/components/checker'
import type { DependencyContainer } from '@/lib/injection'
import type { I18n } from '@/i18n'
import type { Api, ApiHandlers } from '@/api'
import type { CommentNode } from '@/comment'

/**
 * Artalk Context
 */
export interface Context extends EventManager<EventPayloadMap>, DependencyContainer<Services> {
  /**
   * The root element of Artalk
   *
   * @deprecated Use `getEl()` instead
   */
  $root: HTMLElement

  /** Get the root element */
  getEl(): HTMLElement

  /**
   * Inject a dependency object
   *
   * @deprecated Use `inject()` instead
   */
  get<T extends keyof Services>(key: T): Services[T]

  /**
   * Get config object
   *
   * @deprecated Use `getConf()` and `updateConf()` instead
   */
  conf: Config

  /** Get the config */
  getConf(): Config

  /** Update the config */
  updateConf(conf: ConfigPartial): void

  /** Watch the config */
  watchConf<T extends (keyof Config)[]>(
    keys: T,
    effect: (val: Pick<Config, T[number]>) => void,
  ): void

  /** Get the marked instance */
  getMarked(): Marked | undefined

  /** Set dark mode */
  setDarkMode(darkMode: boolean | 'auto'): void

  /** Translate i18n message */
  $t(key: keyof I18n, args?: { [key: string]: string }): string

  /** Get HTTP API client */
  getApi(): Api

  /** Get HTTP API handlers */
  getApiHandlers(): ApiHandlers

  /** Get Data Manager */
  getData(): DataManager

  /** Get User Manager */
  getUser(): UserManager

  /** Fetch comments */
  fetch(params: Partial<ListFetchParams>): void

  /** Reload comments */
  reload(): void

  /** Destroy */
  destroy(): void

  /** Goto the first comment of the list */
  listGotoFirst(): void

  /** Get the comment data list */
  getComments(): CommentData[]

  /** Get the comment node list */
  getCommentNodes(): CommentNode[]

  /**
   * Get the comment data list
   * @deprecated Use `getComments()` instead
   */
  getCommentDataList(): CommentData[]

  /**
   * Get the comment node list
   * @deprecated Use `getCommentNodes()` instead
   */
  getCommentList(): CommentNode[]

  /** Reply to a comment */
  replyComment(commentData: CommentData, $comment: HTMLElement): void

  /** Edit a comment */
  editComment(commentData: CommentData, $comment: HTMLElement): void

  /** Show loading of the editor */
  editorShowLoading(): void

  /** Hide loading of the editor */
  editorHideLoading(): void

  /** Show notify of the editor */
  editorShowNotify(msg: string, type: NotifyLevel): void

  /** Reset the state of the editor */
  editorResetState(): void

  /** Show the sidebar */
  showSidebar(payload?: SidebarShowPayload): void

  /** Hide the sidebar */
  hideSidebar(): void

  /** Check captcha */
  checkCaptcha(payload: CheckerCaptchaPayload): Promise<void>

  /** Check admin */
  checkAdmin(payload: CheckerPayload): Promise<void>
}

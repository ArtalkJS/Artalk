import type { TInjectedServices } from '@/service'
import type { CheckerCaptchaPayload, CheckerPayload } from '@/components/checker'
import type { EventManagerFuncs } from '@/lib/event-manager'
import type { I18n } from '@/i18n'
import type { Api, ApiHandlers } from '@/api'
import type { CommentNode } from '@/comment'
import type {
  SidebarShowPayload,
  EventPayloadMap,
  ArtalkConfig,
  CommentData,
  DataManagerApi,
  ListFetchParams,
  NotifyLevel,
} from '.'

/**
 * Artalk Context
 */
export interface ContextApi extends EventManagerFuncs<EventPayloadMap> {
  /** Artalk 根元素对象 */
  $root: HTMLElement

  /** 依赖注入函数 */
  inject<K extends keyof TInjectedServices>(depName: K, obj: TInjectedServices[K]): void

  /** 获取依赖对象 */
  get<K extends keyof TInjectedServices>(depName: K): TInjectedServices[K]

  /** 配置对象 */
  // TODO: 修改为 getConf() 和 setConf() 并且返回拷贝而不是引用
  conf: ArtalkConfig

  /** marked 依赖对象 */
  getMarked(): any | undefined

  /** 获取 API 以供 HTTP 请求 */
  getApi(): Api

  /** Get API handlers */
  getApiHandlers(): ApiHandlers

  /** 获取数据管理器对象 */
  getData(): DataManagerApi

  /** 评论回复 */
  replyComment(commentData: CommentData, $comment: HTMLElement): void

  /** 编辑评论 */
  editComment(commentData: CommentData, $comment: HTMLElement): void

  /** 获取评论数据 */
  fetch(params: Partial<ListFetchParams>): void

  /** 重载评论数据 */
  reload(): void

  /** 列表滚动到第一个评论的位置 */
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

  /** 显示侧边栏 */
  showSidebar(payload?: SidebarShowPayload): void

  /** 隐藏侧边栏 */
  hideSidebar(): void

  /** 编辑器 - 显示加载 */
  editorShowLoading(): void

  /** 编辑器 - 隐藏加载 */
  editorHideLoading(): void

  /** 编辑器 - 显示提示消息 */
  editorShowNotify(msg: string, type: NotifyLevel): void

  /** 评论框 - 复原状态 */
  editorResetState(): void

  /** 验证码检测 */
  checkCaptcha(payload: CheckerCaptchaPayload): Promise<void>

  /** 管理员检测 */
  checkAdmin(payload: CheckerPayload): Promise<void>

  /** i18n 翻译 */
  $t(key: keyof I18n, args?: { [key: string]: string }): string

  /** 设置夜间模式 */
  setDarkMode(darkMode: boolean | 'auto'): void

  /** 获取配置 */
  getConf(): ArtalkConfig

  /** 获取挂载元素 */
  getEl(): HTMLElement

  /** 更新配置 */
  updateConf(conf: Partial<ArtalkConfig>): void

  /** 监听配置更新 */
  watchConf<T extends (keyof ArtalkConfig)[]>(
    keys: T,
    effect: (val: Pick<ArtalkConfig, T[number]>) => void,
  ): void
}

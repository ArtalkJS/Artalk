import type { TInjectedServices } from '@/service'
import type { CheckerCaptchaPayload, CheckerPayload } from '@/components/checker'
import type { EventManagerFuncs } from '@/lib/event-manager'
import type { TMarked } from '@/lib/marked'
import type { I18n } from '@/i18n'
import type Api from '@/api'
import type { SidebarShowPayload, EventPayloadMap, ArtalkConfig, CommentData, DataManagerApi, ListFetchParams, NotifyLevel } from '.'

/**
 * Context 接口
 *
 * (面向接口的编程)
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
  getMarked(): TMarked | undefined

  /** 获取 API 以供 HTTP 请求 */
  getApi(): Api

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
  checkCaptcha(payload: CheckerCaptchaPayload): void

  /** 管理员检测 */
  checkAdmin(payload: CheckerPayload): void

  /** i18n 翻译 */
  $t(key: keyof I18n, args?: {[key: string]: string}): string

  /** 设置夜间模式 */
  setDarkMode(darkMode: boolean|'auto'): void

  /** 更新配置 */
  updateConf(conf: Partial<ArtalkConfig>): void
}

import type { MarkedOptions } from 'marked'
import type { CommentData } from './data'
import type { Editor } from './editor'
import type { Context } from './context'
import type { I18n } from '@/i18n'

export interface Config {
  /** Element selector or Element to mount the Artalk */
  el: string | HTMLElement

  /** Unique page identifier */
  pageKey: string

  /** Title of the page */
  pageTitle: string

  /** Server address */
  server: string

  /** Site name */
  site: string

  /** Placeholder text for the comment input box */
  placeholder: string

  /** Text to display when there are no comments */
  noComment: string

  /** Text for the send button */
  sendBtn: string

  /** Movable comment box (display below the comment to be replied) */
  editorTravel: boolean

  /** Emoticons settings */
  emoticons: object | any[] | string | false

  /** Gravatar avatar settings */
  gravatar: {
    /** API endpoint */
    mirror: string
    /** API parameters */
    params: string
  }

  /** Avatar URL generator function */
  avatarURLBuilder?: (comment: CommentData) => string

  /** Pagination settings */
  pagination: {
    /** Number of comments to fetch per request */
    pageSize: number
    /** "Read more" mode */
    readMore: boolean
    /** Automatically load more comments when scrolled to the bottom */
    autoLoad: boolean
  }

  /** Height limit configuration */
  heightLimit: {
    /** Maximum height for comment content */
    content: number

    /** Maximum height for child comments */
    children: number

    /** Whether the content is scrollable */
    scrollable: boolean
  }

  /** Voting feature for comments */
  vote: boolean

  /** Downvote button for comments */
  voteDown: boolean

  /** Page Vote Widget */
  pageVote:
    | {
        /** Up Vote Button Selector */
        upBtnEl: string

        /** Down Vote Button Selector */
        downBtnEl: string

        /** Up Vote Count Selector */
        upCountEl: string

        /** Down Vote Count Selector */
        downCountEl: string

        /** Active class name if the vote is already cast */
        activeClass: string
      }
    | boolean

  /** Preview feature for comments */
  preview: boolean

  /** Selector for the element binding to the comment count */
  countEl: string

  /** Selector for the element binding to the page views (PV) count */
  pvEl: string

  /** Attribute name for the PageKey in statistics components */
  statPageKeyAttr: string

  /** Dark mode settings */
  darkMode: boolean | 'auto'

  /** Request timeout (in seconds) */
  reqTimeout: number

  /** Flat mode for comment display */
  flatMode: boolean | 'auto'

  /** Maximum number of levels for nested comments */
  nestMax: number

  /** Sorting order for nested comments */
  nestSort: 'DATE_ASC' | 'DATE_DESC'

  /** Display UA badge (user agent badge) */
  uaBadge: boolean

  /** Show sorting dropdown for comment list */
  listSort: boolean

  /** Enable image upload feature */
  imgUpload: boolean

  /** Image uploader function */
  imgUploader?: (file: File) => Promise<string>

  /** Image lazy load mode */
  imgLazyLoad: false | 'native' | 'data-src'

  /** Enable version check */
  versionCheck: boolean

  /**
   * Use remote configuration (from the backend server)
   *
   * @deprecated
   * The `useBackendConf` is always `true` and planned to be removed in the future.
   *
   * Please use `preferRemoteConf` to control the priority of the remote and local configuration.
   *
   * @default true
   */
  useBackendConf: boolean

  /**
   * Prefer to use the local configuration if available
   *
   * @note
   * If `true`, the local config will be used as a fallback (remote config first).
   *
   * If `false`, the local config will override the remote config (local config first).
   *
   * @default false
   */
  preferRemoteConf: boolean

  /** Localization settings */
  locale: I18n | string

  /** Backend API version (system data, not allowed for user modification) */
  apiVersion: string

  /** URLs for plugin scripts */
  pluginURLs?: string[]

  /** Replacers for the marked (Markdown parser) */
  markedReplacers?: ((raw: string) => string)[]

  /** Options for the marked (Markdown parser) */
  markedOptions?: MarkedOptions

  /** Modifier for list fetch request parameters */
  listFetchParamsModifier?: (params: any) => void

  /**
   * Custom date formatter
   * @param date - The Date object to format
   * @returns Formatted date string
   */
  dateFormatter?: (date: Date) => string

  /** List unread highlight (enable by default in artalk-sidebar) */
  listUnreadHighlight: boolean

  /** The relative element for scrolling (useful if artalk is in a scrollable container) */
  scrollRelativeTo?: () => HTMLElement

  /** Page view increment when comment list is loaded */
  pvAdd: boolean

  /** Immediately fetch comments when Artalk instance is initialized */
  fetchCommentsOnInit: boolean

  /** Callback before submitting a comment */
  beforeSubmit?: (editor: Editor, next: () => void) => void
}

type DeepPartial<T> = {
  [K in keyof T]?: T[K] extends object ? DeepPartial<T[K]> : T[K]
}

export type ConfigPartial = DeepPartial<Config>

// Alias for backward compatibility
export type ArtalkConfig = Config

export interface ConfigManager {
  watchConf: Context['watchConf']
  get: () => Config
  update: (conf: ConfigPartial) => void
}

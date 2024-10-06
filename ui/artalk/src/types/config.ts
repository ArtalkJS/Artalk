import type { MarkedOptions } from 'marked'
import type { CommentData } from './data'
import type { EditorApi } from './editor'
import type { I18n } from '@/i18n'

export interface ArtalkConfig {
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
  imgLazyLoad?: 'native' | 'data-src'

  /** Enable version check */
  versionCheck: boolean

  /** Use backend configuration */
  useBackendConf: boolean

  /** Localization settings */
  locale: I18n | string

  /** Backend API version (system data, not allowed for user modification) */
  apiVersion?: string

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

  /** Custom configuration modifier for remote configuration (referenced by artalk-sidebar) */
  // TODO: Consider merging list-related configuration into a single object, or flatten everything for simplicity and consistency
  remoteConfModifier?: (conf: DeepPartial<ArtalkConfig>) => void

  /** List unread highlight (enable by default in artalk-sidebar) */
  listUnreadHighlight?: boolean

  /** The relative element for scrolling (useful if artalk is in a scrollable container) */
  scrollRelativeTo?: () => HTMLElement

  /** Page view increment when comment list is loaded */
  pvAdd?: boolean

  /** Callback before submitting a comment */
  beforeSubmit?: (editor: EditorApi, next: () => void) => void
}

type DeepPartial<T> = {
  [K in keyof T]?: T[K] extends object ? DeepPartial<T[K]> : T[K]
}

export type ArtalkConfigPartial = DeepPartial<ArtalkConfig>

/**
 * Local User Data (in localStorage)
 *
 * @note Keep flat for easy handling
 */
export interface LocalUser {
  /** Username (aka. Nickname) */
  name: string

  /** Email */
  email: string

  /** Link (aka. Website) */
  link: string

  /** Token (for authorization) */
  token: string

  /** Admin flag */
  is_admin: boolean
}

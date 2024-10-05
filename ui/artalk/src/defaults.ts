import type { ArtalkConfig } from '@/types'

type RequiredExcept<T, K extends keyof T> = Required<Omit<T, K>> & Pick<T, K>
type FunctionKeys<T> = Exclude<
  { [K in keyof T]: NonNullable<T[K]> extends (...args: any[]) => any ? K : never }[keyof T],
  undefined
>
type ExcludedKeys = FunctionKeys<ArtalkConfig>

const defaults: RequiredExcept<ArtalkConfig, ExcludedKeys> = {
  el: '',
  pageKey: '',
  pageTitle: '',
  server: '',
  site: '',

  placeholder: '',
  noComment: '',
  sendBtn: '',
  darkMode: false,
  editorTravel: true,

  flatMode: 'auto',
  nestMax: 2,
  nestSort: 'DATE_ASC',

  emoticons: 'https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json',

  vote: true,
  voteDown: false,
  uaBadge: true,
  listSort: true,
  preview: true,
  countEl: '.artalk-comment-count',
  pvEl: '.artalk-pv-count',
  statPageKeyAttr: 'data-page-key',

  gravatar: {
    mirror: 'https://www.gravatar.com/avatar/',
    params: 'sha256=1&d=mp&s=240',
  },

  pagination: {
    pageSize: 20,
    readMore: true,
    autoLoad: true,
  },

  heightLimit: {
    content: 300,
    children: 400,
    scrollable: false,
  },

  pvAdd: true,
  imgUpload: true,
  imgLazyLoad: 'native',
  reqTimeout: 15000,
  versionCheck: true,
  useBackendConf: true,
  listUnreadHighlight: false,

  locale: 'en',
  apiVersion: '',
  pluginURLs: [],
  markedReplacers: [],
  markedOptions: {},
}

if (ARTALK_LITE) {
  defaults.vote = false
  defaults.uaBadge = false
  defaults.emoticons = false
  defaults.preview = false
}

export default defaults

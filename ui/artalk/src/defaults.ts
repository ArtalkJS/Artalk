import type { Config } from '@/types'

export const Defaults: Readonly<RequiredExcept<Config, ExcludedKeys>> = {
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

  emoticons: ARTALK_LITE
    ? false
    : 'https://cdn.jsdelivr.net/gh/ArtalkJS/Emoticons/grps/default.json',

  pageVote: true,

  vote: ARTALK_LITE ? false : true,
  voteDown: false,
  uaBadge: ARTALK_LITE ? false : true,
  listSort: true,
  preview: ARTALK_LITE ? false : true,
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

  imgUpload: true,
  imgLazyLoad: false,
  reqTimeout: 15000,
  versionCheck: true,
  useBackendConf: true,
  preferRemoteConf: false,
  listUnreadHighlight: false,
  pvAdd: true,
  fetchCommentsOnInit: true,

  locale: 'en',
  apiVersion: '',
  pluginURLs: [],
  markedReplacers: [],
  markedOptions: {},
}

type RequiredExcept<T, K extends keyof T> = Required<Omit<T, K>> & Pick<T, K>
type FunctionKeys<T> = Exclude<
  { [K in keyof T]: NonNullable<T[K]> extends (...args: any[]) => any ? K : never }[keyof T],
  undefined
>
type ExcludedKeys = FunctionKeys<Config>

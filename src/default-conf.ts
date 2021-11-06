import ArtalkConfig from "~/types/artalk-config";
import emoticons from './assets/emoticons.json'

const defaultConf: ArtalkConfig = {
  el: '',
  pageKey: '',
  server: '',
  site: '',

  placeholder: '键入内容...',
  noComment: '「此时无声胜有声」',
  sendBtn: '发送评论',
  darkMode: false,

  emoticons,

  vote: true,
  voteDown: false,
  uaBadge: true,
  flatMode: false,
  maxNesting: 3,
  gravatar: {
    default: 'mp',
    mirror: 'https://sdn.geekzu.org/avatar/',
  },

  readMore: {
    pageSize: 15,
    autoLoad: true,
  },

  reqTimeout: 15000,
  versionCheck: true,
}

export default defaultConf

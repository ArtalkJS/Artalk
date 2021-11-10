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
  flatMode: 'auto',
  maxNesting: 3,
  gravatar: {
    default: 'mp',
    mirror: 'https://sdn.geekzu.org/avatar/',
  },

  pagination: {
    pageSize: 15,
    readMore: true,
    autoLoad: true,
  },

  heightLimit: {
    content: 200,
    children: 300,
  },

  reqTimeout: 15000,
  versionCheck: true,
}

export default defaultConf

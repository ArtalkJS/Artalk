import ArtalkConfig from "~/types/artalk-config";
import emoticons from './assets/emoticons.json'

const defaultConf: ArtalkConfig = {
  el: '',
  placeholder: '来啊，快活啊 ( ゜- ゜)',
  noComment: '快来成为第一个评论的人吧~',
  sendBtn: '发送评论',
  defaultAvatar: 'mp',
  pageKey: '',
  server: '',
  site: '',
  emoticons,
  gravatar: {
    cdn: 'https://sdn.geekzu.org/avatar/'
  },
  voteDown: false,
  darkMode: false,
  reqTimeout: 15000,
  flatMode: false,
  maxNesting: 3,
  versionCheck: true,
}

export default defaultConf

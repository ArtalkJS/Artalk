import './css/main.less'
import marked from 'marked'
import hanabi from 'hanabi'
import Editor from './components/Editor'
import List from './components/List'
import Sidebar from './components/Sidebar'
import Ui from './utils/ui'
import { ArtalkConfig } from '~/types/artalk-config'

/* global ARTALK_VERSION */

const defaultOpts: ArtalkConfig = {
  el: '',
  placeholder: '来啊，快活啊 ( ゜- ゜)',
  noComment: '快来成为第一个评论的人吧~',
  sendBtn: '发送评论',
  defaultAvatar: 'mp',
  pageKey: '',
  serverUrl: '',
  emoticons: require('./assets/emoticons.json'),
  gravatar: {
    cdn: 'https://gravatar.loli.net/avatar/'
  }
}

// eslint-disable-next-line
export let ArtalkInstance: Artalk = null

export default class Artalk {
  public el: HTMLElement

  public ui: Ui
  public editor: Editor
  public list: List
  public sidebar: Sidebar

  public user: {
    nick: string|null,
    email: string|null,
    link: string|null,
    password: string|null,
    isAdmin: boolean
  }

  constructor (public conf: ArtalkConfig) {
    // Version Information
    console.log(`\n %c `
      + `Artalk v${ARTALK_VERSION} %c 一款简洁有趣的自托管评论系统 \n\n%c`
      + `> https://artalk.js.org\n`
      + `> https://github.com/qwqcode/Artalk\n`
      + `> https://qwqaq.com\n`,
      'color: #FFF; background: #1DAAFF; padding:5px 0;', 'color: #FFF; background: #656565; padding:5px 0;', '')

    ArtalkInstance = this

    // Options
    this.conf = Object.assign(defaultOpts, this.conf)

    // Main Element
    try {
      this.el = document.querySelector(this.conf.el)
      if (this.el === null) {
        throw Error(`Sorry, Target element "${this.conf.el}" was not found.`)
      }
    } catch (e) {
      console.error(e)
      throw new Error('Artalk config `el` error')
    }

    this.el.classList.add('artalk')

    // Components
    this.ui = new Ui()
    this.initUser()
    this.initMarked()
    this.editor = new Editor()
    this.list = new List()
    this.sidebar = new Sidebar()
  }

  /** 初始化用户数据 */
  initUser () {
    // 从 localStorage 导入
    const localUser = JSON.parse(window.localStorage.getItem('ArtalkUser') || '{}')
    this.user = {
      nick: localUser.nick || '',
      email: localUser.email || '',
      link: localUser.link || '',
      password: localUser.password || '',
      isAdmin: localUser.isAdmin || false
    }
  }

  /** 保存用户到 localStorage 中 */
  saveUser () {
    window.localStorage.setItem('ArtalkUser', JSON.stringify(this.user))
  }

  /** 是否已填写基本用户信息 */
  checkHasBasicUserInfo () {
    return !!this.user.nick && !!this.user.email
  }

  public marked: (src: string, callback?: (error: any, parseResult: string) => void) => string

  private initMarked () {
    const renderer = new marked.Renderer()
    const linkRenderer = renderer.link
    renderer.link = (href, title, text) => {
      const html = linkRenderer.call(renderer, href, title, text)
      return html.replace(/^<a /, '<a target="_blank" rel="nofollow" ')
    }

    const nMarked = marked
    nMarked.setOptions({
      renderer,
      highlight: (code) => {
        return hanabi(code)
      },
      pedantic: false,
      gfm: true,
      tables: true,
      breaks: true,
      sanitize: true, // 净化
      smartLists: true,
      smartypants: true,
      xhtml: false
    })

    this.marked = nMarked
  }

  public request (action: string, data: object, before: () => void, after: () => void, success: (msg: string, data: any) => void, error: (msg: string, data: any) => void) {
    before()

    data = { action, ...data }
    const formData = new FormData()
    Object.keys(data).forEach(key => formData.set(key, data[key]))

    const xhr = new XMLHttpRequest()
    xhr.timeout = 5000
    xhr.open('POST', this.conf.serverUrl, true)

    xhr.onload = () => {
      after()
      if (xhr.status >= 200 && xhr.status < 400) {
        const respData = JSON.parse(xhr.response)
        if (respData.success) {
          success(respData.msg, respData.data)
        } else {
          error(respData.msg, respData.data)
        }
      } else {
        error(`服务器响应错误 Code: ${xhr.status}`, {})
      }
    };

    xhr.onerror = () => {
      after()
      error('网络错误', {})
    };

    xhr.send(formData)
  }
}

import './style/main.less'
import Checker from './lib/checker'
import Editor from './components/Editor'
import List from './components/List'
import Sidebar from './components/Sidebar'
import { GetLayerWrap } from './components/Layer'
import * as Ui from './lib/ui'
import * as Utils from './lib/utils'
import { ArtalkConfig } from '~/types/artalk-config'
import Context from './Context'
import emoticons from './assets/emoticons.json'
import Constant from './Constant'

const defaultOpts: ArtalkConfig = {
  el: '',
  placeholder: '来啊，快活啊 ( ゜- ゜)',
  noComment: '快来成为第一个评论的人吧~',
  sendBtn: '发送评论',
  defaultAvatar: 'mp',
  pageKey: '',
  serverUrl: '',
  emoticons,
  gravatar: {
    cdn: 'https://sdn.geekzu.org/avatar/'
  },
  darkMode: false,
}

export default class Artalk {
  public conf: ArtalkConfig
  public ctx: Context
  public el: HTMLElement
  public readonly contextID: number = new Date().getTime() // 实例唯一 ID
  public checker: Checker
  public editor: Editor
  public list: List
  public sidebar: Sidebar
  public comments: Comment[] = []

  constructor (conf: ArtalkConfig) {
    // Version Information
    console.log(`\n %c `
      + `Artalk v${ARTALK_VERSION} %c 一款简洁有趣的可拓展评论系统 \n\n%c`
      + `> https://artalk.js.org\n`
      + `> https://github.com/ArtalkJS/Artalk\n`,
      'color: #FFF; background: #1DAAFF; padding:5px 0;', 'color: #FFF; background: #656565; padding:5px 0;', '')

    // Options
    this.conf = { ...defaultOpts, ...conf }
    this.conf.serverUrl = this.conf.serverUrl.replace(/\/$/, '')

    // Main Element
    try {
      const el = document.querySelector<HTMLElement>(this.conf.el)
      if (el === null) {
        throw Error(`Sorry, Target element "${this.conf.el}" was not found.`)
      }
      this.el = el
    } catch (e) {
      console.error(e)
      throw new Error('Artalk config `el` error')
    }

    // Create context
    this.ctx = new Context(this.el, this.conf)

    this.el.classList.add('artalk')
    this.el.setAttribute('atk-run-id', this.contextID.toString())

    // 若该元素中 artalk 已装载
    if (this.el.innerHTML.trim() !== '') this.el.innerHTML = ''

    // Components
    this.initDarkMode()

    this.checker = new Checker(this.ctx)

    this.editor = new Editor(this.ctx)
    this.list = new List(this.ctx)
    this.sidebar = new Sidebar(this.ctx)

    this.el.appendChild(this.editor.el)
    this.el.appendChild(this.list.el)
    this.el.appendChild(this.sidebar.el)

    // 请求获取评论
    this.list.reqComments()

    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.list.checkGoToCommentByUrlHash()
    })

    // 仅管理员显示控制
    this.ctx.addEventListener('check-admin-show-el', () => {
      this.el.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach((itemEl: HTMLElement) => {
        if (this.ctx.user.data.isAdmin)
          itemEl.classList.remove('atk-hide')
        else
          itemEl.classList.add('atk-hide')
      })
    })

    this.ctx.addEventListener('user-changed', () => {
      this.ctx.dispatchEvent('check-admin-show-el')
    })
  }

  /** 暗黑模式 - 初始化 */
  initDarkMode() {
    if (this.conf.darkMode) {
      this.el.classList.add(Constant.DARK_MODE_CLASSNAME)
    } else {
      this.el.classList.remove(Constant.DARK_MODE_CLASSNAME)
    }

    // for Layer
    const { wrapEl: layerWrapEl } = GetLayerWrap(this.ctx)
    if (layerWrapEl) {
      if (this.conf.darkMode) {
        layerWrapEl.classList.add(Constant.DARK_MODE_CLASSNAME)
      } else {
        layerWrapEl.classList.remove(Constant.DARK_MODE_CLASSNAME)
      }
    }
  }

  /** 暗黑模式 - 设定 */
  setDarkMode(darkMode: boolean) {
    this.ctx.conf.darkMode = darkMode
    this.initDarkMode()
  }

  /** 暗黑模式 - 开启 */
  openDarkMode() {
    this.setDarkMode(true)
  }

  /** 暗黑模式 - 关闭 */
  closeDarkMode() {
    this.setDarkMode(false)
  }
}

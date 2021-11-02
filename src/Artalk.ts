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
import { EventPayloadMap, Handler } from '~/types/event'

const defaultOpts: ArtalkConfig = {
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
  darkMode: false,
  reqTimeout: 15000,
  flatMode: false,
  maxNesting: 3,
}

export default class Artalk {
  public ctx: Context
  public conf: ArtalkConfig
  public el: HTMLElement
  public readonly contextID: number = new Date().getTime() // 实例唯一 ID
  public checker: Checker
  public editor: Editor
  public list: List
  public sidebar: Sidebar
  public comments: Comment[] = []

  constructor (conf: ArtalkConfig) {
    // Version Information
    // console.log(`\n %c `
    //   + `Artalk v${ARTALK_VERSION} %c 一款简洁有趣的可拓展评论系统 \n\n%c`
    //   + `> https://artalk.js.org\n`
    //   + `> https://github.com/ArtalkJS/Artalk\n`,
    //   'color: #FFF; background: #1DAAFF; padding:5px 0;', 'color: #FFF; background: #656565; padding:5px 0;', '')

    // Options
    this.conf = { ...defaultOpts, ...conf }
    this.conf.server = this.conf.server.replace(/\/$/, '')

    // Default `pageKey` conf
    if (!this.conf.pageKey) {
      // TODO 自动获取和 atk_comment query 冲突
      // eslint-disable-next-line prefer-destructuring
      this.conf.pageKey = window.location.href.split('#')[0]
    }

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
    this.ctx.on('check-admin-show-el', () => {
      const items: HTMLElement[] = []

      this.el.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))
      // for layer
      const { wrapEl: layerWrapEl } = GetLayerWrap(this.ctx)
      if (layerWrapEl)
        layerWrapEl.querySelectorAll<HTMLElement>(`[atk-only-admin-show]`).forEach(item => items.push(item))

      items.forEach((itemEl: HTMLElement) => {
        if (this.ctx.user.data.isAdmin)
          itemEl.classList.remove('atk-hide')
        else
          itemEl.classList.add('atk-hide')
      })
    })

    this.ctx.on('user-changed', () => {
      this.ctx.trigger('check-admin-show-el')
      this.ctx.trigger('list-refresh-ui')
    })
  }

  /** 暗黑模式 - 初始化 */
  public initDarkMode() {
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
  public setDarkMode(darkMode: boolean) {
    this.ctx.conf.darkMode = darkMode
    this.initDarkMode()
  }

  /** 暗黑模式 - 开启 */
  public openDarkMode() {
    this.setDarkMode(true)
  }

  /** 暗黑模式 - 关闭 */
  public closeDarkMode() {
    this.setDarkMode(false)
  }

  public on<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>): void {
    this.ctx.on(name, handler, 'external')
  }

  public off<K extends keyof EventPayloadMap>(name: K, handler: Handler<EventPayloadMap[K]>): void {
    this.ctx.off(name, handler, 'external')
  }

  public trigger<K extends keyof EventPayloadMap>(name: K, payload?: EventPayloadMap[K]): void {
    this.ctx.trigger(name, payload, 'external')
  }
}

import './css/main.less'
import marked from 'marked'
import hanabi from 'hanabi'
import User from './components/User'
import Checker from './components/Checker'
import Editor from './components/Editor'
import List from './components/List'
import Comment from './components/Comment'
import Sidebar from './components/Sidebar'
import Ui from './utils/ui'
import Utils from './utils'
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
    cdn: 'https://cdn.v2ex.com/gravatar/'
  },
  darkMode: false,
}

export default class Artalk {
  public conf: ArtalkConfig
  public el: HTMLElement
  public readonly runId: number = new Date().getTime() // 实例唯一 ID

  public ui: Ui
  public user: User
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
    this.el.setAttribute('artalk-run-id', this.runId.toString())

    // 若该元素中 artalk 已装载
    if (this.el.innerHTML.trim() !== '') this.el.innerHTML = ''

    // Components
    this.ui = new Ui(this)
    this.user = new User(this)
    this.checker = new Checker(this)
    this.initMarked()
    this.editor = new Editor(this)
    this.list = new List(this)
    this.sidebar = new Sidebar(this)

    // 请求获取评论
    this.list.reqComments()

    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.checkGoToCommentByUrlHash()
    })
  }

  /** 遍历操作 Comment (包括子评论) */
  public eachComment (commentList: Comment[], action: (comment?: Comment, levelList?: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  /** 查找评论项 */
  public findComment (id: number) {
    let comment: Comment|null = null

    this.eachComment(this.comments, (item) => {
      if (item.data.id === id) {
        comment = item
        return false
      }
      return true
    })

    return comment
  }

  /** 获取评论总数 */
  public getCommentCount (): number {
    let count = 0
    this.eachComment(this.comments, () => { count++ })
    return count
  }

  /** 删除评论 */
  public deleteComment (comment: number|Comment) {
    let findComment: Comment
    if (typeof comment === 'number') {
      findComment = this.findComment(comment)
      if (!findComment) throw Error(`未找到评论 ${comment}`)
    } else findComment = comment

    findComment.getElem().remove()
    this.eachComment(this.comments, (item, levelList) => {
      if (item === findComment) {
        levelList.splice(levelList.indexOf(item), 1)
        return false
      }
      return true
    })
  }

  /** 清空所有评论 */
  public clearComments () {
    this.list.commentsWrapEl.innerHTML = ''
    this.list.data = undefined
    this.comments = []
  }

  /** 公共请求 */
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

  /** 跳到评论项位置 - 根据 `location.hash` */
  public checkGoToCommentByUrlHash () {
    let commentId: number = Number(Utils.getLocationParmByName('artalk_comment'))
    if (!commentId) {
      const match = window.location.hash.match(/#artalk-comment-([0-9]+)/)
      if (!match || !match[1] || Number.isNaN(Number(match[1]))) return
      commentId = Number(match[1])
    }

    const comment = this.findComment(commentId)
    if (!comment && this.list.hasMoreComments) {
      this.list.readMore()
      return
    }
    if (!comment) { return }

    this.ui.scrollIntoView(comment.getElem(), false)
    setTimeout(() => {
      comment.getElem().classList.add('artalk-flash-once')
    }, 800)
  }

  public marked: (src: string, callback?: (error: any, parseResult: string) => void) => string

  /** 初始化 Marked */
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
}

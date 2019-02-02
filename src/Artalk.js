import './css/main.scss'
import Editor from './components/Editor'
import List from './components/List'
import $ from 'jquery'
import marked from 'marked'
import hanabi from 'hanabi'
import Sidebar from './components/Sidebar'
import Ui from './js/ui'

/* global ARTALK_VERSION */

const defaultOpts = {
  el: '',
  placeholder: '来啊，快活啊 (/ω＼)',
  noComment: '快来成为第一个评论的人吧~',
  defaultAvatar: 'mp',
  pageSize: 50,
  pageKey: '',
  serverUrl: '',
  emoticons: require('./assets/emoticons.json'),
  gravatar: {
    cdn: 'https://gravatar.loli.net/avatar/'
  }
}

class Artalk {
  constructor (opts) {
    // Copyright
    console.log(`\n %c Artalk v${ARTALK_VERSION} %c 一款简洁有趣的自托管评论系统 \n\n%c> https://artalk.js.org\n> https://github.com/qwqcode/Artalk\n> https://qwqaq.com \n`, 'color: #FFF; background: #1DAAFF; padding:5px 0;', 'color: #FFF; background: #656565; padding:5px 0;', '')

    // Options
    this.opts = Object.assign(defaultOpts, opts)

    // Main Element
    this.el = $(this.opts.el)
    if (!this.el.length) {
      throw Error(`Sorry, Target element "${this.opts.el}" was not found.`)
    }
    this.el.addClass('artalk')

    // Components
    this.ui = new Ui(this)
    this.initMarked()
    this.editor = new Editor(this)
    this.list = new List(this)
    this.sidebar = new Sidebar(this)
  }

  initMarked () {
    let renderer = new marked.Renderer()
    const linkRenderer = renderer.link
    renderer.link = function (href, title, text) {
      const html = linkRenderer.call(renderer, href, title, text)
      return html.replace(/^<a /, '<a target="_blank" rel="nofollow" ')
    }

    this.marked = marked
    this.marked.setOptions({
      renderer: renderer,
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
  }

  request (action, data, before, after, success, error) {
    $.ajax({
      type: 'POST',
      url: this.opts.serverUrl,
      data: Object.assign({ action: action }, data),
      dataType: 'json',
      beforeSend: () => {
        before()
      },
      success: (obj) => {
        after()
        if (obj.success) {
          success(obj.msg, obj.data)
        } else {
          error(obj.msg, obj.data)
        }
      },
      error: () => {
        after()
        error('网络错误', {})
      }
    })
  }
}

export default Artalk

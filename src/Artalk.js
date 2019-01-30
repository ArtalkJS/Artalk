import './css/main.scss'
import Editor from './js/Editor'
import List from './js/List'
import $ from 'jquery'
import marked from 'marked'
import hanabi from 'hanabi'

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
    console.log('\n %c Artalk %c 一款有趣有后端的评论系统 \n\n%c> https://github.com/qwqcode/Artalk\n> https://qwqaq.com \n', 'color: #FFF; background: #1DAAFF; padding:5px 0;', 'color: #FFF; background: #656565; padding:5px 0;', '')

    // 配置
    this.opts = defaultOpts
    this.opts = Object.assign(this.opts, opts)

    this.init()
  }

  init () {
    this.el = $(this.opts.el)
    if (!this.el.length) {
      throw Error(`Sorry, Target element "${this.opts.el}" was not found.`)
    }
    this.el.addClass('artalk')

    this.initMarked()
    this.editor = new Editor(this)
    this.list = new List(this)
  }

  initMarked () {
    this.marked = marked

    this.marked.setOptions({
      renderer: new marked.Renderer(),
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

  showLoading (parentElem) {
    if (!parentElem) {
      parentElem = this.el
    }
    let loadingElem = $(parentElem).find('.artalk-loading')
    if (!loadingElem.length) {
      loadingElem = $(`<div class="artalk-loading" style="display: none;"><div class="artalk-loading-spinner"><svg viewBox="25 25 50 50"><circle cx="50" cy="50" r="20" fill="none" stroke-width="2" stroke-miterlimit="10"></circle></svg></div></div>`).appendTo(parentElem)
    }
    loadingElem.css('display', '')
  }

  hideLoading (parentElem) {
    if (!parentElem) {
      parentElem = this.el
    }
    let loadingElem = $(parentElem).find('.artalk-loading')
    loadingElem.css('display', 'none')
  }

  showErrorDialog () {
    let elem = this.el.find('.artalk-error-dialog')
    if (!elem.length) {

    }
  }

  setGlobalError (msg) {
    let elem = this.el.find('.artalk-error-layer')
    if (!elem.length) {
      elem = $('<div class="artalk-error-layer">Artalk Error<span class="artalk-error-text"></span></div>').appendTo(this.el)
    } else {
      if (msg === null) {
        elem.remove()
      }
    }
    elem.find('.artalk-error-text').text(msg)
  }

  scrollToView (elem) {
    if (this.isScrolledIntoView(elem)) {
      return
    }

    $('html,body').animate({
      scrollTop: $(elem).offset().top - $(window).height() / 2
    }, 500)
  }

  isScrolledIntoView (elem) {
    let docViewTop = $(window).scrollTop()
    let docViewBottom = docViewTop + $(window).height()

    let elemTop = $(elem).offset().top
    let elemBottom = elemTop + $(elem).height()

    return ((elemBottom <= docViewBottom) && (elemTop >= docViewTop))
  }
}

export default Artalk

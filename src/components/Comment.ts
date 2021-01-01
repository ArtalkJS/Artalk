import '../css/comment.less'
import List from './List'
import UADetect from '../utils/detect'
import { CommentData } from '~/types/artalk-data'
import Artalk from '../Artalk'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

export default class Comment extends ArtalkContext {
  public elem: HTMLElement
  public mainEl: HTMLElement
  public bodyEl: HTMLElement
  public contentEl: HTMLElement
  public childrenEl: HTMLElement
  public actionsEl: HTMLElement

  public parent: Comment|null
  public nestedNum: number
  private readonly maxNestedNo = 3 // 最多嵌套层数
  public children: Comment[] = []

  constructor (artalk: Artalk, public list: List, public data: CommentData) {
    super(artalk)

    this.data = { ...this.data }
    this.data.date = this.data.date.replace(/-/g, '/') // 解决 Safari 日期解析 NaN 问题

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层

    this.renderElem()
  }

  private renderElem () {
    this.elem = Utils.createElement(require('../templates/Comment.ejs')(this))
    this.mainEl = this.elem.querySelector('.artalk-comment-main')
    this.bodyEl = this.elem.querySelector('.artalk-body')
    this.contentEl = this.bodyEl.querySelector('.artalk-content')
    this.actionsEl = this.elem.querySelector('.artalk-comment-actions')
    this.childrenEl = null

    // 已折叠的评论
    const contentShowBtn = this.mainEl.querySelector('.artalk-collapsed .artalk-show-btn')
    if (contentShowBtn) {
      contentShowBtn.addEventListener('click', () => {
        if (this.contentEl.classList.contains('artalk-hide')) {
          this.contentEl.classList.remove('artalk-hide')
          this.artalk.ui.playFadeInAnim(this.contentEl)
          contentShowBtn.innerHTML = '收起内容'
        } else {
          this.contentEl.classList.add('artalk-hide')
          contentShowBtn.innerHTML = '查看内容'
        }
      })
    }

    this.initActionBtn()

    return this.elem
  }

  private refreshUI () {
    const originalEl = this.elem
    const newEl = this.renderElem()
    originalEl.replaceWith(newEl) // 替换 document 中的的 elem
    this.playFadeInAnim()

    // 重建子评论元素
    this.artalk.eachComment(this.children, (child) => {
      child.parent.getChildrenEl().appendChild(child.renderElem())
      child.playFadeInAnim()
    })
  }

  initActionBtn () {
    // 绑定回复按钮事件
    const replyBtn = this.actionsEl.querySelector('[data-comment-action="reply"]') as HTMLElement
    if (replyBtn) {
      replyBtn.addEventListener('click', () => {
        this.artalk.editor.setReply(this)
      })
    }

    // 绑定折叠按钮事件
    const collapseBtn = this.actionsEl.querySelector('[data-comment-action="collapse"]') as HTMLElement
    if (collapseBtn) {
      collapseBtn.addEventListener('click', () => {
        this.adminCollapse(collapseBtn)
      })
    }

    // 绑定删除按钮事件
    const delBtn = this.actionsEl.querySelector('[data-comment-action="delete"]') as HTMLElement
    if (delBtn) {
      delBtn.addEventListener('click', () => {
        this.adminDelete(delBtn)
      })
    }
  }

  getIsRoot () {
    return this.parent === null
  }

  getChildren () {
    return this.children
  }

  putChild (childC: Comment) {
    childC.parent = this
    childC.nestedNum = this.nestedNum + 1 // 嵌套层数 +1
    this.children.push(childC)

    this.getChildrenEl().appendChild(childC.getElem())
    childC.playFadeInAnim()
  }

  getChildrenEl () {
    if (this.childrenEl === null) {
      // console.log(this.nestedNo)
      if (this.nestedNum < this.maxNestedNo) {
        this.childrenEl = Utils.createElement('<div class="artalk-comment-children"></div>')
        this.mainEl.appendChild(this.childrenEl)
      } else {
        this.childrenEl = this.parent.getChildrenEl()
      }
    }
    return this.childrenEl
  }

  getParent () {
    return this.parent
  }

  getElem () {
    return this.elem
  }

  getData () {
    return this.data
  }

  getGravatarUrl () {
    return `${this.artalk.conf.gravatar.cdn}${this.data.email_encrypted}?d=${encodeURIComponent(this.artalk.conf.defaultAvatar)}&s=80`
  }

  getContentMarked () {
    return this.artalk.marked(this.data.content)
  }

  getDateFormatted () {
    return Utils.timeAgo(new Date(this.data.date))
  }

  getUserUaBrowser () {
    const info = UADetect(this.data.ua)
    return `${info.browser} ${info.version}`
  }

  getUserUaOS () {
    const info = UADetect(this.data.ua)
    return `${info.os} ${info.osVersion}`
  }

  /** 渐出动画 */
  playFadeInAnim () {
    this.artalk.ui.playFadeInAnim(this.elem)
  }

  /** 管理员 - 评论折叠 */
  adminCollapse (btnElem: HTMLElement) {
    if (btnElem.classList.contains('artalk-in-process')) return // 若正在折叠中
    const btnTextOrg = btnElem.innerText
    const isCollapse = !this.data.is_collapsed
    this.artalk.request('CommentCollapse', {
      id: this.data.id,
      nick: this.artalk.user.data.nick,
      email: this.artalk.user.data.email,
      password: this.artalk.user.data.password,
      is_collapsed: Number(isCollapse)
    }, () => {
      btnElem.classList.add('artalk-in-process')
      btnElem.innerText = isCollapse ? '折叠中...' : '展开中...'
    }, () => {
    }, (msg, data) => {
      btnElem.classList.remove('artalk-in-process')
      this.data.is_collapsed = data.is_collapsed
      this.artalk.eachComment([this], (item) => {
        item.data.is_allow_reply = !data.is_collapsed // 禁止回复
      })
      this.refreshUI()
      this.artalk.ui.playFadeInAnim(this.bodyEl)
      this.list.refreshUI()
    }, (msg, data) => {
      btnElem.classList.add('artalk-error')
      btnElem.innerText = isCollapse ? '折叠失败' : '展开失败'
      setTimeout(() => {
        btnElem.innerText = btnTextOrg
        btnElem.classList.remove('artalk-error')
        btnElem.classList.remove('artalk-in-process')
      }, 2000)
    })
  }

  /** 管理员 - 评论删除 */
  adminDelete (btnElem: HTMLElement) {
    if (btnElem.classList.contains('artalk-in-process')) return // 若正在删除中

    // 删除确认
    const btnClicked = Number(btnElem.getAttribute('data-btn-clicked') || 1)
    if (btnClicked < 2) {
      if (btnClicked === 1) {
        const btnTextOrg = btnElem.innerText
        btnElem.innerText = '确认删除'
        setTimeout(() => {
          btnElem.innerText = btnTextOrg
          btnElem.setAttribute('data-btn-clicked', '')
        }, 2000)
        btnElem.setAttribute('data-btn-clicked', String(btnClicked+1))
      }
      return
    }
    const btnTextOrg = btnElem.innerText
    this.artalk.request('CommentDel', {
      id: this.data.id,
      nick: this.artalk.user.data.nick,
      email: this.artalk.user.data.email,
      password: this.artalk.user.data.password
    }, () => {
      btnElem.classList.add('artalk-in-process')
      btnElem.innerText = '删除中...'
    }, () => {
    }, (msg, data) => {
      btnElem.innerText = btnTextOrg
      this.artalk.deleteComment(this)
      this.list.data.total -= 1 // 评论数 -1
      this.list.refreshUI() // 刷新 list
      btnElem.classList.remove('artalk-in-process')
    }, (msg, data) => {
      btnElem.classList.add('artalk-error')
      btnElem.innerText = '删除失败'
      setTimeout(() => {
        btnElem.innerText = btnTextOrg
        btnElem.classList.remove('artalk-error')
        btnElem.classList.remove('artalk-in-process')
      }, 2000)
    })
  }
}

import '../style/comment.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import UADetect from '../lib/detect'
import { CommentData } from '~/types/artalk-data'
import CommentHTML from './html/comment.html?raw'

export default class Comment extends Component {
  public data: CommentData

  public mainEl!: HTMLElement
  public bodyEl!: HTMLElement
  public contentEl!: HTMLElement
  public childrenEl!: HTMLElement|null
  public actionsEl!: HTMLElement

  public parent: Comment|null
  public nestedNum: number
  private readonly maxNestingNum = 3 // 最多嵌套层数
  public children: Comment[] = []

  constructor (ctx: Context, data: CommentData) {
    super(ctx)

    this.data = { ...data }
    this.data.date = this.data.date.replace(/-/g, '/') // 解决 Safari 日期解析 NaN 问题

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层

    this.renderElem()
  }

  private renderElem () {
    this.el = Utils.createElement(CommentHTML)
    this.mainEl = this.el.querySelector('.atk-comment-main')!
    this.bodyEl = this.el.querySelector('.atk-body')!
    this.contentEl = this.bodyEl.querySelector('.atk-content')!
    this.actionsEl = this.el.querySelector('.atk-comment-actions')!
    this.childrenEl = null

    // 填入内容
    this.el.setAttribute('data-comment-id', `${this.data.id}`)
    this.el.querySelector('.atk-avatar a')!.setAttribute('href', this.data.link)
    this.el.querySelector('.atk-avatar img')!.setAttribute('src', this.getGravatarUrl())

    const nickEl = this.el.querySelector<HTMLLinkElement>('.atk-nick > a')!
    nickEl.innerText = this.data.nick
    nickEl.href = this.data.link

    const badgeEl = this.el.querySelector<HTMLElement>('.atk-badge')!
    if (this.data.badge) {
      badgeEl.innerText = this.data.badge.name || '徽章'
      if (this.data.badge.color)
        badgeEl.style.backgroundColor = this.data.badge.color
    } else {
      badgeEl.remove()
    }

    this.el.querySelector<HTMLElement>('.atk-date')!.innerText = this.getDateFormatted()
    this.el.querySelector<HTMLElement>('.atk-ua.ua-browser')!.innerText = this.getUserUaBrowser()
    this.el.querySelector<HTMLElement>('.atk-ua.ua-os')!.innerText = this.getUserUaOS()

    // 内容 & 折叠
    if (!this.data.is_collapsed)
      this.contentEl.innerHTML = this.getContentMarked()

    if (this.data.is_collapsed) {
      this.contentEl.classList.add('atk-hide', 'atk-type-collapsed')
      const collapsedInfoEl = Utils.createElement(`
      <div class="atk-collapsed">
        <span class="atk-text">该评论已被系统或管理员折叠</span>
        <span class="atk-show-btn">查看内容</span>
      </div>`)
      this.contentEl.insertAdjacentElement('afterbegin', collapsedInfoEl)

      const contentShowBtn = collapsedInfoEl.querySelector('.atk-show-btn')!
      contentShowBtn.addEventListener('click', () => {
        if (this.contentEl.classList.contains('atk-hide')) {
          this.contentEl.innerHTML = this.getContentMarked()
          this.contentEl.classList.remove('atk-hide')
          Ui.playFadeInAnim(this.contentEl)
          contentShowBtn.innerHTML = '收起内容'
        } else {
          this.contentEl.innerHTML = ''
          this.contentEl.classList.add('atk-hide')
          contentShowBtn.innerHTML = '查看内容'
        }
      })
    }

    this.initActionBtn()

    return this.el
  }

  private eachComment (commentList: Comment[], action: (comment: Comment, levelList: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  public refreshUI () {
    const originalEl = this.el
    const newEl = this.renderElem()
    originalEl.replaceWith(newEl) // 替换 document 中的的 elem
    this.playFadeInAnim()

    // 重建子评论元素
    this.eachComment(this.children, (child) => {
      child.parent?.getChildrenEl().appendChild(child.renderElem())
      child.playFadeInAnim()
    })
  }

  initActionBtn () {
    // 绑定回复按钮事件
    if (this.data.is_allow_reply) {
      const replyBtn = Utils.createElement(`<span>回复</span>`)
      this.actionsEl.append(replyBtn)
      replyBtn.addEventListener('click', () => {
        this.ctx.dispatchEvent('editor-reply', this.data)
      })
    }

    // 管理员操作按钮

    // 绑定折叠按钮事件
    const collapseBtn = Utils.createElement(`<span atk-only-admin-show>${this.data.is_collapsed ? '取消折叠' : '折叠'}</span>`)
    this.actionsEl.append(collapseBtn)
    collapseBtn.addEventListener('click', () => {
      this.adminCollapse(collapseBtn)
    })

    // 绑定删除按钮事件
    const delBtn = Utils.createElement(`<span atk-only-admin-show>删除</span>`)
    this.actionsEl.append(delBtn)
    delBtn.addEventListener('click', () => {
      this.adminDelete(delBtn)
    })
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

    this.getChildrenEl().appendChild(childC.getEl())
    childC.playFadeInAnim()
  }

  getChildrenEl () {
    if (this.childrenEl === null) {
      // console.log(this.nestedNum)
      if (this.nestedNum < this.maxNestingNum) {
        this.childrenEl = Utils.createElement('<div class="atk-comment-children"></div>')
        this.mainEl.appendChild(this.childrenEl)
      } else if (this.parent) {
          this.childrenEl = this.parent.getChildrenEl()
      }
    }
    return this.childrenEl
  }

  getParent () {
    return this.parent
  }

  getEl () {
    return this.el
  }

  getData () {
    return this.data
  }

  getGravatarUrl () {
    return `${this.ctx.conf.gravatar?.cdn || ''}${this.data.email_encrypted}?d=${encodeURIComponent(this.ctx.conf.defaultAvatar || '')}&s=80`
  }

  getContentMarked () {
    return Utils.marked(this.ctx, this.data.content)
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
    Ui.playFadeInAnim(this.el)
  }

  /** 管理员 - 评论折叠 */
  adminCollapse (btnElem: HTMLElement) {
    if (btnElem.classList.contains('atk-in-process')) return // 若正在折叠中
    const btnTextOrg = btnElem.innerText
    const isCollapse = !this.data.is_collapsed

    // TODO: 评论折叠
    // this.artalk.request('CommentCollapse', {
    //   id: this.data.id,
    //   nick: this.ctx.user.data.nick,
    //   email: this.ctx.user.data.email,
    //   token: this.ctx.user.data.token,
    //   is_collapsed: Number(isCollapse)
    // }, () => {
    //   btnElem.classList.add('atk-in-process')
    //   btnElem.innerText = isCollapse ? '折叠中...' : '展开中...'
    // }, () => {
    // }, (msg, data) => {
    //   btnElem.classList.remove('atk-in-process')
    //   this.data.is_collapsed = data.is_collapsed
    //   this.eachComment([this], (item) => {
    //     item.data.is_allow_reply = !data.is_collapsed // 禁止回复
    //   })
    //   this.refreshUI()
    //   Ui.playFadeInAnim(this.bodyEl)
    //   this.list.refreshUI()
    // }, (msg, data) => {
    //   btnElem.classList.add('atk-error')
    //   btnElem.innerText = isCollapse ? '折叠失败' : '展开失败'
    //   setTimeout(() => {
    //     btnElem.innerText = btnTextOrg
    //     btnElem.classList.remove('atk-error')
    //     btnElem.classList.remove('atk-in-process')
    //   }, 2000)
    // })
  }

  /** 管理员 - 评论删除 */
  adminDelete (btnElem: HTMLElement) {
    if (btnElem.classList.contains('atk-in-process')) return // 若正在删除中

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

    // TODO: 评论删除
    // this.artalk.request('CommentDel', {
    //   id: this.data.id,
    //   nick: this.ctx.user.data.nick,
    //   email: this.ctx.user.data.email,
    //   token: this.ctx.user.data.token
    // }, () => {
    //   btnElem.classList.add('atk-in-process')
    //   btnElem.innerText = '删除中...'
    // }, () => {
    // }, (msg, data) => {
    //   btnElem.innerText = btnTextOrg
    //   this.artalk.deleteComment(this)
    //   this.list.data.total -= 1 // 评论数 -1
    //   this.list.refreshUI() // 刷新 list
    //   btnElem.classList.remove('atk-in-process')
    // }, (msg, data) => {
    //   btnElem.classList.add('atk-error')
    //   btnElem.innerText = '删除失败'
    //   setTimeout(() => {
    //     btnElem.innerText = btnTextOrg
    //     btnElem.classList.remove('atk-error')
    //     btnElem.classList.remove('atk-in-process')
    //   }, 2000)
    // })
  }
}

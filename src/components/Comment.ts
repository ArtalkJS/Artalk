import '../style/comment.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import UADetect from '../lib/detect'
import { CommentData } from '~/types/artalk-data'
import CommentHTML from './html/comment.html?raw'
import Api from '../lib/api'

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

  public replyTo?: CommentData // 回复对象（flatMode 用）

  public afterRender?: () => void

  constructor (ctx: Context, data: CommentData) {
    super(ctx)

    this.data = { ...data }
    this.data.date = this.data.date.replace(/-/g, '/') // 解决 Safari 日期解析 NaN 问题

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层
  }

  public renderElem () {
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

    const dateEL = this.el.querySelector<HTMLElement>('.atk-date')!
    dateEL.innerText = this.getDateFormatted()
    dateEL.setAttribute('data-atk-comment-date', String(+new Date(this.data.date)))
    this.el.querySelector<HTMLElement>('.atk-ua.ua-browser')!.innerText = this.getUserUaBrowser()
    this.el.querySelector<HTMLElement>('.atk-ua.ua-os')!.innerText = this.getUserUaOS()

    // 内容 & 折叠
    if (!this.data.is_collapsed) {
      this.contentEl.innerHTML = this.getContentMarked()
    } else {
      this.contentEl.classList.add('atk-hide', 'atk-type-collapsed')
      const collapsedInfoEl = Utils.createElement(`
      <div class="atk-collapsed">
        <span class="atk-text">该评论已被系统或管理员折叠</span>
        <span class="atk-show-btn">查看内容</span>
      </div>`)
      this.bodyEl.insertAdjacentElement('beforeend', collapsedInfoEl)

      const contentShowBtn = collapsedInfoEl.querySelector('.atk-show-btn')!
      contentShowBtn.addEventListener('click', (e) => {
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

        e.stopPropagation() // 防止穿透
      })
    }

    // 显示回复的对象
    if (this.replyTo) {
      const replyToEl = Utils.createElement(`
      <div class="atk-reply-to">
        <div class="atk-meta">回复 <span class="atk-nick"></span>:</div>
        <div class="atk-content"></div>
      </div>`)
      replyToEl.querySelector<HTMLElement>('.atk-nick')!.innerText = `@${this.replyTo.nick}`
      replyToEl.querySelector<HTMLElement>('.atk-content')!.innerHTML = Utils.marked(this.ctx, this.replyTo.content)
      this.bodyEl.prepend(replyToEl)
    }

    // 显示待审核状态
    if (this.data.is_pending) {
      const pendingEl = Utils.createElement(`<div class="atk-pending">审核中，仅本人可见。</div>`)
      this.bodyEl.prepend(pendingEl)
    }

    this.initActionBtn()

    if (this.afterRender) this.afterRender()

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
      const replyBtn = Utils.createElement(`<span data-atk-action="comment-reply">回复</span>`)
      this.actionsEl.append(replyBtn)
      replyBtn.addEventListener('click', (e) => {
        this.ctx.dispatchEvent('editor-reply', this.data)
        e.stopPropagation() // 防止穿透
      })
    }

    // 管理员操作按钮

    // 绑定折叠按钮事件
    const collapseBtn = Utils.createElement(`<span atk-only-admin-show>${this.data.is_collapsed ? '取消折叠' : '折叠'}</span>`)
    this.actionsEl.append(collapseBtn)
    collapseBtn.addEventListener('click', (e) => {
      this.adminEdit('collapsed', collapseBtn)
      e.stopPropagation() // 防止穿透
    })

    // 绑定待审核按钮事件
    const pendingBtn = Utils.createElement(`<span atk-only-admin-show>${this.data.is_pending ? '待审' : '已审'}</span>`)
    this.actionsEl.append(pendingBtn)
    pendingBtn.addEventListener('click', (e) => {
      this.adminEdit('pending', pendingBtn)
      e.stopPropagation() // 防止穿透
    })

    // 绑定删除按钮事件
    const delBtn = Utils.createElement(`<span atk-only-admin-show>删除</span>`)
    this.actionsEl.append(delBtn)
    delBtn.addEventListener('click', (e) => {
      this.adminDelete(delBtn)
      e.stopPropagation() // 防止穿透
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
  adminEdit (type: 'collapsed'|'pending', btnElem: HTMLElement) {
    if (btnElem.classList.contains('atk-in-process')) return // 若正在折叠中
    const btnTextOrg = btnElem.innerText

    btnElem.classList.add('atk-in-process')
    btnElem.innerText = '修改中...'

    const params: any = {
      id: this.data.id,
    }
    if (type === 'collapsed') {
      params.is_collapsed = !this.data.is_collapsed
    } else if (type === 'pending') {
      params.is_pending = !this.data.is_pending
    }

    new Api(this.ctx).commentEdit(params).then((comment) => {
      btnElem.classList.remove('atk-in-process')
      this.data = comment
      this.refreshUI()
      Ui.playFadeInAnim(this.bodyEl)
      this.ctx.dispatchEvent('list-refresh-ui')
    }).catch((err) => {
      console.error(err)
      btnElem.classList.add('atk-error')
      btnElem.innerText = '修改失败'
      setTimeout(() => {
        btnElem.innerText = btnTextOrg
        btnElem.classList.remove('atk-error')
        btnElem.classList.remove('atk-in-process')
      }, 2000)
    })
  }

  public onDelete?: (comment: Comment) => void

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

    btnElem.classList.add('atk-in-process')
    btnElem.innerText = '删除中...'
    new Api(this.ctx).commentDel(this.data.id)
      .then(() => {
        btnElem.innerText = btnTextOrg
        btnElem.classList.remove('atk-in-process')
        if (this.onDelete) this.onDelete(this)
      })
      .catch((e) => {
        console.error(e)
        btnElem.classList.add('atk-error')
        btnElem.innerText = '删除失败'
        setTimeout(() => {
          btnElem.innerText = btnTextOrg
          btnElem.classList.remove('atk-error')
          btnElem.classList.remove('atk-in-process')
        }, 2000)
      })
  }
}

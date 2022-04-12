import '../style/comment.less'

import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import UADetect from '../lib/detect'
import { CommentData } from '~/types/artalk-data'
import CommentHTML from './html/comment.html?raw'
import Api from '../api'
import ActionBtn from './action-btn'

export default class Comment extends Component {
  public data: CommentData

  public $main!: HTMLElement
  public $header!: HTMLElement
  public $body!: HTMLElement
  public $content!: HTMLElement
  public $children!: HTMLElement|null
  public $actions!: HTMLElement
  public voteBtnUp?: ActionBtn
  public voteBtnDown?: ActionBtn

  public parent: Comment|null
  public nestedNum: number
  private maxNestingNum: number // 最多嵌套层数
  public children: Comment[] = []

  public replyTo?: CommentData // 回复对象（flatMode 用）
  public $replyTo?: HTMLElement

  public afterRender?: () => void

  private unread = false
  private openable = false
  private openURL?: string
  public openEvt?: () => void

  public onReplyBtnClick?: Function

  constructor(ctx: Context, data: CommentData) {
    super(ctx)

    // 最大嵌套数
    this.maxNestingNum = ctx.conf.maxNesting || 3

    this.data = { ...data }
    this.data.date = this.data.date.replace(/-/g, '/') // 解决 Safari 日期解析 NaN 问题

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层
  }

  /** 渲染 UI */
  public render() {
    this.$el = Utils.createElement(CommentHTML)
    this.$main = this.$el.querySelector('.atk-main')!
    this.$header = this.$el.querySelector('.atk-header')!
    this.$body = this.$el.querySelector('.atk-body')!
    this.$content = this.$body.querySelector('.atk-content')!
    this.$actions = this.$el.querySelector('.atk-actions')!
    this.$children = null

    this.$el.setAttribute('data-comment-id', `${this.data.id}`)

    this.renderCheckUnread()
    this.renderCheckClickable()

    this.renderAvatar()
    this.renderHeader()
    this.renderContent()
    this.renderReplyTo()
    this.renderPending()
    this.renderActionBtn()

    if (this.afterRender) this.afterRender()

    return this.$el
  }

  //#region Renders
  private renderCheckUnread() {
    if (this.unread) this.$el.classList.add('atk-unread')
    else this.$el.classList.remove('atk-unread')
  }

  private renderCheckClickable() {
    if (this.openable) {
      this.$el.classList.add('atk-openable')
    } else {
      this.$el.classList.remove('atk-openable')
    }

    this.$el.addEventListener('click', (evt) => {
      if (this.openable && this.openURL) {
        evt.preventDefault()
        window.open(this.openURL)
      }
      if (this.openEvt)
        this.openEvt()
    })
  }

  private renderAvatar() {
    const $avatar = this.$el.querySelector<HTMLElement>('.atk-avatar')!
    const $avatarImg = Utils.createElement<HTMLImageElement>('<img />')
    $avatarImg.src = this.getGravatarUrl()
    if (this.data.link) {
      const $avatarA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
      $avatarA.href = this.data.link
      $avatarA.append($avatarImg)
      $avatar.append($avatarA)
    } else {
      $avatar.append($avatarImg)
    }
  }

  private renderHeader() {
    const $nick = this.$el.querySelector<HTMLElement>('.atk-nick')!
    if (this.data.link) {
      const $nickA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
      $nickA.innerText = this.data.nick
      $nickA.href = this.data.link
      $nick.append($nickA)
    } else {
      $nick.innerText = this.data.nick
    }

    const $badge = this.$el.querySelector<HTMLElement>('.atk-badge')!
    if (this.data.badge_name) {
      $badge.innerText = this.data.badge_name
      if (this.data.badge_color)
        $badge.style.backgroundColor = this.data.badge_color
    } else {
      $badge.remove()
    }

    if (this.data.is_pinned) {
      const $pinnedBadge = Utils.createElement(`<span class="atk-item atk-pinned-badge">置顶</span>`) // 置顶徽章
      $nick.insertAdjacentElement('afterend', $pinnedBadge)
    }

    const $date = this.$el.querySelector<HTMLElement>('.atk-date')!
    $date.innerText = this.getDateFormatted()
    $date.setAttribute('data-atk-comment-date', String(+new Date(this.data.date)))

    if (this.conf.uaBadge) {
      const $uaWrap = Utils.createElement(`<span class="atk-ua-wrap"></span>`)
      const $uaBrowser = Utils.createElement(`<span class="atk-ua ua-browser"></span>`)
      const $usOS = Utils.createElement(`<span class="atk-ua ua-os"></span>`)
      $uaBrowser.innerText = this.getUserUaBrowser()
      $usOS.innerText = this.getUserUaOS()
      $uaWrap.append($uaBrowser)
      $uaWrap.append($usOS)
      this.$header.append($uaWrap)
    }
  }

  private renderContent() {
    // 内容 & 折叠
    if (!this.data.is_collapsed) {
      this.$content.innerHTML = this.getContentMarked()
      return
    }

    this.$content.classList.add('atk-hide', 'atk-type-collapsed')
    const collapsedInfoEl = Utils.createElement(`
      <div class="atk-collapsed">
        <span class="atk-text">该评论已被系统或管理员折叠</span>
        <span class="atk-show-btn">查看内容</span>
      </div>`)
    this.$body.insertAdjacentElement('beforeend', collapsedInfoEl)

    const contentShowBtn = collapsedInfoEl.querySelector('.atk-show-btn')!
    contentShowBtn.addEventListener('click', (e) => {
      e.stopPropagation() // 防止穿透

      if (this.$content.classList.contains('atk-hide')) {
        this.$content.innerHTML = this.getContentMarked()
        this.$content.classList.remove('atk-hide')
        Ui.playFadeInAnim(this.$content)
        contentShowBtn.innerHTML = '收起内容'
      } else {
        this.$content.innerHTML = ''
        this.$content.classList.add('atk-hide')
        contentShowBtn.innerHTML = '查看内容'
      }
    })
  }

  // 回复的对象
  private renderReplyTo() {
    if (!this.replyTo) return

    this.$replyTo = Utils.createElement(`
      <div class="atk-reply-to">
        <div class="atk-meta">回复 <span class="atk-nick"></span>:</div>
        <div class="atk-content"></div>
      </div>`)
    this.$replyTo.querySelector<HTMLElement>('.atk-nick')!.innerText = `@${this.replyTo.nick}`
    let replyContent = Utils.marked(this.ctx, this.replyTo.content)
    if (this.replyTo.is_collapsed) replyContent = '[已折叠]'
    this.$replyTo.querySelector<HTMLElement>('.atk-content')!.innerHTML = replyContent
    this.$body.prepend(this.$replyTo)
  }

  // 待审核状态
  private renderPending() {
    if (!this.data.is_pending) return

    const pendingEl = Utils.createElement(`<div class="atk-pending">审核中，仅本人可见。</div>`)
    this.$body.prepend(pendingEl)
  }

  /** 初始化评论操作按钮 */
  private renderActionBtn() {
    // 投票功能
    if (this.ctx.conf.vote) {
      // 赞同按钮
      this.voteBtnUp = new ActionBtn(() => `赞同 (${this.data.vote_up || 0})`).appendTo(this.$actions)
      this.voteBtnUp.setClick(() => {
        this.vote('up')
      })

      // 反对按钮
      if (this.ctx.conf.voteDown) {
        this.voteBtnDown = new ActionBtn(() => `反对 (${this.data.vote_down || 0})`).appendTo(this.$actions)
        this.voteBtnDown.setClick(() => {
          this.vote('down')
        })
      }
    }

    // 绑定回复按钮事件
    if (this.data.is_allow_reply) {
      const replyBtn = Utils.createElement(`<span>回复</span>`)
      this.$actions.append(replyBtn)
      replyBtn.addEventListener('click', (e) => {
        e.stopPropagation() // 防止穿透
        if (!this.onReplyBtnClick) {
          this.ctx.trigger('editor-reply', {data: this.data, $el: this.$el})
        } else {
          this.onReplyBtnClick()
        }
      })
    }

    // 绑定折叠按钮事件
    const collapseBtn = new ActionBtn({
      text: () => (this.data.is_collapsed ? '取消折叠' : '折叠'),
      adminOnly: true
    })
    collapseBtn.appendTo(this.$actions)
    collapseBtn.setClick(() => {
      this.adminEdit('collapsed', collapseBtn)
    })

    // 绑定待审核按钮事件
    const pendingBtn = new ActionBtn({
      text: () => (this.data.is_pending ? '待审' : '已审'),
      adminOnly: true
    })
    pendingBtn.appendTo(this.$actions)
    pendingBtn.setClick(() => {
      this.adminEdit('pending', pendingBtn)
    })

    // 绑定删除按钮事件
    const delBtn = new ActionBtn({
      text: '删除',
      confirm: true,
      confirmText: '确认删除',
      adminOnly: true,
    })
    delBtn.appendTo(this.$actions)
    delBtn.setClick(() => {
      this.adminDelete(delBtn)
    })

    // 绑定置顶按钮事件
    const pinnedBtn = new ActionBtn({
      text: () => (this.data.is_pinned ? '取消置顶' : '置顶'),
      adminOnly: true
    })
    pinnedBtn.appendTo(this.$actions)
    pinnedBtn.setClick(() => {
      this.adminEdit('pinned', pendingBtn)
    })
  }
  //#endregion

  /** 刷新评论 UI */
  public refreshUI() {
    const originalEl = this.$el
    const newEl = this.render()
    originalEl.replaceWith(newEl) // 替换 document 中的的 elem
    this.playFadeInAnim()

    // 重建子评论元素
    this.eachComment(this.children, (child) => {
      child.parent?.getChildrenEl().appendChild(child.render())
      child.playFadeInAnim()
    })

    this.ctx.trigger('comments-loaded')
  }

  /** 遍历评论 */
  private eachComment(commentList: Comment[], action: (comment: Comment, levelList: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  getIsRoot() {
    return this.parent === null
  }

  getChildren() {
    return this.children
  }

  putChild(childC: Comment) {
    childC.parent = this
    childC.nestedNum = this.nestedNum + 1 // 嵌套层数 +1
    this.children.push(childC)

    this.getChildrenEl().appendChild(childC.getEl())
    childC.playFadeInAnim()

    // 内容限高
    childC.checkHeightLimitArea('content')
  }

  getChildrenEl() {
    if (this.$children === null) {
      // console.log(this.nestedNum)
      if (this.nestedNum < this.maxNestingNum) {
        this.$children = Utils.createElement('<div class="atk-comment-children"></div>')
        this.$main.appendChild(this.$children)
      } else if (this.parent) {
          this.$children = this.parent.getChildrenEl()
      }
    }

    return this.$children
  }

  getParent() {
    return this.parent
  }

  getParents() {
    const parents: Comment[] = []
    const once = (c: Comment) => {
      if (c.parent) {
        parents.push(c.parent)
        once(c.parent)
      }
    }

    once(this)
    return parents
  }

  getEl() {
    return this.$el
  }

  getData() {
    return this.data
  }

  getGravatarUrl() {
    return Utils.getGravatarURL(this.ctx, this.data.email_encrypted)
  }

  getContentMarked() {
    return Utils.marked(this.ctx, this.data.content)
  }

  getDateFormatted() {
    return Utils.timeAgo(new Date(this.data.date))
  }

  getUserUaBrowser() {
    const info = UADetect(this.data.ua)
    return `${info.browser} ${info.version}`
  }

  getUserUaOS() {
    const info = UADetect(this.data.ua)
    return `${info.os} ${info.osVersion}`
  }

  /** 渐出动画 */
  playFadeInAnim() {
    Ui.playFadeInAnim(this.$el)
  }

  /** 投票操作 */
  vote(type: 'up'|'down') {
    const actionBtn = type === 'up' ? this.voteBtnUp : this.voteBtnDown

    new Api(this.ctx).vote(this.data.id, `comment_${type}`)
    .then((v) => {
      this.data.vote_up = v.up
      this.data.vote_down = v.down
      this.voteBtnUp?.updateText()
      this.voteBtnDown?.updateText()
    })
    .catch((err) => {
      actionBtn?.setError(`投票失败`)
      console.log(err)
    })
  }

  /** 管理员 - 评论状态修改 */
  adminEdit(type: 'collapsed'|'pending'|'pinned', btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在修改中

    btnElem.setLoading(true, '修改中...')

    // 克隆并修改当前数据
    const modify = { ...this.data }
    if (type === 'collapsed') {
      modify.is_collapsed = !modify.is_collapsed
    } else if (type === 'pending') {
      modify.is_pending = !modify.is_pending
    } else if (type === 'pinned') {
      modify.is_pinned = !modify.is_pinned
    }

    new Api(this.ctx).commentEdit(modify).then((comment) => {
      btnElem.setLoading(false)

      // 刷新当前 Comment UI
      this.data = comment
      this.refreshUI()
      Ui.playFadeInAnim(this.$body)

      // 刷新 List UI
      this.ctx.trigger('list-refresh-ui')
    }).catch((err) => {
      console.error(err)
      btnElem.setError('修改失败')
    })
  }

  public onDelete?: (comment: Comment) => void

  /** 管理员 - 评论删除 */
  adminDelete(btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在删除中

    btnElem.setLoading(true, '删除中...')
    new Api(this.ctx).commentDel(this.data.id, this.data.site_name)
      .then(() => {
        btnElem.setLoading(false)
        if (this.onDelete) this.onDelete(this)
      })
      .catch((e) => {
        console.error(e)
        btnElem.setError('删除失败')
      })
  }

  public setUnread(val: boolean) {
    this.unread = val
    if (this.unread) this.$el.classList.add('atk-unread')
    else this.$el.classList.remove('atk-unread')
  }

  public setOpenURL(url: string) {
    if (!url) {
      this.openable = false
      this.$el.classList.remove('atk-openable')
    }

    this.openable = true
    this.openURL = url
    this.$el.classList.add('atk-openable')
  }

  /** 内容限高检测 */
  checkHeightLimit() {
    this.checkHeightLimitArea('content') // 评论内容限高
    this.checkHeightLimitArea('children') // 子评论部分限高（嵌套模式）
  }

  /** 目标内容限高检测 */
  checkHeightLimitArea(area: 'children'|'content') {
    // 参数准备
    const childrenMaxH = this.ctx.conf.heightLimit?.children
    const contentMaxH = this.ctx.conf.heightLimit?.content

    if (area === 'children' && !childrenMaxH) return
    if (area === 'content' && !contentMaxH) return

    // 限高
    let maxHeight: number
    if (area === 'children') maxHeight = childrenMaxH!
    if (area === 'content') maxHeight = contentMaxH!

    // 检测指定元素
    const checkEl = ($el?: HTMLElement|null) => {
      if (!$el) return

      // 是否超过高度
      if (Utils.getHeight($el) > maxHeight) {
        this.heightLimitAdd($el, maxHeight)
      }
    }

    // 执行限高检测
    if (area === 'children') {
      checkEl(this.$children)
    } else if (area === 'content') {
      checkEl(this.$content)
      checkEl(this.$replyTo)

      // 若有图片 · 图片加载完后再检测一次
      Utils.onImagesLoaded(this.$content, () => {
        checkEl(this.$content)
      })
      if (this.$replyTo) {
        Utils.onImagesLoaded(this.$replyTo, () => {
          checkEl(this.$replyTo)
        })
      }
    }
  }

  // 操作 · 取消限高
  heightLimitRemove($el: HTMLElement) {
    if (!$el) return
    if (!$el.classList.contains('atk-height-limit')) return

    $el.classList.remove('atk-height-limit')
    Array.from($el.children).forEach((e) => {
      if (e.classList.contains('atk-height-limit-btn')) e.remove()
    })
    $el.style.height = ''
    $el.style.overflow = ''
  }

  // 操作 · 内容限高
  heightLimitAdd($el: HTMLElement, maxHeight: number) {
    if (!$el) return
    if ($el.classList.contains('atk-height-limit')) return

    $el.classList.add('atk-height-limit')
    $el.style.height = `${maxHeight}px`
    $el.style.overflow = 'hidden'
    const $hideMoreOpenBtn = Utils.createElement(`<div class="atk-height-limit-btn">阅读更多</span>`)
    $hideMoreOpenBtn.onclick = (e) => {
      e.stopPropagation()
      this.heightLimitRemove($el)

      // 子评论数等于 1，直接取消限高
      const children = this.getChildren()
      if (children.length === 1) children[0].heightLimitRemove(children[0].$content)
    }
    $el.append($hideMoreOpenBtn)
  }
}

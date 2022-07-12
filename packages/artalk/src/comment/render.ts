import '../style/comment.less'

import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import ActionBtn from '../components/action-btn'
import CommentHTML from './comment.html?raw'
import Comment from './comment'

export default class CommentRender {
  private comment: Comment

  private get ctx() { return this.comment.ctx }
  private get data() { return this.comment.getData() }
  private get cConf() { return this.comment.getConf() }

  public $el!: HTMLElement
  public $main!: HTMLElement
  public $header!: HTMLElement
  public $headerNick!: HTMLElement
  public $headerBadgeWrap!: HTMLElement
  public $body!: HTMLElement
  public $content!: HTMLElement
  private $childrenWrap!: HTMLElement|null
  public $actions!: HTMLElement
  public voteBtnUp?: ActionBtn
  public voteBtnDown?: ActionBtn

  public $replyTo?: HTMLElement // 回复评论内容 (平铺下显示)
  public $replyAt?: HTMLElement // 回复 AT（层级嵌套下显示）

  public constructor(comment: Comment) {
    this.comment = comment
  }

  public render() {
    this.$el = Utils.createElement(CommentHTML)

    this.$main = this.$el.querySelector('.atk-main')!
    this.$header = this.$el.querySelector('.atk-header')!
    this.$body = this.$el.querySelector('.atk-body')!
    this.$content = this.$body.querySelector('.atk-content')!
    this.$actions = this.$el.querySelector('.atk-actions')!

    this.$el.setAttribute('data-comment-id', `${this.data.id}`)

    this.renderAvatar()
    this.renderHeader()
    this.renderContent()
    this.renderReplyAt()
    this.renderReplyTo()
    this.renderPending()
    this.renderActions()

    this.recoveryChildrenWrap()

    return this.$el
  }

  /** 初始化 - 评论头像 */
  private renderAvatar() {
    const $avatar = this.$el.querySelector<HTMLElement>('.atk-avatar')!
    const $avatarImg = Utils.createElement<HTMLImageElement>('<img />')
    $avatarImg.src = this.comment.getGravatarURL()
    if (this.data.link) {
      const $avatarA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
      $avatarA.href = Utils.isValidURL(this.data.link) ? this.data.link : `https://${this.data.link}`
      $avatarA.append($avatarImg)
      $avatar.append($avatarA)
    } else {
      $avatar.append($avatarImg)
    }
  }

  /** 初始化 - 评论信息 */
  private renderHeader() {
    this.renderHeader_Nick()
    this.renderHeader_VerifyBadge()
    this.renderHeader_Date()
    this.renderHeader_UABadge()
  }

  private renderHeader_Nick() {
    this.$headerNick = this.$el.querySelector<HTMLElement>('.atk-nick')!

    if (this.data.link) {
      const $nickA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
      $nickA.innerText = this.data.nick
      $nickA.href = Utils.isValidURL(this.data.link) ? this.data.link : `https://${this.data.link}`
      this.$headerNick.append($nickA)
    } else {
      this.$headerNick.innerText = this.data.nick
    }
  }

  private renderHeader_VerifyBadge() {
    this.$headerBadgeWrap = this.$el.querySelector<HTMLElement>('.atk-badge-wrap')!
    this.$headerBadgeWrap.innerHTML = ''

    const badgeText = this.data.badge_name
    const badgeColor = this.data.badge_color
    if (badgeText) {
      const $badge = Utils.createElement(`<span class="atk-badge"></span>`)
      $badge.innerText = badgeText.replace('管理员', this.ctx.$t('admin')) // i18n patch
      $badge.style.backgroundColor = badgeColor || ''
      this.$headerBadgeWrap.append($badge)
    }

    if (this.data.is_pinned) {
      const $pinnedBadge = Utils.createElement(`<span class="atk-pinned-badge">${this.ctx.$t('pin')}</span>`) // 置顶徽章
      this.$headerBadgeWrap.append($pinnedBadge)
    }
  }

  private renderHeader_Date() {
    const $date = this.$el.querySelector<HTMLElement>('.atk-date')!
    $date.innerText = this.comment.getDateFormatted()
    $date.setAttribute('data-atk-comment-date', String(+new Date(this.data.date)))
  }

  private renderHeader_UABadge() {
    if (this.ctx.conf.uaBadge) {
      let $uaWrap = this.$header.querySelector('atk-ua-wrap')
      if (!$uaWrap) {
        $uaWrap = Utils.createElement(`<span class="atk-ua-wrap"></span>`)
        this.$header.append($uaWrap)
      }

      $uaWrap.innerHTML = ''
      const $uaBrowser = Utils.createElement(`<span class="atk-ua ua-browser"></span>`)
      const $usOS = Utils.createElement(`<span class="atk-ua ua-os"></span>`)
      const uaInfo = this.comment.getUserUA()
      $uaBrowser.innerText = uaInfo.browser
      $usOS.innerText = uaInfo.os
      $uaWrap.append($uaBrowser)
      $uaWrap.append($usOS)
    }
  }

  /** 初始化 - 评论内容 */
  private renderContent() {
    if (!this.data.is_collapsed) {
      this.$content.innerHTML = this.comment.getContentMarked()
      this.$content.classList.remove('atk-hide', 'atk-collapsed')
      return
    }

    // 内容 & 折叠
    this.$content.classList.add('atk-hide', 'atk-type-collapsed')
    const collapsedInfoEl = Utils.createElement(`
      <div class="atk-collapsed">
        <span class="atk-text">${this.ctx.$t('collapsedMsg')}</span>
        <span class="atk-show-btn">${this.ctx.$t('expand')}</span>
      </div>`)
    this.$body.insertAdjacentElement('beforeend', collapsedInfoEl)

    const contentShowBtn = collapsedInfoEl.querySelector('.atk-show-btn')!
    contentShowBtn.addEventListener('click', (e) => {
      e.stopPropagation() // 防止穿透

      if (this.$content.classList.contains('atk-hide')) {
        this.$content.innerHTML = this.comment.getContentMarked()
        this.$content.classList.remove('atk-hide')
        Ui.playFadeInAnim(this.$content)
        contentShowBtn.innerHTML = this.ctx.$t('collapse')
      } else {
        this.$content.innerHTML = ''
        this.$content.classList.add('atk-hide')
        contentShowBtn.innerHTML = this.ctx.$t('expand')
      }
    })
  }

  /** 初始化 - 层级嵌套模式显示 At */
  private renderReplyAt() {
    if (this.cConf.isFlatMode || this.data.rid === 0) return // not 平铺模式 或 根评论
    if (!this.cConf.replyTo) return

    this.$replyAt = Utils.createElement(`<span class="atk-item atk-reply-at"><span class="atk-arrow"></span><span class="atk-nick"></span></span>`)
    this.$replyAt.querySelector<HTMLElement>('.atk-nick')!.innerText = `${this.cConf.replyTo.nick}`
    this.$replyAt.onclick = () => { this.comment.getActions().goToReplyComment() }

    this.$headerBadgeWrap.insertAdjacentElement('afterend', this.$replyAt)
  }

  /** 初始化 - 回复的对象 */
  private renderReplyTo() {
    if (!this.cConf.isFlatMode) return // 仅平铺模式显示
    if (!this.cConf.replyTo) return

    this.$replyTo = Utils.createElement(`
      <div class="atk-reply-to">
        <div class="atk-meta">${this.ctx.$t('reply')} <span class="atk-nick"></span>:</div>
        <div class="atk-content"></div>
      </div>`)
    const $nick = this.$replyTo.querySelector<HTMLElement>('.atk-nick')!
    $nick.innerText = `@${this.cConf.replyTo.nick}`
    $nick.onclick = () => { this.comment.getActions().goToReplyComment() }
    let replyContent = Utils.marked(this.ctx, this.cConf.replyTo.content)
    if (this.cConf.replyTo.is_collapsed) replyContent = `[${this.ctx.$t('collapsed')}]`
    this.$replyTo.querySelector<HTMLElement>('.atk-content')!.innerHTML = replyContent
    this.$body.prepend(this.$replyTo)
  }

  /** 初始化 - 待审核状态 */
  private renderPending() {
    if (!this.data.is_pending) return

    const pendingEl = Utils.createElement(`<div class="atk-pending">${this.ctx.$t('pendingMsg')}</div>`)
    this.$body.prepend(pendingEl)
  }

  /** 初始化 - 评论操作按钮 */
  private renderActions() {
    this.renderActions_Vote()
    this.renderActions_Reply()

    // 管理员操作
    this.renderActions_Collapse()
    this.renderActions_Moderator()
    this.renderActions_Pin()
    this.renderActions_Edit()
    this.renderActions_Del()
  }

  // 操作按钮 - 投票
  private renderActions_Vote() {
    if (!this.ctx.conf.vote) return // 关闭投票功能

    // 赞同按钮
    this.voteBtnUp = new ActionBtn(this.ctx, () => `${this.ctx.$t('voteUp')} (${this.data.vote_up || 0})`).appendTo(this.$actions)
    this.voteBtnUp.setClick(() => {
      this.comment.getActions().vote('up')
    })

    // 反对按钮
    if (this.ctx.conf.voteDown) {
      this.voteBtnDown = new ActionBtn(this.ctx, () => `${this.ctx.$t('voteDown')} (${this.data.vote_down || 0})`).appendTo(this.$actions)
      this.voteBtnDown.setClick(() => {
        this.comment.getActions().vote('down')
      })
    }
  }

  // 操作按钮 - 回复
  private renderActions_Reply() {
    if (!this.data.is_allow_reply) return // 不允许回复

    const replyBtn = Utils.createElement(`<span>${this.ctx.$t('reply')}</span>`)
    this.$actions.append(replyBtn)
    replyBtn.addEventListener('click', (e) => {
      e.stopPropagation() // 防止穿透
      if (!this.cConf.onReplyBtnClick) {
        this.ctx.replyComment(this.data, this.$el)
      } else {
        this.cConf.onReplyBtnClick()
      }
    })
  }

  // 操作按钮 - 折叠
  private renderActions_Collapse() {
    const collapseBtn = new ActionBtn(this.ctx, {
      text: () => (this.data.is_collapsed ? this.ctx.$t('expand') : this.ctx.$t('collapse')),
      adminOnly: true
    })
    collapseBtn.appendTo(this.$actions)
    collapseBtn.setClick(() => {
      this.comment.getActions().adminEdit('collapsed', collapseBtn)
    })
  }

  // 操作按钮 - 审核
  private renderActions_Moderator() {
    const pendingBtn = new ActionBtn(this.ctx, {
      text: () => (this.data.is_pending ? this.ctx.$t('pending') : this.ctx.$t('approved')),
      adminOnly: true
    })
    pendingBtn.appendTo(this.$actions)
    pendingBtn.setClick(() => {
      this.comment.getActions().adminEdit('pending', pendingBtn)
    })
  }

  // 操作按钮 - 置顶
  private renderActions_Pin() {
    const pinnedBtn = new ActionBtn(this.ctx, {
      text: () => (this.data.is_pinned ? this.ctx.$t('unpin') : this.ctx.$t('pin')),
      adminOnly: true
    })
    pinnedBtn.appendTo(this.$actions)
    pinnedBtn.setClick(() => {
      this.comment.getActions().adminEdit('pinned', pinnedBtn)
    })
  }

  // 操作按钮 - 编辑
  private renderActions_Edit() {
    const editBtn = new ActionBtn(this.ctx, {
      text: this.ctx.$t('edit'),
      adminOnly: true
    })
    editBtn.appendTo(this.$actions)
    editBtn.setClick(() => {
      this.ctx.editComment(this.data, this.$el)
    })
  }

  // 操作按钮 - 删除
  private renderActions_Del() {
    const delBtn = new ActionBtn(this.ctx, {
      text: this.ctx.$t('delete'),
      confirm: true,
      confirmText: this.ctx.$t('deleteConfirm'),
      adminOnly: true,
    })
    delBtn.appendTo(this.$actions)
    delBtn.setClick(() => {
      this.comment.getActions().adminDelete(delBtn)
    })
  }

  /** 内容限高检测 */
  public checkHeightLimit() {
    this.checkHeightLimitArea('content') // 评论内容限高
    this.checkHeightLimitArea('children') // 子评论部分限高（嵌套模式）
  }

  /** 目标内容限高检测 */
  public checkHeightLimitArea(area: 'children'|'content') {
    // 参数准备
    const childrenMaxH = this.ctx.conf.heightLimit.children
    const contentMaxH = this.ctx.conf.heightLimit.content

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
      checkEl(this.$childrenWrap)
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

  /** 移除限高 */
  private heightLimitRemove($el: HTMLElement) {
    if (!$el) return
    if (!$el.classList.contains('atk-height-limit')) return

    $el.classList.remove('atk-height-limit')
    Array.from($el.children).forEach((e) => {
      if (e.classList.contains('atk-height-limit-btn')) e.remove()
    })
    $el.style.height = ''
    $el.style.overflow = ''
  }

  /** 子评论区域移除限高 */
  public heightLimitRemoveForChildren() {
    if (!this.$childrenWrap) return
    this.heightLimitRemove(this.$childrenWrap)
  }

  /** 内容限高区域新增 */
  private heightLimitAdd($el: HTMLElement, maxHeight: number) {
    if (!$el) return
    if ($el.classList.contains('atk-height-limit')) return

    $el.classList.add('atk-height-limit')
    $el.style.height = `${maxHeight}px`
    $el.style.overflow = 'hidden'
    const $hideMoreOpenBtn = Utils.createElement(`<div class="atk-height-limit-btn">${this.ctx.$t('readMore')}</span>`)
    $hideMoreOpenBtn.onclick = (e) => {
      e.stopPropagation()
      this.heightLimitRemove($el)

      // 子评论数等于 1，直接取消限高
      const children = this.comment.getChildren()
      if (children.length === 1) children[0].getRender().heightLimitRemove(children[0].getRender().$content)
    }
    $el.append($hideMoreOpenBtn)
  }

  /** 渐出动画 */
  playFadeAnim() {
    Ui.playFadeInAnim(this.comment.getRender().$el)
  }

  /** 渐出动画 · 评论内容区域 */
  playFadeAnimForBody() {
    Ui.playFadeInAnim(this.comment.getRender().$body)
  }

  /** 获取子评论 Wrap */
  public getChildrenWrap() {
    return this.$childrenWrap
  }

  /** 初始化子评论区域 Wrap */
  public renderChildrenWrap() {
    if (!this.$childrenWrap) {
      this.$childrenWrap = Utils.createElement('<div class="atk-comment-children"></div>')
      this.$main.append(this.$childrenWrap)
    }

    return this.$childrenWrap
  }

  /** 恢复原有的子评论区域 Wrap */
  public recoveryChildrenWrap() {
    if (this.$childrenWrap) {
      this.$main.append(this.$childrenWrap)
    }
  }

  /** 设置已读 */
  public setUnread(val: boolean) {
    if (val) this.$el.classList.add('atk-unread')
    else this.$el.classList.remove('atk-unread')
  }

  /** 设置为可点击的评论 */
  public setOpenable(val: boolean) {
    if (val) this.$el.classList.add('atk-openable')
    else this.$el.classList.remove('atk-openable')
  }

  /** 设置点击评论打开置顶 URL */
  public setOpenURL(url: string) {
    this.setOpenable(true)
    this.$el.onclick = (evt) => {
      evt.preventDefault()
      window.open(url)

      if (this.cConf.openEvt) this.cConf.openEvt()
    }
  }

  /** 设置点击评论时的操作 */
  public setOpenAction(action: () => void) {
    this.setOpenable(true)
    this.$el.onclick = (evt) => {
      evt.preventDefault()
      action()
    }
  }
}

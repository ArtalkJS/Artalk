import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

import CommentHTML from './comment.html?raw'
import Comment from './comment'
import RenderCtx from './render-ctx'
import loadRenders from './renders'

export default class CommentRender extends RenderCtx {
  public constructor(comment: Comment) {
    super(comment)
  }

  public render() {
    // init ui elements
    this.$el = Utils.createElement(CommentHTML)

    this.$main = this.$el.querySelector('.atk-main')!
    this.$header = this.$el.querySelector('.atk-header')!
    this.$body = this.$el.querySelector('.atk-body')!
    this.$content = this.$body.querySelector('.atk-content')!
    this.$actions = this.$el.querySelector('.atk-actions')!

    this.$el.setAttribute('data-comment-id', `${this.data.id}`)

    // call all renders
    loadRenders(this)

    this.recoveryChildrenWrap()

    return this.$el
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

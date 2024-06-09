import ActionBtn from '../components/action-btn'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

import * as HeightLimit from './height-limit'
import CommentHTML from './comment.html?raw'
import loadRenders from './renders'
import type { CommentNode } from '.'

export default class Render {
  public comment: CommentNode

  public get data() {
    return this.comment.getData()
  }
  public get opts() {
    return this.comment.getOpts()
  }

  public $el!: HTMLElement
  public $main!: HTMLElement
  public $header!: HTMLElement
  public $headerNick!: HTMLElement
  public $headerBadgeWrap!: HTMLElement
  public $body!: HTMLElement
  public $content!: HTMLElement
  public $childrenWrap!: HTMLElement | null
  public $actions!: HTMLElement
  public voteBtnUp?: ActionBtn
  public voteBtnDown?: ActionBtn

  public $replyTo?: HTMLElement // 回复评论内容 (平铺下显示)
  public $replyAt?: HTMLElement // 回复 AT（层级嵌套下显示）

  public constructor(comment: CommentNode) {
    this.comment = comment
  }

  /**
   * Render the comment ui
   *
   * If comment data is updated, call this method to re-render the comment ui.
   * The method will be called multiple times, so it should be idempotent.
   *
   * Renders may add event listeners to the comment ui, so it should be called only once or override the original.
   * Please be aware of the memory leak caused by the event listener.
   */
  public render() {
    // init ui elements
    this.$el = Utils.createElement(CommentHTML)

    this.$main = this.$el.querySelector('.atk-main')!
    this.$header = this.$el.querySelector('.atk-header')!
    this.$body = this.$el.querySelector('.atk-body')!
    this.$content = this.$body.querySelector('.atk-content')!
    this.$actions = this.$el.querySelector('.atk-actions')!

    this.$el.setAttribute('id', `atk-comment-${this.data.id}`)

    // call all renders
    loadRenders(this)

    // append children wrap if exists
    if (this.$childrenWrap) {
      this.$main.append(this.$childrenWrap)
    }

    return this.$el
  }

  /** 内容限高检测 */
  public checkHeightLimit() {
    const conf = this.opts.heightLimit
    if (!conf || !conf.content || !conf.children) return // 关闭限高

    const contentMaxH = conf.content
    const childrenMaxH = conf.children

    HeightLimit.check(
      {
        afterExpandBtnClick: () => {
          // 子评论数仅有 1，直接取消限高
          const children = this.comment.getChildren()
          if (children.length === 1)
            HeightLimit.disposeHeightLimit(children[0].getRender().$content)
        },
        scrollable: conf.scrollable,
      },
      [
        // 评论内容限高
        { el: this.$content, max: contentMaxH, imgCheck: true },
        { el: this.$replyTo, max: contentMaxH, imgCheck: true },
        // 子评论区域限高（仅嵌套模式）
        { el: this.$childrenWrap, max: childrenMaxH, imgCheck: false },
      ],
    )
  }

  /** 子评论区域移除限高 */
  public heightLimitRemoveForChildren() {
    if (!this.$childrenWrap) return
    HeightLimit.disposeHeightLimit(this.$childrenWrap)
  }

  /** 渐出动画 */
  playFadeAnim() {
    Ui.playFadeInAnim(this.comment.getRender().$el)
  }

  /** 渐出动画 · 评论内容区域 */
  playFadeAnimForBody() {
    Ui.playFadeInAnim(this.comment.getRender().$body)
  }

  /** Perform the flash animation */
  playFlashAnim() {
    // Make sure the class is removed
    this.$el.classList.remove('atk-flash-once')
    window.setTimeout(() => {
      // Add the class to perform the animation
      this.$el.classList.add('atk-flash-once')
    }, 150)
  }

  /** 获取子评论 Wrap */
  public getChildrenWrap() {
    if (!this.$childrenWrap) {
      // if not exists, create a new one
      this.$childrenWrap = Utils.createElement('<div class="atk-comment-children"></div>')
      this.$main.append(this.$childrenWrap)
    }
    return this.$childrenWrap
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
      evt.stopPropagation()
      window.open(url)
    }
  }

  /** 设置点击评论时的操作 */
  public setOpenAction(action: () => void) {
    this.setOpenable(true)
    this.$el.onclick = (evt) => {
      evt.stopPropagation()
      action()
    }
  }
}

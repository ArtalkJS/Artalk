import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'

import CommentHTML from './comment.html?raw'
import Comment from './comment'
import RenderCtx from './render-ctx'
import loadRenders from './renders'
import * as HeightLimit from './height-limit'

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

    this.$el.setAttribute('id', `atk-comment-${this.data.id}`)

    // call all renders
    loadRenders(this)

    this.recoveryChildrenWrap()

    return this.$el
  }

  /** 内容限高检测 */
  public checkHeightLimit() {
    const conf = this.ctx.conf.heightLimit
    if (!conf || !conf.content || !conf.children) return // 关闭限高

    const contentMaxH = conf.content
    const childrenMaxH = conf.children

    HeightLimit.check({
      postExpandBtnClick: () => {
        // 子评论数仅有 1，直接取消限高
        const children = this.comment.getChildren()
        if (children.length === 1) HeightLimit.disposeHeightLimit(children[0].getRender().$content)
      },
      scrollable: conf.scrollable
    }, [
      // 评论内容限高
      { el: this.$content, max: contentMaxH, imgContains: true },
      { el: this.$replyTo, max: contentMaxH, imgContains: true },
      // 子评论区域限高（仅嵌套模式）
      { el: this.$childrenWrap, max: childrenMaxH, imgContains: false }
    ])
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

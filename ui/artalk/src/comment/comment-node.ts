import type { CommentData, ArtalkConfig, ContextApi } from '@/types'
import $t from '@/i18n'
import type { Api } from '../api'
import * as Ui from '../lib/ui'
import * as Utils from '../lib/utils'
import marked from '../lib/marked'
import UADetect from '../lib/detect'
import CommentUI from './render'
import CommentActions from './actions'

export interface CommentOptions {
  // Hooks
  onAfterRender?: () => void
  onDelete?: Function

  /** The comment being replied to (linked comment) */
  replyTo?: CommentData

  // Referenced from ArtalkConfig
  flatMode: boolean
  vote: boolean
  voteDown: boolean
  uaBadge: boolean
  nestMax: number
  gravatar: ArtalkConfig['gravatar']
  heightLimit: ArtalkConfig['heightLimit']
  avatarURLBuilder: ArtalkConfig['avatarURLBuilder']
  scrollRelativeTo: ArtalkConfig['scrollRelativeTo']

  // TODO: Move to plugin folder and remove from core
  getApi: () => Api
  replyComment: ContextApi['replyComment']
  editComment: ContextApi['editComment']
}

export default class CommentNode {
  $el?: HTMLElement

  private renderInstance: CommentUI
  private actionInstance: CommentActions

  private data: CommentData
  private opts: CommentOptions

  private parent: CommentNode | null
  private children: CommentNode[] = []

  private nestCurt: number // 当前嵌套层数

  constructor(data: CommentData, opts: CommentOptions) {
    this.opts = opts
    this.data = { ...data }
    this.data.date = this.data.date.replace(/-/g, '/') // 解决 Safari 日期解析 NaN 问题

    this.parent = null
    this.nestCurt = 1 // 现在已嵌套 n 层

    this.actionInstance = new CommentActions(this)
    this.renderInstance = new CommentUI(this)
  }

  /** 渲染 UI */
  public render() {
    const newEl = this.renderInstance.render()

    if (this.$el) this.$el.replaceWith(newEl)
    this.$el = newEl
    // Please be aware of the memory leak, the $el may be replaced multiple times.
    // If somewhere else has a reference to the old $el, it will cause a memory leak.
    // So it's limited to use the $el reference by `getEl()`.
    // The `getEl()` will always return the latest $el.

    if (this.opts.onAfterRender) this.opts.onAfterRender()
  }

  /** 获取评论操作实例对象 */
  public getActions() {
    return this.actionInstance
  }

  /** 获取评论渲染器实例对象 */
  public getRender() {
    return this.renderInstance
  }

  /** 获取评论数据 */
  public getData() {
    return this.data
  }

  /** 设置数据 */
  public setData(data: CommentData) {
    this.data = data

    this.render()
    this.getRender().playFadeAnimForBody()
  }

  /** 获取父评论 */
  public getParent() {
    return this.parent
  }

  /** 获取所有子评论 */
  public getChildren() {
    return this.children
  }

  /** 获取当前嵌套层数 */
  public getNestCurt() {
    return this.nestCurt
  }

  /** 判断是否为根评论 */
  public getIsRoot() {
    return this.data.rid === 0
  }

  /** 获取评论 ID */
  public getID() {
    return this.data.id
  }

  /** 置入子评论 */
  public putChild(childNode: CommentNode, insertMode: 'append' | 'prepend' = 'append') {
    childNode.parent = this
    childNode.nestCurt = this.nestCurt + 1 // 嵌套层数 +1
    this.children.push(childNode)

    const $childrenWrap = this.getChildrenWrapEl()
    const $childComment = childNode.getEl()
    if (insertMode === 'append') $childrenWrap.append($childComment)
    else if (insertMode === 'prepend') $childrenWrap.prepend($childComment)

    childNode.getRender().playFadeAnim()
    childNode.getRender().checkHeightLimit() // 内容限高
  }

  /** 获取存放子评论的元素对象 */
  public getChildrenWrapEl(): HTMLElement {
    // console.log(this.nestCurt)
    if (this.nestCurt >= this.opts.nestMax) {
      return this.parent!.getChildrenWrapEl()
    }
    return this.getRender().getChildrenWrap()
  }

  /** 获取所有父评论 */
  public getParents() {
    const flattenParents: CommentNode[] = []
    let parent = this.parent
    while (parent) {
      flattenParents.push(parent)
      parent = parent.getParent()
    }
    return flattenParents
  }

  /**
   * Get the element of the comment
   *
   * The `getEl()` will always return the latest $el after calling `render()`.
   * Please be aware of the memory leak if you use the $el reference directly.
   */
  public getEl() {
    if (!this.$el) throw new Error('comment element not initialized before `getEl()`')
    return this.$el
  }

  /**
   * Focus on the comment
   *
   * Scroll to the comment and perform flash animation
   */
  focus() {
    if (!this.$el) throw new Error('comment element not initialized before `focus()`')

    // 若父评论存在 “子评论部分” 限高，取消限高
    this.getParents().forEach((p) => {
      p.getRender().heightLimitRemoveForChildren()
    })

    // Scroll to comment
    this.scrollIntoView()

    // Perform flash animation
    this.getRender().playFlashAnim()
  }

  scrollIntoView() {
    this.$el &&
      Ui.scrollIntoView(this.$el, false, this.opts.scrollRelativeTo && this.opts.scrollRelativeTo())
  }

  /**
   * Remove the comment node
   */
  remove() {
    this.$el?.remove()
  }

  /** 获取 Gravatar 头像 URL */
  public getGravatarURL() {
    return Utils.getGravatarURL({
      mirror: this.opts.gravatar.mirror,
      params: this.opts.gravatar.params,
      emailMD5: this.data.email_encrypted,
    })
  }

  /** 获取评论 markdown 解析后的内容 */
  public getContentMarked() {
    return marked(this.data.content)
  }

  /** 获取格式化后的日期 */
  public getDateFormatted() {
    return Utils.timeAgo(new Date(this.data.date), $t)
  }

  /** 获取用户 UserAgent 信息 */
  public getUserUA() {
    const info = UADetect(this.data.ua)
    return {
      browser: `${info.browser} ${info.version}`,
      os: `${info.os} ${info.osVersion}`,
    }
  }

  /** 获取配置 */
  public getOpts() {
    return this.opts
  }
}

import { CommentData } from '~/types/artalk-data'
import Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import marked from '../lib/marked'
import UADetect from '../lib/detect'
import CommentUI from './render'
import CommentActions from './actions'

export interface CommentConf {
  isUnread?: boolean
  openURL?: string
  isFlatMode: boolean
  replyTo?: CommentData
  afterRender?: () => void
  openEvt?: () => void
  onReplyBtnClick?: Function
  onDelete?: Function
}

export default class Comment extends Component {
  private renderInstance: CommentUI
  private actionInstance: CommentActions

  private data: CommentData
  private cConf: CommentConf

  private parent: Comment|null
  private children: Comment[] = []

  private nestCurt: number // 当前嵌套层数
  private nestMax: number  // 最大嵌套层数

  constructor(ctx: Context, data: CommentData, conf: CommentConf) {
    super(ctx)

    // 最大嵌套数
    this.nestMax = ctx.conf.nestMax || 3

    this.cConf = conf
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

    if (this.cConf.afterRender) this.cConf.afterRender()
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
  public putChild(childC: Comment, insertMode: 'append'|'prepend' = 'append') {
    childC.parent = this
    childC.nestCurt = this.nestCurt + 1 // 嵌套层数 +1

    this.children.push(childC)

    const $children = this.getChildrenEl()
    if (insertMode === 'append') $children.append(childC.getEl())
    else if (insertMode === 'prepend') $children.prepend(childC.getEl())

    childC.getRender().playFadeAnim()

    // 内容限高
    childC.getRender().checkHeightLimit()
  }

  /** 获取存放子评论的元素对象 */
  public getChildrenEl(): HTMLElement {
    let $children = this.getRender().getChildrenWrap()

    if (!$children) {
      // console.log(this.nestCurt)
      if (this.nestCurt < this.nestMax) {
        $children = this.getRender().renderChildrenWrap()
      } else {
        $children = this.parent!.getChildrenEl()
      }
    }

    return $children
  }

  /** 获取所有父评论 */
  public getParents() {
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

  /** 获取评论元素对象 */
  public getEl() {
    return this.$el
  }

  /** 获取 Gravatar 头像 URL */
  public getGravatarURL() {
    return Utils.getGravatarURL(this.ctx, this.data.email_encrypted)
  }

  /** 获取评论 markdown 解析后的内容 */
  public getContentMarked() {
    return marked(this.ctx, this.data.content)
  }

  /** 获取格式化后的日期 */
  public getDateFormatted() {
    return Utils.timeAgo(new Date(this.data.date), this.ctx)
  }

  /** 获取用户 UserAgent 信息 */
  public getUserUA() {
    const info = UADetect(this.data.ua)
    return {
      browser: `${info.browser} ${info.version}`,
      os: `${info.os} ${info.osVersion}`
    }
  }

  /** 获取配置 */
  public getConf() {
    return this.cConf
  }
}

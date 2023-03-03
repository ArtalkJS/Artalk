import Comment from './comment'
import ActionBtn from '../components/action-btn'

export default class RenderCtx {
  public comment: Comment

  public get ctx() { return this.comment.ctx }
  public get data() { return this.comment.getData() }
  public get conf() { return this.comment.conf }
  public get cConf() { return this.comment.getConf() }

  public $el!: HTMLElement
  public $main!: HTMLElement
  public $header!: HTMLElement
  public $headerNick!: HTMLElement
  public $headerBadgeWrap!: HTMLElement
  public $body!: HTMLElement
  public $content!: HTMLElement
  public $childrenWrap!: HTMLElement|null
  public $actions!: HTMLElement
  public voteBtnUp?: ActionBtn
  public voteBtnDown?: ActionBtn

  public $replyTo?: HTMLElement // 回复评论内容 (平铺下显示)
  public $replyAt?: HTMLElement // 回复 AT（层级嵌套下显示）

  public constructor(comment: Comment) {
    this.comment = comment
  }
}

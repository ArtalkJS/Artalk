import { ListData, CommentData } from '~/types/artalk-data'
import Context from '~/types/context'
import type ArtalkConfig from '~/types/artalk-config'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment from '../comment'
import PgHolder, { TPgMode } from './paginator'
import * as ListNest from './list-nest'
import * as ListUi from './list-ui'
import { handleBackendRefConf } from '../config'

interface ListOptions {
  /** Flat mode */
  flatMode?: boolean

  /** Pagination mode */
  pageMode?: TPgMode

  /** Page size */
  pageSize?: number

  /** 监听指定元素上的滚动 */
  scrollListenerAt?: HTMLElement

  /** 翻页归位到指定元素 */
  repositionAt?: HTMLElement

  /** 启用列表未读高亮 */
  unreadHighlight?: boolean

  /** Sort condition in nest mode */
  nestSortBy?: ListNest.SortByType

  /** Text to show when no comment */
  noCommentText?: string

  // 一些 Hook 函数
  // ----------------
  renderComment?: (comment: Comment) => void
  paramsEditor?: (params: any) => void
  onAfterLoad?: (data: ListData) => void
}

// TODO public 的可配置成员字段放到一个 Options 对象，而不是直接暴露，很混乱
export default class ListLite extends Component {
  /** The options of List */
  protected options: ListOptions = {}
  getOptions() {
    return this.options
  }

  /** 列表评论集区域元素 */
  protected $commentsWrap: HTMLElement

  /** 最后一次请求得到的列表数据 */
  protected data?: ListData

  /** 是否处于加载中状态 */
  protected isLoading: boolean = false

  /** 配置是否已加载 */
  private confLoaded = false

  /** 嵌套模式下的排序方式 */
  private getNestSortBy(): ListNest.SortByType {
    return this.options.nestSortBy || this.ctx.conf.nestSort || 'DATE_ASC'
  }

  /** 平铺模式 */
  private getFlatMode(): boolean {
    if (this.options.flatMode !== undefined)
      return this.options.flatMode

    // 配置开启平铺模式
    if (this.ctx.conf.flatMode === true || Number(this.ctx.conf.nestMax) <= 1)
      return true

    // 自动判断启用平铺模式
    if (this.ctx.conf.flatMode === 'auto' && window.matchMedia("(max-width: 768px)").matches)
      return true

    return false
  }

  /** 分页方式持有者 */
  public pgHolder?: PgHolder

  /** 分页方式 */
  private getPageMode(): TPgMode {
    return this.options.pageMode || (this.conf.pagination.readMore ? 'read-more' : 'pagination')
  }

  private setPageMode(mode: TPgMode) {
    this.options.pageMode = mode
    this.pgHolder?.setMode(mode)
  }

  /** 每页数量 (每次请求获取量) */
  private getPageSize(): number {
    return this.options.pageSize || this.conf.pagination.pageSize
  }

  constructor (ctx: Context, options: ListOptions = {}) {
    super(ctx)

    this.options = options

    // 初始化元素
    this.$el = Utils.createElement(
    `<div class="atk-list-lite">
      <div class="atk-list-comments-wrap"></div>
    </div>`)
    this.$commentsWrap = this.$el.querySelector('.atk-list-comments-wrap')!

    // 评论时间自动更新
    window.setInterval(() => {
      this.$el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach(el => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)), this.ctx)
      })
    }, 30 * 1000) // 30s 更新一次

    // 事件监听
    this.ctx.on('conf-loaded', () => {
    })
  }

  public getData() {
    return this.data
  }

  public clearData() {
    this.data = undefined
  }

  public getLoading() {
    return this.isLoading
  }

  public getCommentsWrapEl() {
    return this.$commentsWrap
  }

  /** 设置加载状态 */
  public setLoading(val: boolean, isFirstLoad: boolean = false) {
    this.isLoading = val
    if (isFirstLoad) {
      Ui.setLoading(val, this.$el)
      return
    }

    this.pgHolder?.setLoading(val)
  }

  /** 评论获取 */
  public async fetchComments(offset: number) {
    if (this.isLoading) return

    const isFirstLoad = (offset === 0)
    const setLoading = (val: boolean) => this.setLoading(val, isFirstLoad)

    // 加载动画
    setLoading(true)

    // 事件通知（开始加载评论）
    this.ctx.trigger('list-load')

    // 清空评论（按钮加载更多的第一页、每次加载分页页面）
    const pageMode = this.getPageMode()
    if ((pageMode === 'read-more' && isFirstLoad) || pageMode === 'pagination') {
      this.ctx.clearAllComments()
    }

    // 请求评论数据
    let listData: ListData
    try {
      // 执行请求
      listData = await this.ctx.getApi().comment
        .get(offset, this.getPageSize(), this.getFlatMode(), this.options.paramsEditor)
    } catch (e: any) {
      this.onError(e.msg || String(e), offset, e.data)
      throw e
    } finally {
      setLoading(false)
    }

    // 清除原有错误
    Ui.setError(this.$el, null)

    // 装载数据
    try {
      this.onLoad(listData, offset)
    } catch (e: any) {
      this.onError(String(e), offset)
      throw e
    } finally {
      setLoading(false)
    }
  }

  protected onLoad(data: ListData, offset: number) {
    this.data = data

    // 装载后端提供的配置
    this.loadConf(data)

    // 导入数据
    this.importComments(data.comments)

    // 分页功能
    this.refreshPagination(offset, (this.getFlatMode() ? data.total : data.total_roots))

    // 加载后事件
    this.refreshUI()

    // 更新页面数据
    this.ctx.updatePage(data.page)

    // 未读消息提示功能
    this.ctx.updateUnreadList(data.unread || [])

    // 事件触发，列表已加载
    this.ctx.trigger('list-loaded')

    // Hook 函数调用
    this.options.onAfterLoad && this.options.onAfterLoad(data)
  }

  private loadConf(data: ListData) {
    if (!this.confLoaded) { // 仅应用一次配置
      let conf: Partial<ArtalkConfig> = {
        apiVersion: data.api_version.version
      }

      // reference conf from backend
      if (this.conf.useBackendConf) {
        conf = { ...conf, ...handleBackendRefConf(data.conf.frontend_conf) }
      }

      this.ctx.updateConf(conf)
      this.confLoaded = true
    }
  }

  /** 分页模式 */
  private refreshPagination(offset: number, total: number) {
    // 初始化
    if (!this.pgHolder) {
      this.pgHolder = new PgHolder({
        list: this,
        mode: this.getPageMode(),
        pageSize: this.getPageSize(),
        total
      })
    }

    // 更新
    this.pgHolder?.update(offset, total)
  }

  /** 错误处理 */
  protected onError(msg: any, offset: number, errData?: any) {
    if (!this.confLoaded) {
      this.ctx.updateConf({})
    }

    msg = String(msg)
    console.error(msg)

    // 加载更多按钮显示错误
    if (offset !== 0 && this.getPageMode() === 'read-more') {
      this.pgHolder?.showErr(this.$t('loadFail'))
      return
    }

    // 显示错误对话框
    Ui.setError(this.$el, ListUi.renderErrorDialog(this, msg, errData))
  }

  /** 刷新界面 */
  public refreshUI() {
    ListUi.refreshUI(this)
  }

  /** 重新加载列表 */
  public reload() {
    this.fetchComments(0)
  }

  /** 创建新评论 */
  protected createComment(cData: CommentData, ctxData?: CommentData[]) {
    if (!ctxData) ctxData = this.ctx.getCommentDataList()

    const comment = new Comment(this.ctx, cData, {
      isFlatMode: this.getFlatMode(),
      afterRender: () => {
        this.options.renderComment && this.options.renderComment(comment)
      },
      onDelete: (c: Comment) => {
        this.ctx.deleteComment(c)
        this.refreshUI()
      },
      replyTo: (cData.rid ? ctxData.find(c => c.id === cData.rid) : undefined)
    })

    // 渲染元素
    comment.render()

    // 放入 comment 总表中
    this.ctx.getCommentList().push(comment)

    return comment
  }

  /** 导入评论 */
  public importComments(srcData: CommentData[]) {
    if (this.getFlatMode()) {
      srcData.forEach((commentData: CommentData) => {
        this.putCommentFlatMode(commentData, srcData, 'append')
      })
    } else {
      this.importCommentsNest(srcData)
    }
  }

  // 导入评论 · 嵌套模式
  private importCommentsNest(srcData: CommentData[]) {
    // 遍历 root 评论
    const rootNodes = ListNest.makeNestCommentNodeList(srcData, this.getNestSortBy(), this.conf.nestMax)
    rootNodes.forEach((rootNode: ListNest.CommentNode) => {
      const rootC = this.createComment(rootNode.comment, srcData)

      // 显示并播放渐入动画
      this.$commentsWrap.appendChild(rootC.getEl())
      rootC.getRender().playFadeAnim()

      // 加载子评论
      const that = this
      ;(function loadChildren(parentC: Comment, parentNode: ListNest.CommentNode) {
        parentNode.children.forEach((node: ListNest.CommentNode) => {
          const childD = node.comment
          const childC = that.createComment(childD, srcData)

          // 插入到父评论中
          parentC.putChild(childC)

          // 递归加载子评论
          loadChildren(childC, node)
        })
      })(rootC, rootNode)

      // 限高检测
      rootC.getRender().checkHeightLimit()
    })
  }

  /** 导入评论 · 平铺模式 */
  private putCommentFlatMode(cData: CommentData, ctxData: CommentData[], insertMode: 'append'|'prepend') {
    if (cData.is_collapsed) cData.is_allow_reply = false
    const comment = this.createComment(cData, ctxData)

    // 可见评论添加到界面
    // 注：不可见评论用于显示 “引用内容”
    if (cData.visible) {
      if (insertMode === 'append') this.$commentsWrap.append(comment.getEl())
      if (insertMode === 'prepend') this.$commentsWrap.prepend(comment.getEl())
      comment.getRender().playFadeAnim()
    }

    // 平铺评论插入后 · 内容限高检测
    comment.getRender().checkHeightLimit()

    return comment
  }

  /** 新增评论 · 首部添加 */
  public insertComment(commentData: CommentData) {
    if (!this.getFlatMode()) {
      // 嵌套模式
      const comment = this.createComment(commentData)

      if (commentData.rid === 0) {
        // root评论 新增
        this.$commentsWrap.prepend(comment.getEl())
      } else {
        // 子评论 新增
        const parent = this.ctx.findComment(commentData.rid)
        if (parent) {
          parent.putChild(comment, (this.getNestSortBy() === 'DATE_ASC' ? 'append' : 'prepend'))

          // 若父评论存在 “子评论部分” 限高，取消限高
          comment.getParents().forEach((p) => {
            p.getRender().heightLimitRemoveForChildren()
          })
        }
      }

      comment.getRender().checkHeightLimit()

      Ui.scrollIntoView(comment.getEl()) // 滚动到可以见
      comment.getRender().playFadeAnim() // 播放评论渐出动画
    } else {
      // 平铺模式
      const comment = this.putCommentFlatMode(commentData, this.ctx.getCommentDataList(), 'prepend')
      Ui.scrollIntoView(comment.getEl()) // 滚动到可见
    }

    // 评论数增加 1
    if (this.data) this.data.total += 1

    // 评论新增后
    this.refreshUI()
    this.ctx.trigger('list-loaded')
    this.ctx.trigger('list-inserted', commentData)
  }

  /** 更新评论 */
  public updateComment(commentData: CommentData) {
    const comment = this.ctx.findComment(commentData.id)
    comment && comment.setData(commentData)
  }
}

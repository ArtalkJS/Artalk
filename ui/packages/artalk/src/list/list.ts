import { ListData, CommentData } from '~/types/artalk-data'
import Context from '~/types/context'
import type ArtalkConfig from '~/types/artalk-config'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment from '../comment/comment'
import type { ListOptions } from './options'
import PgHolder, { TPgMode } from './paginator'
import * as ListNest from './nest'
import ListHTML from './list.html?raw'
import ListLayout from './layout'
import { handleBackendRefConf } from '../config'

export default class List extends Component {
  /** The options of List */
  protected options: ListOptions = {}

  protected layout: ListLayout

  /** 列表评论集区域元素 */
  $commentsWrap?: HTMLElement

  /** 最后一次请求得到的列表数据 */
  protected data?: ListData

  /** 是否处于加载中状态 */
  protected isLoading: boolean = false

  /** 配置是否已加载 */
  private confLoaded = false

  /** 分页方式持有者 */
  public pgHolder?: PgHolder

  constructor (ctx: Context, options: ListOptions = {}) {
    super(ctx)

    this.options = options

    this.initBaseEl()

    // init layout manager
    this.layout = new ListLayout({
      $commentsWrap: this.$commentsWrap!,
      nestSortBy: this.getNestSortBy(),
      nestMax: this.ctx.conf.nestMax,
      flatMode: this.getFlatMode(),
      createComment: (d, c) => this.createComment(d, c),
      findComment: (id) => this.ctx.findComment(id),
      getCommentDataList: () => this.ctx.getCommentDataList()
    })

    // 事件监听
    this.ctx.on('conf-loaded', () => {
    })

    this.ctx.on('list-loaded', () => {
      // 防止评论框被吞
      this.ctx.editorResetState()
    })
  }

  getOptions() {
    return this.options
  }

  /** 嵌套模式下的排序方式 */
  getNestSortBy(): ListNest.SortByType {
    return this.options.nestSortBy || this.ctx.conf.nestSort || 'DATE_ASC'
  }

  /** 平铺模式 */
  getFlatMode(): boolean {
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

  /** 分页方式 */
  getPageMode(): TPgMode {
    return this.options.pageMode || (this.conf.pagination.readMore ? 'read-more' : 'pagination')
  }

  /** 每页数量 (每次请求获取量) */
  getPageSize(): number {
    return this.options.pageSize || this.conf.pagination.pageSize
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
    return this.$commentsWrap!
  }

  private initBaseEl() {
    this.$el = Utils.createElement(
      `<div class="atk-list-lite">
        <div class="atk-list-comments-wrap"></div>
      </div>`)
    this.$commentsWrap = this.$el.querySelector('.atk-list-comments-wrap')!

    if (!this.options.liteMode) {
      const el = Utils.createElement(ListHTML)
      el.querySelector('.atk-list-body')!.append(this.$el) // 把 list 的 $el 变为子元素
      this.$el = el

      this.options.repositionAt = this.$el // 更新翻页归位元素
    }
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
      this.clearAllComments()
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
    this.loadPagination(offset, (this.getFlatMode() ? data.total : data.total_roots))

    // 更新页面数据
    this.ctx.updatePage(data.page)

    // 未读消息提示功能
    this.ctx.updateUnreadList(data.unread || [])

    // 事件触发，列表已加载
    this.ctx.trigger('list-loaded')
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
  private loadPagination(offset: number, total: number) {
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
    this.ctx.trigger('list-error', {
      msg, data: errData
    })
  }

  /** 重新加载列表 */
  public reload() {
    this.fetchComments(0)
  }

  /** 创建新评论 */
  private createComment(comment: CommentData, ctxComments: CommentData[]): Comment {
    const instance = new Comment(this.ctx, comment, {
      isFlatMode: this.getFlatMode(),
      afterRender: () => {
        const renderCommentFn = this.options.renderComment
        renderCommentFn && renderCommentFn(instance)
      },
      onDelete: (c: Comment) => {
        this.deleteComment(c.getID())
      },
      replyTo: (comment.rid ? ctxComments.find(c => c.id === comment.rid) : undefined)
    })

    // 渲染元素
    instance.render()

    // 放入 comment 总表中
    this.ctx.getCommentList().push(instance)

    return instance
  }

  /** 导入评论 */
  public importComments(srcData: CommentData[]) {
    this.layout.import(srcData)
  }

  /** 新增评论 · 首部添加 */
  public insertComment(commentData: CommentData) {
    this.layout.insert(commentData)

    // 评论数增加 1
    if (this.data) this.data.total += 1

    // 评论新增后
    this.ctx.trigger('list-loaded')
    this.ctx.trigger('list-inserted', commentData)
  }

  /** 更新评论 */
  public updateComment(commentData: CommentData) {
    const comment = this.ctx.findComment(commentData.id)
    comment && comment.setData(commentData)

    this.ctx.trigger('list-loaded')
  }

  /** 删除评论 */
  public deleteComment(id: number) {
    const comment = this.ctx.findComment(id)
    if (!comment) throw Error(`Comment ${id} cannot be found`)

    comment.getEl().remove()

    const list = this.ctx.getCommentList()
    list.splice(list.indexOf(comment), 1)

    // 评论数减 1
    if (this.data) this.data.total -= 1

    // 评论删除后
    this.ctx.trigger('list-loaded')
    this.ctx.trigger('list-deleted', comment.getData())
  }

  /** 删除全部评论 */
  public clearAllComments() {
    this.getCommentsWrapEl().innerHTML = ''
    this.clearData()

    this.ctx.clearCommentList()
    this.ctx.trigger('list-loaded')
  }
}

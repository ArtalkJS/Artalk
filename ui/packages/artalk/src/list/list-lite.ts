import { ListData, CommentData, NotifyData } from '~/types/artalk-data'
import Context from '~/types/context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Comment from '../comment'
import PgHolder, { TPgMode } from './paginator'
import * as ListNest from './list-nest'
import * as ListUi from './list-ui'
import { version as ARTALK_VERSION } from '../../package.json'
import { handleBackendRefConf } from '../config'

export default class ListLite extends Component {
  /** 列表评论集区域元素 */
  protected $commentsWrap: HTMLElement

  /** 最后一次请求得到的列表数据 */
  protected data?: ListData

  /** 是否处于加载中状态 */
  protected isLoading: boolean = false

  /** 配置是否已加载 */
  private confLoaded = false

  /** 无评论显示文字 */
  public noCommentText?: string

  /** 嵌套模式下的排序方式 */
  private _nestSortBy?: ListNest.SortByType
  public get nestSortBy() {
    return this._nestSortBy || this.ctx.conf.nestSort || 'DATE_ASC'
  }
  public set nestSortBy(val: ListNest.SortByType) {
    this._nestSortBy = val
  }

  /** 平铺模式 */
  private _flatMode?:boolean
  public get flatMode() {
    if (this._flatMode !== undefined)
      return this._flatMode

    // 配置开启平铺模式
    if (this.ctx.conf.flatMode === true || Number(this.ctx.conf.nestMax) <= 1)
      return true

    // 自动判断启用平铺模式
    if (this.ctx.conf.flatMode === 'auto' && window.matchMedia("(max-width: 768px)").matches)
      return true

    return false
  }
  public set flatMode(val: boolean) {
    this._flatMode = val
  }

  /** 分页方式 */
  public _pageMode?: TPgMode
  public get pageMode() {
    return this._pageMode || (this.conf.pagination.readMore ? 'read-more' : 'pagination')
  }
  public set pageMode(val: 'pagination'|'read-more') {
    this._pageMode = val
    this.pgHolder?.setMode(this._pageMode)
  }

  /** 分页方式持有者 */
  public pgHolder?: PgHolder

  /** 每页数量 (每次请求获取量) */
  private _pageSize?: number
  public get pageSize() {
    return this._pageSize || this.conf.pagination.pageSize
  }
  public set pageSize(val: number) {
    this._pageSize = val
  }

  /** 监听指定元素上的滚动 */
  public scrollListenerAt?: HTMLElement
  /** 翻页归位到指定元素 */
  public repositionAt?: HTMLElement

  // 一些 Hook 函数
  public renderComment?: (comment: Comment) => void
  public paramsEditor?: (params: any) => void
  public onAfterLoad?: (data: ListData) => void

  /** 未读标记 */
  protected unread: NotifyData[] = []
  public unreadHighlight?: boolean // 高亮未读

  constructor (ctx: Context) {
    super(ctx)

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
    if ((this.pageMode === 'read-more' && isFirstLoad) || this.pageMode === 'pagination') {
      this.ctx.clearAllComments()
    }

    // 请求评论数据
    let listData: ListData
    try {
      // 执行请求
      listData = await this.ctx.getApi().comment
        .get(offset, this.pageSize, this.flatMode, this.paramsEditor)
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
    if (!this.confLoaded) {
      // 仅应用一次配置
      const backendRefConf = handleBackendRefConf(data.conf.frontend_conf)
      if (this.conf.useBackendConf) this.ctx.updateConf(backendRefConf)
      else this.ctx.updateConf({}) // 让事件监听 `on('conf-loaded')` 有效，与前者保持相同生命周期环节
      this.confLoaded = true
    }

    // 前后端版本一致性检测
    if (this.ctx.conf.versionCheck && ListUi.versionCheckDialog(this, ARTALK_VERSION, data.api_version.version)) return

    // 导入数据
    this.importComments(data.comments)

    // 分页功能
    this.refreshPagination(offset, (this.flatMode ? data.total : data.total_roots))

    // 加载后事件
    this.refreshUI()

    // 未读消息提示功能
    this.ctx.updateNotifies(data.unread || [])

    // 事件触发，列表已加载
    this.ctx.trigger('list-loaded')

    // Hook 函数调用
    if (this.onAfterLoad) this.onAfterLoad(data)
  }

  /** 分页模式 */
  private refreshPagination(offset: number, total: number) {
    // 初始化
    if (!this.pgHolder) {
      this.pgHolder = new PgHolder({
        list: this,
        mode: this.pageMode,
        pageSize: this.pageSize,
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
    if (offset !== 0 && this.pageMode === 'read-more') {
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
      isFlatMode: this.flatMode,
      afterRender: () => {
        if (this.renderComment) this.renderComment(comment)
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

  /** 导入评论 · 通过请求数据 */
  public importComments(srcData: CommentData[]) {
    if (this.flatMode) {
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
    const rootNodes = ListNest.makeNestCommentNodeList(srcData, this.nestSortBy, this.conf.nestMax)
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
    if (!this.flatMode) {
      // 嵌套模式
      const comment = this.createComment(commentData)

      if (commentData.rid === 0) {
        // root评论 新增
        this.$commentsWrap.prepend(comment.getEl())
      } else {
        // 子评论 新增
        const parent = this.ctx.findComment(commentData.rid)
        if (parent) {
          parent.putChild(comment, (this.nestSortBy === 'DATE_ASC' ? 'append' : 'prepend'))

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
    if (comment) {
      comment.setData(commentData)
    }
  }

  /** 更新未读数据 */
  public updateUnread(notifies: NotifyData[]) {
    this.unread = notifies

    // 高亮评论
    if (this.unreadHighlight === true) {
      this.ctx.getCommentList().forEach((comment) => {
        const notify = this.unread.find(o => o.comment_id === comment.getID())
        if (notify) {
          comment.getRender().setUnread(true)
          comment.getRender().setOpenAction(() => {
            window.open(notify.read_link)
            this.unread = this.unread.filter(o => o.comment_id !== comment.getID()) // remove
            this.ctx.updateNotifies(this.unread)
          })
        } else {
          comment.getRender().setUnread(false)
        }
      })
    }
  }
}

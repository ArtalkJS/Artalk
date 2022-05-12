import { ListData, CommentData, NotifyData, ApiVersionData } from '~/types/artalk-data'
import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Api from '../api'
import Comment from './comment'
import Pagination from './pagination'
import ReadMoreBtn from './read-more-btn'
import * as ListNest from './list-nest'
import { backendMinVersion } from '../../package.json'

export default class ListLite extends Component {
  protected $commentsWrap: HTMLElement

  protected commentList: Comment[] = [] // Note: 无层级结构 + 无须排列
  protected get commentDataList() { return this.commentList.map(c => c.data) }

  protected data?: ListData
  protected isLoading: boolean = false

  public noCommentText: string // 无评论时显示

  /** 平铺模式 */
  public flatMode = false

  /** 嵌套模式 */
  public nestSortBy!: ListNest.SortByType

  /** 分页 */
  public pageMode: 'pagination'|'read-more' = 'pagination'
  public pageSize: number = 20 // 每次请求获取量
  public scrollListenerAt?: HTMLElement // 监听指定元素上的滚动
  public repositionAt?: HTMLElement // 翻页归位到指定元素

  protected pagination?: Pagination   // 分页条
  protected readMoreBtn?: ReadMoreBtn // 阅读更多
  protected autoLoadScrollEvent?: any // 自动加载滚动事件

  public renderComment?: (comment: Comment) => void
  public paramsEditor?: (params: any) => void
  public onAfterLoad?: (data: ListData) => void

  // 未读标记
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

    // 评论为空时显示字符
    this.noCommentText = ctx.conf.noComment || ctx.$t('noComment')

    // 嵌套排序方式
    this.nestSortBy = this.ctx.conf.nestSort || 'DATE_ASC'

    // 评论时间自动更新
    window.setInterval(() => {
      this.$el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach(el => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)), this.ctx)
      })
    }, 30 * 1000) // 30s 更新一次

    // 事件监听
    this.ctx.on('unread-update', (data) => (this.updateUnread(data.notifies)))
  }

  /** 评论获取 */
  public async fetchComments(offset: number) {
    if (this.isLoading) return

    // 加载动画
    const showLoading = () => {
      this.isLoading = true
      if (offset === 0) Ui.showLoading(this.$el)
      else if (this.pageMode === 'read-more') this.readMoreBtn!.setLoading(true)
      else if (this.pageMode === 'pagination') this.pagination!.setLoading(true)
    }
    const hideLoading = () => {
      this.isLoading = false
      if (offset === 0) Ui.hideLoading(this.$el)
      else if (this.pageMode === 'read-more') this.readMoreBtn!.setLoading(false)
      else if (this.pageMode === 'pagination') this.pagination!.setLoading(false)
    }
    showLoading()

    // 事件通知（开始加载评论）
    this.ctx.trigger('comments-load')

    // 清空评论（加载按钮）
    if (this.pageMode === 'read-more' && offset === 0) { this.clearAllComments() }

    // 请求评论数据
    let listData: ListData
    try {
      listData = await new Api(this.ctx).get(offset, this.pageSize, this.flatMode, this.paramsEditor)
    } catch (e: any) {
      this.onError(e.msg || String(e), offset, e.data)
      throw e
    } finally {
      hideLoading()
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
      hideLoading()
    }
  }

  protected onLoad(data: ListData, offset: number) {
    // 清空评论（分页条）
    if (this.pageMode === 'pagination') { this.clearAllComments() }

    this.data = data

    // 版本检测
    const feMinVersion = data.api_version?.fe_min_version || '0.0.0'
    if (this.ctx.conf.versionCheck && this.versionCheck('frontend', feMinVersion, ARTALK_VERSION)) return
    if (this.ctx.conf.versionCheck && this.versionCheck('backend', backendMinVersion, data.api_version?.version)) return

    // 图片上传功能
    if (data.conf && typeof data.conf.img_upload === "boolean") {
      this.ctx.conf.imgUpload = data.conf.img_upload
    }

    // 导入数据
    this.importComments(data.comments)

    // 分页方式
    this.refreshPagination(offset, (this.flatMode ? data.total : data.total_roots)) // 初始化

    // 加载后事件
    this.refreshUI()

    this.ctx.trigger('unread-update', { notifies: data.unread || [] })
    this.ctx.trigger('comments-loaded')
    this.ctx.trigger('conf-updated')

    if (this.onAfterLoad) this.onAfterLoad(data)
  }

  /** 分页模式 */
  private refreshPagination(offset: number, total: number) {
    const modePagination = (this.pageMode === 'pagination')
    const modeReadMoreBtn = (this.pageMode === 'read-more')
    const initialized = (modePagination) ? !!this.pagination : !!this.readMoreBtn

    // 初始化
    if (!initialized) {
      this.initPagination()
    }

    // 更新
    if (modePagination) this.pagination!.update(offset, total)
    if (modeReadMoreBtn) this.readMoreBtn!.update(offset, total)
  }

  private initPagination() {
    // 加载更多按钮
    if (this.pageMode === 'read-more') {
      this.readMoreBtn = new ReadMoreBtn({
        pageSize: this.pageSize,
        onClick: async (o) => {
          await this.fetchComments(o)
        },
        text: this.ctx.$t('loadMore'),
      })
      this.$el.append(this.readMoreBtn.$el)

      // 滚动到底部自动加载
      if (this.conf.pagination?.autoLoad) {
        // 添加滚动事件监听
        const at = this.scrollListenerAt || document
        if (this.autoLoadScrollEvent) at.removeEventListener('scroll', this.autoLoadScrollEvent) // 解除原有
        this.autoLoadScrollEvent = () => {
          if (this.pageMode !== 'read-more') return
          if (!this.readMoreBtn) return

          if (!this.readMoreBtn.hasMore) return
          if (this.isLoading) return

          const $target = this.$el.querySelector<HTMLElement>('.atk-list-comments-wrap > .atk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
          if (!$target) return
          if (Ui.isVisible($target, this.scrollListenerAt)) {
            this.readMoreBtn.click() // 自动点击加载更多按钮
          }
        }
        at.addEventListener('scroll', this.autoLoadScrollEvent)
      }
    }

    // 分页条
    if (this.pageMode === 'pagination') {
      this.pagination = new Pagination((!this.flatMode ? this.data!.total_roots : this.data!.total), {
        pageSize: this.pageSize,
        onChange: async (o) => {
          if (this.ctx.conf.editorTravel === true) {
            this.ctx.trigger('editor-travel-back') // 防止评论框被吞
          }

          await this.fetchComments(o)

          // 滚动到第一个评论的位置
          if (this.repositionAt) {
            const at = this.scrollListenerAt || window
            at.scroll({
              top: this.repositionAt ? Utils.getOffset(this.repositionAt).top : 0,
              left: 0,
            })
          }
        }
      })

      this.$el.append(this.pagination.$el)
    }
  }

  /** 错误处理 */
  protected onError(msg: any, offset: number, errData?: any) {
    msg = String(msg)
    console.error(msg)

    // 加载更多按钮显示错误
    if (offset !== 0 && this.pageMode === 'read-more') {
      this.readMoreBtn?.showErr(this.$t('loadFail'))
      return
    }

    // 显示错误对话框
    const $err = Utils.createElement(`<span>${msg}，${this.$t('listLoadFailMsg')}<br/></span>`)

    const $retryBtn = Utils.createElement(`<span style="cursor:pointer;">${this.$t('listRetry')}</span>`)
    $retryBtn.onclick = () => (this.fetchComments(0))
    $err.appendChild($retryBtn)

    const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
    $err.appendChild(adminBtn)
    if (!this.ctx.user.data.isAdmin) adminBtn.classList.add('atk-hide')

    let sidebarView = ''

    // 找不到站点错误，打开侧边栏并填入创建站点表单
    if (errData?.err_no_site) {
      const viewLoadParam = {
        create_name: this.ctx.conf.site,
        create_urls: `${window.location.protocol}//${window.location.host}`
      }
      // TODO 真的是飞鸽传书啊
      sidebarView = `sites|${JSON.stringify(viewLoadParam)}`
    }

    adminBtn.onclick = () => (this.ctx.trigger('sidebar-show', {
      view: sidebarView
    }))

    Ui.setError(this.$el, $err)
  }

  /** 刷新界面 */
  public refreshUI() {
    // 无评论
    const isNoComment = this.commentList.length <= 0
    let $noComment = this.$commentsWrap.querySelector<HTMLElement>('.atk-list-no-comment')

    if (isNoComment) {
      if (!$noComment) {
        $noComment = Utils.createElement('<div class="atk-list-no-comment"></div>')
        $noComment.innerHTML = this.noCommentText
        this.$commentsWrap.appendChild($noComment)
      }
    } else {
      $noComment?.remove()
    }

    // 仅管理员显示控制
    this.ctx.trigger('check-admin-show-el')
  }

  /** 创建新评论 */
  protected createComment(cData: CommentData, ctxData?: CommentData[]) {
    if (!ctxData) ctxData = this.commentDataList

    const comment = new Comment(this.ctx, cData)
    comment.flatMode = this.flatMode
    comment.afterRender = () => {
      if (this.renderComment) this.renderComment(comment)
    }
    comment.onDelete = (c) => {
      this.deleteComment(c)
      this.refreshUI()
    }

    // 子评论查找回复对象
    comment.replyTo = (cData.rid ? ctxData.find(c => c.id === cData.rid) : undefined)

    // 渲染元素
    comment.render()

    // 放入 comment 总表中
    this.commentList.push(comment)

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
      rootC.playFadeInAnim()

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
      rootC.checkHeightLimit()
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
      comment.playFadeInAnim()
    }

    // 平铺评论插入后 · 内容限高检测
    comment.checkHeightLimit()

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
        const parent = this.findComment(commentData.rid)
        if (parent) {
          parent.putChild(comment, (this.nestSortBy === 'DATE_ASC' ? 'append' : 'prepend'))

          // 若父评论存在 “子评论部分” 限高，取消限高
          comment.getParents().forEach((p) => {
            if (p.$children) p.heightLimitRemove(p.$children)
          })
        }
      }

      comment.checkHeightLimit()

      Ui.scrollIntoView(comment.getEl()) // 滚动到可以见
      comment.playFadeInAnim() // 播放评论渐出动画
    } else {
      // 平铺模式
      const comment = this.putCommentFlatMode(commentData, this.commentDataList, 'prepend')
      Ui.scrollIntoView(comment.getEl()) // 滚动到可见
    }

    // 评论数增加 1
    if (this.data) this.data.total += 1

    // 评论新增后
    this.refreshUI()
    this.ctx.trigger('comments-loaded')
  }

  /** 查找评论项 */
  protected findComment(id: number): Comment|undefined {
    return this.commentList.find(c => c.data.id === id)
  }

  /** 删除评论 */
  public deleteComment(_comment: number|Comment) {
    let comment: Comment
    if (typeof _comment === 'number') {
      const findComment = this.findComment(_comment)
      if (!findComment) throw Error(`Comment ${_comment} cannot be found`)
      comment = findComment
    } else comment = _comment

    comment.getEl().remove()
    this.commentList.splice(this.commentList.indexOf(comment), 1)

    if (this.data) this.data.total -= 1 // 评论数减 1

    this.refreshUI()
  }

  /** 清空所有评论 */
  public clearAllComments() {
    this.$commentsWrap.innerHTML = ''
    this.data = undefined
    this.commentList = []
  }

  /** 更新未读数据 */
  public updateUnread(notifies: NotifyData[]) {
    this.unread = notifies

    // 高亮评论
    if (this.unreadHighlight === true) {
      this.commentList.forEach((comment) => {
        const notify = this.unread.find(o => o.comment_id === comment.data.id)
        if (notify) {
          comment.setUnread(true)
          comment.setOpenURL(notify.read_link)
          comment.openEvt = () => {
            this.unread = this.unread.filter(o => o.comment_id !== comment.data.id) // remove
            this.ctx.trigger('unread-update', {
              notifies: this.unread
            })
          }
        } else {
          comment.setUnread(false)
        }
      })
    }
  }

  /** 版本检测 */
  public versionCheck(name: 'frontend'|'backend', needVersion: string, curtVersion: string): boolean {
    const needUpdate = Utils.versionCompare(needVersion, curtVersion) === 1
    if (needUpdate) {
      // 需要更新
      const errEl = Utils.createElement(`<div>Artalk ${this.$t(name)}版本已过时，请更新以获得完整体验<br/>`
      + `如果你是管理员，请前往 “<a href="https://artalk.js.org/" target="_blank">官方文档</a>” 获得帮助`
      + `<br/><br/>`
      + `<span style="color: var(--at-color-meta);">当前${this.$t(name)}版本 ${curtVersion}，需求版本 >= ${needVersion}</span><br/><br/>`
      + `</div>`)
      const ignoreBtn = Utils.createElement('<span style="cursor:pointer;">忽略</span>')
      ignoreBtn.onclick = () => {
        Ui.setError(this.$el.parentElement!, null)
        this.ctx.conf.versionCheck = false
        this.ctx.trigger('conf-updated')
        this.fetchComments(0)
      }
      errEl.append(ignoreBtn)
      Ui.setError(this.$el.parentElement!, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
    }

    return needUpdate
  }
}

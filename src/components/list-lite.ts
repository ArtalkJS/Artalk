import Context from '../context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Api from '../api'
import Comment from './comment'
import { ListData, CommentData, NotifyData, ApiVersionData } from '~/types/artalk-data'
import Pagination from './pagination'
import ReadMoreBtn from './read-more-btn'

export default class ListLite extends Component {
  private $parent: HTMLElement

  private $commentsWrap: HTMLElement
  public comments: Comment[] = []

  public data?: ListData
  private pageSize: number = 15 // 每次请求获取量

  public noCommentText: string
  public renderComment?: (comment: Comment) => void
  public paramsEditor?: (params: any) => void
  public onAfterLoad?: (data: ListData) => void

  // 状态标识
  private isLoading: boolean = false
  public isFirstLoad = true

  // 平铺模式
  public flatMode?: boolean

  // 分页方式
  public pageMode: 'pagination'|'read-more' = 'pagination'
  public pagination?: Pagination   // 分页条
  public readMoreBtn?: ReadMoreBtn // 阅读更多
  public autoLoadScrollEvent?: any // 自动加载滚动事件
  public autoLoadListenerAt?: HTMLElement // 监听指定元素上的滚动

  // 未读标记
  public unread: NotifyData[] = []
  public unreadHighlight = false

  constructor (ctx: Context, $parent: HTMLElement) {
    super(ctx)

    // 初始化元素
    this.$parent = $parent
    this.$el = Utils.createElement(
    `<div class="atk-list-lite">
      <div class="atk-list-comments-wrap"></div>
    </div>`)
    this.$commentsWrap = this.$el.querySelector('.atk-list-comments-wrap')!

    // 初始化配置
    this.pageSize = this.conf.pagination?.pageSize || this.pageSize
    this.noCommentText = this.conf.noComment || '无评论'

    // 评论时间自动更新
    window.setInterval(() => {
      this.$el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach(el => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)))
      })
    }, 30 * 1000) // 30s 更新一次

    // 事件监听
    this.ctx.on('unread-update', (data) => (this.updateUnread(data.notifies)))
  }

  /** 评论获取 */
  public async fetchComments(offset: number = 0) {
    // 清空评论（加载按钮）
    if (this.pageMode === 'read-more' && offset === 0) { this.clearAllComments() }

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

    // 请求评论数据
    let listData: ListData
    try {
      listData = await new Api(this.ctx).get(offset, this.flatMode, this.paramsEditor)
    } catch (e: any) {
      this.onError(e.msg || String(e))
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
      this.onError(String(e))
      throw e
    } finally {
      hideLoading()
    }

    this.isFirstLoad = false
  }

  public onLoad(data: ListData, offset: number) {
    // 清空评论（分页条）
    if (this.pageMode === 'pagination') { this.clearAllComments() }

    this.data = data

    // 版本检测
    if (this.ctx.conf.versionCheck && this.versionCheck(data.api_version)) return

    // 导入数据
    this.importComments(data.comments)

    // 分页方式
    if (this.isFirstLoad) this.initPageMode() // 初始化
    if (this.pageMode === 'pagination') this.pagination!.update(offset, (!this.flatMode ? data.total_roots : data.total) || 0)
    if (this.pageMode === 'read-more') this.readMoreBtn!.update(offset, (!this.flatMode ? data.total_roots : data.total) || 0)

    // 加载后事件
    this.refreshUI()

    this.ctx.trigger('unread-update', { notifies: data.unread || [] })
    this.ctx.trigger('comments-loaded')

    if (this.onAfterLoad) this.onAfterLoad(data)
  }

  /** 初始化分页模式 */
  initPageMode() {
    // 解除滚动事件监听
    if (this.autoLoadScrollEvent) {
      const at = this.autoLoadListenerAt || document
      at.removeEventListener('scroll', this.autoLoadScrollEvent)
    }

    // 加载更多按钮
    if (this.pageMode === 'read-more') {
      const readMoreBtn = new ReadMoreBtn({
        pageSize: this.pageSize,
        onClick: async (offset) => {
          await this.fetchComments(offset)
        },
      })
      if (this.readMoreBtn) this.readMoreBtn.$el.replaceWith(readMoreBtn.$el)
      else this.$el.append(readMoreBtn.$el)
      this.readMoreBtn = readMoreBtn

      // 滚动到底部自动加载
      if (this.conf.pagination?.autoLoad) {
        this.autoLoadScrollEvent = () => {
          if (this.pageMode !== 'read-more') return
          if (!this.readMoreBtn) return

          if (!this.readMoreBtn.hasMore) return
          if (this.isLoading) return

          const $target = this.$el.querySelector<HTMLElement>('.atk-list-comments-wrap > .atk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
          if (!$target) return
          if (Ui.isVisible($target, this.autoLoadListenerAt)) {
            this.readMoreBtn.click() // 自动点击加载更多按钮
          }
        }
        const at = this.autoLoadListenerAt || document
        at.addEventListener('scroll', this.autoLoadScrollEvent)
      }
    }

    // 分页条
    if (this.pageMode === 'pagination') {
      const pagination = new Pagination((!this.flatMode ? this.data!.total_roots : this.data!.total), {
        pageSize: this.pageSize,
        onChange: async (offset) => {
          if (this.ctx.conf.editorTravel === true) {
            this.ctx.trigger('editor-travel-back') // 防止评论框被吞
          }

          await this.fetchComments(offset)

          // 滚动到第一个评论的位置
          if (this.$parent) {
            const at = this.autoLoadListenerAt || window
            at.scroll({
              top: (at === window && this.$parent) ? Utils.getOffset(this.$parent).top : 0,
              left: 0,
            })
          }
        }
      })

      if (this.pagination) this.pagination.$el.replaceWith(pagination.$el)
      else this.$el.append(pagination.$el)
      this.pagination = pagination
    }
  }

  /** 错误处理 */
  public onError(msg: any) {
    msg = String(msg)
    console.error(msg)

    // 加载更多按钮显示错误
    if (!this.isFirstLoad && this.pageMode === 'read-more') {
      this.readMoreBtn?.showErr(`获取失败`)
      return
    }

    // 显示错误对话框
    const $err = Utils.createElement(`<span>${msg}，无法获取评论列表数据<br/></span>`)

    const $retryBtn = Utils.createElement('<span style="cursor:pointer;">点击重新获取</span>')
    $retryBtn.onclick = () => (this.fetchComments())
    $err.appendChild($retryBtn)

    const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
    adminBtn.onclick = () => (this.ctx.trigger('sidebar-show'))
    if (!this.ctx.user.data.isAdmin) adminBtn.classList.add('atk-hide')
    $err.appendChild(adminBtn)

    Ui.setError(this.$el, $err)
  }

  /** 刷新界面 */
  public refreshUI() {
    // 无评论
    const isNoComment = this.comments.length <= 0
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
  public createComment(data: CommentData) {
    const comment = new Comment(this.ctx, data)
    comment.afterRender = () => {
      if (this.renderComment) this.renderComment(comment)
    }
    comment.onDelete = (c) => {
      this.deleteComment(c)
      this.refreshUI()
    }

    return comment
  }

  /** 导入评论 · 通过请求数据 */
  public importComments(srcData: CommentData[]) {
    if (this.flatMode) {
      srcData.forEach((commentData: CommentData) => {
        this.putCommentFlatMode(commentData, srcData, 'append')
      })
    } else {
      this.importCommentsNesting(srcData)
    }
  }

  // 导入评论 · 嵌套模式
  private importCommentsNesting(srcData: CommentData[]) {
    // 查找并导入所有子评论
    const loadChildren = (parentC: Comment) => {
      const children = srcData.filter(o => o.rid === parentC.data.id)
      children.forEach((childData: CommentData) => {
        const childC = this.createComment(childData)
        childC.render()

        // 插入到父评论中
        parentC.putChild(childC)

        // 递归加载子评论
        loadChildren(childC)
      })
    }

    // 遍历 root 评论
    const rootComments = srcData.filter((o) => o.rid === 0)
    rootComments.forEach((rootData: CommentData) => {
      const rootC = this.createComment(rootData)
      rootC.render()

      // 显示并播放渐入动画
      this.$commentsWrap.appendChild(rootC.getEl())
      rootC.playFadeInAnim()

      // 评论放入 comments 总表中
      this.comments.push(rootC)

      // 加载子评论
      loadChildren(rootC)

      // 限高检测
      rootC.checkHeightLimit()
    })
  }

  /** 导入评论 · 平铺模式 */
  private putCommentFlatMode(cData: CommentData, srcData: CommentData[], insertMode: 'append'|'prepend') {
    if (cData.is_collapsed) cData.is_allow_reply = false
    const comment = this.createComment(cData)
    if (cData.rid !== 0) {
      const rComment = srcData.find(o => o.id === cData.rid)
      if (rComment) comment.replyTo = rComment
    }
    comment.render()

    // 将评论导入 comments 总表中
    if (insertMode === 'append') this.comments.push(comment)
    if (insertMode === 'prepend') this.comments.unshift(comment)

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
      comment.render()

      if (commentData.rid === 0) {
        // root评论 新增
        this.$commentsWrap.prepend(comment.getEl())
        this.comments.unshift(comment)
      } else {
        // 子评论 新增
        const parent = this.findComment(commentData.rid)
        if (parent) {
          parent.putChild(comment)

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
      const comment = this.putCommentFlatMode(commentData, this.comments.map(c => c.data), 'prepend')
      Ui.scrollIntoView(comment.getEl()) // 滚动到可见
    }

    // 评论数增加 1
    if (this.data) this.data.total += 1

    // 评论新增后
    this.refreshUI()
    this.ctx.trigger('comments-loaded')
  }

  /** 遍历操作 Comment (包括子评论) */
  public eachComment(commentList: Comment[], action: (comment: Comment, levelList: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  /** 查找评论项 */
  public findComment(id: number, src?: Comment[]): Comment|null {
    if (!src) src = this.comments
    let comment: Comment|null = null
    this.eachComment(src, (item) => {
      if (item.data.id === id) {
        comment = item
        return false
      }
      return true
    })

    return comment
  }

  /** 删除评论 */
  public deleteComment(comment: number|Comment) {
    let findComment: Comment|null
    if (typeof comment === 'number') {
      findComment = this.findComment(comment)
      if (!findComment) throw Error(`未找到评论 ${comment}`)
    } else findComment = comment

    findComment.getEl().remove()
    this.eachComment(this.comments, (item, levelList) => {
      if (item === findComment) {
        levelList.splice(levelList.indexOf(item), 1)
        return false
      }
      return true
    })

    if (this.data) this.data.total -= 1 // 评论数减 1

    this.refreshUI()
  }

  /** 清空所有评论 */
  public clearAllComments() {
    this.$commentsWrap.innerHTML = ''
    this.data = undefined
    this.comments = []
  }

  /** 更新未读数据 */
  public updateUnread(notifies: NotifyData[]) {
    this.unread = notifies

    // 高亮评论
    if (this.unreadHighlight) {
      this.eachComment(this.comments, (comment) => {
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

  /** 前端版本检测 */
  public versionCheck(versionData: ApiVersionData): boolean {
    const needVersion = versionData?.fe_min_version || '0.0.0'
    const needUpdate = Utils.versionCompare(needVersion, ARTALK_VERSION) === 1
    if (needUpdate) {
      // 需要更新
      const errEl = Utils.createElement(`<div>前端 Artalk 版本已过时，请更新以获得完整体验<br/>`
      + `若您是站点管理员，请前往 “<a href="https://artalk.js.org/" target="_blank">官方文档</a>” 获取帮助`
      + `<br/><br/>`
      + `<span style="color: var(--at-color-meta);">前端版本 ${ARTALK_VERSION}，需求版本 >= ${needVersion}</span><br/><br/>`
      + `</div>`)
      const ignoreBtn = Utils.createElement('<span style="cursor:pointer;">忽略</span>')
      ignoreBtn.onclick = () => {
        Ui.setError(this.ctx, null)
        this.ctx.conf.versionCheck = false
        this.fetchComments(0)
      }
      errEl.append(ignoreBtn)
      Ui.setError(this.ctx, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
    }

    return needUpdate
  }
}

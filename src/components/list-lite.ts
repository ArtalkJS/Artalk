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
  private offset: number = 0
  public type?: 'all'|'mentions'|'mine'|'pending'

  public noCommentText: string
  public renderComment?: (comment: Comment) => void
  public paramsEditor?: (params: any) => void
  public onAfterLoad?: (data: ListData) => void

  private isLoading: boolean = false
  public isFirstLoad = true

  public flatMode?: boolean

  public pageMode: 'pagination'|'read-more' = 'pagination'
  public pagination?: Pagination
  public readMoreBtn?: ReadMoreBtn
  public autoLoadScrollEvent?: any
  public autoLoadListenerAt?: HTMLElement

  public unread: NotifyData[] = []
  public unreadHighlight = false

  constructor (ctx: Context, $parent: HTMLElement) {
    super(ctx)

    this.$parent = $parent
    this.$el = Utils.createElement(
    `<div class="atk-list-lite">
      <div class="atk-list-comments-wrap"></div>
    </div>`)
    this.$commentsWrap = this.$el.querySelector('.atk-list-comments-wrap')!

    // 查看更多
    this.pageSize = this.conf.pagination ? (this.conf.pagination.pageSize || this.pageSize) : this.pageSize

    this.noCommentText = this.conf.noComment || '无评论'

    // 评论时间自动更新
    window.setInterval(() => {
      this.$el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach(el => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)))
      })
    }, 30 * 1000) // 30s 更新一次

    this.ctx.on('unread-update', (data) => (this.updateUnread(data.notifies)))
  }

  public async reqComments(offset: number = 0) {
    if (offset === 0 && this.pageMode !== 'pagination') { this.clearAllComments() }

    // set loading
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
    this.ctx.trigger('comments-load')

    let listData: ListData
    try {
      listData = await new Api(this.ctx).get(offset, this.type, this.flatMode, this.paramsEditor)
    } catch (e: any) {
      this.onError(e.msg || String(e))
      throw e
    } finally {
      hideLoading()
    }

    // version check
    if (this.ctx.conf.versionCheck) {
      const needUpdate = this.apiVersionCheck(listData.api_version || {})
      if (needUpdate) return
    }

    // load data
    this.offset = offset

    try {
      this.onLoad(listData, offset)

      if (this.onAfterLoad) {
        this.onAfterLoad(listData)
      }
    } catch (e: any) {
      this.onError(String(e))
      throw e
    } finally {
      hideLoading()
    }
  }

  public onLoad(data: ListData, offset: number) {
    Ui.setError(this.$el, null)

    if (this.pageMode === 'pagination') {
      this.clearAllComments()
    }

    this.data = data
    this.importComments(data.comments)

    // onLoad 时初始化
    if (this.isFirstLoad) {
      this.onLoadInit()
    }

    if (this.pageMode === 'pagination') {
      this.pagination!.update(offset, this.data?.total_parents || 0)
    }
    if (this.pageMode === 'read-more') {
      if (this.hasMoreComments) this.readMoreBtn!.show()
      else this.readMoreBtn!.hide()
    }

    this.ctx.trigger('unread-update', { notifies: data.unread || [] })

    this.isFirstLoad = false
  }

  /** 初始化加载更多 */
  onLoadInit() {
    if (this.autoLoadScrollEvent) {
      const at = this.autoLoadListenerAt || document
      at.removeEventListener('scroll', this.autoLoadScrollEvent)
    }

    if (this.pageMode === 'read-more') {
      // 阅读更多按钮
      const readMoreBtn = new ReadMoreBtn({
        pageSize: this.pageSize,
        total: 0,
        onClick: async () => {
          const offset = this.offset + this.pageSize
          await this.reqComments(offset)
        },
      })
      if (this.readMoreBtn) this.readMoreBtn.$el.replaceWith(readMoreBtn.$el)
      else this.$el.append(readMoreBtn.$el)
      this.readMoreBtn = readMoreBtn

      // 滚动到底部自动加载
      if (this.conf.pagination?.autoLoad) {
        this.autoLoadScrollEvent = () => {
          if (this.pageMode !== 'read-more') return
          if (!this.hasMoreComments) return
          if (this.isLoading) return

          const $target = this.$el.querySelector<HTMLElement>('.atk-list-comments-wrap > .atk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
          if (!$target) return
          if (Ui.isVisible($target, this.autoLoadListenerAt)) {
            // 加载更多
            this.readMoreBtn!.click()
          }
        }
        const at = this.autoLoadListenerAt || document
        at.addEventListener('scroll', this.autoLoadScrollEvent)
      }
    } else if (this.pageMode === 'pagination') {
      // 分页条
      const pagination = new Pagination(this.parentCommentsCount, {
        pageSize: this.pageSize,
        onChange: async (offset) => {
          await this.reqComments(offset)
          // 滚动到第一个评论的位置
          if (this.$parent) {
            let topPos = 0
            if (!this.autoLoadListenerAt && this.$parent) {
              topPos = Utils.getOffset(this.$parent).top
            }
            const at = this.autoLoadListenerAt || window
            at.scroll({
              top: topPos,
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

  public onError(msg: any) {
    msg = String(msg)
    console.error(msg)
    if (this.isFirstLoad || this.pageMode === 'pagination') {
      const errEl = Utils.createElement(`<span>${msg}，无法获取评论列表数据<br/></span>`)
      const retryBtn = Utils.createElement('<span style="cursor:pointer;">点击重新获取</span>')
      retryBtn.onclick = () => {
        this.reqComments(this.offset)
      }
      errEl.appendChild(retryBtn)
      const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
      adminBtn.onclick = () => {
        this.ctx.trigger('sidebar-show')
      }
      if (!this.ctx.user.data.isAdmin) {
        adminBtn.classList.add('atk-hide')
      }
      errEl.appendChild(adminBtn)
      Ui.setError(this.$el, errEl)
    } else {
      this.readMoreBtn?.showErr(`获取失败`)
    }
  }

  /** 刷新界面 */
  public refreshUI() {
    // 评论为空界面
    const noComment = this.comments.length <= 0
    let noCommentEl = this.$commentsWrap.querySelector<HTMLElement>('.atk-list-no-comment')

    if (noComment) {
      if (!noCommentEl) {
        noCommentEl = Utils.createElement('<div class="atk-list-no-comment"></div>')
        this.$commentsWrap.appendChild(noCommentEl)
        noCommentEl.innerHTML = this.noCommentText
      }
    }

    if (!noComment && noCommentEl)
      noCommentEl.remove()

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
  public importComments(rawData: CommentData[]) {
    // 查找并导入所有子评论
    const queryImportChildren = (parentC: Comment) => {
      const children = rawData.filter(o => o.rid === parentC.data.id)
      if (children.length === 0) return

      children.forEach((itemData: CommentData) => {
        itemData.is_allow_reply = parentC.data.is_allow_reply
        const childC = this.createComment(itemData)
        childC.render()
        parentC.putChild(childC)
        queryImportChildren(childC) // 递归
      })
    }

    // 开始处理 rawData
    if (!this.flatMode) {
      rawData.filter((o) => o.rid === 0).forEach((rootCommentData: CommentData) => {
        if (rootCommentData.is_collapsed) rootCommentData.is_allow_reply = false
        const rootComment = this.createComment(rootCommentData)
        rootComment.render()
        this.comments.push(rootComment) // 将评论导入 comments 总表中

        this.$commentsWrap.appendChild(rootComment.getEl())
        rootComment.playFadeInAnim()

        queryImportChildren(rootComment)
      })
    } else {
      // 平铺模式
      rawData.forEach((commentData: CommentData) => {
        this.putCommentFlatMode(commentData, rawData, 'append')
      })
    }

    this.eachComment(this.comments, (c) => {
      this.checkMoreHide(c)
    })

    this.refreshUI()
    this.ctx.trigger('comments-loaded')
  }

  /** 导入评论 · 平铺模式 */
  private putCommentFlatMode(commentItem: CommentData, comments: CommentData[], insertMode: 'append'|'prepend') {
    if (commentItem.is_collapsed) commentItem.is_allow_reply = false
    const comment = this.createComment(commentItem)
    if (commentItem.rid !== 0) {
      const rComment = comments.find(o => o.id === commentItem.rid)
      if (rComment) comment.replyTo = rComment
    }
    comment.render()

    // 将评论导入 comments 总表中
    if (insertMode === 'append') {
      this.comments.push(comment)
    } else {
      this.comments.unshift(comment)
    }

    if (commentItem.visible) {
      if (insertMode === 'append') {
        this.$commentsWrap.appendChild(comment.getEl())
      } else {
        this.$commentsWrap.prepend(comment.getEl())
      }
      comment.playFadeInAnim()
    }

    this.checkMoreHide(comment)
  }

  /** 插入评论 · 首部添加 */
  public insertComment(commentData: CommentData) {
    if (!this.flatMode) {
      const comment = this.createComment(commentData)
      comment.render()

      if (commentData.rid !== 0) {
        this.findComment(commentData.rid)?.putChild(comment)
      } else {
        this.$commentsWrap.prepend(comment.getEl())
        this.comments.unshift(comment)
      }

      Ui.scrollIntoView(comment.getEl()) // 滚动到可以见
      comment.playFadeInAnim() // 播放评论渐出动画

      this.checkMoreHide(comment)
    } else {
      this.putCommentFlatMode(commentData, this.comments.map(c => c.data), 'prepend')
    }

    if (this.data) this.data.total += 1 // 评论数增加 1
    this.refreshUI() // 更新 list 界面
    this.ctx.trigger('comments-loaded')
  }

  checkMoreHide(c: Comment) {
    const childrenH = this.ctx.conf.heightLimit?.children
    const contentH = this.ctx.conf.heightLimit?.content
    const isChildrenLimit = typeof childrenH === 'number' && childrenH > 0
    const isContentLimit = typeof contentH === 'number' && contentH > 0

    // 子评论内容过多隐藏
    if (isChildrenLimit && c.getIsRoot()) {
      c.checkMoreHide(c.$children, childrenH || 300)
    }

    // 评论内容过多隐藏
    if (isContentLimit) {
      c.checkMoreHide(c.$content, contentH || 200)
      if (c.$replyTo) c.checkMoreHide(c.$replyTo, contentH || 200) // 平铺模式回复内容
    }
  }

  /** 获取评论总数 (包括子评论) */
  get commentsCount(): number {
    let count = 0
    this.eachComment(this.comments, () => { count++ })
    return count
  }

  /** 是否还有更多的评论 */
  get hasMoreComments(): boolean {
    if (!this.data) return false
    return this.data.total_parents > (this.offset + this.pageSize)
  }

  get parentCommentsCount() {
    return this.comments.length
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

  /** 获取评论总数 */
  public getCommentCount(): number {
    let count = 0
    this.eachComment(this.comments, () => { count++ })
    return count
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

    this.refreshUI()
  }

  /** 清空所有评论 */
  public clearAllComments() {
    this.$commentsWrap.innerHTML = ''
    this.data = undefined
    this.comments = []
  }

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

  public apiVersionCheck(versionData: ApiVersionData): boolean {
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
        this.reqComments(0)
      }
      errEl.append(ignoreBtn)
      Ui.setError(this.ctx, errEl, '<span class="atk-warn-title">Artalk Warn</span>')
    }

    return needUpdate
  }
}

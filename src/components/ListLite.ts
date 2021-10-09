import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Api from '../lib/api'
import Comment from './Comment'
import { ListData, CommentData } from '~/types/artalk-data'

export default class ListLite extends Component {
  public el: HTMLElement
  private commentsWrapEl: HTMLElement
  private comments: Comment[] = []

  public data?: ListData
  private pageSize: number = 15 // 每次请求获取量
  private offset: number = 0
  public type?: 'all'|'mentions'|'mine'|'pending'

  private readMoreEl: HTMLElement
  private readMoreLoadingEl: HTMLElement
  private readMoreTextEl: HTMLElement

  public noCommentText: string
  public renderComment?: (comment: Comment) => void
  public paramsEditor?: (params: any) => void

  private isLoading: boolean = false
  public isFirstLoad = true

  public flatMode = false

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(`
    <div class="atk-list-lite">
      <div class="atk-list-comments-wrap"></div>
      <div class="atk-list-read-more" style="display: none;">
        <div class="atk-loading-icon" style="display: none;"></div>
        <span class="atk-text">查看更多</span>
      </div>
    </div>
    `)
    this.commentsWrapEl = this.el.querySelector('.atk-list-comments-wrap')!

    // 查看更多
    this.pageSize = this.conf.readMore ? (this.conf.readMore.pageSize || this.pageSize) : this.pageSize
    this.readMoreEl = this.el.querySelector('.atk-list-read-more')!
    this.readMoreLoadingEl = this.readMoreEl.querySelector('.atk-loading-icon')!
    this.readMoreTextEl = this.readMoreEl.querySelector('.atk-text')!
    this.readMoreEl.addEventListener('click', () => {
      this.readMore()
    })

    this.noCommentText = this.conf.noComment || '无评论'

    // 评论时间自动更新
    setInterval(() => {
      this.el.querySelectorAll<HTMLElement>('[data-atk-comment-date]').forEach(el => {
        const date = el.getAttribute('data-atk-comment-date')
        el.innerText = Utils.timeAgo(new Date(Number(date)))
      })
    }, 30 * 1000) // 30s 更新一次
  }

  public async reqComments(offset: number = 0) {
    if (offset === 0) {
      this.clearComments()
    }

    // set loading
    this.isLoading = true
    if (offset === 0) Ui.showLoading(this.el)
    else this.readMoreBtnSetLoading(true)

    const hideLoading = () => {
      // hide loading
      this.isLoading = false
      if (offset === 0) Ui.hideLoading(this.el)
      else this.readMoreBtnSetLoading(false)
    }

    let listData: ListData
    try {
      listData = await new Api(this.ctx).get(offset, this.type, this.paramsEditor)
    } catch (e: any) {
      this.onError(e.msg || String(e))
      throw e
    } finally {
      hideLoading()
    }

    // load data
    this.offset = offset

    try {
      this.onLoad(listData)
    } catch (e: any) {
      this.onError(String(e))
      throw e
    } finally {
      hideLoading()
    }
  }

  public onLoad (data: ListData) {
    this.data = data
    Ui.setError(this.el, null)
    this.importComments(data.comments)

    // 查看更多按钮
    if (this.hasMoreComments) {
      this.showReadMoreBtn()
    } else {
      this.hideReadMoreBtn()
    }

    // 滚动到底部自动加载
    if (this.isFirstLoad && this.hasMoreComments) {
      this.initScrollBottomAutoLoad()
    }

    this.isFirstLoad = false
  }

  public onError (msg: any) {
    msg = String(msg)
    console.error(msg)
    if (this.isFirstLoad) {
      const errEl = Utils.createElement(`<span>${msg}，无法获取评论列表数据<br/></span>`)
      const retryBtn = Utils.createElement('<span style="cursor:pointer;">点击重新获取</span>')
      retryBtn.onclick = () => {
        this.reqComments(0)
      }
      errEl.appendChild(retryBtn)
      const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
      adminBtn.onclick = () => {
        this.ctx.dispatchEvent('sidebar-show')
      }
      if (!this.ctx.user.data.isAdmin) {
        adminBtn.classList.add('atk-hide')
      }
      errEl.appendChild(adminBtn)
      Ui.setError(this.el, errEl)
    } else {
      this.readMoreBtnShowErr(`${msg} 获取失败`)
    }
  }

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

  /** 导入评论 - 通过请求数据 */
  public importComments (rawData: CommentData[]) {
    // 查找并导入所有子评论
    const queryImportChildren = (parentC: Comment) => {
      const children = rawData.filter(o => o.rid === parentC.data.id)
      if (children.length === 0) return

      children.forEach((itemData: CommentData) => {
        itemData.is_allow_reply = parentC.data.is_allow_reply
        const childC = this.createComment(itemData)
        childC.renderElem()
        parentC.putChild(childC)
        queryImportChildren(childC) // 递归
      })
    }

    // 开始处理 rawData
    if (!this.flatMode) {
      rawData.filter((o) => o.rid === 0).forEach((rootCommentData: CommentData) => {
        if (rootCommentData.is_collapsed) rootCommentData.is_allow_reply = false
        const rootComment = this.createComment(rootCommentData)
        rootComment.renderElem()
        this.comments.push(rootComment) // 将评论导入 comments 总表中

        this.commentsWrapEl.appendChild(rootComment.getEl())
        rootComment.playFadeInAnim()

        queryImportChildren(rootComment)
      })
    } else {
      rawData.forEach((commentData: CommentData) => {
        if (commentData.is_collapsed) commentData.is_allow_reply = false
        const comment = this.createComment(commentData)
        if (commentData.rid !== 0) {
          const rComment = rawData.find(o => o.id === commentData.rid)
          if (rComment) comment.replyTo = rComment
        }
        comment.renderElem()

        this.comments.push(comment) // 将评论导入 comments 总表中

        if (commentData.visible) {
          this.commentsWrapEl.appendChild(comment.getEl())
          comment.playFadeInAnim()
        }
      })
    }

    this.refreshUI()
  }

  public insertComment (commentData: CommentData) {
    const comment = this.createComment(commentData)
    comment.renderElem()

    if (commentData.rid !== 0) {
      this.findComment(commentData.rid)?.putChild(comment)
    } else {
      this.commentsWrapEl.prepend(comment.getEl())
      this.comments.unshift(comment)
    }

    Ui.scrollIntoView(comment.getEl()) // 滚动到可以见
    comment.playFadeInAnim() // 播放评论渐出动画
    if (this.data) this.data.total += 1 // 评论数增加 1
    this.refreshUI() // 更新 list 界面
  }

  /** 刷新界面 */
  public refreshUI () {
    // 评论为空界面
    const noComment = this.comments.length <= 0
    let noCommentEl = this.commentsWrapEl.querySelector<HTMLElement>('.atk-list-no-comment')

    if (noComment) {
      if (!noCommentEl) {
        noCommentEl = Utils.createElement('<div class="atk-list-no-comment"></div>')
        this.commentsWrapEl.appendChild(noCommentEl)
        noCommentEl.innerHTML = this.noCommentText
      }
    }

    if (!noComment && noCommentEl)
      noCommentEl.remove()

    // 仅管理员显示控制
    this.ctx.dispatchEvent('check-admin-show-el')
  }

  /** 获取评论总数 (包括子评论) */
  public getListCommentCount (): number {
    if (!this.data || !this.data.total) return 0
    return Number(this.data.total || '0')
  }

  /** 是否还有更多的评论 */
  get hasMoreComments (): boolean {
    if (!this.data) return false
    return this.data.total_parents > (this.offset + this.pageSize)
  }

  /** 阅读更多操作 */
  readMore () {
    const offset = this.offset + this.pageSize
    this.reqComments(offset)
  }

  /** 阅读更多按钮 - 显示 */
  showReadMoreBtn () {
    this.readMoreEl.style.display = ''
  }

  /** 阅读更多按钮 - 隐藏 */
  hideReadMoreBtn () {
    this.readMoreEl.style.display = 'none'
  }

  /** 阅读更多按钮 - 显示加载 */
  readMoreBtnSetLoading (isLoading: boolean) {
    this.readMoreLoadingEl.style.display = isLoading ? '' : 'none'
    this.readMoreTextEl.style.display = isLoading ? 'none' : ''
  }

  /** 阅读更多按钮 - 显示错误 */
  readMoreBtnShowErr (errMsg: string) {
    this.readMoreBtnSetLoading(false)

    const readMoreTextOrg = this.readMoreTextEl.innerText
    this.readMoreTextEl.innerText = errMsg
    this.readMoreEl.classList.add('atk-err')
    setTimeout(() => {
      this.readMoreTextEl.innerText = readMoreTextOrg
      this.readMoreEl.classList.remove('atk-err')
    }, 2000) // 2s后错误提示复原
  }

  /** 初始化滚动到底部自动查看更多（若开启） */
  initScrollBottomAutoLoad () {
    if (!this.conf.readMore) return
    if (!this.conf.readMore.autoLoad) return

    document.addEventListener('scroll', () => {
      const targetEl = this.el.querySelector<HTMLElement>('.atk-list-comments-wrap > .atk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
      if (!targetEl) return
      if (!this.hasMoreComments) return
      if (this.isLoading) return
      if (Ui.isVisible(targetEl)) {
        // 加载更多
        this.readMore()
      }
    })
  }

  /** 遍历操作 Comment (包括子评论) */
  public eachComment (commentList: Comment[], action: (comment: Comment, levelList: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  /** 查找评论项 */
  public findComment (id: number, src?: Comment[]): Comment|null {
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
  public getCommentCount (): number {
    let count = 0
    this.eachComment(this.comments, () => { count++ })
    return count
  }

  /** 删除评论 */
  public deleteComment (comment: number|Comment) {
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
  }

  /** 清空所有评论 */
  public clearComments () {
    this.commentsWrapEl.innerHTML = ''
    this.data = undefined
    this.comments = []
  }
}

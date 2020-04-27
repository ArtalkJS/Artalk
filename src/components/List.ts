import '../css/list.less'
import Comment from './Comment'
import { ListData } from '~/types/artalk-data'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

export default class List extends ArtalkContext {
  public el: HTMLElement
  public commentsWrapEl: HTMLElement
  public comments: Comment[] = []

  public readMoreEl: HTMLElement
  public readMoreLoadingEl: HTMLElement
  public readMoreTextEl: HTMLElement

  public data: ListData;
  public readonly limit = 5; // 每次请求获取量

  constructor () {
    super()

    this.el = Utils.createElement(require('../templates/List.ejs')(this))
    this.artalk.el.appendChild(this.el)

    this.commentsWrapEl = this.el.querySelector('.artalk-list-comments-wrap')
    this.reqComments()

    this.el.querySelector('[data-action="open-sidebar"]').addEventListener('click', () => {
      this.artalk.sidebar.show()
    })

    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.checkRedirectByUrlHash()
    })

    // 查看更多
    this.readMoreEl = this.el.querySelector('.artalk-list-read-more')
    this.readMoreLoadingEl = this.readMoreEl.querySelector('.artalk-loading-icon')
    this.readMoreTextEl = this.readMoreEl.querySelector('.artalk-text')

    this.readMoreEl.addEventListener('click', () => {
      this.readMore()
    })
  }

  /** 拉取评论 */
  reqComments (offset: Number = 0) {
    if (offset === 0) {
      this.clearComments();
    }

    this.artalk.request('CommentGet', {
      page_key: this.artalk.conf.pageKey,
      limit: this.limit, // 获取评论数
      offset, // 偏移量
    }, () => {
      if (offset === 0) this.artalk.ui.showLoading()
      else this.readMoreBtnSetLoading(true)
    }, () => {
      if (offset === 0) this.artalk.ui.hideLoading()
      else this.readMoreBtnSetLoading(false)
    }, (msg, data: ListData) => {
      this.data = { ...data }
      this.artalk.ui.setGlobalError(null)
      this.importCommentsByReqObj(data.comments)
      // 查看更多按钮
      if (this.hasMoreComments) this.showReadMoreBtn()
      else this.hideReadMoreBtn()

      this.checkRedirectByUrlHash()
    }, (msg, data) => {
      if (offset === 0) {
        const errEl = Utils.createElement(`<span>${msg}，无法获取评论列表数据<br/></span>`)
        const retryBtn = Utils.createElement('<span style="cursor:pointer">点击重新获取</span>')
        retryBtn.addEventListener('click', () => {
          this.reqComments(0)
        })
        errEl.appendChild(retryBtn)
        this.artalk.ui.setGlobalError(errEl)
      } else {
        this.readMoreBtnShowErr(`${msg} 获取失败`)
      }
    })
  }

  /** 导入评论 - 通过请求数据 */
  importCommentsByReqObj (rawData: any[]) {
    if (!Array.isArray(rawData)) { throw new Error('putCommentsByObj 出错：参数非数组') }

    const newCommentList: Comment[] = []
    rawData.forEach((commentData) => {
      if (commentData.id === 0) {
        throw new Error('黑人问号 ??? Comment 的 ID 怎么可能是 0 ?')
      }
      if (commentData.rid === 0) {
        newCommentList.push(new Comment(this, commentData))
      }
    })

    // 查找并导入所有子评论
    const queryChildren = (parentComment: Comment) => {
      rawData.forEach((commentData) => {
        if (commentData.rid === parentComment.data.id) {
          const cComment = new Comment(this, commentData)
          parentComment.setChild(cComment)
          cComment.playFadeInAnim()

          // 是否存在子评论
          const hasChild = (parentId: number) => {
            return rawData.find(o => o.rid === parentId) !== null
          }

          // 继续 递归查找子评论
          if (hasChild(cComment.data.id)) {
            queryChildren(cComment)
          }
        }
      })
    }

    // 查找并导入子评论
    newCommentList.forEach((comment) => {
      queryChildren(comment)
    })

    // 将评论导入 this.comments 总表中
    this.comments.push(...newCommentList)

    // 装载评论元素
    newCommentList.forEach((comment) => {
      this.commentsWrapEl.appendChild(comment.getElem())
      comment.playFadeInAnim()
    })

    this.refreshUI()
  }

  /** 添加评论项 */
  putComment (comment: Comment) {
    this.commentsWrapEl.prepend(comment.getElem())
    this.comments.unshift(comment)

    comment.playFadeInAnim()
    this.refreshUI()
  }

  /** 查找评论项 */
  findComment (id: number) {
    let comment: Comment|null = null

    const findCommentInList = (commentList: Comment[]) => {
      commentList.every((item) => {
        if (comment !== null) return false
        if (item.data.id === id) {
          comment = item
        }
        findCommentInList(item.getChildren())
        return true
      })
    }

    findCommentInList(this.comments)
    return comment
  }

  /** 清空所有评论 */
  clearComments () {
    this.commentsWrapEl.innerHTML = ''
    this.comments = []
    this.data = undefined
  }

  /** 刷新界面 */
  refreshUI () {
    (this.el.querySelector('.artalk-comment-count-num') as HTMLElement).innerText = this.data ? String(this.data.total || 0) : '0'

    let noCommentElem = this.commentsWrapEl.querySelector('.artalk-no-comment') as HTMLElement
    if (this.comments.length <= 0 && !noCommentElem) {
      noCommentElem = Utils.createElement('<div class="artalk-no-comment"></div>')
      noCommentElem.innerText = this.artalk.conf.noComment
      this.commentsWrapEl.appendChild(noCommentElem)
    }
    if (this.comments.length > 0 && noCommentElem !== null) {
        noCommentElem.remove()
    }
  }

  /** 跳到评论项位置 - 根据 `location.hash` */
  checkRedirectByUrlHash () {
    let commentId: number = Number(Utils.getLocationParmByName('artalk_comment'))
    if (!commentId) {
      const match = window.location.hash.match(/#artalk-comment-([0-9]+)/)
      if (!match || !match[1] || Number.isNaN(Number(match[1]))) return
      commentId = Number(match[1])
    }

    const comment = this.findComment(commentId)
    if (!comment && this.hasMoreComments) {
      this.readMore()
      return
    }
    if (!comment) { return }

    this.artalk.ui.scrollIntoView(comment.getElem())
    setTimeout(() => {
      comment.getElem().classList.add('artalk-flash-once')
    }, 800)
  }

  /** 是否还有更多的评论 */
  get hasMoreComments (): boolean {
    if (!this.data) return false
    return this.data.total_parents > (this.data.offset + this.data.limit)
  }

  /** 阅读更多操作 */
  readMore () {
    const offset = this.data.offset + this.limit
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
    this.showReadMoreBtn()
    this.readMoreLoadingEl.style.display = isLoading ? '' : 'none'
    this.readMoreTextEl.style.display = isLoading ? 'none' : ''
  }

  /** 阅读更多按钮 - 显示错误 */
  readMoreBtnShowErr (errMsg: string) {
    this.readMoreBtnSetLoading(false)

    const readMoreTextOrg = this.readMoreTextEl.innerText
    this.readMoreTextEl.innerText = errMsg
    this.readMoreEl.classList.add('artalk-err')
    setTimeout(() => {
      this.readMoreTextEl.innerText = readMoreTextOrg
      this.readMoreEl.classList.remove('artalk-err')
    }, 2000) // 2s后错误提示复原
  }
}

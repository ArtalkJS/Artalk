import md5 from 'md5'
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

  public data: ListData
  public reqPageSize: number = 15 // 每次请求获取量

  public isLoading: boolean = false

  public openSidebarBtnEl: HTMLElement

  constructor () {
    super()

    this.el = Utils.createElement(require('../templates/List.ejs')(this))
    this.artalk.el.appendChild(this.el)

    this.commentsWrapEl = this.el.querySelector('.artalk-list-comments-wrap')

    // 侧边栏呼出按钮
    this.openSidebarBtnEl = this.el.querySelector('[data-action="open-sidebar"]')
    this.openSidebarBtnEl.addEventListener('click', () => {
      this.artalk.sidebar.show()
    })

    // 锚点快速跳转评论
    window.addEventListener('hashchange', () => {
      this.checkRedirectByUrlHash()
    })

    // 查看更多
    this.reqPageSize = this.artalk.conf.readMore ? (this.artalk.conf.readMore.pageSize || this.reqPageSize) : this.reqPageSize
    this.readMoreEl = this.el.querySelector('.artalk-list-read-more')
    this.readMoreLoadingEl = this.readMoreEl.querySelector('.artalk-loading-icon')
    this.readMoreTextEl = this.readMoreEl.querySelector('.artalk-text')

    this.readMoreEl.addEventListener('click', () => {
      this.readMore()
    })

    // 请求获取评论
    this.reqComments()
  }

  /** 拉取评论 */
  reqComments (offset: Number = 0) {
    if (offset === 0) {
      this.clearComments();
    }

    this.artalk.request('CommentGet', {
      page_key: this.artalk.conf.pageKey,
      limit: this.reqPageSize, // 获取评论数
      offset, // 偏移量
    }, () => {
      this.isLoading = true
      if (offset === 0) this.artalk.ui.showLoading()
      else this.readMoreBtnSetLoading(true)
    }, () => {
      this.isLoading = false
      if (offset === 0) this.artalk.ui.hideLoading()
      else this.readMoreBtnSetLoading(false)
    }, (msg, data: ListData) => {
      this.data = { ...data }
      this.artalk.ui.setGlobalError(null)
      this.importCommentsByReqObj(data.comments)
      // 查看更多按钮
      if (this.hasMoreComments) this.showReadMoreBtn()
      else this.hideReadMoreBtn()
      // 锚点跳转
      this.checkRedirectByUrlHash()
      // 滚动到底部自动加载
      if (offset === 0 && this.hasMoreComments) {
        this.initScrollBottomAutoLoad()
      }
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

    this.eachComment(this.comments, (item) => {
      if (item.data.id === id) {
        comment = item
        return false
      }
      return true
    })

    return comment
  }

  /** 遍历操作 Comment (包括子评论) */
  eachComment (commentList: Comment[], action: (comment?: Comment, levelList?: Comment[]) => boolean|void) {
    if (commentList.length === 0) return
    commentList.every((item) => {
      if (action(item, commentList) === false) return false
      this.eachComment(item.getChildren(), action)
      return true
    })
  }

  /** 删除评论 */
  deleteComment (comment: number|Comment) {
    let findComment: Comment
    if (typeof comment === 'number') {
      findComment = this.findComment(comment)
      if (!findComment) throw Error(`未找到评论 ${comment}`)
    } else findComment = comment

    findComment.getElem().remove()
    this.eachComment(this.comments, (item, levelList) => {
      if (item === findComment) {
        levelList.splice(levelList.indexOf(item), 1)
        return false
      }
      return true
    })
  }

  /** 清空所有评论 */
  clearComments () {
    this.commentsWrapEl.innerHTML = ''
    this.comments = []
    this.data = undefined
  }

  /** 刷新界面 */
  refreshUI () {
    (this.el.querySelector('.artalk-comment-count-num') as HTMLElement).innerText = this.getCommentCount().toString()

    let noCommentElem = this.commentsWrapEl.querySelector('.artalk-no-comment') as HTMLElement
    if (this.comments.length <= 0 && !noCommentElem) {
      noCommentElem = Utils.createElement('<div class="artalk-no-comment"></div>')
      noCommentElem.innerText = this.artalk.conf.noComment
      this.commentsWrapEl.appendChild(noCommentElem)
    }
    if (this.comments.length > 0 && noCommentElem !== null) {
        noCommentElem.remove()
    }

    // 已输入个人信息
    if (!!this.artalk.user.nick && !!this.artalk.user.email) {
      this.openSidebarBtnEl.style.display = ''
    } else {
      this.openSidebarBtnEl.style.display = 'none'
    }

    // 是管理员
    if (this.artalk.user.isAdmin) {
      this.el.querySelectorAll('[data-only-admin-show]').forEach((itemEl: HTMLElement) => {
        itemEl.classList.remove('artalk-hide')
      })
    } else {
      this.el.querySelectorAll('[data-only-admin-show]').forEach((itemEl: HTMLElement) => {
        itemEl.classList.add('artalk-hide')
      })
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

    this.artalk.ui.scrollIntoView(comment.getElem(), false)
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
    const offset = this.data.offset + this.reqPageSize
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

  /** 初始化滚动到底部自动查看更多（若开启） */
  initScrollBottomAutoLoad () {
    if (!this.artalk.conf.readMore) return
    if (!this.artalk.conf.readMore.autoLoad) return

    document.addEventListener('scroll', () => {
      const targetEl = this.el.querySelector('.artalk-list-comments-wrap > .artalk-comment-wrap:nth-last-child(3)') // 获取倒数第3个评论元素
      if (!targetEl) return
      if (!this.hasMoreComments) return
      if (this.isLoading) return
      if (this.artalk.ui.isVisible(targetEl as HTMLElement)) {
        // 加载更多
        this.readMore()
      }
    })
  }

  /** 根据请求数据判断 nick 是否为管理员 */
  checkNickEmailIsAdmin (nick: string, email: string) {
    if (!this.data || !this.data.admin_nicks || !this.data.admin_encrypted_emails) return false

    return (this.data.admin_nicks.indexOf(nick) !== -1)
      && (this.data.admin_encrypted_emails.find(o => String(o).toLowerCase() === String(md5(email)).toLowerCase()))
  }

  /** 获取评论总数 (包括子评论) */
  getCommentCount (): number {
    let count = 0
    this.eachComment(this.comments, () => {
      count++
    })

    // 尝试从请求数据中获取
    if (this.data && typeof this.data.total !== 'undefined')
      if (Number(this.data.total) > count)
        return Number(this.data.total)

    return count
  }
}

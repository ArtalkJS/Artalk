import '../style/list.less'

import Context from '../Context'
import Component from '../lib/component'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import Api from '../lib/api'
import Comment from './Comment'
import { ListData, CommentData } from '~/types/artalk-data'
import ListHTML from './html/list.html?raw'

export default class List extends Component {
  public el: HTMLElement
  public commentsWrapEl: HTMLElement
  public comments: Comment[] = []

  public data?: ListData
  public pageSize: number = 15 // 每次请求获取量
  public offset: number = 0

  public readMoreEl: HTMLElement
  public readMoreLoadingEl: HTMLElement
  public readMoreTextEl: HTMLElement

  public isLoading: boolean = false

  public closeCommentBtnEl!: HTMLElement
  public openSidebarBtnEl!: HTMLElement

  public isFirstLoad = true

  constructor (ctx: Context) {
    super(ctx)

    this.el = Utils.createElement(ListHTML)
    this.commentsWrapEl = this.el.querySelector('.atk-list-comments-wrap')!

    // 操作按钮
    this.initListActionBtn()

    // 查看更多
    this.pageSize = this.conf.readMore ? (this.conf.readMore.pageSize || this.pageSize) : this.pageSize
    this.readMoreEl = this.el.querySelector('.atk-list-read-more')!
    this.readMoreLoadingEl = this.readMoreEl.querySelector('.atk-loading-icon')!
    this.readMoreTextEl = this.readMoreEl.querySelector('.atk-text')!
    this.readMoreEl.addEventListener('click', () => {
      this.readMore()
    })

    this.el.querySelector<HTMLElement>('.atk-copyright')!.innerHTML = `Powered By <a href="https://artalk.js.org" target="_blank" title="Artalk v${ARTALK_VERSION}">Artalk</a>`

    this.ctx.addEventListener('list-load', this.onLoad) // 装载事件
    this.ctx.addEventListener('list-error', this.onError) // 错误事件
    this.ctx.addEventListener('list-clear', () => (this.clearComments())) // 清空评论
    this.ctx.addEventListener('list-refresh-ui', () => (this.refreshUI()))
    this.ctx.addEventListener('list-import', (data) => (this.importComments(data)))
    this.ctx.addEventListener('list-insert', (data) => (this.insertComment(data)))
  }

  public async reqComments(offset: number = 0) {
    if (offset === 0) {
      this.clearComments()
    }

    // set loading
    this.isLoading = true
    if (offset === 0) Ui.showLoading(this.ctx)
    else this.readMoreBtnSetLoading(true)

    try {
      const listData = await new Api(this.ctx).get(offset)

      // load data
      this.offset = offset
      this.onLoad(listData)
    } catch (e: any) {
      this.onError(e.msg || String(e))
    } finally {
      // hide loading
      this.isLoading = false
      if (offset === 0) Ui.hideLoading(this.ctx)
      else this.readMoreBtnSetLoading(false)
    }
  }

  public onLoad (data: ListData) {
    this.data = data
    Ui.setGlobalError(this.ctx, null)
    this.importComments(data.comments)

    // 查看更多按钮
    if (this.hasMoreComments) {
      this.showReadMoreBtn()
    } else {
      this.hideReadMoreBtn()
    }

    // 检测锚点跳转
    this.checkGoToCommentByUrlHash()

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
      const retryBtn = Utils.createElement('<span style="cursor:pointer">点击重新获取</span>')
      retryBtn.addEventListener('click', () => {
        this.reqComments(0)
      })
      errEl.appendChild(retryBtn)
      Ui.setGlobalError(this.ctx, errEl)
    } else {
      this.readMoreBtnShowErr(`${msg} 获取失败`)
    }
  }

  /** 导入评论 - 通过请求数据 */
  private importComments (rawData: CommentData[]) {
    // 查找并导入所有子评论
    const queryImportChildren = (parentC: Comment) => {
      const children = rawData.filter(o => o.rid === parentC.data.id)
      if (children.length === 0) return

      children.forEach((itemData: CommentData) => {
        itemData.is_allow_reply = parentC.data.is_allow_reply
        const childC = new Comment(this.ctx, itemData)
        parentC.putChild(childC)

        queryImportChildren(childC) // 递归
      })
    }

    // 开始处理 rawData
    rawData.filter((o) => o.rid === 0).forEach((rootCommentData: CommentData) => {
      if (rootCommentData.is_collapsed) rootCommentData.is_allow_reply = false
      const rootComment = new Comment(this.ctx, rootCommentData)
      this.comments.push(rootComment) // 将评论导入 comments 总表中

      this.commentsWrapEl.appendChild(rootComment.getEl())
      rootComment.playFadeInAnim()

      queryImportChildren(rootComment)
    })

    this.refreshUI(true)
  }

  public insertComment (commentData: CommentData) {
    const comment = new Comment(this.ctx, commentData)
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
  public refreshUI (isFirstCall: boolean = false) {
    this.el.querySelector<HTMLElement>('.atk-comment-count-num')!.innerText = String(this.getListCommentCount())

    // 评论为空界面
    let noCommentElem = this.commentsWrapEl.querySelector<HTMLElement>('.atk-no-comment')
    const isNoComment = this.comments.length <= 0

    if (isNoComment && !noCommentElem) {
      noCommentElem = Utils.createElement('<div class="atk-no-comment"></div>')
      this.commentsWrapEl.appendChild(noCommentElem)
      noCommentElem.innerText = this.conf.noComment || '无评论'
    }
    if (!isNoComment && noCommentElem) {
      noCommentElem.remove()
    }

    // 已输入个人信息
    if (!!this.ctx.user.data.nick && !!this.ctx.user.data.email) {
      this.openSidebarBtnEl.classList.remove('atk-hide')
    } else {
      this.openSidebarBtnEl.classList.add('atk-hide')
    }

    // 仅管理员显示控制
    this.ctx.dispatchEvent('check-admin-show-el')

    // 关闭评论
    if (!!this.data && !!this.data.page && this.data.page.admin_only === true) {
      this.ctx.dispatchEvent('editor-close-comment')
      this.closeCommentBtnEl.innerHTML = '打开评论'
    } else if (!isFirstCall) {
      this.ctx.dispatchEvent('editor-open-comment')
      this.closeCommentBtnEl.innerHTML = '关闭评论'
    }
  }

  private initListActionBtn () {

    // 侧边栏呼出按钮
    this.openSidebarBtnEl = this.el.querySelector('[data-action="open-sidebar"]')!
    this.openSidebarBtnEl.addEventListener('click', () => {
      this.ctx.dispatchEvent('sidebar-show')
    })

    // 关闭评论按钮
    this.closeCommentBtnEl = this.el.querySelector('[data-action="admin-close-comment"]')!
    this.closeCommentBtnEl.addEventListener('click', () => {
      if (!this.data) return

      this.adminSetPage({
        admin_only: !this.data.page.admin_only
      })
    })
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

  /** 管理员设置页面信息 */
  public adminSetPage (conf: { admin_only: boolean }) {
    this.ctx.dispatchEvent('editor-show-loading')
    new Api(this.ctx).editPage(this.data?.page.page_key || '', conf.admin_only)
      .then((page) => {
        if (this.data)
          this.data.page = { ...page }
        this.refreshUI()
      })
      .catch(err => {
        this.ctx.dispatchEvent('editor-notify', { msg: `修改页面数据失败：${err.msg || String(err)}`, type: 'e'})
      })
      .finally(() => {
        this.ctx.dispatchEvent('editor-hide-loading')
      })
  }

  /** 跳到评论项位置 - 根据 `location.hash` */
  public checkGoToCommentByUrlHash () {
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

    Ui.scrollIntoView(comment.getEl(), false)
    setTimeout(() => {
      comment.getEl().classList.add('atk-flash-once')
    }, 800)
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
  public findComment (id: number): Comment|null {
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

import '../css/comment.scss'
import List from './List'
import UADetect from '../utils/detect'
import { CommentData } from '~/types/artalk-data'
import ArtalkContext from '../ArtalkContext'
import Utils from '../utils'

export default class Comment extends ArtalkContext {
  public elem: HTMLElement
  public contentEl: HTMLElement
  public childrenEl: HTMLElement

  public parent: Comment|null
  public nestedNum: number
  private readonly maxNestedNo = 3 // 最多嵌套层数
  public children: Comment[] = []

  constructor (public list: List, public data: CommentData) {
    super()

    this.data = { ...this.data }

    this.elem = Utils.createElement(require('../templates/Comment.ejs')(this))
    this.contentEl = this.elem.querySelector('.artalk-content')

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层
    this.childrenEl = null

    // 绑定回复按钮事件
    this.elem.querySelector('[data-comment-action="reply"]').addEventListener('click', () => {
      this.artalk.editor.setReply(this)
    })
  }

  setChild (comment: Comment) {
    this.children.push(comment)
    this.getChildrenEl().appendChild(comment.getElem())
    comment.parent = this
    comment.nestedNum = this.nestedNum + 1 // 嵌套层数 +1
  }

  getChildren () {
    return this.children
  }

  getChildrenEl () {
    if (this.childrenEl === null) {
      // console.log(this.nestedNo)
      if (this.nestedNum < this.maxNestedNo) {
        this.childrenEl = Utils.createElement('<div class="artalk-comment-children"></div>')
        this.contentEl.appendChild(this.childrenEl)
      } else {
        this.childrenEl = this.parent.getChildrenEl()
      }
    }
    return this.childrenEl
  }

  getParent () {
    return this.parent
  }

  getElem () {
    return this.elem
  }

  getData () {
    return this.data
  }

  getGravatarUrl () {
    return `${this.artalk.conf.gravatar.cdn}${this.data.email_encrypted}?d=${encodeURIComponent(this.artalk.conf.defaultAvatar)}&s=80`
  }

  getContentMarked () {
    return this.artalk.marked(this.data.content)
  }

  padWithZeros (vNumber: number, width: number) {
    let numAsString = vNumber.toString()
    while (numAsString.length < width) {
      numAsString = `0${numAsString}`
    }
    return numAsString
  }

  dateFormat (date: Date) {
    const vDay = this.padWithZeros(date.getDate(), 2)
    const vMonth = this.padWithZeros(date.getMonth() + 1, 2)
    const vYear = this.padWithZeros(date.getFullYear(), 2)
    // var vHour = padWithZeros(date.getHours(), 2);
    // var vMinute = padWithZeros(date.getMinutes(), 2);
    // var vSecond = padWithZeros(date.getSeconds(), 2);
    return `${vYear}-${vMonth}-${vDay}`
  }

  timeAgo (date: Date) {
    try {
      const oldTime = date.getTime()
      const currTime = new Date().getTime()
      const diffValue = currTime - oldTime

      const days = Math.floor(diffValue / (24 * 3600 * 1000))
      if (days === 0) {
        // 计算相差小时数
        const leave1 = diffValue % (24 * 3600 * 1000) // 计算天数后剩余的毫秒数
        const hours = Math.floor(leave1 / (3600 * 1000))
        if (hours === 0) {
          // 计算相差分钟数
          const leave2 = leave1 % (3600 * 1000) // 计算小时数后剩余的毫秒数
          const minutes = Math.floor(leave2 / (60 * 1000))
          if (minutes === 0) {
            // 计算相差秒数
            const leave3 = leave2 % (60 * 1000) // 计算分钟数后剩余的毫秒数
            const seconds = Math.round(leave3 / 1000)
            return `${seconds} 秒前`
          }
          return `${minutes} 分钟前`
        }
        return `${hours} 小时前`
      }
      if (days < 0) return '刚刚'

      if (days < 8) {
        return `${days} 天前`
      }

      return this.dateFormat(date)
    } catch (error) {
      console.error(error)
      return ' - '
    }
  }

  getUserUaBrowser () {
    const info = UADetect(this.data.ua)
    return `${info.browser} ${info.version}`
  }

  getUserUaOS () {
    const info = UADetect(this.data.ua)
    return `${info.os} ${info.osVersion}`
  }
}

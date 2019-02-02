import $ from 'jquery'
import '../css/comment.scss'
import UADetect from '../js/detect.js'

export default class Comment {
  constructor (list, data) {
    this.artalk = list.artalk
    this.data = {
      id: null,
      content: null,
      nick: null,
      email: null,
      email_encrypted: null,
      link: null,
      rid: null,
      ua: null,
      date: null
    }

    this.data = Object.assign(this.data, data)
    this.elem = $(require('../templates/Comment.ejs')(this))
    this.contentEl = this.elem.find('.artalk-content')

    this.parent = null
    this.nestedNum = 1 // 现在已嵌套 n 层
    this.maxNestedNo = 3 // 最多嵌套 n 层
    this.children = []
    this.childrenEl = null

    // 绑定回复按钮事件
    this.elem.find('[data-comment-action="reply"]').click(() => {
      this.artalk.editor.setReply(this)
    })
  }

  setChild (comment) {
    this.children.push(comment)
    $(comment.getElem()).appendTo(this.getChildrenEl())
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
        this.childrenEl = $('<div class="artalk-comment-children"></div>').appendTo(this.contentEl)
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
    return `${this.artalk.opts.gravatar.cdn}${this.data.email_encrypted}?d=${encodeURIComponent(this.artalk.opts.defaultAvatar)}&s=80`
  }

  getContentMarked () {
    return this.artalk.marked(this.data.content)
  }

  padWithZeros (vNumber, width) {
    var numAsString = vNumber.toString()
    while (numAsString.length < width) {
      numAsString = '0' + numAsString
    }
    return numAsString
  }

  dateFormat (date) {
    var vDay = this.padWithZeros(date.getDate(), 2)
    var vMonth = this.padWithZeros(date.getMonth() + 1, 2)
    var vYear = this.padWithZeros(date.getFullYear(), 2)
    // var vHour = padWithZeros(date.getHours(), 2);
    // var vMinute = padWithZeros(date.getMinutes(), 2);
    // var vSecond = padWithZeros(date.getSeconds(), 2);
    return `${vYear}-${vMonth}-${vDay}`
  }

  timeAgo (date) {
    try {
      var oldTime = date.getTime()
      var currTime = new Date().getTime()
      var diffValue = currTime - oldTime

      var days = Math.floor(diffValue / (24 * 3600 * 1000))
      if (days === 0) {
        // 计算相差小时数
        var leave1 = diffValue % (24 * 3600 * 1000) // 计算天数后剩余的毫秒数
        var hours = Math.floor(leave1 / (3600 * 1000))
        if (hours === 0) {
          // 计算相差分钟数
          var leave2 = leave1 % (3600 * 1000) // 计算小时数后剩余的毫秒数
          var minutes = Math.floor(leave2 / (60 * 1000))
          if (minutes === 0) {
            // 计算相差秒数
            var leave3 = leave2 % (60 * 1000) // 计算分钟数后剩余的毫秒数
            var seconds = Math.round(leave3 / 1000)
            return seconds + ' 秒前'
          }
          return minutes + ' 分钟前'
        }
        return hours + ' 小时前'
      }
      if (days < 0) return '刚刚'

      if (days < 8) {
        return days + ' 天前'
      } else {
        return this.dateFormat(date)
      }
    } catch (error) {
      console.error(error)
    }
  }

  getUserUaBrowser () {
    let info = UADetect(this.data.ua)
    return `${info.browser} ${info.version}`
  }

  getUserUaOS () {
    let info = UADetect(this.data.ua)
    return `${info.os} ${info.osVersion}`
  }
}

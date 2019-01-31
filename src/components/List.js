import $ from 'jquery'
import Comment from './Comment.js'
import '../css/list.scss'

class List {
  constructor (artalk) {
    this.artalk = artalk

    this.el = $(require('../templates/List.ejs')(this)).appendTo(this.artalk.el)
    this.initComments()
  }

  initComments () {
    this.comments = []
    this.commentsWrapEl = this.el.find('.artalk-list-comments-wrap')

    $.ajax({
      type: 'POST',
      url: this.artalk.opts.serverUrl,
      data: {
        action: 'CommentGet',
        page_key: this.artalk.opts.pageKey
      },
      dataType: 'json',
      beforeSend: () => {
        this.artalk.showLoading()
      },
      success: (obj) => {
        this.artalk.hideLoading()
        this.putAllCommentByObj(obj.data.comments)
      },
      error: () => {
        this.artalk.hideLoading()
        this.artalk.setGlobalError('网络错误，无法获取评论列表数据')
      }
    })
  }

  putAllCommentByObj (rawData) {
    let comments = []
    for (let i in rawData) {
      let comment = rawData[i]
      if (comment.id === 0) {
        throw Error('黑人问号 ??? Comment 的 ID 怎么可能是 0 ?')
      }
      if (comment.rid === 0) {
        comments.push(new Comment(this, comment))
      }
    }

    let isChildExistByParentId = (parentId) => {
      for (let i in rawData) {
        if (rawData[i].rid === parentId) {
          return true
        }
      }
      return false
    }

    // 导入子评论
    for (let i in comments) {
      let comment = comments[i]

      // 查找所有子评论
      let queryChildren = (parentComment) => {
        for (let cI in rawData) {
          let cData = rawData[cI]
          if (cData.rid === parentComment.data.id) {
            let cComment = new Comment(this, cData)
            parentComment.setChild(cComment)

            // 递归查找子评论的子评论
            if (isChildExistByParentId(cComment.data.id)) {
              queryChildren(cComment)
            }
          }
        }
      }

      queryChildren(comment)
    }

    this.comments = comments
    for (let i in this.comments) {
      let comment = this.comments[i]
      $(comment.getElem()).appendTo(this.commentsWrapEl)
    }

    this.updateIndicator()
  }

  putOneComment (comment) {
    $(comment.getElem()).prependTo(this.commentsWrapEl)
    this.comments.unshift(comment)
    this.updateIndicator()
  }

  updateIndicator () {
    this.el.find('.artalk-comment-count-num').text(this.comments.length)

    let noCommentElem = this.commentsWrapEl.find('.artalk-no-comment')
    if (this.comments.length <= 0) {
      if (!noCommentElem.length) {
        noCommentElem = $('<div class="artalk-no-comment"></div>').text(this.artalk.opts.noComment)
        noCommentElem.appendTo(this.commentsWrapEl)
      }
    } else {
      noCommentElem.remove()
    }
  }
}

export default List

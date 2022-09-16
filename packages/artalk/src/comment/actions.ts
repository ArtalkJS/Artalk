import Comment from './comment'
import ActionBtn from '../components/action-btn'

export default class CommentActions {
  private comment: Comment

  private get ctx() { return this.comment.ctx }
  private get data() { return this.comment.getData() }
  private get cConf() { return this.comment.getConf() }

  public constructor(comment: Comment) {
    this.comment = comment
  }

  /** 投票操作 */
  public vote(type: 'up'|'down') {
    const actionBtn = (type === 'up') ? this.comment.getRender().voteBtnUp : this.comment.getRender().voteBtnDown

    this.ctx.getApi().comment.vote(this.data.id, `comment_${type}`)
    .then((v) => {
      this.data.vote_up = v.up
      this.data.vote_down = v.down
      this.comment.getRender().voteBtnUp?.updateText()
      this.comment.getRender().voteBtnDown?.updateText()
    })
    .catch((err) => {
      actionBtn?.setError(this.ctx.$t('voteFail'))
      console.log(err)
    })
  }

  /** 管理员 - 评论状态修改 */
  public adminEdit(type: 'collapsed'|'pending'|'pinned', btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在修改中

    btnElem.setLoading(true, `${this.ctx.$t('editing')}...`)

    // 克隆并修改当前数据
    const modify = { ...this.data }
    if (type === 'collapsed') {
      modify.is_collapsed = !modify.is_collapsed
    } else if (type === 'pending') {
      modify.is_pending = !modify.is_pending
    } else if (type === 'pinned') {
      modify.is_pinned = !modify.is_pinned
    }

    this.ctx.getApi().comment.commentEdit(modify).then((data) => {
      btnElem.setLoading(false)

      // 刷新当前 Comment UI
      this.comment.setData(data)

      // 刷新 List UI
      this.ctx.listRefreshUI()
    }).catch((err) => {
      console.error(err)
      btnElem.setError(this.ctx.$t('editFail'))
    })
  }

  /** 管理员 - 评论删除 */
  public adminDelete(btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在删除中

    btnElem.setLoading(true, `${this.ctx.$t('deleting')}...`)
    this.ctx.getApi().comment.commentDel(this.data.id, this.data.site_name)
      .then(() => {
        btnElem.setLoading(false)
        if (this.cConf.onDelete) this.cConf.onDelete(this.comment)
      })
      .catch((e) => {
        console.error(e)
        btnElem.setError(this.ctx.$t('deleteFail'))
      })
  }

  /** 快速跳转到该评论 */
  public goToReplyComment() {
    const origHash = window.location.hash
    const modifyHash = `#atk-comment-${this.data.rid}`

    window.location.hash = modifyHash
    if (modifyHash === origHash) window.dispatchEvent(new Event('hashchange')) // 强制触发事件
  }
}

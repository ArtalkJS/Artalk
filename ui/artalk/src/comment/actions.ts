import ActionBtn from '@/components/action-btn'
import $t from '@/i18n'
import type { CommentNode } from '.'

export default class CommentActions {
  private comment: CommentNode

  private get data() { return this.comment.getData() }
  private get opts() { return this.comment.getOpts() }
  private getApi() { return this.comment.getOpts().getApi() }

  public constructor(comment: CommentNode) {
    this.comment = comment
  }

  /** 投票操作 */
  public vote(type: 'up'|'down') {
    const actionBtn = (type === 'up') ? this.comment.getRender().voteBtnUp : this.comment.getRender().voteBtnDown

    this.getApi().votes.vote(`comment_${type}`, this.data.id, { ...this.getApi().getUserFields() })
    .then((res) => {
      this.data.vote_up = res.data.up
      this.data.vote_down = res.data.down
      this.comment.getRender().voteBtnUp?.updateText()
      this.comment.getRender().voteBtnDown?.updateText()
    })
    .catch((err) => {
      actionBtn?.setError($t('voteFail'))
      console.log(err)
    })
  }

  /** 管理员 - 评论状态修改 */
  public adminEdit(type: 'collapsed'|'pending'|'pinned', btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在修改中

    btnElem.setLoading(true, `${$t('editing')}...`)

    // 克隆并修改当前数据
    const modify = { ...this.data }
    if (type === 'collapsed') {
      modify.is_collapsed = !modify.is_collapsed
    } else if (type === 'pending') {
      modify.is_pending = !modify.is_pending
    } else if (type === 'pinned') {
      modify.is_pinned = !modify.is_pinned
    }

    this.getApi().comments.updateComment(this.data.id, {
      ...modify,
    }).then((res) => {
      btnElem.setLoading(false)

      // 刷新当前 Comment UI
      this.comment.setData(res.data)
    }).catch((err) => {
      console.error(err)
      btnElem.setError($t('editFail'))
    })
  }

  /** 管理员 - 评论删除 */
  public adminDelete(btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在删除中

    btnElem.setLoading(true, `${$t('deleting')}...`)
    this.getApi().comments.deleteComment(this.data.id)
      .then(() => {
        btnElem.setLoading(false)
        if (this.opts.onDelete) this.opts.onDelete(this.comment)
      })
      .catch((e) => {
        console.error(e)
        btnElem.setError($t('deleteFail'))
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

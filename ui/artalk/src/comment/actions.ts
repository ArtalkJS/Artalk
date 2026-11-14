import type { CommentNode } from '.'
import $t from '@/i18n'
import ActionBtn from '@/components/action-btn'

export default class CommentActions {
  private comment: CommentNode
  private voting = false

  private get data() {
    return this.comment.getData()
  }
  private get opts() {
    return this.comment.getOpts()
  }
  private getApi() {
    return this.opts.getApi()
  }

  public constructor(comment: CommentNode) {
    this.comment = comment
  }

  /** 投票操作 */
  public vote(choice: 'up' | 'down') {
    if (this.voting || this.comment.isRemoved()) return
    this.voting = true
    const actionBtn =
      choice === 'up' ? this.comment.getRender().voteBtnUp : this.comment.getRender().voteBtnDown
    actionBtn?.setLoading(true)

    const api = this.getApi()
    api.votes
      .createVote('comment', this.data.id, choice, { ...api.getUserFields() })
      .then(({ data: res }) => {
        if (this.comment.isRemoved()) return
        this.data.vote_up = res.up
        this.data.vote_down = res.down
        this.comment.getRender().voteBtnUp?.updateText()
        this.comment.getRender().voteBtnDown?.updateText()
      })
      .catch((err) => {
        if (this.comment.isRemoved()) return
        const render = this.comment.getRender()
        const button = choice === 'up' ? render.voteBtnUp : render.voteBtnDown
        button?.setError($t('voteFail'))
        console.error(err)
      })
      .finally(() => {
        this.voting = false
      })
  }

  /** 管理员 - 评论状态修改 */
  public adminEdit(type: 'collapsed' | 'pending' | 'pinned', btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在修改中

    btnElem.setLoading(true, `${$t('editing')}...`)

    const key = `is_${type}` as const
    this.getApi()
      .comments.updateComment(this.data.id, { ...this.data, [key]: !this.data[key] })
      .then(({ data: res }) => {
        if (this.comment.isRemoved()) return
        btnElem.setLoading(false)

        // 刷新当前 Comment UI
        this.comment.setData(res)
      })
      .catch((err) => {
        if (this.comment.isRemoved()) return
        console.error(err)
        btnElem.setError($t('editFail'))
      })
  }

  /** 管理员 - 评论删除 */
  public adminDelete(btnElem: ActionBtn) {
    if (btnElem.isLoading) return // 若正在删除中

    btnElem.setLoading(true, `${$t('deleting')}...`)
    this.getApi()
      .comments.deleteComment(this.data.id)
      .then(() => {
        if (this.comment.isRemoved()) return
        btnElem.setLoading(false)
        if (this.opts.onDelete) this.opts.onDelete(this.comment)
      })
      .catch((e) => {
        if (this.comment.isRemoved()) return
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

import ActionBtn from '../../components/action-btn'
import type Render from '../render'
import $t from '@/i18n'

/**
 * 评论操作按钮界面
 */
export default function renderActions(r: Render) {
  Object.entries({
    renderVote,
    renderReply,
    // 管理员操作
    renderCollapse,
    renderModerator,
    renderPin,
    renderEdit,
    renderDel,
  }).forEach(([name, render]) => {
    render(r)
  })
}

// 操作按钮 - 投票
function renderVote(r: Render) {
  if (!r.opts.vote) return // 关闭投票功能

  // 赞同按钮
  r.voteBtnUp = new ActionBtn(() => `${$t('voteUp')} (${r.data.vote_up || 0})`).appendTo(r.$actions)
  r.voteBtnUp.setClick(() => {
    r.comment.getActions().vote('up')
  })

  // 反对按钮
  if (r.opts.voteDown) {
    r.voteBtnDown = new ActionBtn(() => `${$t('voteDown')} (${r.data.vote_down || 0})`).appendTo(
      r.$actions,
    )
    r.voteBtnDown.setClick(() => {
      r.comment.getActions().vote('down')
    })
  }
}

// 操作按钮 - 回复
function renderReply(r: Render) {
  if (!r.data.is_allow_reply) return // 不允许回复
  const replyBtn = new ActionBtn({
    text: $t('reply'),
  })
  replyBtn.appendTo(r.$actions)
  replyBtn.setClick(() => {
    r.opts.replyComment(r.data, r.$el)
  })
}

// 操作按钮 - 折叠
function renderCollapse(r: Render) {
  const collapseBtn = new ActionBtn({
    text: () => (r.data.is_collapsed ? $t('expand') : $t('collapse')),
    adminOnly: true,
  })
  collapseBtn.appendTo(r.$actions)
  collapseBtn.setClick(() => {
    r.comment.getActions().adminEdit('collapsed', collapseBtn)
  })
}

// 操作按钮 - 审核
function renderModerator(r: Render) {
  const pendingBtn = new ActionBtn({
    text: () => (r.data.is_pending ? $t('pending') : $t('approved')),
    adminOnly: true,
  })
  pendingBtn.appendTo(r.$actions)
  pendingBtn.setClick(() => {
    r.comment.getActions().adminEdit('pending', pendingBtn)
  })
}

// 操作按钮 - 置顶
function renderPin(r: Render) {
  const pinnedBtn = new ActionBtn({
    text: () => (r.data.is_pinned ? $t('unpin') : $t('pin')),
    adminOnly: true,
  })
  pinnedBtn.appendTo(r.$actions)
  pinnedBtn.setClick(() => {
    r.comment.getActions().adminEdit('pinned', pinnedBtn)
  })
}

// 操作按钮 - 编辑
function renderEdit(r: Render) {
  const editBtn = new ActionBtn({
    text: $t('edit'),
    adminOnly: true,
  })
  editBtn.appendTo(r.$actions)
  editBtn.setClick(() => {
    r.opts.editComment(r.data, r.$el)
  })
}

// 操作按钮 - 删除
function renderDel(r: Render) {
  const delBtn = new ActionBtn({
    text: $t('delete'),
    confirm: true,
    confirmText: $t('deleteConfirm'),
    adminOnly: true,
  })
  delBtn.appendTo(r.$actions)
  delBtn.setClick(() => {
    r.comment.getActions().adminDelete(delBtn)
  })
}

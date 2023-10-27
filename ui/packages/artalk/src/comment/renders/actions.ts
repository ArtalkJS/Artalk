import * as Utils from '../../lib/utils'
import ActionBtn from '../../components/action-btn'
import Render from '../render'

/**
 * 评论操作按钮界面
 */
export default function renderActions(r: Render) {
  Object.entries({
    renderVote, renderReply,
    // 管理员操作
    renderCollapse, renderModerator, renderPin, renderEdit, renderDel
  }).forEach(([name, render]) => {
    render(r)
  })
}


// 操作按钮 - 投票
function renderVote(r: Render) {
  if (!r.ctx.conf.vote) return // 关闭投票功能

  // 赞同按钮
  r.voteBtnUp = new ActionBtn(r.ctx, () => `${r.ctx.$t('voteUp')} (${r.data.vote_up || 0})`).appendTo(r.$actions)
  r.voteBtnUp.setClick(() => {
    r.comment.getActions().vote('up')
  })

  // 反对按钮
  if (r.ctx.conf.voteDown) {
    r.voteBtnDown = new ActionBtn(r.ctx, () => `${r.ctx.$t('voteDown')} (${r.data.vote_down || 0})`).appendTo(r.$actions)
    r.voteBtnDown.setClick(() => {
      r.comment.getActions().vote('down')
    })
  }
}

// 操作按钮 - 回复
function renderReply(r: Render) {
  if (!r.data.is_allow_reply) return // 不允许回复

  const replyBtn = Utils.createElement(`<span>${r.ctx.$t('reply')}</span>`)
  r.$actions.append(replyBtn)
  replyBtn.addEventListener('click', (e) => {
    e.stopPropagation() // 防止穿透
    if (!r.cConf.onReplyBtnClick) {
      r.ctx.replyComment(r.data, r.$el)
    } else {
      r.cConf.onReplyBtnClick()
    }
  })
}

// 操作按钮 - 折叠
function renderCollapse(r: Render) {
  const collapseBtn = new ActionBtn(r.ctx, {
    text: () => (r.data.is_collapsed ? r.ctx.$t('expand') : r.ctx.$t('collapse')),
    adminOnly: true
  })
  collapseBtn.appendTo(r.$actions)
  collapseBtn.setClick(() => {
    r.comment.getActions().adminEdit('collapsed', collapseBtn)
  })
}

// 操作按钮 - 审核
function renderModerator(r: Render) {
  const pendingBtn = new ActionBtn(r.ctx, {
    text: () => (r.data.is_pending ? r.ctx.$t('pending') : r.ctx.$t('approved')),
    adminOnly: true
  })
  pendingBtn.appendTo(r.$actions)
  pendingBtn.setClick(() => {
    r.comment.getActions().adminEdit('pending', pendingBtn)
  })
}

// 操作按钮 - 置顶
function renderPin(r: Render) {
  const pinnedBtn = new ActionBtn(r.ctx, {
    text: () => (r.data.is_pinned ? r.ctx.$t('unpin') : r.ctx.$t('pin')),
    adminOnly: true
  })
  pinnedBtn.appendTo(r.$actions)
  pinnedBtn.setClick(() => {
    r.comment.getActions().adminEdit('pinned', pinnedBtn)
  })
}

// 操作按钮 - 编辑
function renderEdit(r: Render) {
  const editBtn = new ActionBtn(r.ctx, {
    text: r.ctx.$t('edit'),
    adminOnly: true
  })
  editBtn.appendTo(r.$actions)
  editBtn.setClick(() => {
    r.ctx.editComment(r.data, r.$el)
  })
}

// 操作按钮 - 删除
function renderDel(r: Render) {
  const delBtn = new ActionBtn(r.ctx, {
    text: r.ctx.$t('delete'),
    confirm: true,
    confirmText: r.ctx.$t('deleteConfirm'),
    adminOnly: true,
  })
  delBtn.appendTo(r.$actions)
  delBtn.setClick(() => {
    r.comment.getActions().adminDelete(delBtn)
  })
}

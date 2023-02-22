import * as Utils from '../../lib/utils'
import ActionBtn from '../../components/action-btn'
import RenderCtx from '../render-ctx'

/**
 * 评论操作按钮界面
 */
export default function renderActions(ctx: RenderCtx) {
  Object.entries({
    renderVote, renderReply,
    // 管理员操作
    renderCollapse, renderModerator, renderPin, renderEdit, renderDel
  }).forEach(([name, render]) => {
    render(ctx)
  })
}


// 操作按钮 - 投票
function renderVote(ctx: RenderCtx) {
  if (!ctx.ctx.conf.vote) return // 关闭投票功能

  // 赞同按钮
  ctx.voteBtnUp = new ActionBtn(ctx.ctx, () => `${ctx.ctx.$t('voteUp')} (${ctx.data.vote_up || 0})`).appendTo(ctx.$actions)
  ctx.voteBtnUp.setClick(() => {
    ctx.comment.getActions().vote('up')
  })

  // 反对按钮
  if (ctx.ctx.conf.voteDown) {
    ctx.voteBtnDown = new ActionBtn(ctx.ctx, () => `${ctx.ctx.$t('voteDown')} (${ctx.data.vote_down || 0})`).appendTo(ctx.$actions)
    ctx.voteBtnDown.setClick(() => {
      ctx.comment.getActions().vote('down')
    })
  }
}

// 操作按钮 - 回复
function renderReply(ctx: RenderCtx) {
  if (!ctx.data.is_allow_reply) return // 不允许回复

  const replyBtn = Utils.createElement(`<span>${ctx.ctx.$t('reply')}</span>`)
  ctx.$actions.append(replyBtn)
  replyBtn.addEventListener('click', (e) => {
    e.stopPropagation() // 防止穿透
    if (!ctx.cConf.onReplyBtnClick) {
      ctx.ctx.replyComment(ctx.data, ctx.$el)
    } else {
      ctx.cConf.onReplyBtnClick()
    }
  })
}

// 操作按钮 - 折叠
function renderCollapse(ctx: RenderCtx) {
  const collapseBtn = new ActionBtn(ctx.ctx, {
    text: () => (ctx.data.is_collapsed ? ctx.ctx.$t('expand') : ctx.ctx.$t('collapse')),
    adminOnly: true
  })
  collapseBtn.appendTo(ctx.$actions)
  collapseBtn.setClick(() => {
    ctx.comment.getActions().adminEdit('collapsed', collapseBtn)
  })
}

// 操作按钮 - 审核
function renderModerator(ctx: RenderCtx) {
  const pendingBtn = new ActionBtn(ctx.ctx, {
    text: () => (ctx.data.is_pending ? ctx.ctx.$t('pending') : ctx.ctx.$t('approved')),
    adminOnly: true
  })
  pendingBtn.appendTo(ctx.$actions)
  pendingBtn.setClick(() => {
    ctx.comment.getActions().adminEdit('pending', pendingBtn)
  })
}

// 操作按钮 - 置顶
function renderPin(ctx: RenderCtx) {
  const pinnedBtn = new ActionBtn(ctx.ctx, {
    text: () => (ctx.data.is_pinned ? ctx.ctx.$t('unpin') : ctx.ctx.$t('pin')),
    adminOnly: true
  })
  pinnedBtn.appendTo(ctx.$actions)
  pinnedBtn.setClick(() => {
    ctx.comment.getActions().adminEdit('pinned', pinnedBtn)
  })
}

// 操作按钮 - 编辑
function renderEdit(ctx: RenderCtx) {
  const editBtn = new ActionBtn(ctx.ctx, {
    text: ctx.ctx.$t('edit'),
    adminOnly: true
  })
  editBtn.appendTo(ctx.$actions)
  editBtn.setClick(() => {
    ctx.ctx.editComment(ctx.data, ctx.$el)
  })
}

// 操作按钮 - 删除
function renderDel(ctx: RenderCtx) {
  const delBtn = new ActionBtn(ctx.ctx, {
    text: ctx.ctx.$t('delete'),
    confirm: true,
    confirmText: ctx.ctx.$t('deleteConfirm'),
    adminOnly: true,
  })
  delBtn.appendTo(ctx.$actions)
  delBtn.setClick(() => {
    ctx.comment.getActions().adminDelete(delBtn)
  })
}

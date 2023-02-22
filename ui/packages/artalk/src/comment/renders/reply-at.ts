import * as Utils from '../../lib/utils'
import RenderCtx from '../render-ctx'

/**
 * 层级嵌套模式显示 AT 界面
 */
export default function renderReplyAt(ctx: RenderCtx) {
  if (ctx.cConf.isFlatMode || ctx.data.rid === 0) return // not 平铺模式 或 根评论
  if (!ctx.cConf.replyTo) return

  ctx.$replyAt = Utils.createElement(`<span class="atk-item atk-reply-at"><span class="atk-arrow"></span><span class="atk-nick"></span></span>`)
  ctx.$replyAt.querySelector<HTMLElement>('.atk-nick')!.innerText = `${ctx.cConf.replyTo.nick}`
  ctx.$replyAt.onclick = () => { ctx.comment.getActions().goToReplyComment() }

  ctx.$headerBadgeWrap.insertAdjacentElement('afterend', ctx.$replyAt)
}

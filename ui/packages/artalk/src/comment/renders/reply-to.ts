import * as Utils from '../../lib/utils'
import marked from '../../lib/marked'
import RenderCtx from '../render-ctx'

/**
 * 回复的对象界面
 */
export default function renderReplyTo(ctx: RenderCtx) {
  if (!ctx.cConf.isFlatMode) return // 仅平铺模式显示
  if (!ctx.cConf.replyTo) return

  ctx.$replyTo = Utils.createElement(`
    <div class="atk-reply-to">
      <div class="atk-meta">${ctx.ctx.$t('reply')} <span class="atk-nick"></span>:</div>
      <div class="atk-content"></div>
    </div>`)
  const $nick = ctx.$replyTo.querySelector<HTMLElement>('.atk-nick')!
  $nick.innerText = `@${ctx.cConf.replyTo.nick}`
  $nick.onclick = () => { ctx.comment.getActions().goToReplyComment() }
  let replyContent = marked(ctx.ctx, ctx.cConf.replyTo.content)
  if (ctx.cConf.replyTo.is_collapsed) replyContent = `[${ctx.ctx.$t('collapsed')}]`
  ctx.$replyTo.querySelector<HTMLElement>('.atk-content')!.innerHTML = replyContent
  ctx.$body.prepend(ctx.$replyTo)
}

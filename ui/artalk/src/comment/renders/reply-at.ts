import * as Utils from '../../lib/utils'
import type Render from '../render'

/**
 * 层级嵌套模式显示 AT 界面
 */
export default function renderReplyAt(r: Render) {
  if (r.opts.flatMode || r.data.rid === 0) return // not 平铺模式 或 根评论
  if (!r.opts.replyTo) return

  r.$replyAt = Utils.createElement(`<span class="atk-item atk-reply-at"><span class="atk-arrow"></span><span class="atk-nick"></span></span>`)
  r.$replyAt.querySelector<HTMLElement>('.atk-nick')!.innerText = `${r.opts.replyTo.nick}`
  r.$replyAt.onclick = () => { r.comment.getActions().goToReplyComment() }

  r.$headerBadgeWrap.insertAdjacentElement('afterend', r.$replyAt)
}

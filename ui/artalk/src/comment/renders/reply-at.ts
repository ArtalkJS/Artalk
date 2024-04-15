import * as Utils from '../../lib/utils'
import type Render from '../render'

/**
 * Reply user indicator (with arrow icon) in the comment meta header
 */
export default function renderReplyAt(r: Render) {
  if (r.opts.flatMode || r.data.rid === 0) return // if not nested mode or root comment
  if (!r.opts.replyTo) return // if no replyTo data

  r.$replyAt = Utils.createElement(
    `<span class="atk-item atk-reply-at"><span class="atk-arrow"></span><span class="atk-nick"></span></span>`,
  )
  r.$replyAt.querySelector<HTMLElement>('.atk-nick')!.innerText = `${r.opts.replyTo.nick}`
  r.$replyAt.onclick = () => {
    r.comment.getActions().goToReplyComment()
  }

  r.$headerBadgeWrap.insertAdjacentElement('afterend', r.$replyAt)
}

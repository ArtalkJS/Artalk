import * as Utils from '@/lib/utils'
import marked from '@/lib/marked'
import $t from '@/i18n'
import type Render from '../render'

/**
 * 关联评论显示 (被回复的评论)
 */
export default function renderReplyTo(r: Render) {
  if (!r.opts.flatMode) return // 仅平铺模式显示
  if (!r.opts.replyTo) return

  r.$replyTo = Utils.createElement(`
    <div class="atk-reply-to">
      <div class="atk-meta">${$t('reply')} <span class="atk-nick"></span>:</div>
      <div class="atk-content"></div>
    </div>`)

  // Comment author name
  const $nick = r.$replyTo.querySelector<HTMLElement>('.atk-nick')!
  $nick.innerText = `@${r.opts.replyTo.nick}`
  $nick.onclick = () => {
    r.comment.getActions().goToReplyComment()
  }

  // Comment content
  let replyContent = marked(r.opts.replyTo.content)
  if (r.opts.replyTo.is_collapsed) replyContent = `[${$t('collapsed')}]`
  r.$replyTo.querySelector<HTMLElement>('.atk-content')!.innerHTML = replyContent

  // Mount the replyTo element
  r.$body.prepend(r.$replyTo)
}

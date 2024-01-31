import * as Utils from '../../lib/utils'
import marked from '../../lib/marked'
import type Render from '../render'

/**
 * 回复的对象界面
 */
export default function renderReplyTo(r: Render) {
  if (!r.cConf.isFlatMode) return // 仅平铺模式显示
  if (!r.cConf.replyTo) return

  r.$replyTo = Utils.createElement(`
    <div class="atk-reply-to">
      <div class="atk-meta">${r.ctx.$t('reply')} <span class="atk-nick"></span>:</div>
      <div class="atk-content"></div>
    </div>`)
  const $nick = r.$replyTo.querySelector<HTMLElement>('.atk-nick')!
  $nick.innerText = `@${r.cConf.replyTo.nick}`
  $nick.onclick = () => { r.comment.getActions().goToReplyComment() }
  let replyContent = marked(r.cConf.replyTo.content)
  if (r.cConf.replyTo.is_collapsed) replyContent = `[${r.ctx.$t('collapsed')}]`
  r.$replyTo.querySelector<HTMLElement>('.atk-content')!.innerHTML = replyContent
  r.$body.prepend(r.$replyTo)
}

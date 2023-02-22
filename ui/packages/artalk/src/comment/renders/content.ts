import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import RenderCtx from '../render-ctx'

/**
 * 评论内容界面
 */
export default function renderContent(ctx: RenderCtx) {
  if (!ctx.data.is_collapsed) {
    ctx.$content.innerHTML = ctx.comment.getContentMarked()
    ctx.$content.classList.remove('atk-hide', 'atk-collapsed')
    return
  }

  // 内容 & 折叠
  ctx.$content.classList.add('atk-hide', 'atk-type-collapsed')
  const collapsedInfoEl = Utils.createElement(`
    <div class="atk-collapsed">
      <span class="atk-text">${ctx.ctx.$t('collapsedMsg')}</span>
      <span class="atk-show-btn">${ctx.ctx.$t('expand')}</span>
    </div>`)
  ctx.$body.insertAdjacentElement('beforeend', collapsedInfoEl)

  const contentShowBtn = collapsedInfoEl.querySelector('.atk-show-btn')!
  contentShowBtn.addEventListener('click', (e) => {
    e.stopPropagation() // 防止穿透

    if (ctx.$content.classList.contains('atk-hide')) {
      ctx.$content.innerHTML = ctx.comment.getContentMarked()
      ctx.$content.classList.remove('atk-hide')
      Ui.playFadeInAnim(ctx.$content)
      contentShowBtn.innerHTML = ctx.ctx.$t('collapse')
    } else {
      ctx.$content.innerHTML = ''
      ctx.$content.classList.add('atk-hide')
      contentShowBtn.innerHTML = ctx.ctx.$t('expand')
    }
  })
}

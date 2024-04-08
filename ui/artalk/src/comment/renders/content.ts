import $t from '@/i18n'
import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import type Render from '../render'

/**
 * 评论内容界面
 */
export default function renderContent(r: Render) {
  if (!r.data.is_collapsed) {
    r.$content.innerHTML = r.comment.getContentMarked()
    r.$content.classList.remove('atk-hide', 'atk-collapsed')
    return
  }

  // 内容 & 折叠
  r.$content.classList.add('atk-hide', 'atk-type-collapsed')
  const collapsedInfoEl = Utils.createElement(`
    <div class="atk-collapsed">
      <span class="atk-text">${$t('collapsedMsg')}</span>
      <span class="atk-show-btn">${$t('expand')}</span>
    </div>`)
  r.$body.insertAdjacentElement('beforeend', collapsedInfoEl)

  const contentShowBtn = collapsedInfoEl.querySelector<HTMLElement>('.atk-show-btn')!
  contentShowBtn.addEventListener('click', (e) => {
    e.stopPropagation() // 防止穿透

    if (r.$content.classList.contains('atk-hide')) {
      r.$content.innerHTML = r.comment.getContentMarked()
      r.$content.classList.remove('atk-hide')
      Ui.playFadeInAnim(r.$content)
      contentShowBtn.innerText = $t('collapse')
    } else {
      r.$content.innerHTML = ''
      r.$content.classList.add('atk-hide')
      contentShowBtn.innerText = $t('expand')
    }
  })
}

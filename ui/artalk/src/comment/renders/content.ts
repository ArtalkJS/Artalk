import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import type Render from '../render'
import $t from '@/i18n'

/**
 * 评论内容界面
 */
export default function renderContent(r: Render) {
  renderAnnotationQuote(r)

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

/**
 * 划线批注的原文摘录
 *
 * 普通评论没有 anchor，不会渲染；划线评论在管理端 / 侧边栏中可看到原文上下文。
 */
function renderAnnotationQuote(r: Render) {
  const exact = r.data.anchor?.quote?.exact
  if (!exact) return

  const quoteEl = Utils.createElement(`
    <blockquote class="atk-annotation-quote">
      <span class="atk-annotation-quote-label">${$t('annotationQuote')}</span>
      <span class="atk-annotation-quote-text"></span>
      ${r.data.anchor_resolved ? `<span class="atk-annotation-resolved">${$t('annotationResolved')}</span>` : ''}
    </blockquote>`)
  quoteEl.querySelector('.atk-annotation-quote-text')!.textContent = exact
  r.$body.insertBefore(quoteEl, r.$content)
}

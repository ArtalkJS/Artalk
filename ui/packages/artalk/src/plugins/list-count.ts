import ArtalkPlugin from '~/types/plugin'
import * as Utils from '@/lib/utils'
import $t from '@/i18n'

export const ListCount: ArtalkPlugin = (ctx) => {
  const refreshCountNumEl = () => {
    const list = ctx.get('list')
    if (!list) return

    const $count = list.$el.querySelector('.atk-comment-count .atk-text')
    if (!$count) return

    const text = Utils.htmlEncode($t('counter', { count: `${Number(list.getData()?.total) || 0}` }))
    $count.innerHTML = text.replace(/(\d+)/, '<span class="atk-comment-count-num">$1</span>')
  }

  ctx.on('list-loaded', () => {
    refreshCountNumEl()
  })
}

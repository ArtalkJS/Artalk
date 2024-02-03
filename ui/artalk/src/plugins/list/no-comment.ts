import type { ArtalkPlugin } from '@/types'
import * as Utils from '@/lib/utils'
import { sanitize } from '@/lib/sanitizer'

export const NoComment: ArtalkPlugin = (ctx) => {
  ctx.on('list-loaded', (comments) => {
    const list = ctx.get('list')!

    // 无评论
    const isNoComment = comments.length <= 0
    let $noComment = list.getCommentsWrapEl().querySelector<HTMLElement>('.atk-list-no-comment')

    if (isNoComment) {
      if (!$noComment) {
        $noComment = Utils.createElement('<div class="atk-list-no-comment"></div>')

        // sanitize before set innerHTML
        $noComment.innerHTML = sanitize(list.ctx.conf.noComment || list.ctx.$t('noComment'))
        list.getCommentsWrapEl().appendChild($noComment)
      }
    } else {
      $noComment?.remove()
    }
  })
}

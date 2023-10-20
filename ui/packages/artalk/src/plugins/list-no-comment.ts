import type ArtalkPlugin from '~/types/plugin'
import * as Utils from '@/lib/utils'

export const ListNoComment: ArtalkPlugin = (ctx) => {
  ctx.on('list-loaded', () => {
    const list = ctx.get('list')!

    // 无评论
    const isNoComment = list.ctx.getCommentList().length <= 0
    let $noComment = list.getCommentsWrapEl().querySelector<HTMLElement>('.atk-list-no-comment')

    if (isNoComment) {
      if (!$noComment) {
        $noComment = Utils.createElement('<div class="atk-list-no-comment"></div>')

        // TODO POTENTIAL SECURITY RISK: prefer use insane to filter html tags before set innerHTML
        $noComment.innerHTML = list.getOptions().noCommentText || list.ctx.conf.noComment || list.ctx.$t('noComment')
        list.getCommentsWrapEl().appendChild($noComment)
      }
    } else {
      $noComment?.remove()
    }
  })
}

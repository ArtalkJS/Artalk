import * as Utils from '../../lib/utils'
import RenderCtx from '../render-ctx'

/**
 * 评论头像界面
 */
export default function renderAvatar(ctx: RenderCtx) {
  const $avatar = ctx.$el.querySelector<HTMLElement>('.atk-avatar')!
  const $avatarImg = Utils.createElement<HTMLImageElement>('<img />')

  const avatarURLBuilder = ctx.conf.avatarURLBuilder
  $avatarImg.src = avatarURLBuilder ? avatarURLBuilder(ctx.data) : ctx.comment.getGravatarURL()

  if (ctx.data.link) {
    const $avatarA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
    $avatarA.href = Utils.isValidURL(ctx.data.link) ? ctx.data.link : `https://${ctx.data.link}`
    $avatarA.append($avatarImg)
    $avatar.append($avatarA)
  } else {
    $avatar.append($avatarImg)
  }
}

import * as Utils from '../../lib/utils'
import RenderCtx from '../render-ctx'

/**
 * 评论头部界面
 */
export default function renderHeader(ctx: RenderCtx) {
  Object.entries({
    renderNick, renderVerifyBadge, renderDate, renderUABadge
  }).forEach(([name, render]) => {
    render(ctx)
  })
}

function renderNick(ctx: RenderCtx) {
  ctx.$headerNick = ctx.$el.querySelector<HTMLElement>('.atk-nick')!

  if (ctx.data.link) {
    const $nickA = Utils.createElement<HTMLLinkElement>('<a target="_blank" rel="noreferrer noopener nofollow"></a>')
    $nickA.innerText = ctx.data.nick
    $nickA.href = Utils.isValidURL(ctx.data.link) ? ctx.data.link : `https://${ctx.data.link}`
    ctx.$headerNick.append($nickA)
  } else {
    ctx.$headerNick.innerText = ctx.data.nick
  }
}

function renderVerifyBadge(ctx: RenderCtx) {
  ctx.$headerBadgeWrap = ctx.$el.querySelector<HTMLElement>('.atk-badge-wrap')!
  ctx.$headerBadgeWrap.innerHTML = ''

  const badgeText = ctx.data.badge_name
  const badgeColor = ctx.data.badge_color
  if (badgeText) {
    const $badge = Utils.createElement(`<span class="atk-badge"></span>`)
    $badge.innerText = badgeText.replace('管理员', ctx.ctx.$t('admin')) // i18n patch
    $badge.style.backgroundColor = badgeColor || ''
    ctx.$headerBadgeWrap.append($badge)
  }

  if (ctx.data.is_pinned) {
    const $pinnedBadge = Utils.createElement(`<span class="atk-pinned-badge">${ctx.ctx.$t('pin')}</span>`) // 置顶徽章
    ctx.$headerBadgeWrap.append($pinnedBadge)
  }
}

function renderDate(ctx: RenderCtx) {
  const $date = ctx.$el.querySelector<HTMLElement>('.atk-date')!
  $date.innerText = ctx.comment.getDateFormatted()
  $date.setAttribute('data-atk-comment-date', String(+new Date(ctx.data.date)))
}

function renderUABadge(ctx: RenderCtx) {
  if (!ctx.ctx.conf.uaBadge && !ctx.data.ip_region) return

  let $uaWrap = ctx.$header.querySelector('atk-ua-wrap')
  if (!$uaWrap) {
    $uaWrap = Utils.createElement(`<span class="atk-ua-wrap"></span>`)
    ctx.$header.append($uaWrap)
  }

  $uaWrap.innerHTML = ''

  if (ctx.data.ip_region) {
    const $regionBadge = Utils.createElement(`<span class="atk-region-badge"></span>`)
    $regionBadge.innerText = ctx.data.ip_region
    $uaWrap.append($regionBadge)
  }

  if (ctx.ctx.conf.uaBadge) {
    const { browser, os } = ctx.comment.getUserUA()
    if (String(browser).trim()) {
      const $uaBrowser = Utils.createElement(`<span class="atk-ua ua-browser"></span>`)
      $uaBrowser.innerText = browser
      $uaWrap.append($uaBrowser)
    }

    if (String(os).trim()) {
      const $usOS = Utils.createElement(`<span class="atk-ua ua-os"></span>`)
      $usOS.innerText = os
      $uaWrap.append($usOS)
    }
  }
}

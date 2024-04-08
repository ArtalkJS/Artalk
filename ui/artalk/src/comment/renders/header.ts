import $t from '@/i18n'
import * as Utils from '../../lib/utils'
import type Render from '../render'

/**
 * 评论头部界面
 */
export default function renderHeader(r: Render) {
  Object.entries({
    renderNick,
    renderVerifyBadge,
    renderDate,
    renderUABadge,
  }).forEach(([name, render]) => {
    render(r)
  })
}

function renderNick(r: Render) {
  r.$headerNick = r.$el.querySelector<HTMLElement>('.atk-nick')!

  if (r.data.link) {
    const $nickA = Utils.createElement<HTMLLinkElement>(
      '<a target="_blank" rel="noreferrer noopener nofollow"></a>',
    )
    $nickA.innerText = r.data.nick
    $nickA.href = Utils.isValidURL(r.data.link) ? r.data.link : `https://${r.data.link}`
    r.$headerNick.append($nickA)
  } else {
    r.$headerNick.innerText = r.data.nick
  }
}

function renderVerifyBadge(ctx: Render) {
  ctx.$headerBadgeWrap = ctx.$el.querySelector<HTMLElement>('.atk-badge-wrap')!
  ctx.$headerBadgeWrap.innerHTML = ''

  const badgeText = ctx.data.badge_name
  const badgeColor = ctx.data.badge_color
  if (badgeText) {
    const $badge = Utils.createElement(`<span class="atk-badge"></span>`)
    $badge.innerText = badgeText.replace('管理员', $t('admin')) // i18n patch
    $badge.style.backgroundColor = badgeColor || ''
    ctx.$headerBadgeWrap.append($badge)
  } else if (ctx.data.is_verified) {
    const $verifiedBadge = Utils.createElement(
      `<span class="atk-verified-icon" title="${$t('emailVerified')}"></span>`,
    ) // 邮箱验证徽章
    ctx.$headerBadgeWrap.append($verifiedBadge)
  }

  if (ctx.data.is_pinned) {
    const $pinnedBadge = Utils.createElement(`<span class="atk-pinned-badge">${$t('pin')}</span>`) // 置顶徽章
    ctx.$headerBadgeWrap.append($pinnedBadge)
  }
}

function renderDate(ctx: Render) {
  const $date = ctx.$el.querySelector<HTMLElement>('.atk-date')!
  $date.innerText = ctx.comment.getDateFormatted()
  $date.setAttribute('data-atk-comment-date', String(+new Date(ctx.data.date)))
}

function renderUABadge(ctx: Render) {
  if (!ctx.opts.uaBadge && !ctx.data.ip_region) return

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

  if (ctx.opts.uaBadge) {
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

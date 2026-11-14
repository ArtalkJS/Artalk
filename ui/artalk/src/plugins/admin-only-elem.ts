import type { ArtalkPlugin } from '@/types'

export const AdminOnlyElem: ArtalkPlugin = (ctx) => {
  const user = ctx.inject('user')
  const layers = ctx.inject('layers')

  const scanApply = () => {
    applyAdminOnlyEls(
      user.getData().is_admin,
      getAdminOnlyEls({
        $root: ctx.getEl(),
        $list: ctx.inject('list').getEl(),
        $layers: layers.getEl(),
      }),
    )
  }

  ctx.on('list-loaded', () => {
    scanApply()
  })

  ctx.on('user-changed', (user) => {
    scanApply()
  })
}

function getAdminOnlyEls(opts: {
  $root: HTMLElement
  $list: HTMLElement
  $layers: HTMLElement
}): HTMLElement[] {
  const els = new Set<HTMLElement>()
  // The sidebar can move the list outside the widget root. All three roots
  // still belong to this instance, including editor controls inside layers.
  for (const root of [opts.$root, opts.$list, opts.$layers]) {
    root.querySelectorAll<HTMLElement>('[atk-only-admin-show]').forEach((el) => els.add(el))
  }
  return [...els]
}

function applyAdminOnlyEls(isAdmin: boolean, els: HTMLElement[]) {
  els.forEach(($item: HTMLElement) => {
    if (isAdmin) $item.classList.remove('atk-hide')
    else $item.classList.add('atk-hide')
  })
}

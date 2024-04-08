import type { ArtalkPlugin } from '@/types'

export const AdminOnlyElem: ArtalkPlugin = (ctx) => {
  const scanApply = () => {
    applyAdminOnlyEls(
      ctx.get('user').getData().isAdmin,
      getAdminOnlyEls({
        $root: ctx.$root,
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

function getAdminOnlyEls(opts: { $root: HTMLElement }): HTMLElement[] {
  const els: HTMLElement[] = []

  // elements in $root
  opts.$root
    .querySelectorAll<HTMLElement>(`[atk-only-admin-show]`)
    .forEach((item) => els.push(item))

  // TODO: provide a Artalk.conf hook to set whitelist of admin-only elements,
  // and move following code to that hook (move into @artalk/artalk-sidebar)

  // elements in sidebar
  const $sidebarEl = document.querySelector<HTMLElement>('.atk-sidebar')
  if ($sidebarEl)
    $sidebarEl
      .querySelectorAll<HTMLElement>(`[atk-only-admin-show]`)
      .forEach((item) => els.push(item))

  return els
}

function applyAdminOnlyEls(isAdmin: boolean, els: HTMLElement[]) {
  els.forEach(($item: HTMLElement) => {
    if (isAdmin) $item.classList.remove('atk-hide')
    else $item.classList.add('atk-hide')
  })
}

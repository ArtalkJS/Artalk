import type { ArtalkPlugin } from '@/types'
import $t from '@/i18n'

export const SidebarBtn: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  let $openSidebarBtn: HTMLElement | null = null

  const syncByUser = () => {
    if (!$openSidebarBtn) return
    const user = ctx.inject('user').getData()

    // 已输入个人信息
    if (!!user.name && !!user.email) {
      $openSidebarBtn.classList.remove('atk-hide')

      // update button text (normal user or admin)
      const $btnText = $openSidebarBtn.querySelector<HTMLElement>('.atk-text')
      if ($btnText) $btnText.innerText = !user.is_admin ? $t('msgCenter') : $t('ctrlCenter')
    } else {
      $openSidebarBtn.classList.add('atk-hide')
    }
  }

  ctx.watchConf(['locale'], (conf) => {
    $openSidebarBtn = list.getEl().querySelector<HTMLElement>('[data-action="open-sidebar"]')
    if (!$openSidebarBtn) return

    $openSidebarBtn.onclick = () => {
      // use onclick rather than addEventListener to prevent duplicate event
      ctx.showSidebar()
    }

    syncByUser()
  })

  ctx.on('user-changed', (user) => {
    syncByUser()
  })
}

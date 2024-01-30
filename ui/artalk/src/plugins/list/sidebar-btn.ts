import type { ArtalkPlugin } from '@/types'
import $t from '@/i18n'

export const SidebarBtn: ArtalkPlugin = (ctx) => {
  let $openSidebarBtn: HTMLElement|null = null

  const syncByUser = () => {
    if (!$openSidebarBtn) return
    const user = ctx.get('user').getData()

    // 已输入个人信息
    if (!!user.nick && !!user.email) {
      $openSidebarBtn.classList.remove('atk-hide')

      // update button text (normal user or admin)
      const $btnText = $openSidebarBtn.querySelector<HTMLElement>('.atk-text')
      if ($btnText) $btnText.innerText = (!user.isAdmin) ? $t('msgCenter') : $t('ctrlCenter')
    } else {
      $openSidebarBtn.classList.add('atk-hide')
    }
  }

  ctx.watchConf(['locale'], (conf) => {
    const list = ctx.get('list')

    $openSidebarBtn = list.$el.querySelector<HTMLElement>('[data-action="open-sidebar"]')
    if (!$openSidebarBtn) return

    $openSidebarBtn.onclick = () => { // use onclick rather than addEventListener to prevent duplicate event
      ctx.showSidebar()
    }

    syncByUser()
  })

  ctx.on('user-changed', (user) => {
    syncByUser()
  })
}

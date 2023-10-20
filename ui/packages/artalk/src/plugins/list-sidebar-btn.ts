import ArtalkPlugin from '~/types/plugin'
import $t from '@/i18n'

export const ListSidebarBtn: ArtalkPlugin = (ctx) => {
  let $openSidebarBtn: HTMLElement|null = null

  ctx.on('conf-loaded', () => {
    const list = ctx.get('list')
    if (!list) return

    $openSidebarBtn = list.$el.querySelector<HTMLElement>('[data-action="open-sidebar"]')
    if (!$openSidebarBtn) return

    $openSidebarBtn.onclick = () => { // use onclick rather than addEventListener to prevent duplicate event
      ctx.showSidebar()
    }
  })

  ctx.on('user-changed', (user) => {
    if ($openSidebarBtn) {

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
  })
}

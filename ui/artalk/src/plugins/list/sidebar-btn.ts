import type { ArtalkPlugin } from '@/types'
import $t from '@/i18n'

export const SidebarBtn: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')
  const conf = ctx.inject('config')
  let $openSidebarBtn: HTMLElement | null = null

  const syncByUser = () => {
    if (!$openSidebarBtn) return
    const user = ctx.inject('user').getData()

    // 已输入个人信息 且 通知中心未关闭 (管理员始终保留控制中心入口)
    if (!!user.name && !!user.email && (conf.get().notifyCenter || user.is_admin)) {
      $openSidebarBtn.classList.remove('atk-hide')

      // update button text (normal user or admin)
      const $btnText = $openSidebarBtn.querySelector<HTMLElement>('.atk-text')
      if ($btnText) $btnText.innerText = !user.is_admin ? $t('msgCenter') : $t('ctrlCenter')
    } else {
      $openSidebarBtn.classList.add('atk-hide')
    }

    syncHeaderVisibility(user.is_admin)
  }

  // 当评论排序与通知中心同时关闭，且非管理员时，隐藏整个 list-header 容器
  const syncHeaderVisibility = (isAdmin: boolean) => {
    const $header = list.getEl().querySelector<HTMLElement>('.atk-list-header')
    if (!$header) return
    const c = conf.get()
    const shouldHide = !c.listSort && !c.notifyCenter && !isAdmin
    $header.classList.toggle('atk-hide', shouldHide)
  }

  ctx.watchConf(['locale', 'notifyCenter', 'listSort'], (conf) => {
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

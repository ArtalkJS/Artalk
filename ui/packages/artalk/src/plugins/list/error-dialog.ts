import type { ArtalkPlugin, ContextApi } from '~/types'
import * as Utils from '../../lib/utils'
import * as Ui from '../../lib/ui'
import User from '../../lib/user'
import $t from '../../i18n'

export const ErrorDialog: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', () => {
    const list = ctx.get('list')
    if (!list) return

    // clear the original error when a new fetch is triggered
    Ui.setError(list.$el, null)
  })

  ctx.on('list-error', (err) => {
    const list = ctx.get('list')
    if (!list) return
    Ui.setError(list.$el, renderErrorDialog(ctx, err.msg, err.data))
  })
}

export function renderErrorDialog(ctx: ContextApi, errMsg: string, errData?: any): HTMLElement {
  const errEl = Utils.createElement(`<span>${errMsg}，${$t('listLoadFailMsg')}<br/></span>`)

  const $retryBtn = Utils.createElement(`<span style="cursor:pointer;">${$t('listRetry')}</span>`)
  $retryBtn.onclick = () => (ctx.fetch({ offset: 0 }))
  errEl.appendChild($retryBtn)

  const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
  errEl.appendChild(adminBtn)
  if (!User.data.isAdmin) adminBtn.classList.add('atk-hide')

  let sidebarView = ''

  // 找不到站点错误，打开侧边栏并填入创建站点表单
  if (errData?.err_no_site) {
    const viewLoadParam = {
      create_name: ctx.conf.site,
      create_urls: `${window.location.protocol}//${window.location.host}`
    }
    sidebarView = `sites|${JSON.stringify(viewLoadParam)}`
  }

  adminBtn.onclick = () => ctx.showSidebar({
    view: sidebarView as any
  })

  return errEl
}

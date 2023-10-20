import type ArtalkPlugin from '~/types/plugin'
import List from '@/list/list'
import * as Utils from '../lib/utils'
import * as Ui from '../lib/ui'
import User from '../lib/user'
import $t from '../i18n'

export const ListErrorDialog: ArtalkPlugin = (ctx) => {
  ctx.on('list-error', (err) => {
    const list = ctx.get('list')!
    Ui.setError(list.$el, renderErrorDialog(list, err.msg, err.data))
  })
}

export function renderErrorDialog(list: List, errMsg: string, errData?: any): HTMLElement {
  const errEl = Utils.createElement(`<span>${errMsg}，${$t('listLoadFailMsg')}<br/></span>`)

  const $retryBtn = Utils.createElement(`<span style="cursor:pointer;">${$t('listRetry')}</span>`)
  $retryBtn.onclick = () => (list.fetchComments(0))
  errEl.appendChild($retryBtn)

  const adminBtn = Utils.createElement('<span atk-only-admin-show> | <span style="cursor:pointer;">打开控制台</span></span>')
  errEl.appendChild(adminBtn)
  if (!User.data.isAdmin) adminBtn.classList.add('atk-hide')

  let sidebarView = ''

  // 找不到站点错误，打开侧边栏并填入创建站点表单
  if (errData?.err_no_site) {
    const viewLoadParam = {
      create_name: list.ctx.conf.site,
      create_urls: `${window.location.protocol}//${window.location.host}`
    }
    sidebarView = `sites|${JSON.stringify(viewLoadParam)}`
  }

  adminBtn.onclick = () => list.ctx.showSidebar({
    view: sidebarView as any
  })

  return errEl
}

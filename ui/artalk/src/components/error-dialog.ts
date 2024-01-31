import * as Utils from '@/lib/utils'
import * as Ui from '@/lib/ui'
import $t from '@/i18n'

export interface ErrorDialogOptions {
  $err: HTMLElement

  errMsg: string
  errData?: any
  retryFn?: () => void

  onOpenSidebar?: () => void
}

export function showErrorDialog(opts: ErrorDialogOptions) {
  const errEl = Utils.createElement(`<span><span class="error-message"></span><br/><br/></span>`)
  errEl.querySelector<HTMLElement>('.error-message')!.innerText = `${$t('listLoadFailMsg')}\n${opts.errMsg}`

  if (opts.retryFn) {
    const $retryBtn = Utils.createElement(`<span style="cursor:pointer;">${$t('listRetry')}</span>`)
    $retryBtn.onclick = () => opts.retryFn && opts.retryFn()
    errEl.appendChild($retryBtn)
  }

  if (opts.onOpenSidebar) {
    const $openSidebar = Utils.createElement(`<span atk-only-admin-show> | <span style="cursor:pointer;">${$t('openName', { name: $t('ctrlCenter') })}</span></span>`)
    errEl.appendChild($openSidebar)
    $openSidebar.onclick = () => opts.onOpenSidebar && opts.onOpenSidebar()
  }

  Ui.setError(opts.$err, errEl)
}

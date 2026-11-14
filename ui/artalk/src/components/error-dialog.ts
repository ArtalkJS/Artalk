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
  const errEl = Utils.createElement(
    `<span class="atk-error-dialog" role="alert">
      <span class="atk-error-heading"></span>
      <span class="error-message"></span>
      <span class="atk-error-actions"></span>
    </span>`,
  )
  errEl.querySelector<HTMLElement>('.atk-error-heading')!.textContent = $t('listLoadFailMsg')
  errEl.querySelector<HTMLElement>('.error-message')!.textContent = opts.errMsg
  const actions = errEl.querySelector<HTMLElement>('.atk-error-actions')!

  if (opts.retryFn) {
    const $retryBtn = Utils.createElement('<button type="button" class="atk-error-retry"></button>')
    $retryBtn.textContent = $t('listRetry')
    $retryBtn.onclick = () => opts.retryFn?.()
    actions.appendChild($retryBtn)
  }

  if (opts.onOpenSidebar) {
    const $openSidebar = Utils.createElement('<button type="button" atk-only-admin-show></button>')
    $openSidebar.textContent = $t('openName', { name: $t('ctrlCenter') })
    $openSidebar.onclick = () => opts.onOpenSidebar?.()
    actions.appendChild($openSidebar)
  }

  Ui.setError(opts.$err, errEl, '')
}

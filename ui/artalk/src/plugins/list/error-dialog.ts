import * as Ui from '../../lib/ui'
import type { ArtalkPlugin } from '@/types'
import { showErrorDialog } from '@/components/error-dialog'

export const ErrorDialog: ArtalkPlugin = (ctx) => {
  const list = ctx.inject('list')

  ctx.on('list-fetch', () => {
    // clear the original error when a new fetch is triggered
    Ui.setError(list.getEl(), null)
  })

  ctx.on('list-failed', (err) => {
    showErrorDialog({
      $err: list.getEl(),
      errMsg: err.msg,
      errData: err.data,
      retryFn: () => ctx.fetch({ offset: 0 }),
    })
  })
}

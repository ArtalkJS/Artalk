import type { ArtalkPlugin, ContextApi } from '@/types'
import { showErrorDialog } from '@/components/error-dialog'
import * as Ui from '../../lib/ui'

export const ErrorDialog: ArtalkPlugin = (ctx) => {
  ctx.on('list-fetch', () => {
    const list = ctx.get('list')

    // clear the original error when a new fetch is triggered
    Ui.setError(list.$el, null)
  })

  ctx.on('list-failed', (err) => {
    showErrorDialog({
      $err: ctx.get('list').$el,
      errMsg: err.msg,
      errData: err.data,
      retryFn: () => ctx.fetch({ offset: 0 }),
    })
  })
}

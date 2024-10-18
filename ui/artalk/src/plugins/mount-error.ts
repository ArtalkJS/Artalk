import { showErrorDialog } from '@/components/error-dialog'
import type { ArtalkPlugin } from '@/types'

export interface MountErrorOptions {
  err?: { data?: { err_no_site?: boolean }; msg: string }
  onRetry?: () => void
}

export const MountError: ArtalkPlugin<MountErrorOptions> = (ctx, opts = {}) => {
  const list = ctx.inject('list')
  const user = ctx.inject('user')
  const conf = ctx.inject('config')

  const err = opts.err
  if (!err) throw new Error('MountError: `err` is required')

  let sidebarOpenView = ''

  // if response err_no_site, modify the sidebar open view to create site
  if (err.data?.err_no_site) {
    const viewLoadParam = {
      create_name: conf.get().site,
      create_urls: `${window.location.protocol}//${window.location.host}`,
    }
    sidebarOpenView = `sites|${JSON.stringify(viewLoadParam)}`
  }

  showErrorDialog({
    $err: list.getEl(),
    errMsg: err.msg || String(err),
    errData: err.data,
    retryFn: () => opts.onRetry?.(),
    onOpenSidebar: user.getData().is_admin
      ? () =>
          ctx.showSidebar({
            view: sidebarOpenView as any,
          })
      : undefined, // only show open sidebar button when user is admin
  })
}

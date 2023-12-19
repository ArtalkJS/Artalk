import type { ArtalkConfig, ArtalkPlugin, ContextApi } from '~/types'
import { handleConfFormServer } from '@/config'
import { showErrorDialog } from '../components/error-dialog'

export const ConfRemoter: ArtalkPlugin = (ctx) => {
  ctx.on('inited', () => {
    loadConf(ctx)
  })
}

function loadConf(ctx: ContextApi) {
  ctx
    .getApi()
    .system.conf()
    .then((data) => {
      let conf: Partial<ArtalkConfig> = {
        apiVersion: data.version.version, // version info
      }

      // reference conf from backend
      if (ctx.conf.useBackendConf) {
        if (!data.frontend_conf)
          throw new Error(
            'The remote backend does not respond to the frontend conf, but `useBackendConf` conf is enabled',
          )
        conf = { ...conf, ...handleConfFormServer(data.frontend_conf) }
      }

      // apply conf modifier
      ctx.conf.remoteConfModifier && ctx.conf.remoteConfModifier(conf)

      ctx.updateConf(conf)
    })
    .catch((err) => {
      ctx.updateConf({})

      let sidebarOpenView = ''

      // if response err_no_site, modify the sidebar open view to create site
      if (err.data?.err_no_site) {
        const viewLoadParam = {
          create_name: ctx.conf.site,
          create_urls: `${window.location.protocol}//${window.location.host}`,
        }
        sidebarOpenView = `sites|${JSON.stringify(viewLoadParam)}`
      }

      showErrorDialog({
        $err: ctx.get('list').$el,
        errMsg: err.msg || String(err),
        errData: err.data,
        retryFn: () => loadConf(ctx),
        onOpenSidebar: () =>
          ctx.get('user').getData().isAdmin
            ? ctx.showSidebar({
                view: sidebarOpenView as any,
              })
            : undefined, // only show open sidebar button when user is admin
      })

      throw err
    })
    .then(() => {
      // 评论获取
      if (ctx.conf.remoteConfModifier) return // only auto fetch when no remoteConfModifier
      ctx.fetch({ offset: 0 })
    })
    .catch(() => {})
}

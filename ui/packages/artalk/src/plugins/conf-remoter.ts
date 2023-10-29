import type { ArtalkConfig, ArtalkPlugin, ContextApi } from '~/types'
import { handleConfFormServer } from '@/config'
import { showErrorDialog } from '../components/error-dialog'

let confLoaded = false

export const ConfRemoter: ArtalkPlugin = (ctx) => {
  ctx.on('inited', () => {
    if (!confLoaded) loadConf(ctx)
  })
}

function loadConf(ctx: ContextApi) {
  ctx.getApi().system.conf().then((data) => {
    let conf: Partial<ArtalkConfig> = {
      apiVersion: data.version.version, // version info
    }

    // reference conf from backend
    if (ctx.conf.useBackendConf) {
      if (!data.frontend_conf) throw new Error('The remote backend does not respond to the frontend conf, but `useBackendConf` conf is enabled')
      conf = { ...conf, ...handleConfFormServer(data.frontend_conf) }
    }

    ctx.updateConf(conf)
    confLoaded = true
  }).catch((err) => {
    ctx.updateConf({})

    showErrorDialog(ctx, err.msg || String(err), err.data, () => {
      loadConf(ctx)
    })

    throw err
  }).then(() => {
    // 评论获取
    ctx.fetch({ offset: 0 })
  }).catch(() => {})
}
